package ports

import (
	"context"

	"github.com/youruser/yourproject/internal/core/domain"
)

// RBACRepository defines the interface for RBAC data access
type RBACRepository interface {
	// Role operations
	GetRoleByID(ctx context.Context, id string) (*domain.Role, error)
	GetRoleByName(ctx context.Context, name string) (*domain.Role, error)
	GetAllRoles(ctx context.Context) ([]domain.Role, error)
	CreateRole(ctx context.Context, role *domain.Role) error
	UpdateRole(ctx context.Context, role *domain.Role) error
	DeleteRole(ctx context.Context, id string) error

	// Permission operations
	GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error)
	GetPermissionByName(ctx context.Context, name string) (*domain.Permission, error)
	GetAllPermissions(ctx context.Context) ([]domain.Permission, error)
	CreatePermission(ctx context.Context, permission *domain.Permission) error
	DeletePermission(ctx context.Context, id string) error

	// Role-Permission associations
	AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error
	RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error
	GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error)

	// User-Role associations
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RemoveRoleFromUser(ctx context.Context, userID, roleID string) error
	GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error)
	GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error)
	UserHasPermission(ctx context.Context, userID, permissionName string) (bool, error)
	UserHasRole(ctx context.Context, userID, roleName string) (bool, error)
}
