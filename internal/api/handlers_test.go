package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"docgen-service/internal/docgen"
)

// setupTestServer creates a test server for HTTP integration tests
func setupTestServer(t *testing.T) *Server {
	server, err := NewServer("../../assets/shell/template_shell.docx", "../../assets/components/")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	return server
}

func TestGenerateHandler_Success(t *testing.T) {
	server := setupTestServer(t)

	// Create test plan
	plan := docgen.DocumentPlan{
		DocProps: docgen.DocProps{
			Filename: "test_document.docx",
		},
		Body: []docgen.ComponentInstance{
			{
				Component: "DocumentCategoryTitle",
				Props: map[string]interface{}{
					"category_title": "HTTP API TEST",
				},
			},
		},
	}

	// Convert to JSON
	planJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal test plan: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.GenerateHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	expectedContentType := "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	if contentType != expectedContentType {
		t.Errorf("Expected content type %s, got %s", expectedContentType, contentType)
	}

	// Check content disposition
	contentDisposition := w.Header().Get("Content-Disposition")
	if !strings.Contains(contentDisposition, "test_document.docx") {
		t.Errorf("Expected filename in content disposition, got: %s", contentDisposition)
	}

	// Check that response body is not empty
	if w.Body.Len() == 0 {
		t.Error("Response body is empty")
	}

	t.Logf("Generated document size: %d bytes", w.Body.Len())
}

func TestGenerateHandler_InvalidJSON(t *testing.T) {
	server := setupTestServer(t)

	// Create invalid JSON
	invalidJSON := `{"invalid": json}`

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/generate", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.GenerateHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGenerateHandler_EmptyPlan(t *testing.T) {
	server := setupTestServer(t)

	// Create empty plan
	plan := docgen.DocumentPlan{
		DocProps: docgen.DocProps{
			Filename: "empty.docx",
		},
		Body: []docgen.ComponentInstance{}, // Empty body
	}

	// Convert to JSON
	planJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal test plan: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.GenerateHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestGenerateHandler_WrongMethod(t *testing.T) {
	server := setupTestServer(t)

	// Create request with wrong method
	req := httptest.NewRequest(http.MethodGet, "/generate", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.GenerateHandler(w, req)

	// Check response
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGenerateHandler_IntegrationTest(t *testing.T) {
	server := setupTestServer(t)

	// Create comprehensive test plan (similar to our OpticWorks test)
	plan := docgen.DocumentPlan{
		DocProps: docgen.DocProps{
			Filename: "api_integration_test.docx",
		},
		Body: []docgen.ComponentInstance{
			{
				Component: "DocumentCategoryTitle",
				Props: map[string]interface{}{
					"category_title": "HTTP API INTEGRATION TEST",
				},
			},
			{
				Component: "DocumentTitle",
				Props: map[string]interface{}{
					"document_title": "API Integration Test - HTTP Endpoint Validation",
				},
			},
			{
				Component: "DocumentSubject",
				Props: map[string]interface{}{
					"document_subject": "API-TEST-001, Rev 1.0",
				},
			},
			{
				Component: "TestBlock",
				Props: map[string]interface{}{
					"tester_name":     "HTTP Test Suite",
					"test_date":       "9/18/2024",
					"serial_number":   "API-INTEGRATION-001",
					"test_result":     "PASS",
					"additional_info": "All HTTP endpoints validated successfully",
				},
			},
			{
				Component: "AuthorBlock",
				Props: map[string]interface{}{
					"author_name":    "DocGen API Team",
					"company_name":   "DocGen Service",
					"address_line1":  "123 API Drive",
					"address_line2":  "Suite 200",
					"city_state_zip": "San Francisco, CA 94105",
					"phone":          "(555) 123-4567",
					"fax":            "(555) 123-4568",
					"website":        "https://api.docgen.com",
				},
			},
		},
	}

	// Convert to JSON
	planJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal test plan: %v", err)
	}

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/generate", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.GenerateHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Verify document was generated
	if w.Body.Len() == 0 {
		t.Fatal("Generated document is empty")
	}

	t.Logf("Integration test document size: %d bytes", w.Body.Len())
}

func TestHealthHandler(t *testing.T) {
	server := setupTestServer(t)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/health", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.HealthHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse health response: %v", err)
	}

	// Check response fields
	if response["status"] != "healthy" {
		t.Errorf("Expected status healthy, got %v", response["status"])
	}

	if response["service"] != "docgen-service" {
		t.Errorf("Expected service docgen-service, got %v", response["service"])
	}

	// Check that components are loaded
	componentsLoaded, ok := response["components_loaded"].(float64)
	if !ok || componentsLoaded <= 0 {
		t.Errorf("Expected components_loaded > 0, got %v", response["components_loaded"])
	}

	t.Logf("Health check passed: %d components loaded", int(componentsLoaded))
}

func TestComponentsHandler(t *testing.T) {
	server := setupTestServer(t)

	// Create request
	req := httptest.NewRequest(http.MethodGet, "/components", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ComponentsHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected content type application/json, got %s", contentType)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse components response: %v", err)
	}

	// Check that components are listed
	components, ok := response["components"].([]interface{})
	if !ok || len(components) == 0 {
		t.Errorf("Expected components list, got %v", response["components"])
	}

	// Check count matches
	count, ok := response["count"].(float64)
	if !ok || int(count) != len(components) {
		t.Errorf("Expected count to match components length, got count=%v, len=%d", count, len(components))
	}

	t.Logf("Components endpoint returned %d components", len(components))
}

func TestFullHTTPWorkflow(t *testing.T) {
	server := setupTestServer(t)
	mux := server.SetupRoutes()

	// Test health endpoint first
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Health check failed: %d", w.Code)
	}

	// Test components endpoint
	req = httptest.NewRequest(http.MethodGet, "/components", nil)
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Components endpoint failed: %d", w.Code)
	}

	// Test document generation
	plan := docgen.DocumentPlan{
		DocProps: docgen.DocProps{
			Filename: "workflow_test.docx",
		},
		Body: []docgen.ComponentInstance{
			{
				Component: "DocumentCategoryTitle",
				Props: map[string]interface{}{
					"category_title": "WORKFLOW TEST",
				},
			},
		},
	}

	planJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal plan: %v", err)
	}

	req = httptest.NewRequest(http.MethodPost, "/generate", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Document generation failed: %d, body: %s", w.Code, w.Body.String())
	}

	t.Log("Full HTTP workflow test completed successfully")
}