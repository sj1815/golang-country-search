package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/sj1815/golang-country-search/internal/config"
	"github.com/sj1815/golang-country-search/internal/router"
)

// main is the entry point of the application.
func main() {
	// Load application configuration (server port, timeouts, etc.)
	cfg := config.DefaultConfig()

	// Initialize application dependencies (handlers, services, clients, caches)
	deps := config.InitDependencies(cfg)

	// Create HTTP router with all registered routes
	mux := router.NewRouter(deps.CountryHandler)

	server := &http.Server{
		Addr:         cfg.ServerPort,
		Handler:      mux,
		ReadTimeout:  cfg.ServerReadTimeout,
		WriteTimeout: cfg.ServerWriteTimeout,
	}

	// Buffered channel to receive server errors (capacity 1 to avoid blocking)
	serverErrors := make(chan error, 1)

	// Start server in a separate goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		serverErrors <- server.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block and wait for either server error or shutdown signal
	select {
	// Listen for server errors
	case err := <-serverErrors:
		log.Fatalf("Server error: %v", err)
	// Listen for shutdown signal
	case sig := <-shutdown:
		log.Printf("Shutdown signal received: %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Printf("Graceful shutdown failed: %v", err)

			if err := server.Close(); err != nil {
				log.Fatalf("Forced shutdown failed: %v", err)
			}
		}

		log.Println("Server stopped gracefully")
	}
}
