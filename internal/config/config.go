package config

import (
	"time"

	"github.com/sj1815/golang-country-search/internal/cache"
	"github.com/sj1815/golang-country-search/internal/client"
	"github.com/sj1815/golang-country-search/internal/handler"
	"github.com/sj1815/golang-country-search/internal/service"
)

type Config struct {
	ServerPort         string
	HTTPClientTimeout  time.Duration
	ServerReadTimeout  time.Duration
	ServerWriteTimeout time.Duration
	ShutdownTimeout    time.Duration
}

func DefaultConfig() *Config {
	return &Config{
		ServerPort:         ":8000",
		HTTPClientTimeout:  10 * time.Second,
		ServerReadTimeout:  15 * time.Second,
		ServerWriteTimeout: 15 * time.Second,
		ShutdownTimeout:    10 * time.Second,
	}
}

type Dependencies struct {
	CountryHandler *handler.CountryHandler
}

// InitDependencies initializes and returns the application dependencies based on the provided configuration.
func InitDependencies(cfg *Config) *Dependencies {
	countryCache := cache.NewInMemoryCache()
	httpClient := client.NewHTTPClient(cfg.HTTPClientTimeout)
	countryService := service.NewCountryService(httpClient, countryCache)
	countryHandler := handler.NewCountryHandler(countryService)

	return &Dependencies{
		CountryHandler: countryHandler,
	}
}
