package email

import (
	"fmt"
	"log"
	"net/smtp"
	"net/url"

	"github.com/kondohiroki/go-grpc-boilerplate/config"
)

type EmailService struct {
	smtpHost string
	smtpPort int
	username string
	password string
	from     string
}

func NewEmailService(conf *config.Config) EmailService {
	return EmailService{
		smtpHost: conf.Services.Email.SmtpHost,
		smtpPort: conf.Services.Email.SmtpPort,
		username: conf.Services.Email.Username,
		password: url.QueryEscape(conf.Services.Email.Password),
		from:     conf.Services.Email.From,
	}

}

func (s *EmailService) SendEmail(to []string, subject string, body string) error {

	if len(to) == 0 {
		return fmt.Errorf("no recipient")
	}
	// Create the message
	message := []byte("To: " + to[0] + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-version: 1.0;\r\n" +
		"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
		"\r\n" +
		body + "\r\n")

	// Authentication
	auth := smtp.PlainAuth("", s.username, s.password, s.smtpHost)

	// Sending the email
	err := smtp.SendMail(fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort), auth, s.from, to, message)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
