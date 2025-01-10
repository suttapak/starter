package service

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/suttapak/starter/internal/dto"
)

type (
	mailServiceMock struct {
		mock.Mock
	}
)

// NewRequest implements Email.
func (m *mailServiceMock) NewRequest(to []string, subject string) Email {
	return m
}

// ParseTemplate implements Email.
func (m *mailServiceMock) ParseTemplate(ctx context.Context, file string, data any) error {
	args := m.Called()
	return args.Error(0)
}

// ParseVerifyEmailTemplate implements Email.
func (m *mailServiceMock) ParseVerifyEmailTemplate(ctx context.Context, body *dto.VerifyEmailTemplateDataDto) error {
	args := m.Called()
	return args.Error(0)
}

// SendMail implements Email.
func (m *mailServiceMock) SendMail(ctx context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func NewMailServiceMock() *mailServiceMock {
	return &mailServiceMock{}
}
