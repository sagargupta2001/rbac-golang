package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

// GetCountryHandler fetches country details from the external GraphQL API
func (h *APIHandler) GetCountryHandler(w http.ResponseWriter, r *http.Request) {
	// Get the country code from the URL, e.g., "US"
	vars := mux.Vars(r)
	code, ok := vars["code"]
	if !ok {
		respondWithError(w, http.StatusBadRequest, "Country code is required")
		return
	}

	countryDetails, err := h.graphqlSvc.GetCountryDetails(r.Context(), code)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, countryDetails)
}