package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/suttapak/starter/internal/dto"
)

type authServiceMock struct {
	mock.Mock
}

// Login implements Auth.
func (a *authServiceMock) Login(ctx context.Context, body dto.LoginDto) (res *AuthResponse, err error) {
	args := a.Called()
	return args.Get(0).(*AuthResponse), args.Error(1)
}

// Register implements Auth.
func (a *authServiceMock) Register(ctx context.Context, user dto.UserRegisterDto) (res *AuthResponse, err error) {
	args := a.Called()
	return args.Get(0).(*AuthResponse), args.Error(1)
}

// VerifyEmail implements Auth.
func (a *authServiceMock) VerifyEmail(ctx context.Context, body dto.VerifyEmailDto) (res *UserResponse, err error) {
	args := a.Called()
	return args.Get(0).(*UserResponse), args.Error(1)
}

func NewAuthServiceMock() *authServiceMock {
	return &authServiceMock{}
}
