package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/youruser/yourproject/internal/core/domain"
	"github.com/youruser/yourproject/internal/core/ports"
)

// RBACMiddleware holds the RBAC repository for authorization checks
type RBACMiddleware struct {
	rbacRepo ports.RBACRepository
}

// NewRBACMiddleware creates a new RBAC middleware instance
func NewRBACMiddleware(rbacRepo ports.RBACRepository) *RBACMiddleware {
	return &RBACMiddleware{rbacRepo: rbacRepo}
}

// RequireRole checks if the user has a specific role
func (m *RBACMiddleware) RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		ctx := c.Context()
		hasRole, err := m.rbacRepo.UserHasRole(ctx, userID, role)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user role",
			})
		}

		if !hasRole {
			// Check if user is admin (admin has access to everything)
			isAdmin, _ := m.rbacRepo.UserHasRole(ctx, userID, domain.RoleAdmin)
			if !isAdmin {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Insufficient permissions",
				})
			}
		}

		return c.Next()
	}
}

// RequirePermission checks if the user has a specific permission
func (m *RBACMiddleware) RequirePermission(permission string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		ctx := c.Context()
		hasPermission, err := m.rbacRepo.UserHasPermission(ctx, userID, permission)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check user permission",
			})
		}

		if !hasPermission {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}

// RequireAnyRole checks if the user has any of the specified roles
func (m *RBACMiddleware) RequireAnyRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		ctx := c.Context()
		for _, role := range roles {
			hasRole, err := m.rbacRepo.UserHasRole(ctx, userID, role)
			if err != nil {
				continue
			}
			if hasRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// RequireAnyPermission checks if the user has any of the specified permissions
func (m *RBACMiddleware) RequireAnyPermission(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		ctx := c.Context()
		for _, permission := range permissions {
			hasPerm, err := m.rbacRepo.UserHasPermission(ctx, userID, permission)
			if err != nil {
				continue
			}
			if hasPerm {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Insufficient permissions",
		})
	}
}

// RequireAllPermissions checks if the user has all of the specified permissions
func (m *RBACMiddleware) RequireAllPermissions(permissions ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID, ok := c.Locals("user_id").(string)
		if !ok || userID == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User not authenticated",
			})
		}

		ctx := c.Context()
		for _, permission := range permissions {
			hasPerm, err := m.rbacRepo.UserHasPermission(ctx, userID, permission)
			if err != nil || !hasPerm {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Insufficient permissions",
				})
			}
		}

		return c.Next()
	}
}
