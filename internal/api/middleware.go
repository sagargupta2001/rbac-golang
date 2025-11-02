package api

import (
	"context"
	"net/http"
	"rbac/internal/service"
	"rbac/internal/utils"
	"strings"

	"github.com/gorilla/mux"
)

// CtxKey is a custom type for context keys to avoid collisions
type CtxKey string

const (
	// UserIDKey is the key for user ID in context
	UserIDKey CtxKey = "userID"
)

// AuthMiddleware validates the JWT token
func AuthMiddleware(jwtSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")
			if tokenString == authHeader { // No "Bearer " prefix
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			claims, err := utils.ValidateToken(tokenString, jwtSecret)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			// Add user ID to context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RBACMiddleware checks if the user has the required permission
// This is a "factory" that returns a middleware
func RBACMiddleware(rbacSvc service.RBACService, permission string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID, ok := r.Context().Value(UserIDKey).(int64)
			if !ok {
				// This should not happen if AuthMiddleware is applied first
				http.Error(w, "User ID not found in context", http.StatusInternalServerError)
				return
			}

			allowed, err := rbacSvc.CheckPermission(r.Context(), userID, permission)
			if err != nil {
				http.Error(w, "Failed to check permissions", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Forbidden: You do not have the required permission", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}