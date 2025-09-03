package service

import (
	"context"
	"fmt"

	"github.com/suttapak/starter/config"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/idx"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/logger"
)

type (
	Auth interface {
		Register(ctx context.Context, user UserRegisterDto) (res *AuthResponse, err error)
		Login(ctx context.Context, body LoginDto) (res *AuthResponse, err error)
		RefreshToken(ctx context.Context, uId uint) (res *AuthResponse, err error)
		VerifyEmail(ctx context.Context, body VerifyEmailDto) (res *UserResponse, err error)
		SendVerifyEmail(ctx context.Context, userId SendVerifyEmailDto) (err error)
	}
	auth struct {
		// utils
		passwordHelper helpers.Helper
		logger         logger.AppLogger
		// service
		jwtService  JWTService
		mailService Email
		// repo
		userRepo repository.User
		config   *config.Config
	}
	AuthResponse struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}
	UserRegisterDto struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required,min=8"`
		FullName string `json:"full_name" binding:"required"`
	}

	LoginDto struct {
		UserNameEmail string `json:"username" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
	}

	VerifyEmailDto struct {
		Token string `json:"token" form:"token"`
	}
	SendVerifyEmailDto struct {
		UserID uint `json:"user_id"`
	}
)

// RefreshToken implements Auth.
func (a *auth) RefreshToken(ctx context.Context, uId uint) (res *AuthResponse, err error) {
	token, err := a.jwtService.GenerateToken(ctx, uId)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	refreshToken, err := a.jwtService.GenerateRefreshToken(ctx, uId)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
	}
	res = &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	// send res
	return res, err
}

// SendVerifyEmail implements Auth.
func (a *auth) SendVerifyEmail(ctx context.Context, input SendVerifyEmailDto) (err error) {
	// find user email
	userModel, err := a.userRepo.FindById(ctx, nil, input.UserID)
	if err != nil {
		a.logger.Error(err)
		return errs.HandleGorm(err)
	}
	// generate token
	token, err := a.jwtService.GenerateExternalToken(ctx, input.UserID)
	if err != nil {
		a.logger.Error(err)
		return errs.ErrGenerateJWTFail
	}
	// send email
	if err := a.mailService.NewRequest([]string{userModel.Email}, "Verify Email").
		ParseVerifyEmailTemplate(ctx, &VerifyEmailTemplateDataDto{
			Email:           userModel.Email,
			VerifyEmailLink: fmt.Sprintf("%s/api/v1/auth/email/verify?token=%s", a.config.SERVER.HOST_NAME, token),
		}).SendMail(ctx); err != nil {
		a.logger.Error(err)
		return errs.ErrSendEmail
	}
	return
}

func (a auth) VerifyEmail(ctx context.Context, body VerifyEmailDto) (res *UserResponse, err error) {
	email, err := a.jwtService.GetUserIdFromExternalToken(ctx, body.Token)
	if err != nil {
		a.logger.Error(err)
		return nil, errs.ErrGenerateJWTFail
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

func (a auth) Login(ctx context.Context, body LoginDto) (res *AuthResponse, err error) {
	// find user by email or username
	userModel, err := a.userRepo.GetUserByEmailOrUsername(ctx, nil, body.UserNameEmail)
	if err != nil || userModel == nil {
		a.logger.Error(err)
		return nil, errs.ErrUsernameOrPasswordIncorrect
	}

	// check password
	pass, err := a.passwordHelper.CheckPassword(userModel.Password, []byte(body.Password))
	if err != nil || !pass {
		a.logger.Error(err)
		return nil, errs.ErrUsernameOrPasswordIncorrect
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
	res = &AuthResponse{
		Token:        token,
		RefreshToken: refreshToken,
	}
	// send res
	return res, err

}

func (a auth) Register(ctx context.Context, user UserRegisterDto) (res *AuthResponse, err error) {
	var (
		userModel model.User
	)

	// check username in database
	_, duplicateUsername, err := a.userRepo.CheckUsername(ctx, nil, user.Username)
	if duplicateUsername || err != nil {
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

	res = &AuthResponse{
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
	config *config.Config,
) Auth {
	return &auth{
		passwordHelper: passwordHelper,
		logger:         logger,
		userRepo:       userRepo,
		jwtService:     jwtService,
		mailService:    mailService,
		config:         config,
	}
}
