package api

import (
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/token"
	"FiberFinanceAPI/utils"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"time"
)

// Server will handle all request to and from database to client and routes
type Server struct {
	config utils.Config
	repo   db.Repo
	token token.Maker
	logs   *logrus.Logger
	routes *fiber.App
}
// NewServer creates a new Server instance
func NewServer(config utils.Config, repo db.Repo) (Server,error) {
	logrus.WithField("func", "server.go -> NewServer()").Info()
	maker, err := token.NewJWTTokenMaker(config.TokenSymmetricKey);if err != nil {
		return Server{},  fmt.Errorf("cannot create new token %w", err)
	}
	server := Server{
		config: config,
		repo: repo,
		token: maker,
	}
	server.NewLogger(logrus.New())
	server.setupRoutes()
	return server,nil
}
func(server *Server) NewLogger(logs *logrus.Logger) Server {
	server.logs = logs
	return *server

}
// Run runs our Server instance
func (server *Server) Run(address string) error {
	return server.routes.Listen(address)
}
// setupRoutes manages our routes and middleware
func (server *Server) setupRoutes()  {
	server.logs.WithField("func", "server.go -> setupRoutes()").Info()
	server.routes = fiber.New()
	server.routes.Get("/version", server.version)

	v1 := server.routes.Group("/api/v1")
	v1.Post("/users", server.createUser)
	v1.Post("/login", server.loginUser)

	_ = v1.Use(authTokenMiddleWare(server.token))
//	 TODO: Add other urls to the auth routes
}

// errorResponse returns an error response to our users
func errorResponse(status int, message error) fiber.Map {

	return fiber.Map{
		"status": status,
		"error" : message.Error(),
		"time" : time.Now(),
	}
}