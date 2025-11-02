package mysql

import (
	"context"
	"database/sql"
	"rbac/internal/domain"
	"rbac/internal/repository"
)

type mysqlUserRepository struct {
	db repository.DBTX
}

// NewUserRepository creates a new UserRepository
func NewUserRepository(db repository.DBTX) repository.UserRepository {
	return &mysqlUserRepository{db: db}
}

func (r *mysqlUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := "INSERT INTO users (username, email, password_hash) VALUES (?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.PasswordHash)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *mysqlUserRepository) FindByUsername(ctx context.Context, username string) (*domain.User, error) {
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE username = ?"
	row := r.db.QueryRowContext(ctx, query, username)
	
	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {
	query := "SELECT id, username, email, password_hash, created_at FROM users WHERE id = ?"
	row := r.db.QueryRowContext(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlUserRepository) AssignRole(ctx context.Context, userID, roleID int64) error {
	query := "INSERT INTO user_roles (user_id, role_id) VALUES (?, ?)"
	_, err := r.db.ExecContext(ctx, query, userID, roleID)
	return err
}

// GetUserPermissions is the core of our RBAC check
func (r *mysqlUserRepository) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	query := `
		SELECT DISTINCT p.name
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []string
	for rows.Next() {
		var permName string
		if err := rows.Scan(&permName); err != nil {
			return nil, err
		}
		permissions = append(permissions, permName)
	}

	return permissions, nil
}