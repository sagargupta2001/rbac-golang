package repository

import (
	"context"
	"database/sql"
	"rbac/internal/domain"
)

// DBTX is an interface for both *sql.DB and *sql.Tx
type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

// UserRepository defines the methods for interacting with user data
type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByUsername(ctx context.Context, username string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	// RBAC-specific
	AssignRole(ctx context.Context, userID, roleID int64) error
	GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
}

// RoleRepository defines methods for roles and permissions
type RoleRepository interface {
	FindByName(ctx context.Context, name string) (*domain.Role, error)
}

// ProductRepository defines the methods for interacting with product data
type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	FindByID(ctx context.Context, id int64) (*domain.Product, error)
	// Add other CRUD methods (FindAll, Update, Delete) as needed
}
