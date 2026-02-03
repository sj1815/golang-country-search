package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewHTTPClient tests the NewHTTPClient function.
func TestNewHTTPClient(t *testing.T) {
	t.Run("with custom timeout", func(t *testing.T) {
		client := NewHTTPClient(5 * time.Second)

		assert.NotNil(t, client)
		assert.Equal(t, BaseURL, client.baseURL)
		assert.Equal(t, 5*time.Second, client.httpClient.Timeout)
	})

	t.Run("with zero timeout uses default", func(t *testing.T) {
		client := NewHTTPClient(0)

		assert.NotNil(t, client)
		assert.Equal(t, DefaultTimeout, client.httpClient.Timeout)
	})
}

// TestHTTPClient_SearchCountryByName_Success tests the SearchCountryByName method for a successful response.
func TestHTTPClient_SearchCountryByName_Success(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3.1/name/Germany?fullText=true", r.URL.String())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{
            "name": {"common": "Germany", "official": "Federal Republic of Germany"},
            "capital": ["Berlin"],
            "currencies": {"EUR": {"name": "Euro", "symbol": "€"}},
            "population": 83240525
        }]`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	countries, err := client.SearchCountryByName(ctx, "Germany")

	require.NoError(t, err)
	require.Len(t, countries, 1)
	assert.Equal(t, "Germany", countries[0].Name.Common)
	assert.Equal(t, "Berlin", countries[0].Capital[0])
	assert.Equal(t, "€", countries[0].Currencies["EUR"].Symbol)
	assert.Equal(t, 83240525, countries[0].Population)
}

// TestHTTPClient_SearchCountryByName_NotFound tests the SearchCountryByName method for a not found response.
func TestHTTPClient_SearchCountryByName_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	countries, err := client.SearchCountryByName(ctx, "InvalidCountry")

	assert.Error(t, err)
	assert.Nil(t, countries)
	assert.Contains(t, err.Error(), "country not found")
}

// TestHTTPClient_SearchCountryByName_ServerError tests the SearchCountryByName method for a server error response.
func TestHTTPClient_SearchCountryByName_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	countries, err := client.SearchCountryByName(ctx, "Germany")

	assert.Error(t, err)
	assert.Nil(t, countries)
	assert.Contains(t, err.Error(), "unexpected status code")
}

// TestHTTPClient_SearchCountryByName_InvalidJSON tests the SearchCountryByName method for an invalid JSON response.
func TestHTTPClient_SearchCountryByName_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	countries, err := client.SearchCountryByName(ctx, "Germany")

	assert.Error(t, err)
	assert.Nil(t, countries)
	assert.Contains(t, err.Error(), "failed to decode")
}

// TestHTTPClient_SearchCountryByName_ContextCanceled tests the SearchCountryByName method when the context is canceled.
func TestHTTPClient_SearchCountryByName_ContextCanceled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // Simulate slow response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	// Create context that's already canceled
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	countries, err := client.SearchCountryByName(ctx, "Germany")

	assert.Error(t, err)
	assert.Nil(t, countries)
}

// TestHTTPClient_SearchCountryByName_URLEncoding tests that the SearchCountryByName method properly URL-encodes country names.
func TestHTTPClient_SearchCountryByName_URLEncoding(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/v3.1/name/United%20States?fullText=true", r.URL.String())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"name": {"common": "United States"}, "capital": ["Washington, D.C."], "currencies": {"USD": {"symbol": "$"}}, "population": 331000000}]`))
	}))
	defer server.Close()

	client := &HTTPClient{
		baseURL:    server.URL + "/v3.1",
		httpClient: &http.Client{Timeout: 5 * time.Second},
	}

	ctx := context.Background()
	countries, err := client.SearchCountryByName(ctx, "United States")

	require.NoError(t, err)
	assert.Equal(t, "United States", countries[0].Name.Common)
}
