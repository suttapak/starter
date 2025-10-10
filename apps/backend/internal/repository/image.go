package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/suttapak/starter/internal/model"
)

type (
	Image interface {
		Save(ctx context.Context, tx *sqlx.Tx, userId uint, image *model.Image) (*model.Image, error)
		Delete(ctx context.Context, tx *sqlx.Tx, imageId uint) error
	}

	imageSqlx struct {
		db *sqlx.DB
	}
)

func NewImage(db *sqlx.DB) Image {
	return &imageSqlx{db: db}
}

// getDB returns the appropriate database connection (transaction or main DB)
func (i *imageSqlx) getDB(tx *sqlx.Tx) sqlx.ExtContext {
	if tx != nil {
		return tx
	}
	return i.db
}

// Save creates a new image record
func (i *imageSqlx) Save(ctx context.Context, tx *sqlx.Tx, userId uint, image *model.Image) (*model.Image, error) {
	db := i.getDB(tx)

	query := `
		INSERT INTO images (path, url, size, width, height, type, user_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := sqlx.GetContext(ctx, db, image, query,
		image.Path,
		image.Url,
		image.Size,
		image.Width,
		image.Height,
		image.Type,
		userId,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	return image, nil
}

// Delete soft deletes an image
func (i *imageSqlx) Delete(ctx context.Context, tx *sqlx.Tx, imageId uint) error {
	db := i.getDB(tx)

	query := `
		UPDATE images
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`

	result, err := db.ExecContext(ctx, query, imageId)
	if err != nil {
		return fmt.Errorf("failed to delete image: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("image not found")
	}

	return nil
}
