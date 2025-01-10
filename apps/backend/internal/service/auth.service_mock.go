package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/response"
)

type authServiceMock struct {
	mock.Mock
}

// Login implements Auth.
func (a *authServiceMock) Login(ctx context.Context, body dto.LoginDto) (res *response.AuthResponse, err error) {
	args := a.Called()
	return args.Get(0).(*response.AuthResponse), args.Error(1)
}

// Register implements Auth.
func (a *authServiceMock) Register(ctx context.Context, user dto.UserRegisterDto) (res *response.AuthResponse, err error) {
	args := a.Called()
	return args.Get(0).(*response.AuthResponse), args.Error(1)
}

// VerifyEmail implements Auth.
func (a *authServiceMock) VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *response.UserResponse, err error) {
	args := a.Called()
	return args.Get(0).(*response.UserResponse), args.Error(1)
}

func NewAuthServiceMock() *authServiceMock {
	return &authServiceMock{}
}
