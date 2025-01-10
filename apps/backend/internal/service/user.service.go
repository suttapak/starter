package service

import (
	"context"
	"errors"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/logger"

	"gorm.io/gorm"
)

type (
	UserService interface {
		GetUserByUserId(ctx context.Context, uId uint) (res *response.UserResponse, err error)
	}
	userService struct {
		user   repository.User
		logger logger.AppLogger
		help   helpers.Helper
	}
)

// GetUserByUserId implements UserService.
func (u userService) GetUserByUserId(ctx context.Context, uId uint) (res *response.UserResponse, err error) {
	userModel, err := u.user.GetUserByUserId(ctx, nil, uId)
	// not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errs.ErrNotFound
	}
	if err != nil {
		u.logger.Error(err)
		return nil, errs.ErrInternal
	}
	if err := u.help.ParseJson(userModel, &res); err != nil {
		return nil, errs.ErrInternal
	}
	return
}

func newUserService(user repository.User, logger logger.AppLogger, help helpers.Helper) UserService {
	return userService{user: user, logger: logger, help: help}
}
