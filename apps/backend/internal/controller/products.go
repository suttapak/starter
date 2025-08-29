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
	Products interface {
		GetProduct(c *gin.Context)
		GetProducts(c *gin.Context)
		CreateProducts(c *gin.Context)
		UpdateProducts(c *gin.Context)
		DeleteProducts(c *gin.Context)

		UploadProductImages(c *gin.Context)
		DeleteProductImage(c *gin.Context)
	}

	products struct {
		products service.Products
	}
)

// DeleteProductImage implements Products.
func (i *products) DeleteProductImage(c *gin.Context) {
	productId, err := getProductsId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	productImageId, err := getProductImageId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := i.products.DeleteProductImage(c, productId, productImageId); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)

}

// UploadProductImages implements Products.
func (i *products) UploadProductImages(c *gin.Context) {
	uId, err := getProtectUserId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	productId, err := getProductsId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	files, err := getFiles(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	// Call service to upload images
	res, err := i.products.UploadProductImages(c, uId, productId, files)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}

func getProductsId(c *gin.Context) (uint, error) {
	idStr := c.Param("products_id")
	if idStr == "" {
		return 0, errs.ErrNotFound
	}
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errs.ErrBadRequest
	}
	return uint(idInt), nil
}

func (i *products) GetProduct(c *gin.Context) {
	id, err := getProductsId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.products.GetProduct(c, id)
	if err != nil {
		handlerError(c, err)
		return
	}

	handleJsonResponse(c, res)
}
func (i *products) GetProducts(c *gin.Context) {
	id, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	pg, err := helpers.NewPaginate(c)
	if err != nil {
		handlerError(c, err)
		return
	}

	f, err := filter.New[filter.ProductsFilter](c)
	if err != nil {
		handlerError(c, err)
		return
	}
	res, err := i.products.GetProducts(c, id, pg, f)
	if err != nil {
		handlerError(c, err)
		return
	}
	handlePaginationJsonResponse(c, res, pg)
}
func (i *products) CreateProducts(c *gin.Context) {
	teamId, err := getTeamId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.CreateProductRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	res, err := i.products.CreateProducts(c, teamId, &input)
	if err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, res)
}
func (i *products) UpdateProducts(c *gin.Context) {
	id, err := getProductsId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	input := service.UpdateProductRequest{}
	if err := c.ShouldBindJSON(&input); err != nil {
		handlerError(c, errs.ErrBadRequest)
		return
	}
	if err := i.products.UpdateProducts(c, id, &input); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func (i *products) DeleteProducts(c *gin.Context) {
	id, err := getProductsId(c)
	if err != nil {
		handlerError(c, err)
		return
	}
	if err := i.products.DeleteProducts(c, id); err != nil {
		handlerError(c, err)
		return
	}
	handleJsonResponse(c, nil)
}

func NewProducts(
	productsService service.Products,
) Products {
	return &products{
		products: productsService,
	}
}
