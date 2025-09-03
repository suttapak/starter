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
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id		path		integer	true	"Team ID".
//	@Param		products_id	path		integer	true	"Products ID".
//	@Success	201			{object}	Response[any]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/{team_id}/products/{products_id} [delete]
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
//
//	@Tags		products
//	@Accept		multipart/form-data
//	@Produce	json
//	@Param		team_id		path		integer	true	"Team ID".
//	@Param		products_id	path		integer	true	"Products ID".
//	@Param		files		formData	file	true	"Profile image file"
//	@Success	201			{object}	Response[service.ProductResponse]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/{team_id}/products/{products_id}/upload_image [post]
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

// GetProduct
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id		path		integer	true	"Team ID".
//	@Param		products_id	path		integer	true	"Products ID".
//	@Success	200			{object}	Response[service.ProductResponse]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/{team_id}/products/{products_id} [get]
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

// GetProducts
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id	path		integer	true	"Team ID".
//	@Param		page	query		integer	false	"Page".
//	@Param		limit	query		integer	false	"Limit".
//	@Param		code	query		string	false	"Code".
//	@Param		name	query		string	false	"Name".
//	@Param		uom		query		string	false	"UOM".
//	@Success	200		{object}	ResponsePagination[[]service.ProductResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/{team_id}/products [get]
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

// CreateProducts
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id	path		integer							true	"Team ID".
//	@Param		data	body		service.CreateProductRequest	false	"CreateProductRequest".
//	@Success	201		{object}	Response[service.ProductResponse]
//	@Failure	400		{object}	Response[any]
//	@Failure	404		{object}	Response[any]
//	@Failure	500		{object}	Response[any]
//	@Router		/teams/{team_id}/products [post]
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

// UpdateProducts
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id		path		integer							true	"Team ID".
//	@Param		products_id	path		integer							true	"Products ID".
//	@Param		data		body		service.UpdateProductRequest	false	"UpdateProductRequest".
//	@Success	201			{object}	Response[any]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/{team_id}/products/{products_id} [post]
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

// DeleteProducts
//
//	@Tags		products
//	@Accept		json
//	@Produce	json
//	@Param		team_id		path		integer	true	"Team ID".
//	@Param		products_id	path		integer	true	"Products ID".
//	@Success	201			{object}	Response[any]
//	@Failure	400			{object}	Response[any]
//	@Failure	404			{object}	Response[any]
//	@Failure	500			{object}	Response[any]
//	@Router		/teams/{team_id}/products/{products_id} [delete]
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
