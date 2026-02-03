package router

import (
	"net/http"

	"github.com/sj1815/golang-country-search/internal/handler"
)

// NewRouter sets up the HTTP routes for the application.
func NewRouter(countryHandler *handler.CountryHandler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/countries/search", countryHandler.SearchCountry)
	// Additional routes can be added here

	return mux
}
