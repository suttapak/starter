package repository

import (
	"context"

	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
	"gorm.io/gorm"
)

type (
	Products interface {
		GetProduct(ctx context.Context, tx *gorm.DB, id uint) (*model.Product, error)
		GetProducts(ctx context.Context, tx *gorm.DB, teamId uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]model.Product, error)
		CreateProducts(ctx context.Context, tx *gorm.DB, teamId uint, m *CreateProductsRequest) (*model.Product, error)
		UpdateProducts(ctx context.Context, tx *gorm.DB, id uint, m *UpdateProductsRequest) error
		DeleteProducts(ctx context.Context, tx *gorm.DB, id uint) error

		FindProductImage(ctx context.Context, tx *gorm.DB, id uint) (*model.ProductImage, error)
		CreateProductImage(ctx context.Context, tx *gorm.DB, productId uint, imageId uint) error
		DeleteProductImage(ctx context.Context, tx *gorm.DB, productImageId uint) error
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

	products struct {
		db *gorm.DB
	}
)

// FindProductImage implements Products.
func (i *products) FindProductImage(ctx context.Context, tx *gorm.DB, id uint) (*model.ProductImage, error) {
	if tx == nil {
		tx = i.db
	}
	var res model.ProductImage
	err := tx.WithContext(ctx).Where("id = ?", id).Preload("Image").First(&res).Error
	return &res, err
}

// DeleteProductImage implements Products.
func (i *products) DeleteProductImage(ctx context.Context, tx *gorm.DB, productImageId uint) error {
	if tx == nil {
		tx = i.db
	}
	return tx.WithContext(ctx).Where("id = ?", productImageId).Delete(&model.ProductImage{}).Error
}

// CreateProductImage implements Products.
func (i *products) CreateProductImage(ctx context.Context, tx *gorm.DB, productId uint, imageId uint) error {
	if tx == nil {
		tx = i.db
	}
	profileImage := &model.ProductImage{
		ProductID: productId,
		ImageID:   imageId,
	}
	return tx.WithContext(ctx).Create(profileImage).Error

}

func (i *products) GetProduct(ctx context.Context, tx *gorm.DB, id uint) (*model.Product, error) {
	if tx == nil {
		tx = i.db
	}
	var res model.Product
	err := tx.Where("id = ?", id).Preload("ProductProductCategory").
		Preload("ProductProductCategory.ProductCategory").
		Preload("ProductImage.Image").
		First(&res).Error
	return &res, err
}

func (i *products) GetProducts(ctx context.Context, tx *gorm.DB, id uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]model.Product, error) {
	if tx == nil {
		tx = i.db
	}
	tx = tx.WithContext(ctx)

	if f.Name != "" {
		tx = tx.Where("name ILIKE ?", "%"+f.Name+"%")
	}
	if f.Code != "" {
		tx = tx.Where("code ILIKE ?", "%"+f.Code+"%")
	}
	if f.UOM != "" {
		tx = tx.Where("uom ILIKE ?", "%"+f.UOM+"%")
	}

	var res []model.Product
	if pg == nil {
		err := tx.Where("team_id = ?", id).
			Preload("ProductProductCategory").
			Preload("ProductProductCategory.ProductCategory").
			Preload("ProductImage.Image").
			Find(&res).Error
		return res, err
	}
	tx = tx.Model(&model.Product{}).Where("team_id = ?", id).Count(&pg.Count)
	helpers.Paging(pg)
	err := tx.Offset(pg.Offset).Limit(pg.Limit).
		Preload("ProductProductCategory").
		Preload("ProductProductCategory.ProductCategory").
		Preload("ProductImage.Image").
		Find(&res).Error
	return res, err
}

func (i *products) CreateProducts(ctx context.Context, tx *gorm.DB, teamId uint, m *CreateProductsRequest) (*model.Product, error) {
	if tx == nil {
		tx = i.db
	}
	category := []model.ProductProductCategory{}
	for _, v := range m.CategoryID {
		category = append(category, model.ProductProductCategory{
			ProductCategoryID: v,
		})
	}

	input := model.Product{
		TeamID:                 teamId,
		Code:                   m.Code,
		Name:                   m.Name,
		Description:            m.Description,
		UOM:                    m.UOM,
		Price:                  m.Price,
		ProductProductCategory: category,
	}
	err := tx.Create(&input).Error
	return &input, err
}

func (i *products) UpdateProducts(ctx context.Context, tx *gorm.DB, id uint, m *UpdateProductsRequest) error {
	if tx == nil {
		tx = i.db
	}
	// Delete old category
	if err := tx.Where("product_id = ?", id).Delete(&model.ProductProductCategory{}).Error; err != nil {
		return err
	}

	// Insert new category
	category := []model.ProductProductCategory{}
	for _, v := range m.CategoryID {
		category = append(category, model.ProductProductCategory{
			ProductCategoryID: v,
			ProductID:         id,
		})
	}
	if err := tx.Save(&category).Error; err != nil {
		return err
	}

	input := model.Product{
		Code:                   m.Code,
		Name:                   m.Name,
		Description:            m.Description,
		UOM:                    m.UOM,
		Price:                  m.Price,
		ProductProductCategory: category,
	}
	return tx.Model(&model.Product{}).Where("id = ?", id).Updates(&input).Error
}

func (i *products) DeleteProducts(ctx context.Context, tx *gorm.DB, id uint) error {
	if tx == nil {
		tx = i.db
	}
	input := model.Product{}
	err := tx.Where("id = ?", id).Delete(&input).Error
	return err
}

func NewProducts(db *gorm.DB) Products {
	return &products{
		db: db,
	}
}
