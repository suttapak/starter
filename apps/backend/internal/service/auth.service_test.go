package service_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suttapak/starter/errs"
	"github.com/suttapak/starter/helpers"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/model"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/internal/response"
	"github.com/suttapak/starter/internal/service"
	"github.com/suttapak/starter/logger"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	type (
		PeparePassword struct {
			output bool
			err    error
		}
		PepareCheckUser struct {
			output *model.User
			err    error
		}
		PepareJwt struct {
			output string
			err    error
		}
		TestCase struct {
			name            string
			input           dto.LoginDto
			output          *response.AuthResponse
			err             error
			peparePassowrd  PeparePassword
			pepareCheckUser PepareCheckUser
			pepareJwt       PepareJwt
			pepareReJwt     PepareJwt
		}
	)
	testcases := []TestCase{
		{
			name: "TEST_SUCCESS",
			peparePassowrd: PeparePassword{
				output: true,
				err:    nil,
			},
			pepareCheckUser: PepareCheckUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},

			input: dto.LoginDto{
				UserNameEmail: "username",
				Password:      "password",
			},
			output: &response.AuthResponse{
				Token:        "token",
				RefreshToken: "retoken",
			},
			err: nil,
		},
		{
			name: "TEST_USERNOTFOUNd",
			peparePassowrd: PeparePassword{
				output: true,
				err:    nil,
			},
			pepareCheckUser: PepareCheckUser{
				output: nil,
				err:    gorm.ErrRecordNotFound,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},

			input: dto.LoginDto{
				UserNameEmail: "username",
				Password:      "password",
			},
			output: nil,
			err:    errs.ErrUsernameOrPasswordIncorect,
		},
		{
			name: "TEST_PASSWORD_INCORRECT",
			peparePassowrd: PeparePassword{
				output: false,
				err:    errs.ErrHashPassword,
			},
			pepareCheckUser: PepareCheckUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},

			input: dto.LoginDto{
				UserNameEmail: "username",
				Password:      "password",
			},
			output: nil,
			err:    errs.ErrUsernameOrPasswordIncorect,
		},
		{
			name: "TEST_JWT_ERROR",
			peparePassowrd: PeparePassword{
				output: true,
				err:    nil,
			},
			pepareCheckUser: PepareCheckUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			pepareJwt: PepareJwt{
				output: "",
				err:    errs.ErrGenerateJWTFail,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},

			input: dto.LoginDto{
				UserNameEmail: "username",
				Password:      "password",
			},
			output: nil,
			err:    errs.ErrGenerateJWTFail,
		},
		{
			name: "TEST_REJWT_ERROR",
			peparePassowrd: PeparePassword{
				output: true,
				err:    nil,
			},
			pepareCheckUser: PepareCheckUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "",
				err:    errs.ErrGenerateJWTFail,
			},

			input: dto.LoginDto{
				UserNameEmail: "username",
				Password:      "password",
			},
			output: nil,
			err:    errs.ErrGenerateJWTFail,
		},
	}

	ctx := context.Background()
	for _, tc := range testcases {

		t.Run(tc.name, func(t *testing.T) {
			//mail
			mail := service.NewMailServiceMock()
			mailRepo := repository.NewMailRepositoryMock()
			// password helper
			passwordHelper := helpers.NewHelperMock()
			passwordHelper.On("CheckPassword").Return(tc.peparePassowrd.output, tc.peparePassowrd.err)

			// user repo
			userRepo := repository.NewUserRepositoryMock()
			userRepo.On("GetUserByEmailOrUsername").Return(tc.pepareCheckUser.output, tc.pepareCheckUser.err)

			// jwt service
			jwtService := service.NewJwtServiceMock()
			jwtService.On("GenerateToken").Return(tc.pepareJwt.output, tc.pepareJwt.err)
			jwtService.On("GenerateRefreshToken").Return(tc.pepareReJwt.output, tc.pepareReJwt.err)
			// logger
			l := logger.NewLoggerMock()

			// new service
			authService := service.NewAuth(l, userRepo, jwtService, passwordHelper, mail, mailRepo)
			res, err := authService.Login(ctx, tc.input)

			// test
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.output, res)
		})
	}
}

func TestRegister(t *testing.T) {
	type (
		PepareRegister struct {
			output *model.User
			err    error
		}
		PepareCheckUsernameAndEmail struct {
			output1 *model.User
			output2 bool
			err     error
		}

		PeparePassword struct {
			output string
			err    error
		}
		PepareCheckUser struct {
			output *model.User
			err    error
		}
		PepareJwt struct {
			output string
			err    error
		}
		TestCase struct {
			name             string
			input            dto.UserRegisterDto
			output           *response.AuthResponse
			err              error
			peparePassowrd   PeparePassword
			pepareCheckUser  PepareCheckUsernameAndEmail
			pepareCheckEmail PepareCheckUsernameAndEmail
			pepareJwt        PepareJwt
			pepareReJwt      PepareJwt
			pepareRegister   PepareRegister
		}
	)
	testCases := []TestCase{
		{
			name: "TEST_SUCCESS",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: &response.AuthResponse{
				Token:        "token",
				RefreshToken: "retoken",
			},
			err: nil,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_ERROR_DUPLICATE_USERNAME",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrDuplicateUsername,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_ERROR_DUPLICATE_EMAIL",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrDuplicateEmail,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_ERROR_HASHPASSWORD",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrHashPassword,
			peparePassowrd: PeparePassword{
				output: "",
				err:    errs.ErrHashPassword,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_ERR_JWT",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrGenerateJWTFail,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    errs.ErrGenerateJWTFail,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_ERR_RE_JWT",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrGenerateJWTFail,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    errs.ErrGenerateJWTFail,
			},
			pepareRegister: PepareRegister{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
		},
		{
			name: "TEST_REGISTER_FAILS",
			input: dto.UserRegisterDto{
				Username: "username",
				Email:    "email",
				Password: "password",
			},
			output: nil,
			err:    errs.ErrRegisterUsername,
			peparePassowrd: PeparePassword{
				output: "password",
				err:    nil,
			},
			pepareCheckUser: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     nil,
			},
			pepareJwt: PepareJwt{
				output: "token",
				err:    nil,
			},
			pepareReJwt: PepareJwt{
				output: "retoken",
				err:    nil,
			},
			pepareRegister: PepareRegister{
				output: nil,
				err:    gorm.ErrDuplicatedKey,
			},
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mail := service.NewMailServiceMock()
			mailRepo := repository.NewMailRepositoryMock()
			// password
			passwordHelper := helpers.NewHelperMock()
			passwordHelper.On("HashPassword").Return(tc.peparePassowrd.output, tc.peparePassowrd.err)

			// logger
			logger := logger.NewLoggerMock()

			// user
			user := repository.NewUserRepositoryMock()
			user.On("CheckUsername").Return(tc.pepareCheckUser.output1, tc.pepareCheckUser.output2, tc.pepareCheckUser.err)
			user.On("CheckEmail").Return(tc.pepareCheckEmail.output1, tc.pepareCheckEmail.output2, tc.pepareCheckEmail.err)
			user.On("Register").Return(tc.pepareRegister.output, tc.pepareRegister.err)

			// jwt
			jwtService := service.NewJwtServiceMock()
			jwtService.On("GenerateToken").Return(tc.pepareJwt.output, tc.pepareJwt.err)
			jwtService.On("GenerateRefreshToken").Return(tc.pepareReJwt.output, tc.pepareReJwt.err)

			// service
			auth := service.NewAuth(logger, user, jwtService, passwordHelper, mail, mailRepo)
			res, err := auth.Register(ctx, tc.input)
			// test
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.output, res)
		})
	}
}

func TestVerifyEmail(t *testing.T) {
	type (
		PepareCheckUsernameAndEmail struct {
			output1 *model.User
			output2 bool
			err     error
		}

		PepareVerifyUser struct {
			output *model.User
			err    error
		}
		PepareJwt struct {
			output string
			err    error
		}
		TestCase struct {
			name              string
			input             dto.VerifyEmailDto
			output            *response.UserResponse
			err               error
			pepareCheckEmail  PepareCheckUsernameAndEmail
			pepareJwtGetemail PepareJwt
			pepareVerifyEmail PepareVerifyUser
			parseJsonErr      error
		}
	)

	testCases := []TestCase{
		{
			name: "TEST_SUCCESS",
			input: dto.VerifyEmailDto{
				Token: "token",
			},
			output: &response.UserResponse{CommonModel: response.CommonModel{ID: 1}},
			err:    nil,
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareJwtGetemail: PepareJwt{
				output: "email",
				err:    nil,
			},
			pepareVerifyEmail: PepareVerifyUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			parseJsonErr: nil,
		},
		{
			name: "TEST_JWT_ERROR",
			input: dto.VerifyEmailDto{
				Token: "token",
			},
			output: nil,
			err:    errs.ErrGenerateJWTFail,
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareJwtGetemail: PepareJwt{
				output: "",
				err:    errs.ErrGenerateJWTFail,
			},
			pepareVerifyEmail: PepareVerifyUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			parseJsonErr: nil,
		},
		{
			name: "TEST_USER_NOT_EXIST",
			input: dto.VerifyEmailDto{
				Token: "token",
			},
			output: nil,
			err:    errs.ErrBadRequest,
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: false,
				err:     gorm.ErrRecordNotFound,
			},
			pepareJwtGetemail: PepareJwt{
				output: "email",
				err:    nil,
			},
			pepareVerifyEmail: PepareVerifyUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			parseJsonErr: nil,
		},
		{
			name: "TEST_VERIFY_EMAIL_ERROR",
			input: dto.VerifyEmailDto{
				Token: "token",
			},
			output: nil,
			err:    errs.ErrVerifyEmail,
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareJwtGetemail: PepareJwt{
				output: "email",
				err:    nil,
			},
			pepareVerifyEmail: PepareVerifyUser{
				output: nil,
				err:    gorm.ErrDuplicatedKey,
			},
			parseJsonErr: nil,
		},
		{
			name: "TEST_PARSE_JSON_ERROR",
			input: dto.VerifyEmailDto{
				Token: "token",
			},
			output: nil,
			err:    errs.ErrInvalid,
			pepareCheckEmail: PepareCheckUsernameAndEmail{
				output1: nil,
				output2: true,
				err:     nil,
			},
			pepareJwtGetemail: PepareJwt{
				output: "email",
				err:    nil,
			},
			pepareVerifyEmail: PepareVerifyUser{
				output: &model.User{CommonModel: model.CommonModel{ID: 1}},
				err:    nil,
			},
			parseJsonErr: errs.ErrInternal,
		},
	}
	ctx := context.Background()
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mail := service.NewMailServiceMock()
			mailRepo := repository.NewMailRepositoryMock()
			// password
			passwordHelper := helpers.NewHelperMock()
			passwordHelper.On("ParseJson").Return(tc.parseJsonErr)

			// logger
			logger := logger.NewLoggerMock()

			// user
			user := repository.NewUserRepositoryMock()
			user.On("VerifyEmail").Return(tc.pepareVerifyEmail.output, tc.pepareVerifyEmail.err)
			user.On("CheckEmail").Return(tc.pepareCheckEmail.output1, tc.pepareCheckEmail.output2, tc.pepareCheckEmail.err)

			// jwt
			jwtService := service.NewJwtServiceMock()
			jwtService.On("GetEmailFormToken").Return(tc.pepareJwtGetemail.output, tc.pepareJwtGetemail.err)

			// service
			auth := service.NewAuth(logger, user, jwtService, passwordHelper, mail, mailRepo)
			res, err := auth.VerifyEmail(ctx, tc.input)
			// test
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.output, res)
		})
	}
}
