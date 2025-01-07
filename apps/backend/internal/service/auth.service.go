package service

import (
	"context"
	"errors"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/logger"

	"gorm.io/gorm"
)

type (
	Auth interface {
		Register(ctx context.Context, user dto.UserRegisterDto) (res *response.AuthResponse, err error)
		Login(ctx context.Context, body dto.LoginDto) (res *response.AuthResponse, err error)
		VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *response.UserResponse, err error)
	}
	auth struct {
		// utils
		logger logger.AppLogger
		// service
		jwtService  JWTService
		mailService Email
		// repo
		mailRepo repository.MailRepository
		userRepo repository.User
	}
)

func (a auth) VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *response.UserResponse, err error) {
	a.logger.Info(body.Token)
	email, err := a.jwtService.GetEmailFormToken(ctx, body.Token)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	_, have, err := a.userRepo.CheckEmail(ctx, nil, email)
	if !have {
		a.logger.Error(errs.ErrGenerateJWTFail)
		return nil, errs.ErrBadRequest
	}
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrInternal
	}
	userModel, err := a.userRepo.VerifyEmail(ctx, nil, email)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrInternal
	}
	if err := helpers.ParseJson(userModel, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return
}

func (a auth) Login(ctx context.Context, body dto.LoginDto) (res *response.AuthResponse, err error) {
	// find user by email or username
	userModel, err := a.userRepo.GetUserByEmailOrUsername(ctx, nil, body.UserNameEmail)
	if err != nil || userModel == nil {
		a.logger.Error(err)
		return nil, errs.ErrUnauthorized
	}

	// check password
	pass, err := helpers.CheckPassword(userModel.Password, []byte(body.Password))
	if err != nil || !pass {
		a.logger.Error(err)
		return nil, errs.ErrBadRequest
	}
	// generate token
	token, err := a.jwtService.GenerateToken(ctx, userModel.ID)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, userModel.ID)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	res = &response.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	// send res
	return res, err

}

func (a auth) Register(ctx context.Context, user dto.UserRegisterDto) (res *response.AuthResponse, err error) {
	// transaction handler
	tx := a.userRepo.BeginTx()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	var (
		userModel model.User
	)

	// check username in database
	_, dulicateUsername, err := a.userRepo.CheckUsername(ctx, tx, user.Username)
	if dulicateUsername {
		a.logger.Debug("username are duplicate")
		return nil, errs.ErrDuplicateUsername
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error(err)
		return nil, errs.ErrDuplicateUsername
	}

	// check username in database
	_, duplicateEmail, err := a.userRepo.CheckEmail(ctx, tx, user.Username)
	if duplicateEmail {
		a.logger.Debug("email are duplicate")
		return nil, errs.ErrDuplicateEmail
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		a.logger.Error(err)
		return nil, errs.ErrDuplicateEmail
	}

	// hash user password
	user.Password, err = helpers.HashPassword(user.Password)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrHashPassword
	}

	// register user to database
	userModel = model.User{
		Username:  user.Username,
		Password:  user.Password,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		RoleID:    idx.RoleUser,
	}
	register, err := a.userRepo.Register(ctx, tx, userModel)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrRegisterUsername
	}

	// generate token
	token, err := a.jwtService.GenerateToken(ctx, register.ID)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}

	// generate refresh token
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, register.ID)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}

	res = &response.AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	return
}

func NewAuth(logger logger.AppLogger, userRepo repository.User, jwtService JWTService) Auth {
	return &auth{
		logger:     logger,
		userRepo:   userRepo,
		jwtService: jwtService,
	}
}
