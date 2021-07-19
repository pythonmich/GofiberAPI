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
	//TODO: ADD validate to Server Struct
	config utils.Config
	repo   db.Repo
	//validate
	token  auth.Maker
	logs   *utils.StandardLogger
	routes *fiber.App
}

// NewServer creates a new Server instance
func NewServer(config utils.Config, logs *utils.StandardLogger, repo db.Repo) (Server, error) {
	logs.WithField("func", "server.go -> NewServer()").Info()
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
	logs.Info("New Server Created")
	return server, nil
}

// Run runs our Server instance
func (server *Server) Run(address string) error {
	return server.routes.Listen(address)
}

// registerRoutes manages our routes and middleware
func (server *Server) registerRoutes() {
	server.logs.WithField("func", "server.go -> registerRoutes()").Debug()
	server.routes = fiber.New()
	server.routes.Get("/version", server.version)

	permissions := newPermissions(server.repo, server.logs)

	v1 := server.routes.Group("/api/v1")
	v1.Post("/users", server.createUser)
	v1.Post("/login", server.loginUser)

	authRoutes := v1.Use(authTokenMiddleWare(server.token, server.logs))
	// -----TOKENS----
	authRoutes.Post("/refresh", server.refreshToken).Use(permissions.wrap())
	//TODO: Add other urls to the auth routes
	//----ROLES----
	authRoutes.Post("/users/:userID/roles", server.grantRole)
	authRoutes.Delete("/users/:userID/roles", server.revokeRole)
	authRoutes.Get("/users/:userID/role", server.getUserRole)
	authRoutes.Get("/users/:userID/roles", server.listRoles)
	authRoutes.Use(permissions.wrap(admin))

}

// errorResponse returns an error response to our users
func errorResponse(status int, message error) fiber.Map {
	return fiber.Map{
		"status": status,
		"error":  message.Error(),
		"time":   time.Now(),
	}
}
