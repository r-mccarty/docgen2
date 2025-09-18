package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"docgen-service/internal/docgen"
)

// Server holds the HTTP server dependencies
type Server struct {
	engine *docgen.Engine
}

// NewServer creates a new API server with the DocGen engine
func NewServer(shellPath, componentsDir string) (*Server, error) {
	engine, err := docgen.NewEngine(shellPath, componentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create DocGen engine: %w", err)
	}

	return &Server{
		engine: engine,
	}, nil
}

// GenerateHandler handles POST /generate requests
func (s *Server) GenerateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Log request start
	log.Printf("POST /generate - Request started")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("POST /generate - Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse JSON plan
	var plan docgen.DocumentPlan
	if err := json.Unmarshal(body, &plan); err != nil {
		log.Printf("POST /generate - Failed to parse JSON: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Basic validation
	if len(plan.Body) == 0 {
		log.Printf("POST /generate - Empty document plan body")
		http.Error(w, "Document plan body cannot be empty", http.StatusBadRequest)
		return
	}

	// Generate document using existing engine
	result, err := s.engine.Assemble(plan)
	if err != nil {
		log.Printf("POST /generate - Document assembly failed: %v", err)
		http.Error(w, "Failed to generate document", http.StatusInternalServerError)
		return
	}

	// Determine filename
	filename := plan.DocProps.Filename
	if filename == "" {
		filename = "generated_document.docx"
	}
	if !strings.HasSuffix(strings.ToLower(filename), ".docx") {
		filename += ".docx"
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(result)))

	// Write document data
	if _, err := w.Write(result); err != nil {
		log.Printf("POST /generate - Failed to write response: %v", err)
		return
	}

	log.Printf("POST /generate - Document generated successfully: %s (%d bytes)", filename, len(result))
}

// HealthHandler handles GET /health requests
func (s *Server) HealthHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if engine is available and components are loaded
	components := s.engine.GetLoadedComponents()

	response := map[string]interface{}{
		"status": "healthy",
		"service": "docgen-service",
		"components_loaded": len(components),
		"available_components": components,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GET /health - Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// ComponentsHandler handles GET /components requests
func (s *Server) ComponentsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	components := s.engine.GetLoadedComponents()

	response := map[string]interface{}{
		"components": components,
		"count": len(components),
		"note": "Detailed component specifications available in /docs/components/",
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("GET /components - Failed to encode response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// SetupRoutes configures the HTTP routes for the server
func (s *Server) SetupRoutes() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/generate", s.GenerateHandler)
	mux.HandleFunc("/health", s.HealthHandler)
	mux.HandleFunc("/components", s.ComponentsHandler)

	return mux
}