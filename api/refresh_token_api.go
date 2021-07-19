package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// refreshTokenRequest data user sends to server to refresh token
type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	db.LoginSession
}

// TODO: Store RefreshToken In redis in future
//refreshToken to refresh our token
func (server *Server) refreshToken(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "refresh_token_api.go -> refreshToken()").Debug()
	var req refreshTokenRequest
	if err := ctx.BodyParser(&req); err != nil {
		server.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := validateStruct(&req, server.logs); len(errs) > 0 {
		server.logs.Warn("invalid data provided")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs, server.logs))
	}
	server.logs.WithFields(logrus.Fields{
		"deviceID": req.DeviceID,
	}).Debug()

	// We ensure that a new token is not issued until enough time has elapsed
	// In this case, a new token will only be issued if the old token is within
	// 30 seconds of expiry. Otherwise, return a bad request status
	authPayload := ctx.Locals(authorizationPayloadKey).(*auth.AccessPayload)
	if time.Unix(authPayload.EXP, 0).Sub(time.Now()) > 30*time.Second {
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("access token not expired")))
	}
	payload, err := server.token.VerifyRefreshToken(req.RefreshToken)
	if err != nil {
		// if the token has expired the user will need to re-login otherwise we will delete the old refresh token  and create a new pair
		if errors.Is(err, auth.ErrExpiredToken) {
			err = fmt.Errorf("refresh token for user %s is expired", payload.SUB)
			server.logs.WithError(err).Warn(fmt.Errorf("refresh token for user %s is expired", payload.SUB))
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("token is expired login required")))

		}
		server.logs.WithError(err).Warn("refresh token verification failed")
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	args := db.GetSessionsParams{
		UserID:       db.UserID(payload.SUB),
		DeviceID:     req.DeviceID,
		RefreshToken: req.RefreshToken,
	}
	session, err := server.repo.GetSession(ctx.Context(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			server.logs.WithError(err).Warn(err)
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, errors.New("session not found")))
		}
		server.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	// verify if user still exists in our database
	user, err := server.repo.GetUserByID(ctx.Context(), string(session.UserID))
	if err != nil {
		if err == sql.ErrNoRows {
			server.logs.WithError(err).Warn(err)
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		server.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	accessToken, exp, refreshToken, rexp, err := tokenCredentials(user, server)
	if err != nil {
		server.logs.WithError(err).Warn(err.Error())
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}

	arg := db.SaveRefreshTokenParams{
		UserID:       user.ID,
		DeviceID:     req.DeviceID, // value shall be passed currently random string
		RefreshToken: refreshToken,
		ExpiresAt:    rexp,
	}

	if err = server.repo.SaveRefreshToken(ctx.Context(), arg); err != nil {
		server.logs.WithError(err).Warn("unable to save refresh token")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(err)
	}

	resp := ResponseTokens{
		Token: auth.TokenAccess{
			AccessToken:          accessToken,
			RefreshToken:         refreshToken,
			AccessTokenExpiresAt: exp,
		},
		User: user,
	}
	server.logs.WithField("user_id", user.ID).Debug("tokens generated successfully")
	return ctx.Status(http.StatusOK).JSON(resp)
}
