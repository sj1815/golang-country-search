// Package handler provides HTTP handlers for the API.
package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/sj1815/golang-country-search/internal/model"
	"github.com/sj1815/golang-country-search/internal/service"
)

// CountryHandler handles HTTP requests related to countries.
type CountryHandler struct {
	service service.CountryService
}

// NewCountryHandler creates a new instance of CountryHandler.
func NewCountryHandler(service service.CountryService) *CountryHandler {
	return &CountryHandler{
		service: service,
	}
}

// SearchCountry handles the search for a country by name.
func (h *CountryHandler) SearchCountry(w http.ResponseWriter, r *http.Request) {
	// Start timing the request
	// start := time.Now()

	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	countryName := r.URL.Query().Get("name")
	if countryName == "" {
		h.writeError(w, http.StatusBadRequest, "name query parameter is required")
		return
	}

	country, err := h.service.SearchCountry(r.Context(), countryName)
	if err != nil {
		log.Printf("Error searching country: %v", err)
		h.writeError(w, http.StatusNotFound, err.Error())
		return
	}

	// Log response time
	// elapsed := time.Since(start)
	// log.Printf("RESPONSE TIME: Request for '%s' took %v", countryName, elapsed)

	h.writeJSON(w, http.StatusOK, country)
}

// writeJSON writes the given data as a JSON response with the specified status code.
func (h *CountryHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
	}
}

// writeError writes an error response with the specified status code and message.
func (h *CountryHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, model.ErrorResponse{
		Error:   http.StatusText(status),
		Message: message,
	})
}
