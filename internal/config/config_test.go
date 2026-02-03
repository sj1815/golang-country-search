package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	assert.NotNil(t, cfg)
	assert.Equal(t, ":8000", cfg.ServerPort)
	assert.Equal(t, 10*time.Second, cfg.HTTPClientTimeout)
	assert.Equal(t, 15*time.Second, cfg.ServerReadTimeout)
	assert.Equal(t, 15*time.Second, cfg.ServerWriteTimeout)
	assert.Equal(t, 10*time.Second, cfg.ShutdownTimeout)
}

func TestInitDependencies(t *testing.T) {
	cfg := DefaultConfig()

	deps := InitDependencies(cfg)

	assert.NotNil(t, deps)
	assert.NotNil(t, deps.CountryHandler)
}

func TestInitDependencies_WithCustomConfig(t *testing.T) {
	cfg := &Config{
		ServerPort:         ":9000",
		HTTPClientTimeout:  5 * time.Second,
		ServerReadTimeout:  10 * time.Second,
		ServerWriteTimeout: 10 * time.Second,
		ShutdownTimeout:    5 * time.Second,
	}

	deps := InitDependencies(cfg)

	assert.NotNil(t, deps)
	assert.NotNil(t, deps.CountryHandler)
}
