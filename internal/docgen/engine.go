package docgen

import (
	"fmt"

	"docgen-service/internal/validator"
)

// NewEngine creates a new DocGen engine with the loaded shell and components
func NewEngine(shellPath, componentsDir, schemaPath string) (*Engine, error) {
	// Load the shell document
	shell, err := LoadShell(shellPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load shell: %w", err)
	}

	// Load all components
	components, err := LoadComponents(componentsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to load components: %w", err)
	}

	// Initialize the validator
	val, err := validator.New(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	return &Engine{
		shell:      shell,
		components: components,
		validator:  val,
	}, nil
}

// ValidatePlan validates a document plan using the CUE schema
func (e *Engine) ValidatePlan(plan map[string]interface{}) *validator.ValidationResult {
	return e.validator.Validate(plan)
}

// Assemble generates a DOCX document from the given plan
func (e *Engine) Assemble(plan DocumentPlan) ([]byte, error) {
	return e.AssembleDocument(plan)
}

// GetLoadedComponents returns the names of all loaded components
func (e *Engine) GetLoadedComponents() []string {
	var names []string
	for name := range e.components {
		names = append(names, name)
	}
	return names
}