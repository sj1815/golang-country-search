package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/saurabhj/golang-country-search/internal/model"
)

const (
	DefaultTimeout = 10 * time.Second
	BaseURL        = "https://restcountries.com/v3.1"
)

type CountryClient interface {
	SearchCountryByName(ctx context.Context, name string) ([]model.RESTCountryResponse, error)
}

type HTTPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewHTTPClient creates a new instance of HTTPClient with the specified timeout.
func NewHTTPClient(timeout time.Duration) *HTTPClient {
	if timeout == 0 {
		timeout = DefaultTimeout
	}

	return &HTTPClient{
		baseURL: BaseURL,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SearchCountryByName searches for a country by its full name using the REST Countries API.
func (c *HTTPClient) SearchCountryByName(ctx context.Context, name string) ([]model.RESTCountryResponse, error) {
	endpoint := fmt.Sprintf("%s/name/%s?fullText=true", c.baseURL, url.PathEscape(name))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("SearchCountryByName: failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("SearchCountryByName: request execution failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("SearchCountryByName: country not found: %s", name)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("SearchCountryByName: unexpected status code: %d", resp.StatusCode)
	}

	var countries []model.RESTCountryResponse
	if err := json.NewDecoder(resp.Body).Decode(&countries); err != nil {
		return nil, fmt.Errorf("SearchCountryByName: failed to decode response: %w", err)
	}

	return countries, nil
}
