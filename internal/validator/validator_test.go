package validator

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidatorInitialization(t *testing.T) {
	schemaPath := "../../assets/schemas/rules.cue"

	validator, err := New(schemaPath)
	if err != nil {
		t.Fatalf("Failed to initialize validator: %v", err)
	}

	if validator == nil {
		t.Fatal("Validator is nil")
	}
}

func TestValidatorWithValidPlans(t *testing.T) {
	schemaPath := "../../assets/schemas/rules.cue"
	validator, err := New(schemaPath)
	if err != nil {
		t.Fatalf("Failed to initialize validator: %v", err)
	}

	validPlans := []string{
		"../../assets/plans/full_integration_test.json",
		"../../assets/plans/smoke_test_TestBlock.json",
		"../../assets/plans/smoke_test_DocumentTitle.json",
	}

	for _, planPath := range validPlans {
		t.Run("ValidPlan_"+planPath, func(t *testing.T) {
			// Read and parse the plan file
			planData, err := os.ReadFile(planPath)
			if err != nil {
				t.Skipf("Could not read plan file %s: %v", planPath, err)
				return
			}

			var plan map[string]interface{}
			if err := json.Unmarshal(planData, &plan); err != nil {
				t.Fatalf("Failed to parse plan JSON: %v", err)
			}

			// Validate the plan
			result := validator.Validate(plan)
			if !result.Valid {
				t.Errorf("Expected plan to be valid, but got errors: %+v", result.Errors)
			}

			if len(result.Errors) > 0 {
				t.Errorf("Expected no errors for valid plan, but got: %+v", result.Errors)
			}
		})
	}
}

func TestValidatorWithInvalidPlans(t *testing.T) {
	schemaPath := "../../assets/schemas/rules.cue"
	validator, err := New(schemaPath)
	if err != nil {
		t.Fatalf("Failed to initialize validator: %v", err)
	}

	invalidPlans := []struct {
		path        string
		expectError string
	}{
		{
			path:        "../../assets/plans/invalid/missing_required_prop.json",
			expectError: "category_title",
		},
		{
			path:        "../../assets/plans/invalid/bad_subject_format.json",
			expectError: "document_subject",
		},
		{
			path:        "../../assets/plans/invalid/bad_test_result_enum.json",
			expectError: "test_result",
		},
		{
			path:        "../../assets/plans/invalid/missing_title_component.json",
			expectError: "DocumentTitle",
		},
		{
			path:        "../../assets/plans/invalid/bad_test_date_format.json",
			expectError: "test_date",
		},
	}

	for _, test := range invalidPlans {
		t.Run("InvalidPlan_"+test.path, func(t *testing.T) {
			// Read and parse the plan file
			planData, err := os.ReadFile(test.path)
			if err != nil {
				t.Skipf("Could not read plan file %s: %v", test.path, err)
				return
			}

			var plan map[string]interface{}
			if err := json.Unmarshal(planData, &plan); err != nil {
				t.Fatalf("Failed to parse plan JSON: %v", err)
			}

			// Validate the plan
			result := validator.Validate(plan)
			if result.Valid {
				t.Errorf("Expected plan to be invalid, but validation passed")
			}

			if len(result.Errors) == 0 {
				t.Errorf("Expected validation errors, but got none")
			}

			// Check that the expected error is mentioned
			foundExpectedError := false
			for _, validationError := range result.Errors {
				if contains(validationError.Message, test.expectError) || contains(validationError.Path, test.expectError) {
					foundExpectedError = true
					break
				}
			}

			if !foundExpectedError {
				t.Errorf("Expected error containing '%s', but got errors: %+v", test.expectError, result.Errors)
			}

			t.Logf("Validation failed as expected with errors: %+v", result.Errors)
		})
	}
}

func TestValidatorSpecificValidationRules(t *testing.T) {
	schemaPath := "../../assets/schemas/rules.cue"
	validator, err := New(schemaPath)
	if err != nil {
		t.Fatalf("Failed to initialize validator: %v", err)
	}

	testCases := []struct {
		name    string
		plan    map[string]interface{}
		valid   bool
		errText string
	}{
		{
			name: "ValidDocumentSubjectFormat",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentTitle",
						"props": map[string]interface{}{
							"document_title": "Test Title",
						},
					},
					map[string]interface{}{
						"component": "DocumentSubject",
						"props": map[string]interface{}{
							"document_subject": "DOC-1234, Rev A",
						},
					},
				},
			},
			valid: true,
		},
		{
			name: "InvalidDocumentSubjectFormat",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentTitle",
						"props": map[string]interface{}{
							"document_title": "Test Title",
						},
					},
					map[string]interface{}{
						"component": "DocumentSubject",
						"props": map[string]interface{}{
							"document_subject": "DOC-1234, Rev. A", // Invalid: period after "Rev"
						},
					},
				},
			},
			valid:   false,
			errText: "document_subject",
		},
		{
			name: "ValidTestResult",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentTitle",
						"props": map[string]interface{}{
							"document_title": "Test Title",
						},
					},
					map[string]interface{}{
						"component": "TestBlock",
						"props": map[string]interface{}{
							"tester_name":     "John Doe",
							"test_date":       "9/18/2024",
							"serial_number":   "SN-001",
							"test_result":     "FAIL",
							"additional_info": "Some info",
						},
					},
				},
			},
			valid: true,
		},
		{
			name: "InvalidTestResult",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentTitle",
						"props": map[string]interface{}{
							"document_title": "Test Title",
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
			},
			valid:   false,
			errText: "test_result",
		},
		{
			name: "MissingDocumentTitle",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentSubject",
						"props": map[string]interface{}{
							"document_subject": "DOC-1234, Rev A",
						},
					},
				},
			},
			valid:   false,
			errText: "DocumentTitle",
		},
		{
			name: "EmptyRequiredField",
			plan: map[string]interface{}{
				"body": []interface{}{
					map[string]interface{}{
						"component": "DocumentTitle",
						"props": map[string]interface{}{
							"document_title": "", // Empty required field
						},
					},
				},
			},
			valid:   false,
			errText: "document_title",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validator.Validate(tc.plan)

			if result.Valid != tc.valid {
				t.Errorf("Expected valid=%v, got valid=%v. Errors: %+v", tc.valid, result.Valid, result.Errors)
			}

			if !tc.valid && tc.errText != "" {
				foundExpectedError := false
				for _, validationError := range result.Errors {
					if contains(validationError.Message, tc.errText) || contains(validationError.Path, tc.errText) {
						foundExpectedError = true
						break
					}
				}
				if !foundExpectedError {
					t.Errorf("Expected error containing '%s', but got errors: %+v", tc.errText, result.Errors)
				}
			}
		})
	}
}

func TestValidatorNonExistentSchema(t *testing.T) {
	_, err := New("/nonexistent/path/schema.cue")
	if err == nil {
		t.Error("Expected error for non-existent schema path, but got none")
	}
}

func TestValidatorInvalidJSON(t *testing.T) {
	schemaPath := "../../assets/schemas/rules.cue"
	validator, err := New(schemaPath)
	if err != nil {
		t.Fatalf("Failed to initialize validator: %v", err)
	}

	// Test with various invalid data types that should fail JSON parsing or validation
	invalidPlans := []map[string]interface{}{
		nil, // nil plan
		{},  // empty plan
		{
			"body": "not an array", // body should be array
		},
		{
			"body": []interface{}{
				"not an object", // body items should be objects
			},
		},
		{
			"body": []interface{}{
				map[string]interface{}{
					"component": 123, // component should be string
					"props":     map[string]interface{}{},
				},
			},
		},
	}

	for i, plan := range invalidPlans {
		t.Run("InvalidData_"+string(rune(i+'0')), func(t *testing.T) {
			result := validator.Validate(plan)
			if result.Valid {
				t.Errorf("Expected plan to be invalid, but validation passed for plan: %+v", plan)
			}
		})
	}
}

// Helper function to check if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
		 (s[:len(substr)] == substr ||
		  s[len(s)-len(substr):] == substr ||
		  findSubstring(s, substr))))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}