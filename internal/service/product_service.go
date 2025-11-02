package service

import (
	"context"
	"rbac/internal/domain"
	"rbac/internal/repository"
)

type productService struct {
	productRepo repository.ProductRepository
}

// NewProductService creates a new ProductService
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) CreateProduct(ctx context.Context, req domain.CreateProductRequest, userID int64) (*domain.Product, error) {
	product := &domain.Product{
		Name:            req.Name,
		Price:           req.Price,
		CreatedByUserID: userID,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}