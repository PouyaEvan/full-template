package ports

import "context"

type EmailService interface {
	SendEmail(ctx context.Context, to []string, subject string, body string) error
	SendResetPasswordEmail(ctx context.Context, to string, resetToken string) error
}
