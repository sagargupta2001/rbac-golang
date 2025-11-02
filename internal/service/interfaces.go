package service

import (
	"context"
	"rbac/internal/domain"
)

// AuthService handles user registration and login
type AuthService interface {
	Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error)
	Login(ctx context.Context, req domain.LoginRequest) (string, error)
}

// RBACService handles permission checks
type RBACService interface {
	CheckPermission(ctx context.Context, userID int64, requiredPermission string) (bool, error)
}

// ProductService handles product-related business logic
type ProductService interface {
	CreateProduct(ctx context.Context, req domain.CreateProductRequest, userID int64) (*domain.Product, error)
	// GetProduct(ctx context.Context, id int64) (*domain.Product, error)
}
