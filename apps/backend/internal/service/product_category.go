package service

import (
	"context"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

type (
	ProductCategory interface {
		GetProductCategory(ctx context.Context, id uint) (*ProductCategoryResponse, error)
		GetProductCategories(ctx context.Context, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]ProductCategoryResponse, error)
		CreateProductCategory(ctx context.Context, input *CreateProductCategoryRequest) error
		UpdateProductCategory(ctx context.Context, id uint, input *UpdateProductCategoryRequest) error
		DeleteProductCategory(ctx context.Context, id uint) error
	}
	productCategory struct {
		logger          logger.AppLogger
		helper          helpers.Helper
		productCategory repository.ProductCategory
	}

	CreateProductCategoryRequest struct {
		// Add field here
		TeamId uint   `json:"team_id" binding:"required"`
		Name   string `json:"name" binding:"required"`
	}
	UpdateProductCategoryRequest struct {
		// Add field here
		Name string `json:"name" binding:"required"`
	}

	ProductProductCategoryResponse struct {
		CommonModel
		ProductID         uint                    `json:"product_id"`
		ProductCategoryID uint                    `json:"category_id"`
		ProductCategory   ProductCategoryResponse `json:"category"`
	}

	ProductCategoryResponse struct {
		CommonModel
		Name string `json:"name"`
	}
)

func (i *productCategory) GetProductCategory(ctx context.Context, id uint) (*ProductCategoryResponse, error) {
	model, err := i.productCategory.GetProductCategory(ctx, nil, id)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res ProductCategoryResponse
	if err := i.helper.ParseJson(model, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return &res, nil
}

func (i *productCategory) GetProductCategories(ctx context.Context, teamId uint, pg *helpers.Pagination, f *filter.ProductCategoryFilter) ([]ProductCategoryResponse, error) {
	models, err := i.productCategory.GetProductCategories(ctx, nil, teamId, pg, f)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res []ProductCategoryResponse
	if err := i.helper.ParseJson(models, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return res, nil
}

func (i *productCategory) CreateProductCategory(ctx context.Context, input *CreateProductCategoryRequest) error {
	body := repository.CreateProductCategoryRequest(*input)
	if err := i.productCategory.CreateProductCategory(ctx, nil, &body); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}
func (i *productCategory) UpdateProductCategory(ctx context.Context, id uint, input *UpdateProductCategoryRequest) error {
	body := repository.UpdateProductCategoryRequest(*input)
	if err := i.productCategory.UpdateProductCategory(ctx, nil, id, &body); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}
func (i *productCategory) DeleteProductCategory(ctx context.Context, id uint) error {
	if err := i.productCategory.DeleteProductCategory(ctx, nil, id); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

func NewProductCategory(
	logger logger.AppLogger,
	helper helpers.Helper,
	productCategoryRepo repository.ProductCategory,
) ProductCategory {
	return &productCategory{
		logger:          logger,
		helper:          helper,
		productCategory: productCategoryRepo,
	}
}
