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
	_, err := gorm.G[model.Image](tx).Where("id = ?", imageId).Delete(ctx)
	return err
}

// Save implements Image.
func (i *image) Save(ctx context.Context, tx *gorm.DB, userId uint, image *model.Image) (*model.Image, error) {
	if tx == nil {
		tx = i.db
	}
	err := gorm.G[model.Image](tx).Create(ctx, image)
	return image, err
}

func NewImage(db *gorm.DB) Image {
	return &image{
		db: db,
	}
}
