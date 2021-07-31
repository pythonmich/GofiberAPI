package api

import (
	model "FiberFinanceAPI/database/models"
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
	model.UserRole
}

var roleNotFound = errors.New("user does not have any role(s)")

func (s *Server) grantRole(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_role_api.go -> grantRole()").Debug()
	var req roleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}

	if err := ctx.BodyParser(&req); err != nil {
		s.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}

	args := db.RoleParams{
		UserID: model.UserID(userID),
		Role:   req.Role,
	}
	role, err := s.repo.GrantRole(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			s.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres error codes")
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, errors.New("role already allocated to user")))
			}
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	s.logs.WithField("message", "successful").Debug(fmt.Sprintf("role granted for user %s", role.UserID))
	return ctx.Status(http.StatusCreated).JSON(role)
}

func (s *Server) revokeRole(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_role_api.go -> revokeRole()").Debug()
	var req roleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	if err := ctx.BodyParser(&req); err != nil {
		s.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}

	args := db.RoleParams{
		UserID: model.UserID(userID),
		Role:   req.Role,
	}
	err := s.repo.RevokeRole(ctx.Context(), args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, roleNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": "role revoked successfully"})
}

func (s *Server) getUserRole(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_role_api.go -> getUserRole()").Debug()
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}

	role, err := s.repo.GetUserRoleByID(ctx.Context(), model.UserID(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, roleNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.WithField("message", "successful").Debug(fmt.Sprintf("role for user %s", role.UserID))
	return ctx.Status(http.StatusOK).JSON(role)
}

type listRoleRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (s *Server) listRoles(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_role_api.go -> listRoles()").Debug()
	var req listRoleRequest
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	if err := ctx.QueryParser(&req); err != nil {
		s.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}
	s.logs.WithFields(logrus.Fields{"limit": req.PageSize, "offset": (req.PageID - 1) * req.PageSize}).Debug()
	args := db.ListUserRoleParams{
		UserID: model.UserID(userID),
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	roles, err := s.repo.ListUsersByRole(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(roles) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, roleNotFound))
	}
	s.logs.Debug("roles successfully returned")
	return ctx.Status(http.StatusOK).JSON(roles)
}
