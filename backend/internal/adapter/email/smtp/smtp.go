package smtp

import (
	"context"
	"fmt"
	"net/smtp"
	"os"
)

type SMTPAdapter struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

func NewSMTPAdapter() *SMTPAdapter {
	return &SMTPAdapter{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     os.Getenv("SMTP_PORT"),
		Username: os.Getenv("SMTP_USER"),
		Password: os.Getenv("SMTP_PASS"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

func (s *SMTPAdapter) SendEmail(ctx context.Context, to []string, subject string, body string) error {
	auth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: %s\r\n"+
		"\r\n"+
		"%s\r\n", to[0], subject, body))

	addr := fmt.Sprintf("%s:%s", s.Host, s.Port)

	// Note: In production, use a more robust library or worker queue
	return smtp.SendMail(addr, auth, s.From, to, msg)
}

func (s *SMTPAdapter) SendResetPasswordEmail(ctx context.Context, to string, resetToken string) error {
	subject := "Password Reset Request"
	body := fmt.Sprintf("Click here to reset your password: https://myapp.com/reset-password?token=%s", resetToken)
	return s.SendEmail(ctx, []string{to}, subject, body)
}
