package domain

import (
	"errors"
	"regexp"
	"time"
)

var (
	ErrInvalidPhone = errors.New("invalid phone number")
	ErrInvalidOTP   = errors.New("invalid OTP")
)

type User struct {
	ID                   string
	Email                string
	Phone                string
	PasswordHash         string
	IsTwoFactorEnabled   bool
	TwoFactorSecret      string
	TwoFactorBackupCodes []string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

func NewUser(phone string) (*User, error) {
	if !isValidPhone(phone) {
		return nil, ErrInvalidPhone
	}
	return &User{
		Phone:     phone,
		CreatedAt: time.Now(),
	}, nil
}

func isValidPhone(phone string) bool {
	// Simple regex for Iranian mobile numbers
	re := regexp.MustCompile(`^09\d{9}$`)
	return re.MatchString(phone)
}
