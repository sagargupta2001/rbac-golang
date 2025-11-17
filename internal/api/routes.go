package api

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// RegisterRoutes sets up all routes for the application
func (h *APIHandler) RegisterRoutes(router *mux.Router, jwtSecret string) {
	// Create middleware instances
	auth := AuthMiddleware(jwtSecret)
	// Create RBAC middleware for specific permissions
	canCreateProduct := RBACMiddleware(h.rbacSvc, "create_product")
	canReadProduct := RBACMiddleware(h.rbacSvc, "read_product")
	// canDeleteUser := RBACMiddleware(h.rbacSvc, "delete_user") // Example

	// Public routes (Auth)
	router.HandleFunc("/register", h.RegisterHandler).Methods("POST")
	router.HandleFunc("/login", h.LoginHandler).Methods("POST")

	router.HandleFunc("/test-graphql/{code}", h.GetCountryHandler).Methods("GET")

	// Protected routes (Products)
	// We apply middleware in order: Auth (to get user) -> RBAC (to check perm)
	productRouter := router.PathPrefix("/products").Subrouter()
	productRouter.Use(auth) // All product routes require at least login

	// POST /products - Requires 'create_product' permission
	productRouter.HandleFunc("", h.CreateProductHandler).Methods("POST").Handler(
		canCreateProduct(http.HandlerFunc(h.CreateProductHandler)),
	)
	
	// GET /products/{id} - Requires 'read_product' permission
	productRouter.HandleFunc("/{id:[0-9]+}", h.GetProductHandler).Methods("GET").Handler(
		canReadProduct(http.HandlerFunc(h.GetProductHandler)),
	)
	
	// Example of a route only an admin could access
	// adminRouter := router.PathPrefix("/admin").Subrouter()
	// adminRouter.Use(auth, canDeleteUser)
	// adminRouter.HandleFunc("/users/{id}", h.DeleteUserHandler).Methods("DELETE")

	log.Println("Registered API routes")
}

// Dummy handler to satisfy the routes file
func (h *APIHandler) GetProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	// In a real app, you'd call h.productSvc.GetProduct(r.Context(), id)
	respondWithJSON(w, http.StatusOK, map[string]string{"message": "GET product " + id, "status": "ok"})
}