package api

import (
	"FiberFinanceAPI/auth"
	model "FiberFinanceAPI/database/models"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"context"
	"errors"
	"github.com/bluele/gcache"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

/*
	Logic behind this is to have an admin that can grant roles to other and also revoke them.
	A user with an admin role will have freedom in our system as well
*/

type permissionInter interface {
	wrap(permissionType ...permissionType) fiber.Handler
	check(ctx *fiber.Ctx, permissionType ...permissionType) bool
}

type permission struct {
	db.Repo
	cache gcache.Cache
	logs  *utils.StandardLogger
}

func (p *permission) withRoles(payload *auth.AccessPayload, roleFunc func(role model.UserRole) bool) (bool, error) {
	if payload.SUB == "" {
		p.logs.Debug("userID not provided")
		return false, errors.New("userID not provided")
	}
	role, err := p.getRole(model.UserID(payload.SUB))
	p.logs.WithField("role", role).Debug("user role")
	if err != nil {
		p.logs.Warn(err)
		return false, err
	}
	p.logs.Debug("role returned successfully")
	return roleFunc(role), nil
}

func (p *permission) check(ctx *fiber.Ctx, permissionType ...permissionType) bool {
	p.logs.WithField("func", "permissions.go -> check()").Debug()
	for _, permission := range permissionType {
		p.logs.WithField("permission", permission).Debug()
		switch permission {
		case admin:
			payload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
			if allowed, _ := p.withRoles(payload, adminOnly); allowed {
				return true
			}
		case member:
			payload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
			if allowed := memberOnly(payload); allowed {
				return true
			}
		case memberIsTarget:
			payload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
			userID := ctx.Params("userID")
			//TODO: Store our userID in a context so as to make it available each time during requests to avoid repetition when getting userID
			if allowed := memberIsTargetOnly(ctx, model.UserID(userID), payload); allowed {
				return true
			}
		case prospect:
			if allowed := prospects(); allowed {
				return true
			}
		}
	}
	return false
}

func newPermissions(repo db.Repo, logs *utils.StandardLogger) permissionInter {
	p := &permission{
		Repo: repo,
		logs: logs,
	}
	p.cache = gcache.New(200).
		//caching scheme is to remove the least recently used frame when the cache is full
		//and a new page is referenced which is not there in cache
		LRU().
		LoaderExpireFunc(func(key interface{}) (interface{}, *time.Duration, error) {
			userID := key.(model.UserID)
			role, err := repo.GetUserRoleByID(context.Background(), userID)
			if err != nil {
				return nil, nil, err
			}
			expires := 1 * time.Minute
			return role, &expires, nil
		}).
		Build()
	return p
}

// we need a function to call to get user's from cache (if we want to have roles in cache it will get it from database
func (p *permission) getRole(id model.UserID) (model.UserRole, error) {
	p.logs.WithField("func", "permissions.go -> getRole()").Debug()
	role, err := p.cache.Get(id)
	p.logs.WithField("role", role).Debug()
	if err != nil {
		p.logs.WithError(err).Warn()
		return model.UserRole{}, err
	}
	p.logs.Debug("role returned successfully")
	return role.(model.UserRole), nil
}

func (p *permission) wrap(permissionType ...permissionType) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		p.logs.WithField("func", "permissions.go -> wrap()").Debug()
		if allowed := p.check(ctx, permissionType...); !allowed {
			p.logs.Debug("permission denied")
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("user unauthorized, permission denied")))
		}
		return ctx.Next()
	}
}
