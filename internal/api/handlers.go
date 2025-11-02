package api

import (
	"encoding/json"
	"net/http"
	"rbac/internal/domain"
	"rbac/internal/service"
)

// APIHandler holds all services, acting as our dependency injection container
type APIHandler struct {
	authSvc    service.AuthService
	rbacSvc    service.RBACService
	productSvc service.ProductService
}

// NewAPIHandler creates a new APIHandler with all its dependencies
func NewAPIHandler(authSvc service.AuthService, rbacSvc service.RBACService, productSvc service.ProductService) *APIHandler {
	return &APIHandler{
		authSvc:    authSvc,
		rbacSvc:    rbacSvc,
		productSvc: productSvc,
	}
}

// --- Product Handlers ---

// CreateProductHandler handles product creation
func (h *APIHandler) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Get userID from context (set by AuthMiddleware)
	userID, ok := r.Context().Value(UserIDKey).(int64)
	if !ok {
		respondWithError(w, http.StatusInternalServerError, "User ID not found in context")
		return
	}

	product, err := h.productSvc.CreateProduct(r.Context(), req, userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

    // This is the line that returns the JSON body
	respondWithJSON(w, http.StatusCreated, product)
}

// --- Helper Functions ---

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}