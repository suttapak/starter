package repository

import (
	"context"

	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	ProductCategory interface {
		GetProductCategory(ctx context.Context, tx *gorm.DB, id uint) (*model.ProductCategory, error)
		GetProductCategories(ctx context.Context, tx *gorm.DB, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]model.ProductCategory, error)
		CreateProductCategory(ctx context.Context, tx *gorm.DB, m *CreateProductCategoryRequest) error
		UpdateProductCategory(ctx context.Context, tx *gorm.DB, id uint, m *UpdateProductCategoryRequest) error
		DeleteProductCategory(ctx context.Context, tx *gorm.DB, id uint) error
	}

	CreateProductCategoryRequest struct {
		TeamId uint   `json:"teamId"`
		Name   string `json:"name"`
	}

	UpdateProductCategoryRequest struct {
		Name string `json:"name"`
	}

	productCategory struct {
		db *gorm.DB
	}
)

func (i *productCategory) GetProductCategory(ctx context.Context, tx *gorm.DB, id uint) (*model.ProductCategory, error) {
	if tx == nil {
		tx = i.db
	}
	var res model.ProductCategory
	err := tx.Where("id = ?", id).First(&res).Error
	return &res, err
}

func (i *productCategory) GetProductCategories(ctx context.Context, tx *gorm.DB, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]model.ProductCategory, error) {
	if tx == nil {
		tx = i.db
	}

	var res []model.ProductCategory
	tx = tx.Model(&model.ProductCategory{}).Where("team_id = ?", teamId)

	if f.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+f.Name+"%")
	}
	tx = tx.Count(&pg.Count)
	helpers.Paging(pg)
	err := tx.Offset(pg.Offset).Limit(pg.Limit).Find(&res).Error
	return res, err
}

func (i *productCategory) CreateProductCategory(ctx context.Context, tx *gorm.DB, m *CreateProductCategoryRequest) error {
	if tx == nil {
		tx = i.db
	}
	input := model.ProductCategory{
		TeamID: m.TeamId,
		Name:   m.Name,
	}
	err := tx.Create(&input).Error
	return err
}

func (i *productCategory) UpdateProductCategory(ctx context.Context, tx *gorm.DB, id uint, m *UpdateProductCategoryRequest) error {
	if tx == nil {
		tx = i.db
	}
	input := model.ProductCategory{
		Name: m.Name,
	}
	err := tx.Model(&model.ProductCategory{}).Where("id = ?", id).Updates(&input).Error
	return err
}

func (i *productCategory) DeleteProductCategory(ctx context.Context, tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = i.db
	}
	input := model.ProductCategory{}
	err := tx.Where("id = ?", id).Delete(&input).Error
	return err
}

func NewProductCategory(db *gorm.DB) ProductCategory {
	return &productCategory{
		db: db,
	}
}
