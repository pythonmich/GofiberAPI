package api

import (
	"FiberFinanceAPI/auth"
	"FiberFinanceAPI/utils"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

// authTokenMiddleWare server side middleware verify auth on server side
func authTokenMiddleWare(maker auth.Maker, logs *utils.StandardLogger) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		var err error
		logs.WithField("func", "auth_middleware.go -> authTokenMiddleWare()").Debug()
		authHeader := ctx.Get(authorizationHeaderKey)
		if len(authHeader) == 0 {
			err = errors.New("authorization header not provided")
			logs.WithError(err).Warn()
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		logs.WithField("len authHeader", len(authHeader)).Debug()
		// fields is used to split the authorization header by space
		// we expect the results field to have at least 2 element
		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			err = errors.New("invalid authorization header format")
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		authType := strings.ToLower(fields[0])
		logs.WithField("authType", authType).Debug()
		if authType != authorizationTypeBearer {
			err = errors.New("authorization type is not supported by the server")
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		accessToken := fields[1]
		logs.Debug("Token successfully gotten from header")
		payload, err := maker.VerifyAccessToken(accessToken)
		if err != nil {
			logs.WithError(err).Warn()
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		// store payload in context so that other requests can match value
		ctx.Locals(authorizationPayloadKey, payload)
		return ctx.Next()
	}
}
