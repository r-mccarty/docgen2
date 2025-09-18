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
	server, err := NewServer("../../assets/shell/template_shell.docx", "../../assets/components/", "../../assets/schemas/rules.cue")
	if err != nil {
		t.Fatalf("Failed to create test server: %v", err)
	}
	return server
}

func TestGenerateHandler_Success(t *testing.T) {
	server := setupTestServer(t)

	// Create test plan (must include DocumentTitle per schema)
	plan := docgen.DocumentPlan{
		DocProps: docgen.DocProps{
			Filename: "test_document.docx",
		},
		Body: []docgen.ComponentInstance{
			{
				Component: "DocumentTitle",
				Props: map[string]interface{}{
					"document_title": "API Test Document",
				},
			},
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
					"document_subject": "DOC-1001, Rev A",
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
				Component: "DocumentTitle",
				Props: map[string]interface{}{
					"document_title": "Workflow Test Document",
				},
			},
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

// Validation endpoint tests
func TestValidatePlanHandler_ValidPlan(t *testing.T) {
	server := setupTestServer(t)

	// Create valid plan
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "valid_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Valid Test Document",
				},
			},
			map[string]interface{}{
				"component": "DocumentSubject",
				"props": map[string]interface{}{
					"document_subject": "DOC-1234, Rev A",
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
	req := httptest.NewRequest(http.MethodPost, "/validate-plan", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	// Check response fields
	if response["status"] != "valid" {
		t.Errorf("Expected status valid, got %v", response["status"])
	}

	if response["valid"] != true {
		t.Errorf("Expected valid=true, got %v", response["valid"])
	}

	t.Log("Valid plan validation test passed")
}

func TestValidatePlanHandler_InvalidPlan_MissingTitle(t *testing.T) {
	server := setupTestServer(t)

	// Create invalid plan (missing DocumentTitle component)
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "invalid_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentCategoryTitle",
				"props": map[string]interface{}{
					"category_title": "TEST CATEGORY",
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
	req := httptest.NewRequest(http.MethodPost, "/validate-plan", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	// Check response fields
	if response["status"] != "invalid" {
		t.Errorf("Expected status invalid, got %v", response["status"])
	}

	if response["valid"] != false {
		t.Errorf("Expected valid=false, got %v", response["valid"])
	}

	// Check that errors are present
	errors, ok := response["errors"].([]interface{})
	if !ok || len(errors) == 0 {
		t.Errorf("Expected validation errors, got %v", response["errors"])
	}

	t.Logf("Invalid plan validation test passed with %d errors", len(errors))
}

func TestValidatePlanHandler_InvalidPlan_BadSubjectFormat(t *testing.T) {
	server := setupTestServer(t)

	// Create invalid plan (bad DocumentSubject format)
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "invalid_subject_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Valid Title",
				},
			},
			map[string]interface{}{
				"component": "DocumentSubject",
				"props": map[string]interface{}{
					"document_subject": "DOC-1234, Rev. A", // Invalid: period after Rev
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
	req := httptest.NewRequest(http.MethodPost, "/validate-plan", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	// Check response fields
	if response["status"] != "invalid" {
		t.Errorf("Expected status invalid, got %v", response["status"])
	}

	// Check that errors are present
	errors, ok := response["errors"].([]interface{})
	if !ok || len(errors) == 0 {
		t.Errorf("Expected validation errors, got %v", response["errors"])
	}

	t.Logf("Bad subject format validation test passed with %d errors", len(errors))
}

func TestValidatePlanHandler_InvalidJSON(t *testing.T) {
	server := setupTestServer(t)

	// Create invalid JSON
	invalidJSON := `{"invalid": json}`

	// Create request
	req := httptest.NewRequest(http.MethodPost, "/validate-plan", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestValidatePlanHandler_WrongMethod(t *testing.T) {
	server := setupTestServer(t)

	// Create request with wrong method
	req := httptest.NewRequest(http.MethodGet, "/validate-plan", nil)

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status %d, got %d", http.StatusMethodNotAllowed, w.Code)
	}
}

func TestGenerateHandler_ValidationIntegration(t *testing.T) {
	server := setupTestServer(t)

	// Test that invalid plans are rejected by /generate endpoint
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "validation_integration_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Valid Title",
				},
			},
			map[string]interface{}{
				"component": "TestBlock",
				"props": map[string]interface{}{
					"tester_name":     "John Doe",
					"test_date":       "9/18/2024",
					"serial_number":   "SN-001",
					"test_result":     "MAYBE", // Invalid enum value
					"additional_info": "Some info",
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

	// Check response - should be validation error
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}

	// Parse response to verify it's a validation error
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	if response["status"] != "invalid" {
		t.Errorf("Expected validation error response, got %v", response)
	}

	t.Log("Generate endpoint validation integration test passed")
}

func TestValidationEndpointInRoutes(t *testing.T) {
	server := setupTestServer(t)
	mux := server.SetupRoutes()

	// Test that /validate-plan endpoint is properly routed
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "route_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Route Test Document",
				},
			},
		},
	}

	planJSON, err := json.Marshal(plan)
	if err != nil {
		t.Fatalf("Failed to marshal plan: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/validate-plan", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusOK, w.Code, w.Body.String())
	}

	t.Log("Validation endpoint routing test passed")
}

func TestValidatePlanHandler_MultipleDocumentTitles(t *testing.T) {
	server := setupTestServer(t)

	// Create plan with multiple DocumentTitle components (violates compositional rule)
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "multiple_titles_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentCategoryTitle",
				"props": map[string]interface{}{
					"category_title": "TEST CATEGORY",
				},
			},
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "First Title",
				},
			},
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Second Title",
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
	req := httptest.NewRequest(http.MethodPost, "/validate-plan", bytes.NewBuffer(planJSON))
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Call handler
	server.ValidatePlanHandler(w, req)

	// Check response
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}

	// Parse response
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	// Check response fields
	if response["status"] != "invalid" {
		t.Errorf("Expected status invalid, got %v", response["status"])
	}

	if response["valid"] != false {
		t.Errorf("Expected valid=false, got %v", response["valid"])
	}

	// Check that errors are present
	errors, ok := response["errors"].([]interface{})
	if !ok || len(errors) == 0 {
		t.Errorf("Expected validation errors, got %v", response["errors"])
	}

	t.Logf("Multiple DocumentTitle validation test passed with %d errors", len(errors))
}

func TestGenerateHandler_MultipleDocumentTitles(t *testing.T) {
	server := setupTestServer(t)

	// Test that plans with multiple DocumentTitle components are rejected by /generate endpoint
	plan := map[string]interface{}{
		"doc_props": map[string]interface{}{
			"filename": "multiple_titles_generate_test.docx",
		},
		"body": []interface{}{
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "First Title",
				},
			},
			map[string]interface{}{
				"component": "DocumentTitle",
				"props": map[string]interface{}{
					"document_title": "Second Title",
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

	// Check response - should be validation error
	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d. Body: %s", http.StatusBadRequest, w.Code, w.Body.String())
	}

	// Parse response to verify it's a validation error
	var response map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse validation response: %v", err)
	}

	if response["status"] != "invalid" {
		t.Errorf("Expected validation error response, got %v", response)
	}

	t.Log("Generate endpoint multiple DocumentTitle validation test passed")
}