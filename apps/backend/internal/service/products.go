package service

import (
	"context"
	"math"
	"mime/multipart"
	"path/filepath"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/i18n"
	"github.com/suttapak/starter/internal/filter"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

type (
	Products interface {
		GetProduct(ctx context.Context, id uint) (*ProductResponse, error)
		GetProducts(ctx context.Context, id uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]ProductResponse, error)
		CreateProducts(ctx context.Context, teamId uint, input *CreateProductRequest) (*ProductResponse, error)
		UpdateProducts(ctx context.Context, id uint, input *UpdateProductRequest) error
		DeleteProducts(ctx context.Context, id uint) error
		// Upload Product's Images
		UploadProductImages(ctx context.Context, userId, productId uint, files []*multipart.FileHeader) (*ProductResponse, error)
		DeleteProductImage(ctx context.Context, productId, productImageId uint) error
	}
	products struct {
		logger       logger.AppLogger
		helper       helpers.Helper
		db           repository.DatabaseTransaction
		products     repository.Products
		codeService  CodeService
		excel        Excel
		i18n         i18n.I18N
		imageRepo    repository.Image
		imageService ImageFileService
	}

	CreateProductRequest struct {
		Code        string  `json:"code"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		UOM         string  `json:"uom" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		CategoryID  []uint  `json:"category_id"`
	}
	UpdateProductRequest struct {
		Code        string  `json:"code" binding:"required"`
		Name        string  `json:"name" binding:"required"`
		Description string  `json:"description"`
		UOM         string  `json:"uom" binding:"required"`
		Price       float64 `json:"price" binding:"required"`
		CategoryID  []uint  `json:"category_id"`
	}

	ProductResponse struct {
		CommonModel
		TeamID                 uint                             `json:"team_id"`
		Code                   string                           `json:"code" `
		Name                   string                           `json:"name"`
		Description            string                           `json:"description"`
		UOM                    string                           `json:"uom"`
		Price                  int64                            `json:"price"`
		ProductProductCategory []ProductProductCategoryResponse `json:"product_product_category"`
		ProductImage           []ProductImage                   `json:"product_image"`
	}

	ProductImage struct {
		CommonModel
		ProductID uint  `json:"product_id"`
		ImageID   uint  `json:"image_id"`
		Image     Image `json:"image"`
	}
)

// DeleteProductImage implements Products.
func (i *products) DeleteProductImage(ctx context.Context, productId, productImageId uint) error {
	// begin tx
	tx := i.db.BeginTx()
	// get product image
	productImageModel, err := i.products.FindImage(ctx, tx, productImageId)
	if err != nil {
		i.logger.Error(err)
		return errs.ErrInternal
	}
	if productImageModel.ProductID != productId {
		return errs.ErrForbidden
	}
	// remove image
	if err := i.imageRepo.Delete(ctx, tx, productImageModel.ImageID); err != nil {
		i.logger.Error(err)
		i.db.RollbackTx(tx)
		return errs.ErrInternal
	}
	// remove product image
	if err := i.products.DeleteImageById(ctx, tx, productImageId); err != nil {
		i.logger.Error(err)
		i.db.RollbackTx(tx)
		return errs.ErrInternal
	}
	// remove file in disk
	i.db.CommitTx(tx)
	if err := i.imageService.DeleteFile(productImageModel.Image.Path); err != nil {
		i.logger.Error(err)
		// i.products.RollbackTx(tx)
		return errs.ErrInternal
	}
	// return response
	return nil
}

// UploadProductImages implements Products.
func (i *products) UploadProductImages(ctx context.Context, userId, productId uint, files []*multipart.FileHeader) (*ProductResponse, error) {
	success := []*model.Image{}
	rollback := func() {
		for _, img := range success {
			if err := i.imageRepo.Delete(ctx, nil, img.ID); err != nil {
				i.logger.Error(err)
			}
		}
		// delete file in disk
		for _, img := range success {
			if err := i.imageService.DeleteFile(img.Path); err != nil {
				i.logger.Error(err)
			}
		}
	}
	for _, fileHeader := range files {
		// Process each file
		isImage, err := i.imageService.IsImageFromFileHeader(fileHeader)
		if !isImage || err != nil {
			i.logger.Error(err)
			rollback()
			return nil, errs.ErrFileUploadNotImage
		}
		// get image size , width height mine type
		imgStats, err := i.imageService.GetImageStatsFromFileHeader(fileHeader)
		if err != nil {
			i.logger.Error(err)
			rollback()
			return nil, errs.ErrFileImageCanNotGetStats
		}
		// create file path
		imageName := i.imageService.GetUuidFileNameFromFileHeader(fileHeader)
		filePath := filepath.Join(productImagePath, imageName)
		// save file
		if err := i.imageService.SaveFileFromFileHeader(fileHeader, filePath); err != nil {
			i.logger.Error(err)
			rollback()
			return nil, errs.ErrFileImageCanNotSaveToDisk
		}

		m := model.Image{
			Path:   filePath,
			Url:    imageName,
			Size:   imgStats.size,
			Width:  uint(imgStats.width),
			Height: uint(imgStats.height),
			Type:   imgStats.mimeType,
			UserID: userId,
		}
		imageModel, err := i.imageRepo.Save(ctx, nil, userId, &m)
		if err != nil {
			i.logger.Error(err)
			rollback()
			return nil, errs.HandleGorm(err)
		}
		if err := i.products.CreateImage(ctx, nil, productId, imageModel.ID); err != nil {
			i.logger.Error(err)
			rollback()
			return nil, errs.HandleGorm(err)
		}
		success = append(success, imageModel)
	}
	return i.GetProduct(ctx, productId)
}

func (i *products) getDownloadExcelHeader(local i18n.Local) []string {
	header := []string{
		i.i18n.GetMessage(local, "excel_header_download_stock.index"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_code"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_name"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_description"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_price"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_uom"),
		i.i18n.GetMessage(local, "excel_header_download_stock.product_category"),
		i.i18n.GetMessage(local, "excel_header_download_stock.available_qty"),
		i.i18n.GetMessage(local, "excel_header_download_stock.remaining_qty"),
	}

	return header
}

func (i *products) GetProduct(ctx context.Context, id uint) (*ProductResponse, error) {

	model, err := i.products.FindById(ctx, nil, id)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res ProductResponse
	if err := i.helper.ParseJson(model, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return &res, nil
}

func (i *products) GetProducts(ctx context.Context, id uint, pg *helpers.Pagination, f *filter.ProductsFilter) ([]ProductResponse, error) {
	models, err := i.products.FindAll(ctx, nil, id, pg, f)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	var res []ProductResponse
	if err := i.helper.ParseJson(models, &res); err != nil {
		i.logger.Error(err)
		return nil, errs.ErrInvalid
	}
	return res, nil
}

func (i *products) CreateProducts(ctx context.Context, teamId uint, input *CreateProductRequest) (res *ProductResponse, err error) {
	if input.Code == "" {
		code, err := i.codeService.GenerateProductCode(ctx, teamId)
		if err != nil {
			i.logger.Error(err)
			return nil, errs.HandleGorm(err)
		}
		input.Code = code
	}

	// Convert price to cents
	data := repository.CreateProductsRequest{
		Code:        input.Code,
		Name:        input.Name,
		Description: input.Description,
		UOM:         input.UOM,
		Price:       int64(math.Round(input.Price * 100)),
		CategoryID:  input.CategoryID,
	}
	body := repository.CreateProductsRequest(data)
	productModel, err := i.products.Create(ctx, nil, teamId, &body)
	if err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	if err := i.helper.ParseJson(productModel, &res); err != nil {
		i.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	return res, nil
}
func (i *products) UpdateProducts(ctx context.Context, id uint, input *UpdateProductRequest) error {
	data := repository.UpdateProductsRequest{
		Name:        input.Name,
		Description: input.Description,
		UOM:         input.UOM,
		Price:       int64(math.Round(input.Price * 100)),
		CategoryID:  input.CategoryID,
	}

	body := repository.UpdateProductsRequest(data)
	if err := i.products.Save(ctx, nil, id, &body); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}
func (i *products) DeleteProducts(ctx context.Context, id uint) error {
	if err := i.products.DeleteById(ctx, nil, id); err != nil {
		i.logger.Error(err)
		return errs.HandleGorm(err)
	}
	return nil
}

func NewProducts(
	logger logger.AppLogger,
	helper helpers.Helper,
	db repository.DatabaseTransaction,
	productsRepo repository.Products,
	codeService CodeService,
	excel Excel,
	i18n i18n.I18N,
	imageRepo repository.Image,
	imageService ImageFileService,

) Products {
	return &products{
		logger:       logger,
		helper:       helper,
		db:           db,
		products:     productsRepo,
		codeService:  codeService,
		excel:        excel,
		i18n:         i18n,
		imageRepo:    imageRepo,
		imageService: imageService,
	}
}
