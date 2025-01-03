package repository

import (
	"context"

	"gopkg.in/gomail.v2"
)

type (
	MailRepository interface {
		Send(ctx context.Context, message *gomail.Message) error
	}
	mailRepository struct {
		dialer *gomail.Dialer
	}
)

// Send implements MailRepository.
func (m *mailRepository) Send(ctx context.Context, message *gomail.Message) error {
	err := m.dialer.DialAndSend(message)
	return err
}

func newMailRepository(dialer *gomail.Dialer) MailRepository {
	return &mailRepository{dialer}
}
