package api

import (
	"encoding/json"
	"net/http"
	"rbac/internal/domain"
)

func (h *APIHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Basic validation
	if req.Username == "" || req.Password == "" || req.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	user, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

func (h *APIHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req domain.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	token, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, domain.LoginResponse{Token: token})
}