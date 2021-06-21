package api

import (
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/util"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

// Server will handle all request to and from database to client and routes
type Server struct {
	config util.Config
	repo db.Repo
	logs *logrus.Logger
	routes *fiber.App
}
// NewServer creates a new Server instance
func NewServer(config util.Config, repo db.Repo ,logs *logrus.Logger) (Server,error) {
	server := Server{
		config: config,
		repo: repo,
		logs: logs,
	}

	server.setupRoutes()
	return server,nil
}
// Run runs our Server instance
func (server *Server) Run(address string) error {
	return server.routes.Listen(address)
}
// setupRoutes manages our routes and middleware
func (server *Server) setupRoutes()  {
	server.routes = fiber.New()

	server.routes.Get("/version", server.version)

}