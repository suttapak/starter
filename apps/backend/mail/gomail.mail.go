package mail

import (
	"crypto/tls"
	"labostack/config"

	"gopkg.in/gomail.v2"
)

func newMail(cfg *config.Config) (*gomail.Dialer, error) {
	mailer := gomail.NewDialer(
		cfg.MAIL.HOST,
		587,
		cfg.MAIL.USERNAME,
		cfg.MAIL.PASSWORD,
	)
	mailer.TLSConfig = &tls.Config{InsecureSkipVerify: true}
	return mailer, nil
}
