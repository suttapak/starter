package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/mock"
)

type (
	jwtServiceMock struct {
		mock.Mock
	}
)

// CreateVerifyEmailToken implements JWTService.
func (j *jwtServiceMock) CreateVerifyEmailToken(ctx context.Context, email string) (token string, err error) {
	args := j.Called()
	return args.String(0), args.Error(1)
}

// GenerateRefreshToken implements JWTService.
func (j *jwtServiceMock) GenerateRefreshToken(ctx context.Context, userId uint) (token string, err error) {
	args := j.Called()
	return args.String(0), args.Error(1)
}

// GenerateToken implements JWTService.
func (j *jwtServiceMock) GenerateToken(ctx context.Context, userId uint) (token string, err error) {
	args := j.Called()
	return args.String(0), args.Error(1)
}

// GetEmailFormToken implements JWTService.
func (j *jwtServiceMock) GetEmailFormToken(ctx context.Context, token string) (email string, err error) {
	args := j.Called()
	return args.String(0), args.Error(1)
}

// GetUserIdFormToken implements JWTService.
func (j *jwtServiceMock) GetUserIdFormToken(ctx context.Context, token string) (uId uint, err error) {
	args := j.Called()
	return uint(args.Int(0)), args.Error(1)
}

// ParserToken implements JWTService.
func (j *jwtServiceMock) ParserToken(ctx context.Context, t string, secret string) (token *jwt.Token, err error) {
	args := j.Called()
	return args.Get(0).(*jwt.Token), args.Error(1)
}

func NewJwtServiceMock() *jwtServiceMock {
	return &jwtServiceMock{}
}
