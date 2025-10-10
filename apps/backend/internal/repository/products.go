package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
)

type (
	Products interface {
		FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Product, error)
		FindAll(ctx context.Context, tx *sqlx.Tx, teamId uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]model.Product, error)
		Create(ctx context.Context, tx *sqlx.Tx, teamId uint, m *CreateProductsRequest) (*model.Product, error)
		Save(ctx context.Context, tx *sqlx.Tx, id uint, m *UpdateProductsRequest) error
		DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error

		FindImage(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ProductImage, error)
		CreateImage(ctx context.Context, tx *sqlx.Tx, productId uint, imageId uint) error
		DeleteImageById(ctx context.Context, tx *sqlx.Tx, productImageId uint) error
	}

	CreateProductsRequest struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UOM         string `json:"uom"`
		Price       int64  `json:"price"`
		CategoryID  []uint `json:"category_id"`
	}

	UpdateProductsRequest struct {
		Code        string `json:"code"`
		Name        string `json:"name"`
		Description string `json:"description"`
		UOM         string `json:"uom"`
		Price       int64  `json:"price"`
		CategoryID  []uint `json:"category_id"`
	}

	productsSqlx struct {
		db *sqlx.DB
	}
)

func NewProducts(db *sqlx.DB) Products {
	return &productsSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (p *productsSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return p.db
}

// FindById retrieves a product by ID with all relations
func (p *productsSqlx) FindById(ctx context.Context, tx *sqlx.Tx, id uint) (*model.Product, error) {
	db := p.getDB(tx)

	var product model.Product
	query := `
		SELECT id, team_id, code, name, description, uom, price, created_at, updated_at
		FROM products
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &product, query, id)
	if err != nil {
		return nil, err
	}

	// Load ProductProductCategory relations
	p.loadProductCategories(ctx, db, &product)

	p.loadProductImages(ctx, db, &product)

	return &product, nil
}

// FindAll retrieves all products for a team with pagination and filtering
func (p *productsSqlx) FindAll(ctx context.Context, tx *sqlx.Tx, teamId uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]model.Product, error) {
	db := p.getDB(tx)

	// Build count query
	countQuery := `
		SELECT COUNT(*)
		FROM products
		WHERE team_id = $1 AND deleted_at IS NULL
	`
	args := []interface{}{teamId}
	argIndex := 2

	if f != nil {
		if f.Name != "" {
			countQuery += fmt.Sprintf(` AND name ILIKE $%d`, argIndex)
			args = append(args, "%"+f.Name+"%")
			argIndex++
		}
		if f.Code != "" {
			countQuery += fmt.Sprintf(` AND code ILIKE $%d`, argIndex)
			args = append(args, "%"+f.Code+"%")
			argIndex++
		}
		if f.UOM != "" {
			countQuery += fmt.Sprintf(` AND uom ILIKE $%d`, argIndex)
			args = append(args, "%"+f.UOM+"%")
			argIndex++
		}
	}

	// Handle no pagination case
	if pg == nil {
		query := `
			SELECT id, team_id, code, name, description, uom, price, created_at, updated_at
			FROM products
			WHERE team_id = $1 AND deleted_at IS NULL
		`

		// Add filters
		filterArgs := []interface{}{teamId}
		argIdx := 2
		if f != nil {
			if f.Name != "" {
				query += fmt.Sprintf(` AND name ILIKE $%d`, argIdx)
				filterArgs = append(filterArgs, "%"+f.Name+"%")
				argIdx++
			}
			if f.Code != "" {
				query += fmt.Sprintf(` AND code ILIKE $%d`, argIdx)
				filterArgs = append(filterArgs, "%"+f.Code+"%")
				argIdx++
			}
			if f.UOM != "" {
				query += fmt.Sprintf(` AND uom ILIKE $%d`, argIdx)
				filterArgs = append(filterArgs, "%"+f.UOM+"%")
				argIdx++
			}
		}

		query += ` ORDER BY created_at DESC`

		var products []model.Product
		err := sqlx.SelectContext(ctx, db, &products, query, filterArgs...)
		if err != nil {
			return nil, fmt.Errorf("failed to find products: %w", err)
		}

		// Load relations for each product
		for i := range products {
			if err := p.loadProductCategories(ctx, db, &products[i]); err != nil {
				return nil, err
			}
			if err := p.loadProductImages(ctx, db, &products[i]); err != nil {
				return nil, err
			}
		}

		return products, nil
	}

	// Count total records
	if err := sqlx.GetContext(ctx, db, &pg.Count, countQuery, args...); err != nil {
		return nil, fmt.Errorf("failed to count products: %w", err)
	}

	// Apply pagination
	helpers.Paging(pg)

	// Build main query
	query := `
		SELECT id, team_id, code, name, description, uom, price, created_at, updated_at
		FROM products
		WHERE team_id = $1 AND deleted_at IS NULL
	`

	// Re-add filters with same logic
	if f != nil {
		idx := 2
		if f.Name != "" {
			query += fmt.Sprintf(` AND name ILIKE $%d`, idx)
			idx++
		}
		if f.Code != "" {
			query += fmt.Sprintf(` AND code ILIKE $%d`, idx)
			idx++
		}
		if f.UOM != "" {
			query += fmt.Sprintf(` AND uom ILIKE $%d`, idx)
			idx++
		}
	}

	query += ` ORDER BY created_at DESC LIMIT $` + fmt.Sprintf("%d", argIndex) + ` OFFSET $` + fmt.Sprintf("%d", argIndex+1)
	args = append(args, pg.Limit, pg.Offset)

	var products []model.Product
	err := sqlx.SelectContext(ctx, db, &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to find products: %w", err)
	}

	// Load relations for each product
	for i := range products {
		if err := p.loadProductCategories(ctx, db, &products[i]); err != nil {
			return nil, err
		}
		if err := p.loadProductImages(ctx, db, &products[i]); err != nil {
			return nil, err
		}
	}

	return products, nil
}

// Create creates a new product with categories
func (p *productsSqlx) Create(ctx context.Context, tx *sqlx.Tx, teamId uint, m *CreateProductsRequest) (*model.Product, error) {
	db := p.getDB(tx)

	// Insert product
	productQuery := `
		INSERT INTO products (team_id, code, name, description, uom, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	product := &model.Product{
		TeamID:      teamId,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		UOM:         m.UOM,
		Price:       m.Price,
	}

	err := sqlx.GetContext(ctx, db, product, productQuery,
		teamId,
		m.Code,
		m.Name,
		m.Description,
		m.UOM,
		m.Price,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Insert product categories
	if len(m.CategoryID) > 0 {
		for _, categoryID := range m.CategoryID {
			categoryQuery := `
				INSERT INTO product_product_categories (product_id, product_category_id, created_at, updated_at)
				VALUES ($1, $2, NOW(), NOW())
			`
			_, err := db.ExecContext(ctx, categoryQuery, product.ID, categoryID)
			if err != nil {
				return nil, fmt.Errorf("failed to create product category relation: %w", err)
			}
		}
	}

	// Load relations
	if err := p.loadProductCategories(ctx, db, product); err != nil {
		return nil, err
	}

	return product, nil
}

// Save updates an existing product and its categories
func (p *productsSqlx) Save(ctx context.Context, tx *sqlx.Tx, id uint, m *UpdateProductsRequest) error {
	db := p.getDB(tx)

	// Update product
	productQuery := `
		UPDATE products
		SET code = $1, name = $2, description = $3, uom = $4, price = $5, updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, productQuery,
		m.Code,
		m.Name,
		m.Description,
		m.UOM,
		m.Price,
		id,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	// Delete old categories (soft delete not needed for junction table)
	deleteQuery := `DELETE FROM product_product_categories WHERE product_id = $1`
	_, err = db.ExecContext(ctx, deleteQuery, id)
	if err != nil {
		return err
	}

	// Insert new categories
	if len(m.CategoryID) > 0 {
		for _, categoryID := range m.CategoryID {
			categoryQuery := `
				INSERT INTO product_product_categories (product_id, product_category_id, created_at, updated_at)
				VALUES ($1, $2, NOW(), NOW())
			`
			_, err := db.ExecContext(ctx, categoryQuery, id, categoryID)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DeleteById soft deletes a product
func (p *productsSqlx) DeleteById(ctx context.Context, tx *sqlx.Tx, id uint) error {
	db := p.getDB(tx)

	query := `
		UPDATE products
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// FindImage retrieves a product image by ID with the Image relation
func (p *productsSqlx) FindImage(ctx context.Context, tx *sqlx.Tx, id uint) (*model.ProductImage, error) {
	db := p.getDB(tx)

	var productImage model.ProductImage
	query := `
		SELECT id, product_id, image_id, created_at, updated_at
		FROM product_images
		WHERE id = $1 AND deleted_at IS NULL
	`

	err := sqlx.GetContext(ctx, db, &productImage, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("product image not found")
		}
		return nil, fmt.Errorf("failed to find product image: %w", err)
	}

	// Load Image relation
	imageQuery := `
		SELECT id, path, url, size, width, height, type, user_id, created_at, updated_at
		FROM images
		WHERE id = $1 AND deleted_at IS NULL
	`
	image := &model.Image{}
	if err := sqlx.GetContext(ctx, db, image, imageQuery, productImage.ImageID); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("failed to load image: %w", err)
		}
	} else {
		productImage.Image = image
	}

	return &productImage, nil
}

// CreateImage creates a product image relation
func (p *productsSqlx) CreateImage(ctx context.Context, tx *sqlx.Tx, productId uint, imageId uint) error {
	db := p.getDB(tx)

	query := `
		INSERT INTO product_images (product_id, image_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`

	_, err := db.ExecContext(ctx, query, productId, imageId)
	if err != nil {
		return fmt.Errorf("failed to create product image: %w", err)
	}

	return nil
}

// DeleteImageById soft deletes a product image
func (p *productsSqlx) DeleteImageById(ctx context.Context, tx *sqlx.Tx, productImageId uint) error {
	db := p.getDB(tx)

	query := `
		UPDATE product_images
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, productImageId)
	if err != nil {
		return fmt.Errorf("failed to delete product image: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product image not found")
	}

	return nil
}

// loadProductCategories loads ProductProductCategory relations for a product
func (p *productsSqlx) loadProductCategories(ctx context.Context, db sqlx.ExtContext, product *model.Product) error {
	query := `
		SELECT ppc.id, ppc.product_id, ppc.product_category_id, ppc.created_at, ppc.updated_at
		FROM product_product_categories ppc
		WHERE ppc.product_id = $1 AND ppc.deleted_at IS NULL
	`

	var categories []model.ProductProductCategory
	err := sqlx.SelectContext(ctx, db, &categories, query, product.ID)
	if err != nil {
		return err
	}

	// Load ProductCategory for each relation
	for i := range categories {
		categoryQuery := `
			SELECT id, team_id, name, created_at, updated_at
			FROM product_categories
			WHERE id = $1 AND deleted_at IS NULL
		`
		category := &model.ProductCategory{}
		if err := sqlx.GetContext(ctx, db, category, categoryQuery, categories[i].ProductCategoryID); err != nil {
			return err
		} else {
			categories[i].ProductCategory = category
		}
	}

	product.ProductProductCategory = categories
	return nil
}

// loadProductImages loads ProductImage relations for a product
func (p *productsSqlx) loadProductImages(ctx context.Context, db sqlx.ExtContext, product *model.Product) error {
	query := `
		SELECT id, product_id, image_id, created_at, updated_at
		FROM product_images
		WHERE product_id = $1 AND deleted_at IS NULL
	`

	var images []model.ProductImage
	err := sqlx.SelectContext(ctx, db, &images, query, product.ID)
	if err != nil {
		return err
	}

	// Load Image for each ProductImage
	for i := range images {
		imageQuery := `
			SELECT id, path, url, size, width, height, type, user_id, created_at, updated_at
			FROM images
			WHERE id = $1 AND deleted_at IS NULL
		`
		image := &model.Image{}
		if err := sqlx.GetContext(ctx, db, image, imageQuery, images[i].ImageID); err != nil {
			return err
		} else {
			images[i].Image = image
		}
	}

	product.ProductImage = images
	return nil
}
