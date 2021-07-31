package database

import (
	model "FiberFinanceAPI/database/models"
	"context"
	"time"
)

const createCategory = `--name: CreateCategory :one
INSERT INTO categories(parent_id, user_id, name) 
VALUES ($1, $2, $3)
RETURNING category_id, parent_id, user_id, name, created_at, deleted_at
`

type CreateCategoryParams struct {
	ParentID model.CategoryID `json:"parent_id"`
	UserID   model.UserID     `json:"user_id"`
	Name     string           `json:"name"`
}

func (q *Queries) CreateCategory(ctx context.Context, args CreateCategoryParams) (model.Category, error) {
	q.logs.WithField("func", "database/sqlc/category.go -> CreateCategory()").Debug()
	row := q.db.QueryRowContext(ctx, createCategory, args.ParentID, args.UserID, args.Name)
	var category model.Category
	err := row.Scan(
		&category.ID,
		&category.ParentID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.DeletedAt,
	)
	return category, err
}

const updateCategory = `--name: UpdateCategory :one
UPDATE categories SET name = $3
AND parent_id = $2
WHERE category_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING category_id, parent_id, user_id, name, created_at, deleted_at`

type UpdateCategoryParams struct {
	CategoryID model.CategoryID `json:"category_id"`
	ParentID   model.CategoryID `json:"parent_id"`
	Name       string           `json:"name"`
}

func (q *Queries) UpdateCategory(ctx context.Context, args UpdateCategoryParams) (model.Category, error) {
	q.logs.WithField("func", "database/sqlc/category.go -> UpdateCategory()").Debug()
	row := q.db.QueryRowContext(ctx, updateCategory, args.CategoryID, args.ParentID, args.Name)
	var category model.Category
	err := row.Scan(
		&category.ID,
		&category.ParentID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.DeletedAt,
	)
	return category, err
}

const getCategoryByID = `--name: GetCategoryByID :one
SELECT * FROM categories
WHERE category_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
LIMIT 1`

func (q *Queries) GetCategoryByID(ctx context.Context, id model.CategoryID) (model.Category, error) {
	q.logs.WithField("func", "database/sqlc/category.go -> GetCategoryByID()").Debug()
	row := q.db.QueryRowContext(ctx, getCategoryByID, id)
	var category model.Category
	err := row.Scan(
		&category.ID,
		&category.ParentID,
		&category.UserID,
		&category.Name,
		&category.CreatedAt,
		&category.DeletedAt,
	)
	return category, err
}

const listCategory = `--name: ListCategories :many
SELECT * FROM categories
WHERE user_id = $1 AND deleted_at = '0001-01-01 00:00:00Z'
ORDER BY category_id
LIMIT $2 
OFFSET $3`

type ListCategoryParams struct {
	UserID model.UserID `json:"category_id"`
	Limit  int32        `json:"limit"`
	Offset int32        `json:"offset"`
}

func (q *Queries) ListCategories(ctx context.Context, args ListCategoryParams) ([]model.Category, error) {
	q.logs.WithField("func", "database/sqlc/category.go -> ListCategories()").Debug()
	rows, err := q.db.QueryContext(ctx, listCategory, args.UserID, args.Limit, args.Offset)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		q.logs.WithError(err).Warn()
	}()
	var categories []model.Category
	for rows.Next() {
		var category model.Category
		err = rows.Scan(
			&category.ID,
			&category.ParentID,
			&category.UserID,
			&category.Name,
			&category.CreatedAt,
			&category.DeletedAt,
		)
		categories = append(categories, category)
	}
	return categories, err
}

const deleteCategory = `--name: DeleteUser :exec
UPDATE categories SET deleted_at = now()
WHERE category_id = $1
AND deleted_at = '0001-01-01 00:00:00Z'
RETURNING deleted_at;
`

func (q *Queries) DeleteCategory(ctx context.Context, id model.CategoryID) (time.Time, error) {
	q.logs.WithField("func", "database/sqlc/category.go -> DeleteCategory()").Debug()
	row := q.db.QueryRowContext(ctx, deleteCategory, id)
	var category model.Category
	err := row.Scan(
		&category.DeletedAt,
	)
	return category.DeletedAt, err
}
