package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youruser/yourproject/internal/core/domain"
	"github.com/youruser/yourproject/internal/core/ports"
)

type RBACRepository struct {
	db *pgxpool.Pool
}

func NewRBACRepository(db *pgxpool.Pool) ports.RBACRepository {
	return &RBACRepository{db: db}
}

// Role operations

func (r *RBACRepository) GetRoleByID(ctx context.Context, id string) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE id = $1`

	var role domain.Role
	err := r.db.QueryRow(ctx, query, id).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get permissions for the role
	permissions, err := r.GetRolePermissions(ctx, id)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

func (r *RBACRepository) GetRoleByName(ctx context.Context, name string) (*domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles WHERE name = $1`

	var role domain.Role
	err := r.db.QueryRow(ctx, query, name).Scan(
		&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	// Get permissions for the role
	permissions, err := r.GetRolePermissions(ctx, role.ID)
	if err != nil {
		return nil, err
	}
	role.Permissions = permissions

	return &role, nil
}

func (r *RBACRepository) GetAllRoles(ctx context.Context) ([]domain.Role, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM roles ORDER BY name`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *RBACRepository) CreateRole(ctx context.Context, role *domain.Role) error {
	query := `INSERT INTO roles (name, description, created_at, updated_at) VALUES ($1, $2, $3, $4) RETURNING id`

	return r.db.QueryRow(ctx, query, role.Name, role.Description, role.CreatedAt, role.UpdatedAt).Scan(&role.ID)
}

func (r *RBACRepository) UpdateRole(ctx context.Context, role *domain.Role) error {
	query := `UPDATE roles SET name = $1, description = $2, updated_at = $3 WHERE id = $4`

	_, err := r.db.Exec(ctx, query, role.Name, role.Description, role.UpdatedAt, role.ID)
	return err
}

func (r *RBACRepository) DeleteRole(ctx context.Context, id string) error {
	query := `DELETE FROM roles WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

// Permission operations

func (r *RBACRepository) GetPermissionByID(ctx context.Context, id string) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at FROM permissions WHERE id = $1`

	var perm domain.Permission
	err := r.db.QueryRow(ctx, query, id).Scan(
		&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &perm, nil
}

func (r *RBACRepository) GetPermissionByName(ctx context.Context, name string) (*domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at FROM permissions WHERE name = $1`

	var perm domain.Permission
	err := r.db.QueryRow(ctx, query, name).Scan(
		&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &perm, nil
}

func (r *RBACRepository) GetAllPermissions(ctx context.Context) ([]domain.Permission, error) {
	query := `SELECT id, name, description, resource, action, created_at FROM permissions ORDER BY resource, action`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var perm domain.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

func (r *RBACRepository) CreatePermission(ctx context.Context, permission *domain.Permission) error {
	query := `INSERT INTO permissions (name, description, resource, action, created_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	return r.db.QueryRow(ctx, query, permission.Name, permission.Description, permission.Resource, permission.Action, permission.CreatedAt).Scan(&permission.ID)
}

func (r *RBACRepository) DeletePermission(ctx context.Context, id string) error {
	query := `DELETE FROM permissions WHERE id = $1`

	_, err := r.db.Exec(ctx, query, id)
	return err
}

// Role-Permission associations

func (r *RBACRepository) AssignPermissionToRole(ctx context.Context, roleID, permissionID string) error {
	query := `INSERT INTO role_permissions (role_id, permission_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	return err
}

func (r *RBACRepository) RemovePermissionFromRole(ctx context.Context, roleID, permissionID string) error {
	query := `DELETE FROM role_permissions WHERE role_id = $1 AND permission_id = $2`

	_, err := r.db.Exec(ctx, query, roleID, permissionID)
	return err
}

func (r *RBACRepository) GetRolePermissions(ctx context.Context, roleID string) ([]domain.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.resource, p.action, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.resource, p.action`

	rows, err := r.db.Query(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var perm domain.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

// User-Role associations

func (r *RBACRepository) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	query := `INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING`

	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *RBACRepository) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	query := `DELETE FROM user_roles WHERE user_id = $1 AND role_id = $2`

	_, err := r.db.Exec(ctx, query, userID, roleID)
	return err
}

func (r *RBACRepository) GetUserRoles(ctx context.Context, userID string) ([]domain.Role, error) {
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		INNER JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []domain.Role
	for rows.Next() {
		var role domain.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt, &role.UpdatedAt); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	return roles, rows.Err()
}

func (r *RBACRepository) GetUserPermissions(ctx context.Context, userID string) ([]domain.Permission, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.resource, p.action, p.created_at
		FROM permissions p
		INNER JOIN role_permissions rp ON p.id = rp.permission_id
		INNER JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.resource, p.action`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []domain.Permission
	for rows.Next() {
		var perm domain.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description, &perm.Resource, &perm.Action, &perm.CreatedAt); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}

	return permissions, rows.Err()
}

func (r *RBACRepository) UserHasPermission(ctx context.Context, userID, permissionName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			INNER JOIN role_permissions rp ON p.id = rp.permission_id
			INNER JOIN user_roles ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = $1 AND p.name = $2
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, permissionName).Scan(&exists)
	return exists, err
}

func (r *RBACRepository) UserHasRole(ctx context.Context, userID, roleName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM user_roles ur
			INNER JOIN roles r ON ur.role_id = r.id
			WHERE ur.user_id = $1 AND r.name = $2
		)`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID, roleName).Scan(&exists)
	return exists, err
}
