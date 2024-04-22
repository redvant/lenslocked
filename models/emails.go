package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@lenslock.com"
)

type SMTPConfig struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     int    `env:"SMTP_PORT" envDefault:"587"`
	Username string `env:"SMTP_USERNAME,required"`
	Password string `env:"SMTP_PASSWORD,required"`
}

type EmailService struct {
	// DefaultSender is used as the default sender when one isn't
	// provided for an email.
	DefaultSender string

	dialer *mail.Dialer
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}
	return &es
}
