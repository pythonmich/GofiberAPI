package api

import (
	"FiberFinanceAPI/token"
	"errors"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strings"
)

const (
	authorizationHeaderKey = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)
// authTokenMiddleWare server side middleware verify token on server side
func authTokenMiddleWare(maker token.Maker) fiber.Handler{
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get(authorizationHeaderKey); if len(authHeader) == 0 {
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("authorization header not provided")))
		}
		// fields is used to split the authorization header by space
		// we expect the results field to have at least 2 element
		fields := strings.Fields(authHeader); if len(fields) < 2 {
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("invalid authorization header format")))
		}
		authType := strings.ToLower(fields[0])
		if authType != authorizationTypeBearer{
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, errors.New("authorization type is not supported by the server")))
		}

		accessToken := fields[1]
		payload, err := maker.VerifyToken(accessToken); if err != nil{
			status = http.StatusUnauthorized
			return ctx.Status(status).JSON(errorResponse(status, err))
		}
		// store payload in context so that other requests can match value
		ctx.Locals(authorizationPayloadKey, payload)
		return ctx.Next()
	}
}
