package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
)

var status int

var (
	ErrUserExist = errors.New("user with email exists")
	userNotFound = errors.New("user does not exist")
)

// createUserRequest the required credentials to create a user
type createUserRequest struct {
	db.Session
	Email    string `json:"email"  validate:"required,max=155,email"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}

// createUser request to be stored in our database
func (server *Server) createUser(ctx *fiber.Ctx) error {
	// show function name to track errors faster
	server.logs.WithField("func", "users_api.go -> createUser()").Debug()
	// req load parameters
	var req createUserRequest
	if err := ctx.BodyParser(&req); err != nil {
		server.logs.WithError(err).Warn("could not decode parameters")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := validateStruct(&req, server.logs); len(errs) > 0 {
		server.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs, server.logs))
	}
	server.logs.WithFields(logrus.Fields{
		"email": req.Email,
	}).Debug()
	hashPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		server.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	args := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashPassword,
	}
	user, err := server.repo.CreateUser(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			server.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres codes")
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, ErrUserExist))
			}
		}
		server.logs.WithError(err).Warn(err.Error())
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))

	}

	server.logs.WithField("message", "created user successfully").Info("Successful")
	return ctx.Status(http.StatusCreated).JSON(user)
}

type loginUserRequest struct {
	db.LoginSession
	Email    string `json:"email" validate:"required,email,max=155"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}
type ResponseTokens struct {
	Token auth.TokenAccess `json:"token"`
	User  db.User          `json:"user"`
}

// loginUser login our user if email and password are a match
func (server *Server) loginUser(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_api.go -> loginUser()").Debug()
	var req loginUserRequest
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
	user, err := server.repo.GetUserByEmail(ctx.Context(), req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		server.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	// Check Password if matches and is correct
	if err = utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		server.logs.WithError(err).Warn("Invalid password provided by user")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	accessToken, exp, refreshToken, rexp, err := tokenCredentials(user, server)
	if err != nil {
		server.logs.WithError(err).Warn(err.Error())
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	args := db.SaveRefreshTokenParams{
		UserID:       user.ID,
		DeviceID:     req.DeviceID, // value shall be passed currently random string
		RefreshToken: refreshToken,
		ExpiresAt:    rexp,
	}

	if err = server.repo.SaveRefreshToken(ctx.Context(), args); err != nil {
		server.logs.WithError(err).Warn("unable to save refresh token")
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
	server.logs.WithField("user_id", user.ID).Debug("user logged in")
	return ctx.Status(http.StatusOK).JSON(resp)
}

// tokenCredentials returns access token, refresh token and when they all expire
func tokenCredentials(user db.User, server *Server) (accessToken string, accessTokenExp int64, refreshToken string,
	refreshTokenExp int64, err error) {
	// Create Token for the valid user
	accessToken, err = server.token.CreateAccessToken(string(user.ID), server.config.TokenDuration)
	if err != nil {
		server.logs.WithError(err).Warn("cannot create token")
	}
	server.logs.Debug("Server side Access Token Created")
	accessTokenExp, err = server.token.AccessTokenExpiresAt(accessToken)
	if err != nil {
		server.logs.WithError(err).Warn()
	}
	refreshToken, err = server.token.CreateRefreshToken(string(user.ID), server.config.RefreshTokenDuration)
	if err != nil {
		server.logs.WithError(err).Warn()
	}
	server.logs.Debug("Server side Refresh Token Created")
	refreshTokenExp, err = server.token.RefreshTokenExpiresAt(refreshToken)
	if err != nil {
		server.logs.WithError(err).Warn("unable to get Refresh Token expires at")
	}
	server.logs.Debug("Token response successful")
	return
}
