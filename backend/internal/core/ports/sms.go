package ports

import "context"

type SMSGateway interface {
	// SendOTP sends a one-time password to the given phone number
	SendOTP(ctx context.Context, phoneNumber string, code string) error
}
