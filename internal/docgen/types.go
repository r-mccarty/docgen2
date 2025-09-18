package docgen

import "docgen-service/internal/validator"

// DocumentPlan represents the top-level JSON structure for document generation
type DocumentPlan struct {
	DocProps DocProps            `json:"doc_props"`
	Body     []ComponentInstance `json:"body"`
}

// DocProps contains metadata about the document to be generated
type DocProps struct {
	Filename string `json:"filename"`
}

// ComponentInstance represents a single component to be rendered in the document
type ComponentInstance struct {
	Component string                 `json:"component"`
	Props     map[string]interface{} `json:"props"`
}

// InMemoryDocx represents a DOCX file loaded into memory as a map of file paths to content
type InMemoryDocx map[string][]byte

// Engine holds the loaded shell document and component library
type Engine struct {
	shell      InMemoryDocx
	components map[string]string
	validator  *validator.Validator
}