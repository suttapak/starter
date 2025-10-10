package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
)

type (
	ProductCategory interface {
		FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ProductCategory, error)
		FindAll(ctx context.Context, tx *sqlx.Tx, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]model.ProductCategory, error)
		Create(ctx context.Context, tx *sqlx.Tx, m *CreateProductCategoryRequest) error
		Save(ctx context.Context, tx *sqlx.Tx, id uint, m *UpdateProductCategoryRequest) error
		DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error
	}

	CreateProductCategoryRequest struct {
		TeamId uint   `json:"teamId"`
		Name   string `json:"name"`
	}

	UpdateProductCategoryRequest struct {
		Name string `json:"name"`
	}

	productCategorySqlx struct {
		db *sqlx.DB
	}
)

func NewProductCategory(db *sqlx.DB) ProductCategory {
	return &productCategorySqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (p *productCategorySqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return p.db
}

// FindById retrieves a product category by ID
func (p *productCategorySqlx) FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ProductCategory, error) {
	db := p.getDB(tx)

	var category model.ProductCategory
	query := `
		SELECT id, team_id, name, created_at, updated_at
		FROM product_categories
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &category, query, id)
	return &category, err
}

// FindAll retrieves all product categories for a team with pagination
func (p *productCategorySqlx) FindAll(ctx context.Context, tx *sqlx.Tx, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]model.ProductCategory, error) {
	db := p.getDB(tx)

	// Build count query
	countQuery := `
		SELECT COUNT(*)
		FROM product_categories
		WHERE team_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{teamId}

	if f != nil && f.Name != "" {
		countQuery += ` AND name LIKE $2`
		args = append(args, "%"+f.Name+"%")
	}

	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to count product categories: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Build main query
	query := `
		SELECT id, team_id, name, created_at, updated_at
		FROM product_categories
		WHERE team_id = $1 AND deleted_at IS NULL
	`

	if f != nil && f.Name != "" {
		if len(args) == 2 {
			query += ` AND name LIKE $2`
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", len(args)+1) + ` OFFSET $` + fmt.Sprintf("%d", len(args)+2)
	args = append(args, pg.Limit, pg.Offset)

	var categories []model.ProductCategory
	err := sqlx.SelectContext(ctx, db, &categories, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find product categories: %w", err)
	}

	return categories, nil
}

// Create creates a new product category
func (p *productCategorySqlx) Create(ctx context.Context, tx *sqlx.Tx, m *CreateProductCategoryRequest) error {
	db := p.getDB(tx)

	query := `
		INSERT INTO product_categories (team_id, name, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := db.ExecContext(ctx, query, m.TeamId, m.Name)
	if err != nil {
		return fmt.Errorf("failed to create product category: %w", err)
	}

	return nil
}

// Save updates an existing product category
func (p *productCategorySqlx) Save(ctx context.Context, tx *sqlx.Tx, id uint, m *UpdateProductCategoryRequest) error {
	db := p.getDB(tx)

	query := `
		UPDATE product_categories
		SET name = $1, updated_at = NOW()
		WHERE id = $2 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, m.Name, id)
	if err != nil {
		return fmt.Errorf("failed to update product category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product category not found")
	}

	return nil
}

// DeleteById soft deletes a product category
func (p *productCategorySqlx) DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error {
	db := p.getDB(tx)

	query := `
		UPDATE product_categories
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product category: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product category not found")
	}

	return nil
}
