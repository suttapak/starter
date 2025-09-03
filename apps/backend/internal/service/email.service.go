package service

import (
	"bytes"
	"context"
	"html/template"

	"github.com/go-playground/validator/v10"
	"github.com/suttapak/starter/internal/repository"

	"gopkg.in/gomail.v2"
)

type (
	RequestApproveTransactionDto struct {
		Team        string `validate:"required"`
		Code        string `validate:"required"`
		Ref         string
		User        string `validate:"required"`
		RequestDate string `validate:"required"`
		TotalPrice  string `validate:"required"`
		Remark      string
		ApproveURL  string `validate:"required,url"`
	}

	RejectAndApproveTransactionDto struct {
		Code        string `validate:"required"`
		Ref         string
		User        string `validate:"required"`
		RequestDate string `validate:"required"`
		TotalPrice  string `validate:"required"`
		Remark      string
		ApproveURL  string `validate:"required,url"`
	}
	VerifyEmailTemplateDataDto struct {
		Email           string
		VerifyEmailLink string
	}
	InviteTeamMemberTemplateDataDto struct {
		TeamName     string
		JoinTeamLink string
	}
)

func (r *RequestApproveTransactionDto) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}
func (r *RejectAndApproveTransactionDto) Validate() error {
	validate := validator.New()
	return validate.Struct(r)
}

type (
	Email interface {
		NewRequest(to []string, subject string) Email
		SendMail(ctx context.Context) error
		ParseVerifyEmailTemplate(ctx context.Context, body *VerifyEmailTemplateDataDto) Email
		ParseInviteTeamMemberTemplate(ctx context.Context, body *InviteTeamMemberTemplateDataDto) Email
		ParseRequestApproveTransactionTemplate(ctx context.Context, body *RequestApproveTransactionDto) Email
		ParseApproveTransactionTemplate(ctx context.Context, body *RejectAndApproveTransactionDto) Email
		ParseRejectTransactionTemplate(ctx context.Context, body *RejectAndApproveTransactionDto) Email
	}
	email struct {
		to      []string
		subject string
		body    string
		mail    repository.MailRepository
		err     error
	}
)

// ParseApproveTransactionTemplate implements Email.
func (e *email) ParseApproveTransactionTemplate(ctx context.Context, body *RejectAndApproveTransactionDto) Email {
	const (
		templateFile = "./template/mail/approve-transaction.html"
	)
	if err := body.Validate(); err != nil {
		e.err = err
		return e
	}
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// ParseRejectTransactionTemplate implements Email.
func (e *email) ParseRejectTransactionTemplate(ctx context.Context, body *RejectAndApproveTransactionDto) Email {
	const (
		templateFile = "./template/mail/reject-transaction.html"
	)
	if err := body.Validate(); err != nil {
		e.err = err
		return e
	}
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// ParseRequestApproveTransactionTemplate implements Email.
func (e *email) ParseRequestApproveTransactionTemplate(ctx context.Context, body *RequestApproveTransactionDto) Email {
	const (
		templateFile = "./template/mail/request-approve-transaction.html"
	)
	if err := body.Validate(); err != nil {
		e.err = err
		return e
	}
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// ParseInviteTeamMemberTemplate implements Email.
func (e *email) ParseInviteTeamMemberTemplate(ctx context.Context, body *InviteTeamMemberTemplateDataDto) Email {
	const (
		templateFile = "./template/mail/join-team.html"
	)
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// SendMail implements Email.
func (e *email) SendMail(ctx context.Context) error {
	if e.err != nil {
		return e.err
	}
	msg := gomail.NewMessage()
	msg.SetHeader("From", "noreply@labotron.co.th")
	msg.SetHeader("To", e.to...)
	msg.SetHeader("Subject", e.subject)
	msg.SetBody("text/html", e.body)
	if err := e.mail.Send(ctx, msg); err != nil {
		return err
	}
	return nil

}

// ParseVerifyEmailTemplate implements Email.
func (e *email) ParseVerifyEmailTemplate(ctx context.Context, body *VerifyEmailTemplateDataDto) Email {
	const (
		templateFile = "./template/mail/register.html"
	)
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// ParseTemplate implements Email.
func (e *email) parseTemplate(ctx context.Context, file string, data any) error {
	_ = ctx
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
