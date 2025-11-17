package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"rbac/internal/api"
	"rbac/internal/config"
	"rbac/internal/repository/mysql"
	"rbac/internal/service"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

func main() {
	// --- 1. Load Configuration ---
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// --- 2. Initialize Database ---
	db, err := mysql.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// --- 3. Initialize Layers (Dependency Injection) ---

	// Repository Layer
	userRepo := mysql.NewUserRepository(db)
	roleRepo := mysql.NewRoleRepository(db)
	productRepo := mysql.NewProductRepository(db)

	// Service Layer
	authSvc := service.NewAuthService(userRepo, roleRepo, cfg.JWTSecret, cfg.JWTExpirationInHours)
	rbacSvc := service.NewRBACService(userRepo)
	productSvc := service.NewProductService(productRepo)
	graphqlSvc := service.NewGraphQLService()

	// API/Handler Layer
	apiHandler := api.NewAPIHandler(authSvc, rbacSvc, productSvc, graphqlSvc)

	// --- 4. Setup Router & Routes ---
	router := mux.NewRouter()
	apiHandler.RegisterRoutes(router, cfg.JWTSecret)

	// --- 5. Start HTTP Server (with Graceful Shutdown) ---
	srv := &http.Server{
		Addr:    cfg.ServerPort,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	// Wait for interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Create a context with timeout for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}