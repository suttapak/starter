package repository

import (
	"context"

	"github.com/stretchr/testify/mock"
	gomail "gopkg.in/gomail.v2"
)

type (
	mailRepositoryMock struct {
		mock.Mock
	}
)

// Send implements MailRepository.
func (m *mailRepositoryMock) Send(ctx context.Context, message *gomail.Message) error {
	args := m.Called()
	return args.Error(0)
}

func NewMailRepositoryMock() *mailRepositoryMock {
	return &mailRepositoryMock{}
}
