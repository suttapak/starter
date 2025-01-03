package service

import (
	"bytes"
	"context"
	"html/template"
	"github.com/suttapak/starter/internal/dto"
	"github.com/suttapak/starter/internal/repository"

	"gopkg.in/gomail.v2"
)

type (
	Email interface {
		NewRequest(to []string, subject string) Email
		ParseTemplate(ctx context.Context, file string, data any) error
		ParseVerifyEmailTemplate(ctx context.Context, body *dto.VerifyEmailTemplateDataDto) error
		SendMail(ctx context.Context) error
	}
	email struct {
		to      []string
		subject string
		body    string

		mail repository.MailRepository
	}
)

// SendMail implements Email.
func (e *email) SendMail(ctx context.Context) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "matee@labotron.co.th")
	msg.SetHeader("To", e.to...)
	msg.SetHeader("Subject", e.subject)
	msg.SetBody("text/html", e.body)
	if err := e.mail.Send(ctx, msg); err != nil {
		return err
	}
	return nil
}

// ParseVerifyEmailTemplate implements Email.
func (e *email) ParseVerifyEmailTemplate(ctx context.Context, body *dto.VerifyEmailTemplateDataDto) error {
	const (
		templateFile = "./template/mail/register.html"
	)
	if err := e.ParseTemplate(ctx, templateFile, body); err != nil {
		return err
	}
	return nil
}

// ParseTemplate implements Email.
func (e *email) ParseTemplate(ctx context.Context, file string, data any) error {
	t, err := template.ParseFiles(file)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	e.body = buf.String()
	return nil
}

// NewRequest implements Email.
func (e *email) NewRequest(to []string, subject string) Email {
	return &email{
		to:      to,
		subject: subject,
		mail:    e.mail,
	}
}

func newEmailService(mail repository.MailRepository) Email {
	return &email{mail: mail}
}
