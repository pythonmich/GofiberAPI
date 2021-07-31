package api

import (
	model "FiberFinanceAPI/database/models"
	db "FiberFinanceAPI/database/sqlc"
	"FiberFinanceAPI/utils"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	accountNotFound   = errors.New("account(s) not found or deleted")
	accountExists     = errors.New("account already exists")
	accountDeletedMSG = "account successfully deleted at %s"
)

//TODO: Create An API that returns account balances

type createAccountRequest struct {
	AccountName string             `json:"account_name" validate:"required"`
	AccountType model.AccountType  `json:"account_type" validate:"required"`
	Currency    utils.CurrencyCode `json:"currency" validate:"required,currency"`
}

func (s *Server) createAccount(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "accounts_api.go -> createAccount()").Debug()

	var req createAccountRequest
	userID := ctx.Locals("userID").(model.UserID)

	if err := ctx.BodyParser(&req); err != nil {
		s.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}
	// Balance is zero by default initial the user will be required in future to active the account by depositing an amount
	args := db.CreateAccountParams{
		UserID:      userID,
		AccountName: req.AccountName,
		AccountType: req.AccountType,
		Balance:     0,
		Currency:    req.Currency,
	}

	account, err := s.repo.CreateAccount(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			s.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres error codes")
			switch pqErr.Code.Name() {
			case "unique_violation", "foreign_key_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, accountExists))
			}
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("account returned successfully")
	return ctx.Status(http.StatusCreated).JSON(account)
}

func (s *Server) getAccount(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "accounts_api.go -> getAccount()").Debug()
	accountID := ctx.Params("accountID")
	if accountID == "" {
		s.logs.WithField("accountID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("accountID not provided")))
	}

	account, err := s.repo.GetAccountByID(ctx.Context(), model.AccountID(accountID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, accountNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("account returned successfully")

	return ctx.Status(http.StatusOK).JSON(account)

}

type accountBalanceResponse struct {
	AccountID model.AccountID    `json:"account_id"`
	UserID    model.UserID       `json:"user_id"`
	Balance   int64              `json:"balance"`
	Currency  utils.CurrencyCode `json:"currency"`
	Type      model.AccountType  `json:"type"`
}

func newBalanceResponse(account model.Account) accountBalanceResponse {
	return accountBalanceResponse{
		AccountID: account.AccountID,
		UserID:    account.UserID,
		Balance:   account.Balance,
		Currency:  account.Currency,
		Type:      account.Type,
	}
}

func (s *Server) accountBalance(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "accounts_api.go -> getAccount()").Debug()
	accountID := ctx.Params("accountID")
	if accountID == "" {
		s.logs.WithField("accountID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("accountID not provided")))
	}

	account, err := s.repo.GetAccountByID(ctx.Context(), model.AccountID(accountID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, accountNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("account balance returned successfully")
	response := newBalanceResponse(account)
	return ctx.Status(http.StatusOK).JSON(response)
}

type listAccountsRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (s *Server) listAccounts(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "accounts_api.go -> listAccounts()").Debug()
	var req listAccountsRequest
	userID := ctx.Locals("userID").(model.UserID)

	if err := ctx.QueryParser(&req); err != nil {
		s.logs.WithError(err).Warn("cannot decode parameters")
		status = http.StatusUnprocessableEntity
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if errs := s.validate.validateRequests(&req); len(errs) > 0 {
		s.logs.Warn("request data is invalid")
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errs)
	}
	s.logs.WithFields(logrus.Fields{"limit": req.PageSize, "offset": (req.PageID - 1) * req.PageSize}).Debug()
	args := db.ListAccountParams{
		UserID: userID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	accounts, err := s.repo.ListAccounts(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(accounts) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, accountNotFound))
	}
	s.logs.Info("accounts returned successfully")
	return ctx.Status(http.StatusOK).JSON(accounts)
}

func (s *Server) deleteAccount(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "accounts_api.go -> deleteAccount()").Debug()
	accountID := ctx.Params("accountID")
	if accountID == "" {
		s.logs.WithField("accountID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("accountID not provided")))
	}

	deletedAt, err := s.repo.DeleteAccount(ctx.Context(), model.AccountID(accountID))
	if err != nil {
		if err == sql.ErrNoRows {
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, accountNotFound))
		}

		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("account deleted successfully")
	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf(accountDeletedMSG, deletedAt.Format(time.ANSIC))})

}
