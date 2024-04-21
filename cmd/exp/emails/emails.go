package main

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/go-mail/mail/v2"
)

type config struct {
	Host     string `env:"SMTP_HOST,required"`
	Port     int    `env:"SMTP_PORT" envDefault:"587"`
	Username string `env:"SMTP_USERNAME,required"`
	Password string `env:"SMTP_PASSWORD,required"`
}

func main() {
	cfg := config{}
	err := env.Parse(&cfg)
	if err != nil {
		panic(err)
	}

	from := "test@lenslocked.com"
	to := "jon@calhoun.io"
	subject := "This is a test email"
	plaintext := "This is the body of the email"
	html := `<h1>Hello there buddy!</h1><p>This is the email</p><p>Hope you enjoy it</p>`

	msg := mail.NewMessage()
	msg.SetHeader("To", to)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", plaintext)
	msg.AddAlternative("text/html", html)
	// msg.WriteTo(os.Stdout)

	dialer := mail.NewDialer(cfg.Host, cfg.Port, cfg.Username, cfg.Password)
	err = dialer.DialAndSend(msg)
	if err != nil {
		panic(err)
	}
	fmt.Println("Message sent")
}
