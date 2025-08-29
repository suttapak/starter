package service

import (
	"context"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

const (
	profileImagePath        = "public/static/profile/images"
	productImagePath        = "public/static/product/images"
	deletedProductImagePath = "public/static/product/deleted"
)

type (
	UserService interface {
		GetUserByUserId(ctx context.Context, uId uint) (res *UserResponse, err error)
		FindUserByUsername(ctx context.Context, username string) (res []UserResponse, err error)
		CheckUserIsVerifyEmail(ctx context.Context, userId uint) (bool, error)
		// Create profile image receive a image form file header and
		// get image info such as size, width, height, and type and user id who create
		CreateProfileImage(ctx context.Context, userId uint, fileHeader *multipart.FileHeader) (res *UserResponse, err error)
	}
	userService struct {
		user         repository.User
		image        repository.Image
		imageService ImageFileService
		logger       logger.AppLogger
		help         helpers.Helper
	}

	UserResponse struct {
		CommonModel
		Username      string         `json:"username"`
		Email         string         `json:"email"`
		EmailVerifyed bool           `json:"email_verifyed"`
		FullName      string         `json:"full_name"`
		RoleID        uint           `json:"role_id"`
		ProfileImage  []ProfileImage `json:"profile_image"`
		Role          Role           `json:"role"`
	}

	ProfileImage struct {
		CommonModel
		UserID  uint  `json:"user_id"`
		ImageID uint  `json:"image_id"`
		Image   Image `json:"image"`
	}

	Role struct {
		CommonModel
		Name string `json:"name" gorm:"uniqueIndex:idx_role_name"`
	}
	Image struct {
		CommonModel
		Path   string  `json:"path"`
		Url    string  `json:"url"`
		Size   float64 `json:"size"`
		Width  uint    `json:"width"`
		Height uint    `json:"height"`
		Type   string  `json:"type"`
	}
)

// CreateProfileImage implements UserService.
func (a *userService) CreateProfileImage(ctx context.Context, userId uint, fileHeader *multipart.FileHeader) (res *UserResponse, err error) {
	// check is image?
	isImage, err := a.imageService.IsImageFromFileHeader(fileHeader)
	if !isImage || err != nil {
		a.logger.Error(err)
		return nil, errs.ErrFileUploadNotImage
	}
	// get image size , width height mine type
	imgStats, err := a.imageService.GetImageStatsFromFileHeader(fileHeader)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrFileImageCanNotGetStats
	}
	// create file path
	imageName := a.imageService.GetUuidFileNameFromFileHeader(fileHeader)
	filePath := filepath.Join(profileImagePath, imageName)
	// save file
	if err := a.imageService.SaveFileFromFileHeader(fileHeader, filePath); err != nil {
		a.logger.Error(err)
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
	imageModel, err := a.image.CreateImage(ctx, nil, userId, &m)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	// save user image
	if err := a.user.CreateImageProfile(ctx, nil, userId, imageModel.ID); err != nil {
		a.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}
	// create profile image response
	res, err = a.GetUserByUserId(ctx, userId)
	return
}

// CheckUserIsVerifyEmail implements Auth.
func (a *userService) CheckUserIsVerifyEmail(ctx context.Context, userId uint) (bool, error) {
	res, err := a.user.CheckUserIsVerifyEmail(ctx, nil, userId)
	if err != nil {
		a.logger.Error(err)
		return false, errs.HandleGorm(err)
	}
	return res, nil
}

// FindUserByUsername implements UserService.
func (u *userService) FindUserByUsername(ctx context.Context, username string) (res []UserResponse, err error) {

	model, err := u.user.FindByUsername(ctx, nil, strings.ToLower(username))
	if err != nil {
		u.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}

	if err := u.help.ParseJson(model, &res); err != nil {
		return nil, errs.ErrInternal
	}
	return
}

// GetUserByUserId implements UserService.
func (u *userService) GetUserByUserId(ctx context.Context, uId uint) (res *UserResponse, err error) {
	userModel, err := u.user.FindById(ctx, nil, uId)
	if err != nil {
		u.logger.Error(err)
		return nil, errs.HandleGorm(err)
	}

	if err := u.help.ParseJson(userModel, &res); err != nil {
		return nil, errs.ErrInternal
	}
	return
}

func newUserService(
	user repository.User,
	image repository.Image,
	imageService ImageFileService,
	logger logger.AppLogger,
	help helpers.Helper,
) UserService {
	return &userService{
		user:         user,
		image:        image,
		imageService: imageService,
		logger:       logger,
		help:         help,
	}
}
