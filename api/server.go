package api

import (
	"FiberFinanceAPI/auth"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

// Server will handle all request to and from database to client and routes
type Server struct {
	config   utils.Config
	repo     db.Repo
	validate validates
	token    auth.Maker
	logs     *utils.StandardLogger
	routes   *fiber.App
}

// NewServer creates a new Server instance
func NewServer(config utils.Config, logs *utils.StandardLogger, repo db.Repo) (Server, error) {
	logs.WithField("func", "server.go -> NewServer()").Debug()
	maker, err := auth.NewJWTTokenMaker(config.TokenSymmetricKey, config.RefreshTokenSymmetricKey, logs)
	if err != nil {
		return Server{}, fmt.Errorf("cannot create new token %w", err)
	}
	server := Server{
		config: config,
		repo:   repo,
		logs:   logs,
		token:  maker,
	}
	server.registerRoutes()
	server.validate = newValidator(server.logs)
	logs.Debug("New Server Created")
	return server, nil
}

// Run runs our Server instance
func (s *Server) Run(address string) error {
	return s.routes.Listen(address)
}

// errorResponse returns an error response to our users
func errorResponse(status int, message error) fiber.Map {
	return fiber.Map{
		"status": status,
		"error":  message.Error(),
		"time":   time.Now(),
	}
}
