package api

import (
	model "FiberFinanceAPI/database/models"
	db "FiberFinanceAPI/database/sqlc"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var (
	transactionNotFound   = errors.New("transaction(s) not found or deleted")
	transactionDeletedMSG = "transaction successfully deleted at %s"
)

type transactionRequest struct {
	AccountID       model.AccountID       `json:"account_id" validate:"required"`
	CategoryID      model.CategoryID      `json:"category_id" validate:"required"`
	Name            string                `json:"name" validate:"required"`
	TransactionType model.TransactionType `json:"transaction_type" validate:"required"`
	Amount          int64                 `json:"amount" validate:"required"`
	Notes           string                `json:"notes"`
	Date            time.Time             `json:"date" validate:"required"`
}

func (s *Server) createTransaction(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> createTransaction()").Debug()
	var req transactionRequest
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
	args := db.CreateTransactionParams{
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		Date:            req.Date,
	}
	transaction, err := s.repo.CreateTransaction(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("Transaction created successfully")
	return ctx.Status(http.StatusCreated).JSON(transaction)
}

func (s *Server) updateTransaction(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> updateTransaction()").Debug()
	var req transactionRequest
	userID := ctx.Locals("userID").(model.UserID)

	transactionID := model.TransactionID(ctx.Params("transactionID"))
	if transactionID == "" {
		s.logs.WithField("transactionID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("transactionID not provided")))
	}

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
	args := db.UpdateTransactionParams{
		TransactionID:   transactionID,
		UserID:          userID,
		AccountID:       req.AccountID,
		CategoryID:      req.CategoryID,
		Name:            req.Name,
		TransactionType: req.TransactionType,
		Amount:          req.Amount,
		Date:            req.Date,
	}
	transaction, err := s.repo.UpdateTransaction(ctx.Context(), args)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
		}
		s.logs.WithError(err).Warn("could not update transaction")
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("Transaction updated successfully")
	return ctx.Status(http.StatusCreated).JSON(transaction)
}

func (s *Server) getTransaction(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> getTransaction()").Debug()
	transactionID := model.TransactionID(ctx.Params("transactionID"))
	if transactionID == "" {
		s.logs.WithField("transactionID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("transactionID not provided")))
	}
	transaction, err := s.repo.GetTransactionByID(ctx.Context(), transactionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
		}
		s.logs.WithError(err).Warn("could not get transaction")
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("Transaction returned successfully")
	return ctx.Status(http.StatusCreated).JSON(transaction)
}

type listTransactionsRequest struct {
	PageID   int32     `query:"page_id" validate:"required,min=1"`
	PageSize int32     `query:"page_size" validate:"required,min=5,max=10"`
	From     time.Time `json:"from" validate:"required"`
	To       time.Time `json:"to" validate:"required"`
}

func (s *Server) listTransactionsByUserID(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> listTransactionsByUserID()").Debug()
	var req listTransactionsRequest
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
	args := db.ListTxByUserIDParams{
		UserID: userID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
		From:   req.From,
		To:     req.To,
	}
	transactions, err := s.repo.ListTransactionsByUserID(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, errors.New("no transactions available")))
	}
	if len(transactions) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
	}
	s.logs.Info("transactions returned successfully")

	return ctx.Status(http.StatusOK).JSON(transactions)
}

func (s *Server) listTransactionsByAccountID(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> listTransactionsByAccountID()").Debug()
	var req listTransactionsRequest
	accountID := model.AccountID(ctx.Params("accountID"))
	if accountID == "" {
		s.logs.WithField("accountID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("accountID not provided")))
	}
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
	args := db.ListTxByAccountIDParams{
		AccountID: accountID,
		Limit:     req.PageSize,
		Offset:    (req.PageID - 1) * req.PageSize,
		From:      req.From,
		To:        req.To,
	}
	transactions, err := s.repo.ListTransactionsByAccountID(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(transactions) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
	}
	s.logs.Info("transactions returned successfully")

	return ctx.Status(http.StatusOK).JSON(transactions)
}

func (s *Server) listTransactionsByCategoryID(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> listTransactionsByCategoryID()").Debug()
	var req listTransactionsRequest
	categoryID := model.CategoryID(ctx.Params("categoryID"))
	if categoryID == "" {
		s.logs.WithField("categoryID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("categoryID not provided")))
	}
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
	args := db.ListTxByCategoryIDParams{
		CategoryID: categoryID,
		Limit:      req.PageSize,
		Offset:     (req.PageID - 1) * req.PageSize,
		From:       req.From,
		To:         req.To,
	}
	transactions, err := s.repo.ListTransactionsByCategoryID(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(transactions) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
	}
	s.logs.Info("transactions returned successfully")
	return ctx.Status(http.StatusOK).JSON(transactions)
}

func (s *Server) deleteTransaction(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "transactions.go -> deleteTransaction()").Debug()
	transactionID := model.TransactionID(ctx.Params("transactionID"))
	if transactionID == "" {
		s.logs.WithField("transactionID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("transactionID not provided")))
	}

	deletedAt, err := s.repo.DeleteTransaction(ctx.Context(), transactionID)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, transactionNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("transaction deleted successfully")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf(transactionDeletedMSG, deletedAt.Format(time.ANSIC))})
}
