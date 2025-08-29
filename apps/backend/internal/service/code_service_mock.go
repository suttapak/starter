package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/suttapak/starter/internal/model"
)

type codeServiceMock struct {
	mock.Mock
}

// GenerateLotCode implements CodeService.
func (c *codeServiceMock) GenerateLotCode(ctx context.Context, productId uint) (string, error) {
	args := c.Called()
	return args.String(0), args.Error(1)
}

// GenerateProductCode implements CodeService.
func (c *codeServiceMock) GenerateProductCode(ctx context.Context, teamId uint) (string, error) {
	args := c.Called()
	return args.String(0), args.Error(1)
}

// GenerateTransactionCode implements CodeService.
func (c *codeServiceMock) GenerateTransactionCode(ctx context.Context, transactionType model.EntityType, teamId uint) (string, error) {
	args := c.Called()
	return args.String(0), args.Error(1)
}

func NewCodeServiceMock() *codeServiceMock {
	return &codeServiceMock{}
}
