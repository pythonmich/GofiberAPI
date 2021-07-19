package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"context"
	"errors"
	"github.com/bluele/gcache"
	"github.com/gofiber/fiber/v2"
	_ "github.com/gofiber/fiber/v2/middleware/cache"
	"net/http"
	"time"
)

type permissionInter interface {
	wrap(permissionType ...permissionType) fiber.Handler
	check(ctx *fiber.Ctx, permissionType ...permissionType) bool
}

type permission struct {
	db.Repo
	cache gcache.Cache
	logs  *utils.StandardLogger
}

func (p permission) withRoles(payload *auth.AccessPayload, roleFunc func(role *db.UserRole) bool) (bool, error) {
	if payload.SUB != "" {
		return false, errors.New("userID not provided")
	}
	role, err := p.getRole(db.UserID(payload.SUB))
	if err != nil {
		return false, err
	}
	return roleFunc(role), nil
}

func (p *permission) check(ctx *fiber.Ctx, permissionType ...permissionType) bool {
	p.logs.WithField("func", "permissions.go -> check()").Debug()
	payload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	for _, permission := range permissionType {
		switch permission {
		case admin:
			if allowed, _ := p.withRoles(payload, adminOnly); allowed {
				return true
			}
		case member:
			if allowed := memberOnly(payload); allowed {
				return true
			}
		case memberIsTarget:
			userID := ctx.Params("userID")
			if userID == "" {
				p.logs.WithField("userID", "not provided").Debug()
				status = http.StatusBadRequest
				return false
			}
			if allowed := memberIsTargetOnly(db.UserID(userID), payload); allowed {
				return true
			}
		}
	}
	return true
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
			userID := key.(db.UserID)
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
func (p *permission) getRole(id db.UserID) (*db.UserRole, error) {
	role, err := p.cache.Get(id)
	if err != nil {
		p.logs.WithError(err).Warn()
		return &db.UserRole{}, err
	}
	return role.(*db.UserRole), nil
}

func (p *permission) wrap(permissionType ...permissionType) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		p.logs.WithField("func", "permissions.go -> wrap()").Debug()

		if allowed := p.check(ctx, permissionType...); !allowed {
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("user unauthorized, permission denied")))
		}

		//if payload.SUB != userID {
		//	err := errors.New("id does not match")
		//	status = http.StatusUnauthorized
		//	return ctx.Status(status).JSON(errorResponse(status, err))
		//}
		//p.logs.Debug(userID, payload.SUB)

		return ctx.Next()
	}
}
