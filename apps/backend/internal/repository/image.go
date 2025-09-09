package repository

import (
	"context"

	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Image interface {
		Save(ctx context.Context, tx *gorm.DB, userId uint, image *model.Image) (*model.Image, error)
		Delete(ctx context.Context, tx *gorm.DB, imageId uint) error
	}
	image struct {
		db *gorm.DB
	}
)

// Delete implements Image.
func (i *image) Delete(ctx context.Context, tx *gorm.DB, imageId uint) error {
	if tx == nil {
		tx = i.db
	}
	return tx.WithContext(ctx).Delete(&model.Image{}, imageId).Error
}

// Save implements Image.
func (i *image) Save(ctx context.Context, tx *gorm.DB, userId uint, image *model.Image) (*model.Image, error) {
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
