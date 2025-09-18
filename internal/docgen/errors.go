package docgen

import "fmt"

// DocGenError represents errors that occur during document generation
type DocGenError struct {
	Operation string
	Err       error
}

func (e *DocGenError) Error() string {
	return fmt.Sprintf("docgen %s: %v", e.Operation, e.Err)
}

func (e *DocGenError) Unwrap() error {
	return e.Err
}

// NewDocGenError creates a new DocGenError
func NewDocGenError(operation string, err error) *DocGenError {
	return &DocGenError{
		Operation: operation,
		Err:       err,
	}
}

// ComponentNotFoundError represents errors when a requested component is not found
type ComponentNotFoundError struct {
	ComponentName string
}

func (e *ComponentNotFoundError) Error() string {
	return fmt.Sprintf("component not found: %s", e.ComponentName)
}

// ShellLoadError represents errors when loading the shell document
type ShellLoadError struct {
	Path string
	Err  error
}

func (e *ShellLoadError) Error() string {
	return fmt.Sprintf("failed to load shell document from %s: %v", e.Path, e.Err)
}

func (e *ShellLoadError) Unwrap() error {
	return e.Err
}