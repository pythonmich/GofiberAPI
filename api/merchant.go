package api

import (
	model "FiberFinanceAPI/database/models"
	db "FiberFinanceAPI/database/sqlc"
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
	merchantNotFound   = errors.New("merchant(s) not found or deleted")
	merchantDeletedMSG = "merchant successfully deleted at %s"
)

type merchantRequest struct {
	Name string `json:"name" validate:"required"`
}

func (s *Server) createMerchant(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "merchant_api.go -> createMerchant()").Debug()
	var req merchantRequest
	userID := ctx.Locals("userID").(model.UserID)
	s.logs.Debug(userID)
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
	args := db.CreateMerchantParams{
		UserID: userID,
		Name:   req.Name,
	}

	merchant, err := s.repo.CreateMerchant(ctx.Context(), args)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			s.logs.WithField(string(pqErr.Code), pqErr.Code.Name()).Debug("postgres error codes")
			switch pqErr.Code.Name() {
			case "unique_violation":
				status = http.StatusForbidden
				return ctx.Status(status).JSON(errorResponse(status, err))
			}
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	return ctx.Status(http.StatusCreated).JSON(merchant)
}

func (s *Server) getMerchant(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "merchant_api.go -> getMerchant()").Debug()

	merchantID := ctx.Params("merchantID")
	if merchantID == "" {
		s.logs.WithField("merchantID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("merchantID not provided")))
	}

	merchant, err := s.repo.GetMerchantByID(ctx.Context(), model.MerchantID(merchantID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, merchantNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("merchant returned successfully")

	return ctx.Status(http.StatusOK).JSON(merchant)
}

func (s *Server) updateMerchant(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "merchant_api.go -> updateMerchant()").Debug()
	var req merchantRequest
	merchantID := ctx.Params("merchantID")
	if merchantID == "" {
		s.logs.WithField("merchantID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("merchantID not provided")))
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
	args := db.UpdateMerchantParams{
		MerchantID: model.MerchantID(merchantID),
		Name:       req.Name,
	}
	merchant, err := s.repo.UpdateMerchant(ctx.Context(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, merchantNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("merchant updated successfully")
	return ctx.Status(http.StatusOK).JSON(merchant)
}

type listMerchantRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (s *Server) listMerchants(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "merchant_api.go -> listMerchants()").Debug()
	var req listMerchantRequest
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
	args := db.ListMerchantParams{
		UserID: userID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	merchants, err := s.repo.ListMerchants(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(merchants) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, merchantNotFound))
	}
	s.logs.Info("merchants returned successfully")

	return ctx.Status(http.StatusOK).JSON(merchants)
}

func (s *Server) deleteMerchant(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "merchant_api.go -> deleteMerchant()").Debug()
	merchantID := ctx.Params("merchantID")
	if merchantID == "" {
		s.logs.WithField("merchantID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("merchantID not provided")))
	}

	deletedAt, err := s.repo.DeleteMerchant(ctx.Context(), model.MerchantID(merchantID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, merchantNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("merchant deleted successfully")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf(merchantDeletedMSG, deletedAt.Format(time.ANSIC))})
}
