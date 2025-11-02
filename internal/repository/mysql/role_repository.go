package mysql

import (
	"context"
	"database/sql"
	"rbac/internal/domain"
	"rbac/internal/repository"
)

type mysqlRoleRepository struct {
	db repository.DBTX
}

func NewRoleRepository(db repository.DBTX) repository.RoleRepository {
	return &mysqlRoleRepository{db: db}
}

func (r *mysqlRoleRepository) FindByName(ctx context.Context, name string) (*domain.Role, error) {
	query := "SELECT id, name FROM roles WHERE name = ?"
	row := r.db.QueryRowContext(ctx, query, name)

	var role domain.Role
	err := row.Scan(&role.ID, &role.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repository.ErrNotFound
		}
		return nil, err
	}
	return &role, nil
}