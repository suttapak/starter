package service

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/fs"

	"github.com/go-playground/validator/v10"
	"github.com/suttapak/starter/internal/repository"
	"github.com/suttapak/starter/mtemplate"

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

func listFiles(fsys fs.FS, dir string, indent string) {
	entries, err := fs.ReadDir(fsys, dir)
	if err != nil {
		fmt.Printf("error reading %s: %v\n", dir, err)
		return
	}

	for _, e := range entries {
		name := e.Name()
		if e.IsDir() {
			fmt.Printf("%s📁 %s/\n", indent, name)
			// Recursive call for subdirectories
			listFiles(fsys, dir+"/"+name, indent+"  ")
		} else {
			fmt.Printf("%s📄 %s\n", indent, name)
		}
	}
}

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
	}
	email struct {
		to      []string
		subject string
		body    string
		mail    repository.MailRepository
		err     error
	}
)

// ParseInviteTeamMemberTemplate implements Email.
func (e *email) ParseInviteTeamMemberTemplate(ctx context.Context, body *InviteTeamMemberTemplateDataDto) Email {
	const (
		templateFile = "mail/join-team.html"
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
		templateFile = "mail/register.html"
	)
	if err := e.parseTemplate(ctx, templateFile, body); err != nil {
		e.err = err
	}
	return e
}

// ParseTemplate implements Email.
func (e *email) parseTemplate(ctx context.Context, file string, data any) error {
	listFiles(mtemplate.EmailTemplateFS, "mail", "")
	_ = ctx
	t, err := template.ParseFS(mtemplate.EmailTemplateFS, file)
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

func NewEmail(mail repository.MailRepository) Email {
	return &email{mail: mail}
}
