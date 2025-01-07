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
		J struct {
			userId     uint
			token      string
			errToken   error
			reToken    string
			errReToken error
		}
		G struct {
			user *model.User
			err  error
		}
		M struct {
			getUserEmail G
			jwtService   J
		}
		TC struct {
			name        string
			input       dto.LoginDto
			expected    *response.AuthResponse
			errExpected error
			mockup      M
		}
	)
	p, _ := helpers.HashPassword("123412345")
	var testcases = []TC{
		{
			name: "TEST_SUCCESS",
			input: dto.LoginDto{
				UserNameEmail: "suttapak",
				Password:      "123412345",
			},
			expected: &response.AuthResponse{
				Token:        "token",
				RefreshToken: "reToken",
			},
			errExpected: nil,
			mockup: M{
				getUserEmail: G{
					user: &model.User{
						CommonModel: model.CommonModel{
							ID: 1,
						},
						Username: "suttapak",
						Password: p,
					},
					err: nil,
				},
				jwtService: J{
					userId:     1,
					token:      "token",
					errToken:   nil,
					reToken:    "reToken",
					errReToken: nil,
				},
			},
		},
		{
			name: "TEST_PASSWORD_INVALID",
			input: dto.LoginDto{
				UserNameEmail: "suttapak",
				Password:      "1234123456",
			},
			expected:    nil,
			errExpected: errs.ErrBadRequest,
			mockup: M{
				getUserEmail: G{
					user: &model.User{
						CommonModel: model.CommonModel{
							ID: 1,
						},
						Username: "suttapak",
						Password: p,
					},
					err: nil,
				},
				jwtService: J{
					userId:     1,
					token:      "token",
					errToken:   nil,
					reToken:    "reToken",
					errReToken: nil,
				},
			},
		},
		{
			name: "TEST_USERNAME_NOT_FOUND",
			input: dto.LoginDto{
				UserNameEmail: "suttapak",
				Password:      "123412345",
			},
			expected:    nil,
			errExpected: errs.ErrUnauthorized,
			mockup: M{
				getUserEmail: G{
					user: nil,
					err:  gorm.ErrRecordNotFound,
				},
				jwtService: J{
					userId:     1,
					token:      "token",
					errToken:   nil,
					reToken:    "reToken",
					errReToken: nil,
				},
			},
		},
	}
	ctx := context.Background()
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			userRepo := repository.NewUserRepositoryMock()
			userRepo.On("GetUserByEmailOrUsername", ctx, (*gorm.DB)(nil), tc.input.UserNameEmail).Return(tc.mockup.getUserEmail.user, tc.mockup.getUserEmail.err)
			jwtService := service.NewJwtServiceMock()
			jwtService.On("GenerateToken", ctx, tc.mockup.jwtService.userId).Return(tc.mockup.jwtService.token, tc.mockup.jwtService.errToken)
			jwtService.On("GenerateRefreshToken", ctx, tc.mockup.jwtService.userId).Return(tc.mockup.jwtService.reToken, tc.mockup.jwtService.errReToken)
			l := logger.NewLoggerMock()
			authService := service.NewAuth(l, userRepo, jwtService)
			res, err := authService.Login(ctx, tc.input)
			assert.Equal(t, tc.errExpected, err)
			assert.Equal(t, tc.expected, res)
		})
	}
}
