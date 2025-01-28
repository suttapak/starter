package service

import (
	"context"

	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/logger"
)

type (
	Auth interface {
		Register(ctx context.Context, user dto.UserRegisterDto) (res *response.AuthResponse, err error)
		Login(ctx context.Context, body dto.LoginDto) (res *response.AuthResponse, err error)
		VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *response.UserResponse, err error)
	}
	auth struct {
		// utils
		passwordHelper helpers.Helper
		logger         logger.AppLogger
		// service
		jwtService  JWTService
		mailService Email
		// repo
		mailRepo repository.MailRepository
		userRepo repository.User
	}
)

func (a auth) VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *response.UserResponse, err error) {
	email, err := a.jwtService.GetEmailFormToken(ctx, body.Token)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	_, have, err := a.userRepo.CheckEmail(ctx, nil, email)
	if !have || err != nil {
		a.logger.Error(errs.ErrGenerateJWTFail)
		return nil, errs.ErrBadRequest
	}

	userModel, err := a.userRepo.VerifyEmail(ctx, nil, email)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrVerifyEmail
	}
	if err := a.passwordHelper.ParseJson(userModel, &res); err != nil {
		return nil, errs.ErrInvalid
	}
	return
}

func (a auth) Login(ctx context.Context, body dto.LoginDto) (res *response.AuthResponse, err error) {
	// find user by email or username
	userModel, err := a.userRepo.GetUserByEmailOrUsername(ctx, nil, body.UserNameEmail)
	if err != nil || userModel == nil {
		a.logger.Error(err)
		return nil, errs.ErrUsernameOrPasswordIncorect
	}

	// check password
	pass, err := a.passwordHelper.CheckPassword(userModel.Password, []byte(body.Password))
	if err != nil || !pass {
		a.logger.Error(err)
		return nil, errs.ErrUsernameOrPasswordIncorect
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
	var (
		userModel model.User
	)

	// check username in database
	_, dulicateUsername, err := a.userRepo.CheckUsername(ctx, nil, user.Username)
	if dulicateUsername || err != nil {
		a.logger.Debug("username are duplicate")
		return nil, errs.ErrDuplicateUsername
	}

	// check username in database
	_, duplicateEmail, err := a.userRepo.CheckEmail(ctx, nil, user.Email)
	if duplicateEmail || err != nil {
		a.logger.Debug("email are duplicate")
		return nil, errs.ErrDuplicateEmail
	}

	// hash user password
	user.Password, err = a.passwordHelper.HashPassword(user.Password)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrHashPassword
	}

	// register user to database
	userModel = model.User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		FullName: user.FullName,
		RoleID:   idx.RoleUser,
	}
	registered, err := a.userRepo.Register(ctx, nil, userModel)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrRegisterUsername
	}

	// generate token
	token, err := a.jwtService.GenerateToken(ctx, registered.ID)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}

	// generate refresh token
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, registered.ID)
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

func NewAuth(
	logger logger.AppLogger,
	userRepo repository.User,
	jwtService JWTService,
	passwordHelper helpers.Helper,
	mailService Email,
	mailRepo repository.MailRepository,
) Auth {
	return &auth{
		passwordHelper: passwordHelper,
		logger:         logger,
		userRepo:       userRepo,
		jwtService:     jwtService,
		mailService:    mailService,
		mailRepo:       mailRepo,
	}
}
