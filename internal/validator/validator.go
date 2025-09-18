package validator

import (
	"fmt"
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/cuecontext"
	"cuelang.org/go/cue/errors"
	"cuelang.org/go/cue/load"
)

// ValidationError represents a structured validation error
type ValidationError struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

// ValidationResult represents the result of validating a document plan
type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

// Validator handles CUE schema validation for document plans
type Validator struct {
	ctx    *cue.Context
	schema cue.Value
}

// New creates a new validator instance by loading the CUE schema from the specified path
func New(schemaPath string) (*Validator, error) {
	ctx := cuecontext.New()

	// Load the CUE configuration
	buildInstances := load.Instances([]string{schemaPath}, nil)
	if len(buildInstances) == 0 {
		return nil, fmt.Errorf("no CUE instances found at path: %s", schemaPath)
	}

	// Check for load errors
	if err := buildInstances[0].Err; err != nil {
		return nil, fmt.Errorf("failed to load CUE schema: %w", err)
	}

	// Build the schema value
	schema := ctx.BuildInstance(buildInstances[0])
	if err := schema.Err(); err != nil {
		return nil, fmt.Errorf("failed to build CUE schema: %w", err)
	}

	return &Validator{
		ctx:    ctx,
		schema: schema,
	}, nil
}

// Validate validates a document plan against the CUE schema
func (v *Validator) Validate(plan map[string]interface{}) *ValidationResult {
	// Convert the plan to a CUE value
	planValue := v.ctx.Encode(plan)
	if err := planValue.Err(); err != nil {
		return &ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Path:    "root",
					Message: fmt.Sprintf("failed to encode plan: %v", err),
				},
			},
		}
	}

	// Unify the plan with the schema (this applies the constraints)
	unified := v.schema.LookupPath(cue.ParsePath("#DocumentPlan")).Unify(planValue)

	// Validate the unified value
	if err := unified.Validate(cue.Concrete(true)); err != nil {
		return &ValidationResult{
			Valid:  false,
			Errors: extractValidationErrors(err),
		}
	}

	return &ValidationResult{
		Valid:  true,
		Errors: nil,
	}
}

// extractValidationErrors converts CUE errors to structured validation errors
func extractValidationErrors(err error) []ValidationError {
	var validationErrors []ValidationError

	// Handle CUE error list
	if cueErr, ok := err.(errors.Error); ok {
		for _, e := range errors.Errors(cueErr) {
			validationErrors = append(validationErrors, ValidationError{
				Path:    extractPath(e),
				Message: extractMessage(e),
			})
		}
	} else {
		// Fallback for other error types
		validationErrors = append(validationErrors, ValidationError{
			Path:    "unknown",
			Message: err.Error(),
		})
	}

	return validationErrors
}

// extractPath extracts a human-readable path from a CUE error
func extractPath(err errors.Error) string {
	// Try to extract path from error message
	msg := err.Error()

	// Look for patterns like "body[0].props.document_title"
	if strings.Contains(msg, ".") {
		parts := strings.Fields(msg)
		for _, part := range parts {
			if strings.Contains(part, ".") && !strings.Contains(part, " ") {
				// Clean up the part
				cleaned := strings.Trim(part, "():,;")
				if len(cleaned) > 1 {
					return cleaned
				}
			}
		}
	}

	// Look for array index patterns like "body[0]"
	if strings.Contains(msg, "[") && strings.Contains(msg, "]") {
		parts := strings.Fields(msg)
		for _, part := range parts {
			if strings.Contains(part, "[") && strings.Contains(part, "]") {
				cleaned := strings.Trim(part, "():,;")
				if len(cleaned) > 1 {
					return cleaned
				}
			}
		}
	}

	return "root"
}

// extractMessage extracts a clean error message from a CUE error
func extractMessage(err errors.Error) string {
	msg := err.Error()

	// Clean up common CUE error patterns
	cleanMessage := strings.ReplaceAll(msg, "conflicting values", "validation failed:")
	cleanMessage = strings.ReplaceAll(cleanMessage, "incomplete value", "missing required value")

	// Remove CUE-specific formatting
	if idx := strings.Index(cleanMessage, ":"); idx > 0 && idx < 20 {
		cleanMessage = strings.TrimSpace(cleanMessage[idx+1:])
	}

	return cleanMessage
}