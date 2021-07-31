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
	categoryNotFound   = errors.New("category does not exist or deleted")
	categoryDeletedMSG = "category successfully deleted at %s"
)

type createCategoryRequest struct {
	ParentID model.CategoryID `json:"parent_id"`
	Name     string           `json:"name" validate:"required"`
}

func (s *Server) createCategory(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "category_api.go -> createCategory()").Debug()
	var req createCategoryRequest
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
	//s.logs.Debug(reflect.DeepEqual(req.ParentID, nil))
	args := db.CreateCategoryParams{
		ParentID: req.ParentID,
		UserID:   userID,
		Name:     req.Name,
	}

	category, err := s.repo.CreateCategory(ctx.Context(), args)
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
	return ctx.Status(http.StatusCreated).JSON(category)
}

func (s *Server) getCategory(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "category_api.go -> getCategory()").Debug()
	categoryID := ctx.Params("categoryID")
	if categoryID == "" {
		s.logs.WithField("categoryID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("categoryID not provided")))
	}

	category, err := s.repo.GetCategoryByID(ctx.Context(), model.CategoryID(categoryID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, categoryNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("category returned successfully")

	return ctx.Status(http.StatusOK).JSON(category)
}

type updateCategoryRequest struct {
	ParentID model.CategoryID `json:"parent_id" validate:"required"`
	Name     string           `json:"name" validate:"required"`
}

func (s *Server) updateCategory(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "category_api.go -> updateCategory()").Debug()
	var req updateCategoryRequest
	categoryID := ctx.Params("categoryID")
	if categoryID == "" {
		s.logs.WithField("categoryID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("categoryID not provided")))
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
	args := db.UpdateCategoryParams{
		CategoryID: model.CategoryID(categoryID),
		ParentID:   req.ParentID,
		Name:       req.Name,
	}
	category, err := s.repo.UpdateCategory(ctx.Context(), args)
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, categoryNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("category updated successfully")
	return ctx.Status(http.StatusOK).JSON(category)
}

type listCategoryRequest struct {
	PageID   int32 `query:"page_id" validate:"required,min=1"`
	PageSize int32 `query:"page_size" validate:"required,min=5,max=10"`
}

func (s *Server) listCategories(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "category_api.go -> listCategories()").Debug()
	var req listCategoryRequest
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
	args := db.ListCategoryParams{
		UserID: userID,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}
	categories, err := s.repo.ListCategories(ctx.Context(), args)
	if err != nil {
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	if len(categories) == 0 {
		status = http.StatusNotFound
		return ctx.Status(status).JSON(errorResponse(status, categoryNotFound))
	}
	s.logs.Info("categories returned successfully")

	return ctx.Status(http.StatusOK).JSON(categories)
}

func (s *Server) deleteCategory(ctx *fiber.Ctx) error {
	s.logs.WithField("func", "category_api.go -> deleteCategory()").Debug()
	categoryID := ctx.Params("categoryID")
	if categoryID == "" {
		s.logs.WithField("categoryID", "not provided").Debug()
		status = http.StatusBadRequest
		return ctx.Status(status).JSON(errorResponse(status, errors.New("categoryID not provided")))
	}

	deletedAt, err := s.repo.DeleteCategory(ctx.Context(), model.CategoryID(categoryID))
	if err != nil {
		if err == sql.ErrNoRows {
			s.logs.WithError(err).Warn()
			status = http.StatusNotFound
			return ctx.Status(status).JSON(errorResponse(status, categoryNotFound))
		}
		s.logs.WithError(err).Warn()
		status = http.StatusInternalServerError
		return ctx.Status(status).JSON(errorResponse(status, err))
	}
	s.logs.Info("category deleted successfully")

	return ctx.Status(http.StatusOK).JSON(fiber.Map{"message": fmt.Sprintf(categoryDeletedMSG, deletedAt.Format(time.ANSIC))})
}
