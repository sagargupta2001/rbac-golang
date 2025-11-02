package service

import (
	"context"
	"rbac/internal/repository"
)

type rbacService struct {
	userRepo repository.UserRepository
}

// NewRBACService creates a new RBACService
func NewRBACService(userRepo repository.UserRepository) RBACService {
	return &rbacService{userRepo: userRepo}
}

// CheckPermission checks if a user has a specific permission
func (s *rbacService) CheckPermission(ctx context.Context, userID int64, requiredPermission string) (bool, error) {
	permissions, err := s.userRepo.GetUserPermissions(ctx, userID)
	if err != nil {
		return false, err
	}

	// Simple map-based lookup for O(1) average time complexity
	permMap := make(map[string]struct{})
	for _, p := range permissions {
		permMap[p] = struct{}{}
	}

	// Check if the required permission exists in the user's permissions
	if _, ok := permMap[requiredPermission]; ok {
		return true, nil
	}

	return false, nil // Forbidden
}