package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/saurabhj/golang-country-search/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCountryService is a mock implementation of service.CountryService
type MockCountryService struct {
	mock.Mock
}

func (m *MockCountryService) SearchCountry(ctx context.Context, name string) (*model.Country, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Country), args.Error(1)
}

func TestNewCountryHandler(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	assert.NotNil(t, handler)
}

func TestCountryHandler_SearchCountry_Success(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	expectedCountry := &model.Country{
		Name:       "Germany",
		Capital:    "Berlin",
		Currency:   "€",
		Population: 83240525,
	}

	mockService.On("SearchCountry", mock.Anything, "Germany").Return(expectedCountry, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=Germany", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), "Germany")
	assert.Contains(t, rec.Body.String(), "Berlin")
	assert.Contains(t, rec.Body.String(), "€")
	mockService.AssertExpectations(t)
}

func TestCountryHandler_SearchCountry_MissingNameParam(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "name query parameter is required")
	mockService.AssertNotCalled(t, "SearchCountry")
}

func TestCountryHandler_SearchCountry_EmptyNameParam(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), "name query parameter is required")
}

func TestCountryHandler_SearchCountry_NotFound(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	mockService.On("SearchCountry", mock.Anything, "InvalidCountry").Return(nil, errors.New("country not found"))

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=InvalidCountry", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Contains(t, rec.Body.String(), "Not Found")
	mockService.AssertExpectations(t)
}

func TestCountryHandler_SearchCountry_MethodNotAllowed(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/api/countries/search?name=Germany", nil)
			rec := httptest.NewRecorder()

			handler.SearchCountry(rec, req)

			assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
			assert.Contains(t, rec.Body.String(), "method not allowed")
		})
	}
}

func TestCountryHandler_SearchCountry_ServiceError(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	mockService.On("SearchCountry", mock.Anything, "Germany").Return(nil, errors.New("internal error"))

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=Germany", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)
	mockService.AssertExpectations(t)
}

func TestCountryHandler_SearchCountry_ResponseFormat(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	expectedCountry := &model.Country{
		Name:       "India",
		Capital:    "New Delhi",
		Currency:   "₹",
		Population: 1380004385,
	}

	mockService.On("SearchCountry", mock.Anything, "India").Return(expectedCountry, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/countries/search?name=India", nil)
	rec := httptest.NewRecorder()

	handler.SearchCountry(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	// Check JSON structure
	body := rec.Body.String()
	assert.Contains(t, body, `"name":"India"`)
	assert.Contains(t, body, `"capital":"New Delhi"`)
	assert.Contains(t, body, `"currency":"₹"`)
	assert.Contains(t, body, `"population":1380004385`)
}

func TestCountryHandler_WriteJSON(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	rec := httptest.NewRecorder()
	data := map[string]string{"key": "value"}

	handler.writeJSON(rec, http.StatusOK, data)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), `"key":"value"`)
}

func TestCountryHandler_WriteError(t *testing.T) {
	mockService := new(MockCountryService)
	handler := NewCountryHandler(mockService)

	rec := httptest.NewRecorder()

	handler.writeError(rec, http.StatusBadRequest, "test error message")

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, "application/json", rec.Header().Get("Content-Type"))
	assert.Contains(t, rec.Body.String(), `"error":"Bad Request"`)
	assert.Contains(t, rec.Body.String(), `"message":"test error message"`)
}
