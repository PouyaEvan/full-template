package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// Role-Based Access Control
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Assumes Protected() middleware has already run and set the user claims
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)

		userRole := claims["role"].(string)

		if userRole != role && userRole != "admin" { // Admin usually has access to everything
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}
