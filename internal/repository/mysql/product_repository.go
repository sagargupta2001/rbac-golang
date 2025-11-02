package mysql

import (
	"context"
	"rbac/internal/domain"
	"rbac/internal/repository"
)

type mysqlProductRepository struct {
	db repository.DBTX
}

func NewProductRepository(db repository.DBTX) repository.ProductRepository {
	return &mysqlProductRepository{db: db}
}

func (r *mysqlProductRepository) Create(ctx context.Context, product *domain.Product) error {
	query := "INSERT INTO products (name, price, created_by_user) VALUES (?, ?, ?)"
	res, err := r.db.ExecContext(ctx, query, product.Name, product.Price, product.CreatedByUserID)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	product.ID = id
	return nil
}

func (r *mysqlProductRepository) FindByID(ctx context.Context, id int64) (*domain.Product, error) {
	// ... implementation ...
	// (For brevity, you can implement this similar to FindByID in user repo)
	return nil, nil
}