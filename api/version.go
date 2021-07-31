package api

import (
	"FiberFinanceAPI/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
)

// serverVersion represents the current server version
type serverVersion struct {
	Version string `json:"version"`
}

// API route to return our API version
// it will return the version when the server starts which can be used if necessary
func (s *Server) version(ctx *fiber.Ctx) error {
	version := serverVersion{Version: utils.GetVersion(s.config)}
	s.logs.WithFields(logrus.Fields{"version": version.Version}).Info("version response")
	return ctx.Status(http.StatusOK).JSON(version)
}
