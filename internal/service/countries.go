package service

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/saurabhj/golang-country-search/internal/cache"
	"github.com/saurabhj/golang-country-search/internal/client"
	"github.com/saurabhj/golang-country-search/internal/model"
)

type CountryService interface {
	SearchCountry(ctx context.Context, name string) (*model.Country, error)
}

type countryService struct {
	client client.CountryClient
	cache  cache.Cache
}

// NewCountryService creates a new instance of CountryService.
func NewCountryService(client client.CountryClient, cache cache.Cache) CountryService {
	return &countryService{
		client: client,
		cache:  cache,
	}
}

// SearchCountry searches for a country by its name.
func (s *countryService) SearchCountry(ctx context.Context, name string) (*model.Country, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("SearchCountry: country name cannot be empty")
	}

	cacheKey := strings.ToLower(name)

	// Check cache first
	if cached, found := s.cache.Get(cacheKey); found {
		if country, ok := cached.(*model.Country); ok {
			log.Printf("CACHE HIT: Found country in cache: %s", cacheKey)
			return country, nil
		}
	}

	log.Printf("CACHE MISS: Country not in cache, calling API: %s", cacheKey)

	response, err := s.client.SearchCountryByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("SearchCountry: failed to search country by name: %w", err)
	}

	if len(response) == 0 {
		return nil, fmt.Errorf("SearchCountry: no country data found for name: %s", name)
	}

	country := transformToCountry(response[0])

	// Store in cache for future requests
	s.cache.Set(cacheKey, country)
	log.Printf("CACHE SET: Stored country in cache: %s", cacheKey)

	return country, nil
}

// transformToCountry converts a RESTCountryResponse to a Country model.
func transformToCountry(apiResp model.RESTCountryResponse) *model.Country {
	country := &model.Country{
		Name:       apiResp.Name.Common,
		Population: apiResp.Population,
	}

	if len(apiResp.Capital) > 0 {
		country.Capital = apiResp.Capital[0]
	}

	for _, currency := range apiResp.Currencies {
		if currency.Symbol != "" {
			country.Currency = currency.Symbol
			break
		}

		if currency.Name != "" {
			country.Currency = currency.Name
			break
		}
	}

	return country
}
