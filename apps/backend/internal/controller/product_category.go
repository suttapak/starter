package controller

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/service"
)

type (
	ProductCategory interface {
		GetProductCategory(c *gin.Context)
		GetProductCategories(c *gin.Context)
		CreateProductCategory(c *gin.Context)
		UpdateProductCategory(c *gin.Context)
		DeleteProductCategory(c *gin.Context)
	}

	productCategory struct {
		productCategory service.ProductCategory
	}
)

func getProductCategoryId(c *gin.Context) (uint, error) {
	idStr := c.Param("product_category_id")
	if idStr == "" {
		return 0, errs.ErrNotFound
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errs.ErrBadRequest
	}
	return uint(idInt), nil
}

func (i *productCategory) GetProductCategory(c *gin.Context) {
	id, err := getProductCategoryId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.productCategory.GetProductCategory(c, id)
	if err != nil {
		handlerError(c, err)
		return
	}

	handleJsonResponse(c, res)
}
func (i *productCategory) GetProductCategories(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}

	f, err := filter.New[filter.ProductCategoryFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.productCategory.GetProductCategories(c, teamId, pg, f)
	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}
func (i *productCategory) CreateProductCategory(c *gin.Context) {
	input := service.CreateProductCategoryRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	if err := i.productCategory.CreateProductCategory(c, &input); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}
func (i *productCategory) UpdateProductCategory(c *gin.Context) {
	id, err := getProductCategoryId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.UpdateProductCategoryRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	if err := i.productCategory.UpdateProductCategory(c, id, &input); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func (i *productCategory) DeleteProductCategory(c *gin.Context) {
	id, err := getProductCategoryId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := i.productCategory.DeleteProductCategory(c, id); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func NewProductCategory(
	productCategoryService service.ProductCategory,
) ProductCategory {
	return &productCategory{
		productCategory: productCategoryService,
	}
}
