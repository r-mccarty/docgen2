# DocGen2

A declarative, component-based document generation system that renders Microsoft Word documents from JSON plans. DocGen2 implements a "React for Docs" paradigm, treating document creation as a rendering process.

## Quick Start

### HTTP Server (Recommended)
```bash
# Run HTTP server
go run ./cmd/server -server

# Or build and run
go build -o docgen-server ./cmd/server
./docgen-server -server
```

The server will start on port 8080 (configurable via `PORT` environment variable).

### CLI Mode
```bash
# Build CLI
go build -o docgen-cli ./cmd/server

# Run CLI
./docgen-cli -shell assets/shell/template_shell.docx \
             -components assets/components/ \
             -plan assets/plans/test_plan_01.json \
             -output output/generated_document.docx
```

### Docker (Production)
```bash
# Build image
docker build -t docgen-service:latest .

# Run container
docker run -p 8080:8080 docgen-service:latest
```

## HTTP API

### Endpoints

#### `POST /generate`
Generate a Word document from a JSON document plan.

**Request:**
```bash
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{
    "doc_props": {
      "filename": "my_document.docx"
    },
    "body": [
      {
        "component": "DocumentCategoryTitle",
        "props": {
          "category_title": "TEST DOCUMENT"
        }
      }
    ]
  }' \
  --output generated_document.docx
```

**Response:**
- **Content-Type:** `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **Content-Disposition:** `attachment; filename="my_document.docx"`
- **Body:** Binary DOCX file data

#### `GET /health`
Health check endpoint.

```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "docgen-service",
  "components_loaded": 5,
  "available_components": ["DocumentCategoryTitle", "DocumentTitle", ...]
}
```

#### `GET /components`
List available components.

```bash
curl http://localhost:8080/components
```

**Response:**
```json
{
  "components": ["DocumentCategoryTitle", "DocumentTitle", ...],
  "count": 5,
  "note": "Detailed component specifications available in /docs/components/"
}
```

### Environment Variables

- `PORT` - Server port (default: 8080)
- `DOCGEN_SHELL_PATH` - Path to shell document (default: ./assets/shell/template_shell.docx)
- `DOCGEN_COMPONENTS_DIR` - Components directory (default: ./assets/components/)

## Development

### Running Tests
```bash
# Run unit tests
go test ./internal/docgen/...

# Run HTTP API tests
go test ./internal/api/...

# Run all tests with verbose output
go test -v ./...
```

### End-to-End Testing

The project includes comprehensive E2E testing through the CLI interface. To run the full test suite:

#### 1. Setup
```bash
# Ensure output directory exists
mkdir -p ./output
```

#### 2. Individual Component Smoke Tests
Run each component in isolation to verify basic functionality:

```bash
# DocumentCategoryTitle smoke test
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentCategoryTitle.json \
    -output ./output/cli_smoke_test_DocumentCategoryTitle.docx

# DocumentTitle smoke test
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentTitle.json \
    -output ./output/cli_smoke_test_DocumentTitle.docx

# DocumentSubject smoke test
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentSubject.json \
    -output ./output/cli_smoke_test_DocumentSubject.docx

# TestBlock smoke test
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_TestBlock.json \
    -output ./output/cli_smoke_test_TestBlock.docx

# AuthorBlock smoke test
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_AuthorBlock.json \
    -output ./output/cli_smoke_test_AuthorBlock.docx
```

#### 3. Full Integration Test
Test all components working together in a complete document:

```bash
go run ./cmd/server \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/full_integration_test.json \
    -output ./output/cli_full_integration_test.docx
```

#### 4. Verification
After running the tests:
1. Check that all commands complete without error
2. Verify that 6 `.docx` files are generated in `./output/`
3. Open each document to visually verify correct rendering
4. Compare outputs between unit tests (`/test_output/`) and CLI tests (`/output/`)

## Component Library

The system includes 5 production-ready components:

- **DocumentCategoryTitle**: Category header with decorative underline
- **DocumentTitle**: Main document title with metadata integration
- **DocumentSubject**: Document subject/revision line
- **TestBlock**: Test form with 5 input fields
- **AuthorBlock**: Author contact information block

See `/docs/components/` for detailed component documentation and usage examples.

## Architecture

DocGen2 follows a strict separation of concerns:
- **Document Plans (JSON)**: Declarative specifications describing document structure
- **Component Library**: Reusable, parameterizable OpenXML snippets
- **Shell Document**: Minimal `.docx` containing document-wide styles
- **Validation Layer**: CUE schemas for business rule enforcement (to be implemented)

For complete architecture details, see `docs/docgen-vision.md` and `docs/docgen-mvp-spec.md`.
