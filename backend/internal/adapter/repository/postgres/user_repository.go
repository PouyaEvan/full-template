package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/youruser/yourproject/internal/core/domain"
	"github.com/youruser/yourproject/internal/core/ports"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) ports.UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (email, phone, password_hash, is_two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	err := r.db.QueryRow(ctx, query, user.Email, user.Phone, user.PasswordHash, user.IsTwoFactorEnabled, user.TwoFactorSecret, user.TwoFactorBackupCodes, user.CreatedAt, user.UpdatedAt).Scan(&user.ID)
	return err
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, email, phone, password_hash, is_two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at FROM users WHERE email = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PasswordHash,
		&user.IsTwoFactorEnabled, &user.TwoFactorSecret, &user.TwoFactorBackupCodes,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	query := `SELECT id, email, phone, password_hash, is_two_factor_enabled, two_factor_secret, two_factor_backup_codes, created_at, updated_at FROM users WHERE id = $1`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Email, &user.Phone, &user.PasswordHash,
		&user.IsTwoFactorEnabled, &user.TwoFactorSecret, &user.TwoFactorBackupCodes,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `UPDATE users SET email=$1, phone=$2, password_hash=$3, is_two_factor_enabled=$4, two_factor_secret=$5, two_factor_backup_codes=$6, updated_at=$7 WHERE id=$8`

	_, err := r.db.Exec(ctx, query, user.Email, user.Phone, user.PasswordHash, user.IsTwoFactorEnabled, user.TwoFactorSecret, user.TwoFactorBackupCodes, user.UpdatedAt, user.ID)
	return err
}
