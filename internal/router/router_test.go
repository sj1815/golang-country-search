package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saurabhj/golang-country-search/internal/handler"
	"github.com/saurabhj/golang-country-search/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCountryService is a mock implementation of service.CountryService
type MockCountryService struct {
	mock.Mock
}

// SearchCountry is a mock implementation of the SearchCountry method.
func (m *MockCountryService) SearchCountry(ctx context.Context, name string) (*model.Country, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Country), args.Error(1)
}

// TestNewRouter tests the NewRouter function.
func TestNewRouter(t *testing.T) {
	mockService := new(MockCountryService)
	countryHandler := handler.NewCountryHandler(mockService)

	router := NewRouter(countryHandler)

	assert.NotNil(t, router)
}

// TestRouter_CountrySearchRoute tests the /api/countries/search route.
func TestRouter_CountrySearchRoute(t *testing.T) {
	mockService := new(MockCountryService)
	countryHandler := handler.NewCountryHandler(mockService)

	expectedCountry := &model.Country{
		Name:       "Germany",
		Capital:    "Berlin",
		Currency:   "€",
		Population: 83240525,
	}

	mockService.On("SearchCountry", mock.Anything, "Germany").Return(expectedCountry, nil)

	router := NewRouter(countryHandler)

	// Create test server with the router
	server := httptest.NewServer(router)
	defer server.Close()

	// Make request to the route
	resp, err := http.Get(server.URL + "/api/countries/search?name=Germany")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

// TestRouter_UnknownRoute tests an unknown route.
func TestRouter_UnknownRoute(t *testing.T) {
	mockService := new(MockCountryService)
	countryHandler := handler.NewCountryHandler(mockService)

	router := NewRouter(countryHandler)

	server := httptest.NewServer(router)
	defer server.Close()

	resp, err := http.Get(server.URL + "/api/unknown")
	assert.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
