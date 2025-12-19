package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base32"
	"image/png"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPConfig struct {
	Issuer      string
	AccountName string
}

// GenerateTOTPSecret generates a new TOTP key and returns the key object and its QR code as PNG bytes.
func GenerateTOTPSecret(config TOTPConfig) (*otp.Key, []byte, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      config.Issuer,
		AccountName: config.AccountName,
	})
	if err != nil {
		return nil, nil, err
	}

	// Convert to PNG
	var buf bytes.Buffer
	img, err := key.Image(200, 200)
	if err != nil {
		return nil, nil, err
	}
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, nil, err
	}

	return key, buf.Bytes(), nil
}

// ValidateTOTP validates a passcode against a secret.
func ValidateTOTP(passcode string, secret string) bool {
	return totp.Validate(passcode, secret)
}

// GenerateBackupCodes generates n random backup codes.
func GenerateBackupCodes(n int) ([]string, error) {
	codes := make([]string, n)
	for i := 0; i < n; i++ {
		// Generate 5 bytes of random data (enough for 8 base32 chars)
		b := make([]byte, 5)
		_, err := rand.Read(b)
		if err != nil {
			return nil, err
		}
		// Encode to base32
		codes[i] = base32.StdEncoding.EncodeToString(b)
	}
	return codes, nil
}
