package api

import (
	"FiberFinanceAPI/auth"
	model "FiberFinanceAPI/database/models"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var status int

var (
	ErrUserExist   = errors.New("user with email exists")
	userNotFound   = errors.New("user does not exist or deleted")
	userDeletedMSG = "user successfully deleted at %s"
)

// createUserRequest the required credentials to create a user
type createUserRequest struct {
	model.SessionDeviceID
	Email    string `json:"email"  validate:"required,max=155,email"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}

// createUser request to be stored in our database
func (s *Server) createUser(ctx *fiber.Ctx) error {
	// show function name to track errors faster
	s.logs.WithField("func", "users_api.go -> createUser()").Debug()
	// req load parameters
	var req createUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		s.logs.WithError(err).Warn("could not decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}
	s.logs.WithFields(logrus.Fields{
		"email": req.Email,
	}).Debug()
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	args := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashPassword,
	}
	user, err := s.repo.CreateUser(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			s.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres codes")
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, ErrUserExist))
			}
		}
		s.logs.WithError(err).Warn(err.Error())
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))

	}

	s.logs.WithField("message", "created user successfully").Info("Successful")
	return ctx.Status(http.StatusCreated).JSON(user)
}

type loginUserRequest struct {
	model.SessionDeviceID
	Email    string `json:"email" validate:"required,email,max=155"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}
type ResponseTokens struct {
	Token auth.TokenAccess `json:"token"`
	User  model.User       `json:"user"`
}

// loginUser login our user if email and password are a match
func (s *Server) loginUser(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_api.go -> loginUser()").Debug()
	var req loginUserRequest
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
	user, err := s.repo.GetUserByEmail(ctx.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	// Check Password if matches and is correct
	if err = utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		s.logs.WithError(err).Warn("Invalid password provided by user")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	accessToken, exp, refreshToken, rexp, err := s.tokenCredentials(user)
	if err != nil {
		s.logs.WithError(err).Warn(err.Error())
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	args := db.SaveRefreshTokenParams{
		UserID:       user.ID,
		DeviceID:     req.DeviceID, // value shall be passed currently random string
		RefreshToken: refreshToken,
		ExpiresAt:    rexp,
	}

	if err = s.repo.SaveRefreshToken(ctx.Context(), args); err != nil {
		s.logs.WithError(err).Warn("unable to save refresh token")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(err)
	}

	resp := ResponseTokens{
		Token: auth.TokenAccess{
			AccessToken:           accessToken,
			RefreshToken:          refreshToken,
			AccessTokenExpiresAt:  exp,
			RefreshTokenExpiresAt: rexp,
		},
		User: user,
	}
	s.logs.WithField("user_id", user.ID).Debug("user logged in")
	return ctx.Status(http.StatusOK).JSON(resp)
}

func (s *Server) getUserByID(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_api.go -> getUserByID()").Debug()
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}
	user, err := s.repo.GetUserByID(ctx.Context(), model.UserID(userID))
	if err != nil {
		s.logs.WithError(err).Warn()
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Debug("user returned successfully")
	return ctx.Status(http.StatusOK).JSON(user)
}

type listUsersRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (s *Server) listUsers(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_api.go -> listUsers()").Debug()
	var req listUsersRequest
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
	args := db.ListUserParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	users, err := s.repo.ListUsers(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(users) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, userNotFound))
	}
	s.logs.Info("users returned successfully")

	return ctx.Status(http.StatusOK).JSON(users)
}

type changePasswordRequest struct {
	Password string `json:"password" validate:"required"`
}

func (s *Server) changePassword(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_api.go -> changePassword()").Debug()
	var req changePasswordRequest
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
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		s.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	args := db.UpdatePasswordParams{
		UserID:       model.UserID(userID),
		HashPassword: hashPassword,
	}

	changedAt, err := s.repo.UpdatePassword(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Debug("password successfully changed")
	changeAtMSG := fiber.Map{"Message": fmt.Sprintf("Password successfully changed at %s", changedAt.Format(time.ANSIC))}
	return ctx.Status(http.StatusOK).JSON(changeAtMSG)
}

func (s *Server) deleteUser(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "users_api.go -> deleteUser()").Debug()
	userID := ctx.Params("userID")
	if userID == "" {
		s.logs.WithField("userID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("userID not provided")))
	}

	deletedAt, err := s.repo.DeleteUser(ctx.Context(), model.UserID(userID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("user deleted successfully")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf(userDeletedMSG, deletedAt.Format(time.ANSIC))})
}
