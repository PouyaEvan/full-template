package http

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/youruser/yourproject/internal/core/domain"
	"github.com/youruser/yourproject/internal/core/ports"
	"github.com/youruser/yourproject/pkg/auth"
)

type AuthHandler struct {
	SMSGateway ports.SMSGateway
	Redis      *redis.Client
	UserRepo   ports.UserRepository
}

func NewAuthHandler(sms ports.SMSGateway, rdb *redis.Client, userRepo ports.UserRepository) *AuthHandler {
	return &AuthHandler{
		SMSGateway: sms,
		Redis:      rdb,
		UserRepo:   userRepo,
	}
}

func (h *AuthHandler) SendOTP(c *fiber.Ctx) error {
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Generate 6 digit code
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Store in Redis (5 min expiration)
	err := h.Redis.Set(context.Background(), "otp:"+req.Phone, code, 5*time.Minute).Err()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to store OTP"})
	}

	// Send SMS
	err = h.SMSGateway.SendOTP(context.Background(), req.Phone, code)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to send SMS: " + err.Error()})
	}

	return c.JSON(fiber.Map{"message": "OTP sent"})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	type Request struct {
		Phone string `json:"phone"`
		Code  string `json:"code"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check Redis
	storedCode, err := h.Redis.Get(context.Background(), "otp:"+req.Phone).Result()
	if err == redis.Nil {
		return c.Status(400).JSON(fiber.Map{"error": "OTP expired or not found"})
	} else if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Internal error"})
	}

	if storedCode != req.Code {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid OTP"})
	}

	// Clear OTP
	h.Redis.Del(context.Background(), "otp:"+req.Phone)

	// Find or Create User
	dummyEmail := req.Phone + "@example.com"
	user, err := h.UserRepo.GetByEmail(context.Background(), dummyEmail)
	if err != nil {
		// If not found, create
		user = &domain.User{
			Phone:     req.Phone,
			Email:     dummyEmail,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		if err := h.UserRepo.Create(context.Background(), user); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to create user"})
		}
	}

	// Check 2FA
	if user.IsTwoFactorEnabled {
		tempToken, _ := generateToken(user.ID, true)
		return c.JSON(fiber.Map{
			"2fa_required": true,
			"temp_token":   tempToken,
		})
	}

	// Generate final JWT
	token, _ := generateToken(user.ID, false)

	return c.JSON(fiber.Map{"token": token})
}

func (h *AuthHandler) Setup2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	user, err := h.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	key, qrBytes, err := auth.GenerateTOTPSecret(auth.TOTPConfig{
		Issuer:      "YourApp",
		AccountName: user.Email,
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate TOTP secret"})
	}

	// Store secret temporarily or permanently (but disabled)
	user.TwoFactorSecret = key.Secret()
	if err := h.UserRepo.Update(context.Background(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save secret"})
	}

	c.Set("Content-Type", "image/png")
	return c.Send(qrBytes)
}

func (h *AuthHandler) Enable2FA(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	type Request struct {
		Code string `json:"code"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := h.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if !auth.ValidateTOTP(req.Code, user.TwoFactorSecret) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid TOTP code"})
	}

	// Generate backup codes
	backupCodes, err := auth.GenerateBackupCodes(10)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to generate backup codes"})
	}

	user.IsTwoFactorEnabled = true
	user.TwoFactorBackupCodes = backupCodes
	if err := h.UserRepo.Update(context.Background(), user); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to enable 2FA"})
	}

	return c.JSON(fiber.Map{
		"message":      "2FA enabled",
		"backup_codes": backupCodes,
	})
}

func (h *AuthHandler) Verify2FALogin(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	type Request struct {
		Code string `json:"code"`
	}
	var req Request
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := h.UserRepo.GetByID(context.Background(), userID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	}

	if !auth.ValidateTOTP(req.Code, user.TwoFactorSecret) {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid TOTP code"})
	}

	// Generate final JWT
	token, _ := generateToken(user.ID, false)
	return c.JSON(fiber.Map{"token": token})
}

func generateToken(userID string, isTemp bool) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(72 * time.Hour).Unix(),
	}
	if isTemp {
		claims["is_temp"] = true
		claims["exp"] = time.Now().Add(5 * time.Minute).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("your-256-bit-secret"))
}

// OAuth2 Callback Placeholder
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// In a real app, you'd exchange the code for a token using golang.org/x/oauth2
	code := c.Query("code")
	return c.JSON(fiber.Map{"message": "Google Login Successful", "code": code})
}
