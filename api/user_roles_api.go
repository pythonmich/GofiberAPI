package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

type roleRequest struct {
	db.UserRole
}

func (server *Server) grantRole(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_role_api.go -> grantRole()").Debug()
	var req roleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		server.logs.WithField("userId", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	authPayload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	if authPayload.SUB != userID {
		err := errors.New("id does not match")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if err := ctx.BodyParser(&req); err != nil {
		server.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := validateStruct(&req, server.logs); len(errs) > 0 {
		server.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs, server.logs))
	}

	args := db.RoleParams{
		UserID: db.UserID(userID),
		Role:   req.Role,
	}
	role, err := server.repo.GrantRole(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			server.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres codes")
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, errors.New("role already allocated to user")))
			}
		}
		server.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	server.logs.WithField("message", "successful").Debug(fmt.Sprintf("role granted for user %s", role.UserID))
	return ctx.Status(http.StatusCreated).JSON(role)
}

func (server *Server) revokeRole(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_role_api.go -> revokeRole()").Debug()
	var req roleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		server.logs.WithField("userId", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	authPayload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	if authPayload.SUB != userID {
		err := errors.New("id does not match")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if err := ctx.BodyParser(&req); err != nil {
		server.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := validateStruct(&req, server.logs); len(errs) > 0 {
		server.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs, server.logs))
	}

	args := db.RoleParams{
		UserID: db.UserID(userID),
		Role:   req.Role,
	}
	err := server.repo.RevokeRole(ctx.Context(), args)
	if err != nil {
		server.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "role revoked successfully"})
}

func (server *Server) getUserRole(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_role_api.go -> getUserRole()").Debug()
	userID := ctx.Params("userID")
	if userID == "" {
		server.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	authPayload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	if authPayload.SUB != userID {
		err := errors.New("id does not match")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	role, err := server.repo.GetUserRoleByID(ctx.Context(), db.UserID(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			server.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, errors.New("user does not have role")))
		}
		server.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	server.logs.WithField("message", "successful").Debug(fmt.Sprintf("role for user %s", role.UserID))
	return ctx.Status(http.StatusOK).JSON(role)
}

type listRoleRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (server *Server) listRoles(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_role_api.go -> listRoles()").Debug()
	var req listRoleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		server.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	if err := ctx.QueryParser(&req); err != nil {
		server.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := validateStruct(&req, server.logs); len(errs) > 0 {
		server.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs, server.logs))
	}
	authPayload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	server.logs.WithFields(logrus.Fields{"limit": req.PageSize, "offset": (req.PageID - 1) * req.PageSize}).Debug()
	args := db.ListUserRoleParams{
		UserID: db.UserID(authPayload.SUB),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := server.repo.ListUsersByRole(ctx.Context(), args)
	if err != nil {
		server.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(err)
	}
	server.logs.Debug("roles successfully returned")
	return ctx.Status(http.StatusOK).JSON(roles)
}
