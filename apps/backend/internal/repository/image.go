package repository

import (
	"context"

	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Image interface {
		CreateImage(ctx context.Context, tx *gorm.DB, userId uint, image *model.Image) (*model.Image, error)
		DeleteImage(ctx context.Context, tx *gorm.DB, imageId uint) error
	}
	image struct {
		db *gorm.DB
	}
)

// DeleteImage implements Image.
func (i *image) DeleteImage(ctx context.Context, tx *gorm.DB, imageId uint) error {
	if tx == nil {
		tx = i.db
	}
	return tx.WithContext(ctx).Delete(&model.Image{}, imageId).Error
}

// CreateImage implements Image.
func (i *image) CreateImage(ctx context.Context, tx *gorm.DB, userId uint, image *model.Image) (*model.Image, error) {
	if tx == nil {
		tx = i.db
	}
	err := tx.WithContext(ctx).Create(&image).Error
	return image, err
}

func NewImage(db *gorm.DB) Image {
	return &image{
		db: db,
	}
}
