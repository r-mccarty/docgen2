# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DocGen2 is a declarative, component-based document generation system that renders Microsoft Word documents from JSON plans. It implements a "React for Docs" paradigm, treating document creation as a rendering process rather than template modification.

## Core Architecture

The system follows a strict separation of concerns:

- **Document Plans (JSON)**: Declarative specifications describing document structure and content
- **Component Library**: Reusable, parameterizable OpenXML snippets stored as `.component.xml` files
- **Shell Document**: Minimal `.docx` containing document-wide styles and formatting definitions
- **Validation Layer**: CUE schemas enforcing business rules and structural constraints

### Key Directories

- `/docs/`: Comprehensive specifications and workflows
  - `docgen-mvp-spec.md`: Go microservice implementation specification
  - `docgen-vision.md`: Overall architecture and design philosophy
  - `document-plan-spec.md`: JSON document plan specification with component examples
  - `api-endpoints.md`: Complete HTTP API documentation with examples and status
  - `example-component-extraction.md`: AI-assisted component authoring workflow
  - `asset-generation-procedure.md`: Manual component creation workflow
  - `/components/`: Component library documentation with usage guides
- `/assets/`: Document generation assets
  - `/components/`: Parameterized OpenXML component files (`.component.xml`)
  - `/schemas/`: CUE validation schemas (`rules.cue`)
  - `/plans/invalid/`: Test plans for validation failure scenarios
- `/cmd/server/`: HTTP server and CLI application entry point
- `/internal/api/`: HTTP handlers and REST API implementation
- `/internal/docgen/`: Core document generation engine
- `/internal/validator/`: CUE schema validation with structured error reporting

## Development Status

**Current Phase: Business Logic Validation Complete (Milestone 3)**

- ✅ **CLI Renderer**: Functional Go CLI for document generation (Milestone 1 Complete)
- ✅ **Component Library**: 5 production-ready components with comprehensive documentation
- ✅ **HTTP API**: Production-ready HTTP microservice with REST endpoints (Phase 2 Complete)
- ✅ **CUE Validation**: Business logic validation layer with structured error reporting (Milestone 3 Complete)
- ✅ **Docker Containerization**: Multi-stage Dockerfile with distroless base image
- ✅ **Comprehensive Testing**: Unit tests, HTTP integration tests, E2E validation, and validation scenarios

### Available Components

The component library now includes:
- **DocumentCategoryTitle**: Category header with decorative underline
- **DocumentTitle**: Main document title with metadata integration
- **DocumentSubject**: Document subject/revision line
- **TestBlock**: Test form with 5 input fields
- **AuthorBlock**: Author contact information block

See `/docs/components/` for complete specifications and usage examples.

## Component Authoring Workflow

Components are created by:
1. **Isolating**: Extract single visual elements from master templates into scaffold documents
2. **Extracting**: Unzip scaffold `.docx` files to access clean `document.xml`
3. **Parameterizing**: Replace hard-coded text with `{{prop_name}}` placeholders
4. **Saving**: Store as `.component.xml` files in `/assets/components/`

## HTTP Microservice Implementation

The Go microservice implements:
- ✅ **`POST /generate`**: Core endpoint with CUE validation and DOCX generation
- ✅ **`POST /validate-plan`**: Validation-only endpoint for Planner service integration
- ✅ **`GET /health`**: Health check endpoint with component status
- ✅ **`GET /components`**: Component discovery endpoint
- ✅ **CUE business logic validation**: Structured error reporting with field paths
- ✅ **Document assembly engine**: Combining shell + components with prop injection
- ✅ **Docker containerization**: Multi-stage Dockerfile with distroless base
- ✅ **Environment configuration**: PORT, DOCGEN_SHELL_PATH, DOCGEN_COMPONENTS_DIR, DOCGEN_SCHEMA_PATH
- ✅ **Comprehensive testing**: Unit tests, HTTP integration tests, CLI compatibility, validation scenarios
- ✅ **Dual-mode operation**: HTTP server and CLI modes in single binary

Components use `strings.NewReplacer` for prop substitution, and the `etree` library for XML manipulation.

## Key Files to Read First

1. `docs/docgen-mvp-spec.md` - Complete implementation specification
2. `docs/api-endpoints.md` - HTTP API documentation with endpoints and examples
3. `docs/docgen-vision.md` - Architecture philosophy and design rationale
4. `docs/example-component-extraction.md` - AI-assisted component creation workflow

## Component Creation Best Practices

### Lessons from Component Library Development

Based on creating and optimizing 5 production components, follow these practices:

**1. XML Simplification is Critical:**
- Remove ALL revision tracking metadata: `w14:paraId`, `w:rsidR`, `w:rsidRPr`, `w:rsidRDefault`, `w:rsidP`
- Strip `wp14:anchorId`, `wp14:editId` from drawing elements
- Remove `w14:anchorId`, `o:spid` from VML fallback elements
- Clean components are ~60% smaller and more robust

**2. Strategic Namespace Management:**
- Do NOT include `xmlns:` declarations on `<w:p>` elements (handled by shell document)
- Add namespaces only where first needed:
  - `xmlns:mc` on `<mc:AlternateContent>`
  - Drawing namespaces on `<wp:inline>`
  - VML namespaces on `<v:line>` fallback elements
- This prevents XML parser conflicts and reduces brittleness

**3. Component Structure Validation:**
- Verify `{{ prop_name }}` placement in `<w:t>` elements
- Preserve essential paragraph properties (`<w:pPr>`) for styling
- Maintain drawing complexity for visual elements (lines, shapes)
- Keep both modern (`<w:drawing>`) and legacy (`<w:pict>`) fallbacks

**4. Optimization Workflow:**
1. Extract raw XML from scaffold document
2. Remove `<w:document>` and `<w:body>` wrappers
3. Strip all revision IDs and redundant namespaces
4. Add minimal necessary namespaces to specific elements
5. Parameterize text content with `{{ props }}`
6. Verify structure integrity

**5. Standard Document Layout:**
For company documents, follow this vertical component order:
1. Header (to be created) - Document classification, logos
2. DocumentCategoryTitle - Document type identifier
3. DocumentTitle - Main document title
4. DocumentSubject - Document number/revision
5. TestBlock - Test execution details (for test documents)
6. AuthorBlock - Author and company contact information

**6. Component Documentation:**
- Each component must have comprehensive documentation in `/docs/components/`
- Include props specifications, usage examples, and styling notes
- Update document plan specification with new component examples
- Maintain cross-references between plan spec and component docs

## Development Notes

- ✅ **Go Project**: Complete with `go.mod`, `etree` and CUE dependencies, production-ready structure
- ✅ **Cloud Ready**: Designed for GCP Cloud Run deployment with Docker containerization
- ✅ **Production Testing**: Comprehensive test suite validates all functionality including validation scenarios
- ✅ **Semantic Styling**: Components use Word's built-in styles, not direct formatting
- ✅ **CUE Validation**: Business logic validation with structured error reporting (Milestone 3 Complete)

## Usage Examples

### HTTP Server Mode
```bash
# Start server
go run ./cmd/server -server

# Validate document plan
curl -X POST http://localhost:8080/validate-plan \
  -H "Content-Type: application/json" \
  -d @assets/plans/test_plan_01.json

# Generate document via API
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d @assets/plans/test_plan_01.json \
  --output generated.docx
```

### CLI Mode
```bash
# Generate document via CLI
go run ./cmd/server \
  -shell assets/shell/template_shell.docx \
  -components assets/components/ \
  -schema assets/schemas/rules.cue \
  -plan assets/plans/test_plan_01.json \
  -output generated.docx
```

### Docker Deployment
```bash
# Build and run container
docker build -t docgen-service .
docker run -p 8080:8080 docgen-service
```