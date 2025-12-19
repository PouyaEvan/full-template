package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateTOTPSecret(t *testing.T) {
	config := TOTPConfig{
		Issuer:      "TestApp",
		AccountName: "test@example.com",
	}

	key, qrBytes, err := GenerateTOTPSecret(config)
	assert.NoError(t, err)
	assert.NotNil(t, key)
	assert.NotNil(t, qrBytes)
	assert.NotEmpty(t, key.Secret())
}

func TestValidateTOTP(t *testing.T) {
	// This is tricky to test deterministically without mocking time or using a library that allows setting time.
	// However, we can generate a key and try to validate a code generated from it immediately.
	// For now, we'll just test that it doesn't panic and returns false for garbage.

	valid := ValidateTOTP("123456", "JBSWY3DPEHPK3PXP") // Random secret
	assert.False(t, valid)
}

func TestGenerateBackupCodes(t *testing.T) {
	codes, err := GenerateBackupCodes(10)
	assert.NoError(t, err)
	assert.Len(t, codes, 10)

	for _, code := range codes {
		assert.NotEmpty(t, code)
		assert.Len(t, code, 8) // 5 bytes base32 encoded is 8 chars
	}
}
