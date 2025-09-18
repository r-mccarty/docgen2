package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"docgen-service/internal/api"
)

// Config holds the server configuration
type Config struct {
	Port          string
	ShellPath     string
	ComponentsDir string
	SchemaPath    string
}

// LoadConfig loads configuration from environment variables with sensible defaults
func LoadConfig() *Config {
	config := &Config{
		Port:          getEnv("PORT", "8080"),
		ShellPath:     getEnv("DOCGEN_SHELL_PATH", "./assets/shell/template_shell.docx"),
		ComponentsDir: getEnv("DOCGEN_COMPONENTS_DIR", "./assets/components/"),
		SchemaPath:    getEnv("DOCGEN_SCHEMA_PATH", "./assets/schemas/rules.cue"),
	}

	// Validate paths exist
	if _, err := os.Stat(config.ShellPath); os.IsNotExist(err) {
		log.Fatalf("Shell document not found: %s", config.ShellPath)
	}

	if _, err := os.Stat(config.ComponentsDir); os.IsNotExist(err) {
		log.Fatalf("Components directory not found: %s", config.ComponentsDir)
	}

	if _, err := os.Stat(config.SchemaPath); os.IsNotExist(err) {
		log.Fatalf("Schema file not found: %s", config.SchemaPath)
	}

	return config
}

// getEnv returns environment variable value or default if not set
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func runServer() {
	log.Printf("Starting DocGen HTTP Server...")

	// Load configuration
	config := LoadConfig()
	log.Printf("Configuration loaded:")
	log.Printf("  Port: %s", config.Port)
	log.Printf("  Shell: %s", config.ShellPath)
	log.Printf("  Components: %s", config.ComponentsDir)
	log.Printf("  Schema: %s", config.SchemaPath)

	// Create API server
	server, err := api.NewServer(config.ShellPath, config.ComponentsDir, config.SchemaPath)
	if err != nil {
		log.Fatalf("Failed to create API server: %v", err)
	}

	// Setup routes
	mux := server.SetupRoutes()

	// Add basic logging middleware
	handler := loggingMiddleware(mux)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + config.Port,
		Handler: handler,
		// Timeouts
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second, // Longer write timeout for document generation
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on port %s", config.Port)
		log.Printf("Available endpoints:")
		log.Printf("  POST /generate      - Generate document from JSON plan")
		log.Printf("  POST /validate-plan - Validate document plan against schema")
		log.Printf("  GET  /health        - Health check")
		log.Printf("  GET  /components    - List available components")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Printf("Server shutting down...")

	// Create shutdown context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server shutdown complete")
}

// loggingMiddleware adds basic request logging
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Create a response writer wrapper to capture status code
		ww := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Process request
		next.ServeHTTP(ww, r)

		// Log request
		duration := time.Since(start)
		log.Printf("%s %s %d %v %s",
			r.Method,
			r.URL.Path,
			ww.statusCode,
			duration,
			r.RemoteAddr,
		)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}