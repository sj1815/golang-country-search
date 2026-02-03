package service

import (
	"context"
	"errors"
	"testing"

	"github.com/saurabhj/golang-country-search/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCache is a mock implementation of cache.Cache
type MockCache struct {
	mock.Mock
}

// Get is a mock implementation of the Get method.
func (m *MockCache) Get(key string) (interface{}, bool) {
	args := m.Called(key)
	return args.Get(0), args.Bool(1)
}

// Set is a mock implementation of the Set method.
func (m *MockCache) Set(key string, value interface{}) {
	m.Called(key, value)
}

// MockClient is a mock implementation of client.CountryClient
type MockClient struct {
	mock.Mock
}

// SearchCountryByName is a mock implementation of the SearchCountryByName method.
func (m *MockClient) SearchCountryByName(ctx context.Context, name string) ([]model.RESTCountryResponse, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.RESTCountryResponse), args.Error(1)
}

// TestNewCountryService tests the NewCountryService function.
func TestNewCountryService(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	service := NewCountryService(mockClient, mockCache)

	assert.NotNil(t, service)
}

// TestCountryService_SearchCountry_CacheHit tests the SearchCountry method when the country is found in the cache.
func TestCountryService_SearchCountry_CacheHit(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	cachedCountry := &model.Country{
		Name:       "Germany",
		Capital:    "Berlin",
		Currency:   "€",
		Population: 83240525,
	}

	// Cache returns the country
	mockCache.On("Get", "germany").Return(cachedCountry, true)

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "Germany")

	assert.NoError(t, err)
	assert.Equal(t, cachedCountry, country)
	mockCache.AssertExpectations(t)
	mockClient.AssertNotCalled(t, "SearchCountryByName")
}

// TestCountryService_SearchCountry_CacheMiss tests the SearchCountry method when the country is not found in the cache.
func TestCountryService_SearchCountry_CacheMiss(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	apiResponse := []model.RESTCountryResponse{
		{
			Name:       model.CountryName{Common: "Germany", Official: "Federal Republic of Germany"},
			Capital:    []string{"Berlin"},
			Currencies: map[string]model.CurrencyInfo{"EUR": {Name: "Euro", Symbol: "€"}},
			Population: 83240525,
		},
	}

	mockCache.On("Get", "germany").Return(nil, false)
	mockClient.On("SearchCountryByName", mock.Anything, "Germany").Return(apiResponse, nil)
	mockCache.On("Set", "germany", mock.AnythingOfType("*model.Country")).Return()

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "Germany")

	assert.NoError(t, err)
	assert.Equal(t, "Germany", country.Name)
	assert.Equal(t, "Berlin", country.Capital)
	assert.Equal(t, "€", country.Currency)
	assert.Equal(t, 83240525, country.Population)

	mockCache.AssertExpectations(t)
	mockClient.AssertExpectations(t)
}

// TestCountryService_SearchCountry_EmptyName tests the SearchCountry method when the country name is empty.
func TestCountryService_SearchCountry_EmptyName(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "")

	assert.Error(t, err)
	assert.Nil(t, country)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestCountryService_SearchCountry_WhitespaceName tests the SearchCountry method when the country name is whitespace.
func TestCountryService_SearchCountry_WhitespaceName(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "   ")

	assert.Error(t, err)
	assert.Nil(t, country)
	assert.Contains(t, err.Error(), "cannot be empty")
}

// TestCountryService_SearchCountry_ClientError tests the SearchCountry method when the client returns an error.
func TestCountryService_SearchCountry_ClientError(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	mockCache.On("Get", "invalidcountry").Return(nil, false)
	mockClient.On("SearchCountryByName", mock.Anything, "InvalidCountry").Return(nil, errors.New("country not found"))

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "InvalidCountry")

	assert.Error(t, err)
	assert.Nil(t, country)
	assert.Contains(t, err.Error(), "failed to search country")
}

// TestCountryService_SearchCountry_EmptyResponse tests the SearchCountry method when the client returns an empty response.
func TestCountryService_SearchCountry_EmptyResponse(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	mockCache.On("Get", "unknown").Return(nil, false)
	mockClient.On("SearchCountryByName", mock.Anything, "Unknown").Return([]model.RESTCountryResponse{}, nil)

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "Unknown")

	assert.Error(t, err)
	assert.Nil(t, country)
	assert.Contains(t, err.Error(), "no country data found")
}

// TestCountryService_SearchCountry_CaseInsensitiveCache tests the SearchCountry method with case-insensitive cache keys.
func TestCountryService_SearchCountry_CaseInsensitiveCache(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	cachedCountry := &model.Country{Name: "India"}

	mockCache.On("Get", "india").Return(cachedCountry, true)

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "INDIA")

	assert.NoError(t, err)
	assert.Equal(t, "India", country.Name)
}

// TestCountryService_SearchCountry_NoCurrency tests the SearchCountry method when the country has no currency information.
func TestCountryService_SearchCountry_NoCurrency(t *testing.T) {
	mockClient := new(MockClient)
	mockCache := new(MockCache)

	apiResponse := []model.RESTCountryResponse{
		{
			Name:       model.CountryName{Common: "Antarctica"},
			Capital:    []string{},
			Currencies: map[string]model.CurrencyInfo{},
			Population: 0,
		},
	}

	mockCache.On("Get", "antarctica").Return(nil, false)
	mockClient.On("SearchCountryByName", mock.Anything, "Antarctica").Return(apiResponse, nil)
	mockCache.On("Set", "antarctica", mock.AnythingOfType("*model.Country")).Return()

	service := NewCountryService(mockClient, mockCache)
	ctx := context.Background()

	country, err := service.SearchCountry(ctx, "Antarctica")

	assert.NoError(t, err)
	assert.Equal(t, "Antarctica", country.Name)
	assert.Equal(t, "", country.Capital)
	assert.Equal(t, "", country.Currency)
}

// TestTransformToCountry tests the transformToCountry function.
func TestTransformToCountry(t *testing.T) {
	tests := []struct {
		name     string
		input    model.RESTCountryResponse
		expected *model.Country
	}{
		{
			name: "full data",
			input: model.RESTCountryResponse{
				Name:       model.CountryName{Common: "India", Official: "Republic of India"},
				Capital:    []string{"New Delhi"},
				Currencies: map[string]model.CurrencyInfo{"INR": {Name: "Indian rupee", Symbol: "₹"}},
				Population: 1380004385,
			},
			expected: &model.Country{
				Name:       "India",
				Capital:    "New Delhi",
				Currency:   "₹",
				Population: 1380004385,
			},
		},
		{
			name: "no capital",
			input: model.RESTCountryResponse{
				Name:       model.CountryName{Common: "Test"},
				Capital:    []string{},
				Currencies: map[string]model.CurrencyInfo{"USD": {Symbol: "$"}},
				Population: 100,
			},
			expected: &model.Country{
				Name:       "Test",
				Capital:    "",
				Currency:   "$",
				Population: 100,
			},
		},
		{
			name: "currency name fallback",
			input: model.RESTCountryResponse{
				Name:       model.CountryName{Common: "Test"},
				Capital:    []string{"Capital"},
				Currencies: map[string]model.CurrencyInfo{"XYZ": {Name: "Test Dollar", Symbol: ""}},
				Population: 100,
			},
			expected: &model.Country{
				Name:       "Test",
				Capital:    "Capital",
				Currency:   "Test Dollar",
				Population: 100,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := transformToCountry(tt.input)
			assert.Equal(t, tt.expected.Name, result.Name)
			assert.Equal(t, tt.expected.Capital, result.Capital)
			assert.Equal(t, tt.expected.Population, result.Population)
			if tt.expected.Currency != "" {
				assert.NotEmpty(t, result.Currency)
			}
		})
	}
}
