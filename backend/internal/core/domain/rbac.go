package domain

import (
	"time"
)

// Role represents a user role in the system
type Role struct {
	ID          string
	Name        string
	Description string
	Permissions []Permission
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Permission represents a permission that can be assigned to roles
type Permission struct {
	ID          string
	Name        string
	Description string
	Resource    string
	Action      string
	CreatedAt   time.Time
}

// Common role names
const (
	RoleAdmin     = "admin"
	RoleUser      = "user"
	RoleModerator = "moderator"
)

// Common permission actions
const (
	ActionRead   = "read"
	ActionWrite  = "write"
	ActionDelete = "delete"
	ActionAccess = "access"
	ActionManage = "manage"
)

// HasPermission checks if the role has a specific permission
func (r *Role) HasPermission(permissionName string) bool {
	for _, p := range r.Permissions {
		if p.Name == permissionName {
			return true
		}
	}
	return false
}

// HasResourceAccess checks if the role can perform an action on a resource
func (r *Role) HasResourceAccess(resource, action string) bool {
	for _, p := range r.Permissions {
		if p.Resource == resource && p.Action == action {
			return true
		}
	}
	return false
}
