package handler

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/youruser/yourproject/internal/core/ports"
)

type AuthHandler struct {
	SMSGateway ports.SMSGateway
	Redis      *redis.Client
}

func NewAuthHandler(sms ports.SMSGateway, rdb *redis.Client) *AuthHandler {
	return &AuthHandler{
		SMSGateway: sms,
		Redis:      rdb,
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

	// Generate JWT (Mocked here)
	token := "mock-jwt-token-for-" + req.Phone

	return c.JSON(fiber.Map{"token": token})
}

// OAuth2 Callback Placeholder
func (h *AuthHandler) GoogleCallback(c *fiber.Ctx) error {
	// In a real app, you'd exchange the code for a token using golang.org/x/oauth2
	code := c.Query("code")
	return c.JSON(fiber.Map{"message": "Google Login Successful", "code": code})
}
