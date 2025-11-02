package service

import (
	"context"
	"errors"
	"rbac/internal/domain"
	"rbac/internal/repository"
	"rbac/internal/utils"
)

// authService is the implementation of AuthService
type authService struct {
	userRepo      repository.UserRepository
	roleRepo      repository.RoleRepository
	jwtSecret     string
	jwtExpiration int64
}

// NewAuthService creates a new AuthService
func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, jwtSecret string, jwtExp int64) AuthService {
	return &authService{
		userRepo:      userRepo,
		roleRepo:      roleRepo,
		jwtSecret:     jwtSecret,
		jwtExpiration: jwtExp,
	}
}

func (s *authService) Register(ctx context.Context, req domain.RegisterRequest) (*domain.User, error) {
	// Check if user already exists
	_, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err == nil {
		return nil, errors.New("username already taken")
	}
	if err != repository.ErrNotFound {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	// Create user in DB
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Assign default "user" role
	// Note: In a real app, you'd seed the "user" role in your migration
	role, err := s.roleRepo.FindByName(ctx, "user") // Assumes "user" role (ID 2) exists
	if err != nil {
		// This is a fallback. Seed your roles!
		// For this demo, let's assume 'user' role is ID 2 if not found.
		// A better way is to seed 'admin' and 'user' roles in migrations.
		// Let's just log a warning and continue, assigning a hardcoded role ID
		// if we can't find it. This is NOT good practice, but fine for this demo.
		// log.Printf("Warning: 'user' role not found, assigning ID 2. Please seed roles.")
		// We'll skip assignment if role not found, to avoid crash.
		if err == repository.ErrNotFound {
			// This is bad, the user will have NO roles.
			// Let's fix this by seeding.
			return nil, errors.New("default 'user' role not found. Please seed roles in DB.")
		}
		return nil, err
	}
	
	if err := s.userRepo.AssignRole(ctx, user.ID, role.ID); err != nil {
		// Log this, but registration was successful
		// log.Printf("Failed to assign default role to user %d: %v", user.ID, err)
		return nil, errors.New("failed to assign default role")
	}


	user.PasswordHash = "" // Clear password before returning
	return user, nil
}

func (s *authService) Login(ctx context.Context, req domain.LoginRequest) (string, error) {
	user, err := s.userRepo.FindByUsername(ctx, req.Username)
	if err != nil {
		if err == repository.ErrNotFound {
			return "", errors.New("invalid username or password")
		}
		return "", err
	}

	// Check password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return "", errors.New("invalid username or password")
	}

	// Generate JWT
	token, err := utils.GenerateToken(user.ID, s.jwtSecret, s.jwtExpiration)
	if err != nil {
		return "", err
	}

	return token, nil
}