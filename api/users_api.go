package api

import (
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var status int

var (
	ErrUserExist = errors.New("user with email exists")
	userNotFound = errors.New("user does not exist")
)

type createUserRequest struct {
	Email string `json:"email"  validate:"required,max=155,email"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}

type userResponse struct {
	ID db.UserID `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse{
	return  userResponse{
		ID: user.ID,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
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
	if errs := ValidateStruct(&req); len(errs) > 0 {
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs))
	}
	server.logs.WithFields(logrus.Fields{
		"email": req.Email,
	}).Info()
	hashPassword, err := utils.HashPassword(req.Password); if err != nil{
		server.logs.WithError(err).Warn(err)
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	args := db.CreateUserParams{
		Email:        req.Email,
		PasswordHash: hashPassword,
	}
	user, err := server.repo.CreateUser(ctx.Context(), args); if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status,ErrUserExist))
			}
		}
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))

	}

	response := newUserResponse(user)
	server.logs.WithField("message", "created user successfully").Info("Successful")
	return ctx.Status(http.StatusOK).JSON(response)
}

type loginUserRequest struct {
	Email string `json:"email" validate:"required,email,max=155"`
	Password string `json:"password" validate:"required,min=6,max=55"`
}
type loginResp struct {
	AccessToken string `json:"access_token"`
}
func (server Server) loginUser(ctx *fiber.Ctx) error {
	server.logs.WithField("func", "users_api.go -> loginUser()").Debug()
	var req loginUserRequest
	if err := ctx.BodyParser(&req); err != nil{
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := ValidateStruct(&req); len(errs) > 0 {
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(validateResponse(errs))
	}
	user, err := server.repo.GetUserByEmail(ctx.Context(), req.Email); if err!= nil{
		if err == sql.ErrNoRows{
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, userNotFound))
		}
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if err = utils.CheckPassword(req.Password, user.PasswordHash); err != nil{
		status = http.StatusUnauthorized
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	token, err := server.token.CreateToken(string(user.ID), server.config.TokenDuration); if err != nil{
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	resp := loginResp{
		AccessToken: token,
	}
	return ctx.Status(http.StatusOK).JSON(resp)

}