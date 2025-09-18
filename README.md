# DocGen2

A declarative, component-based document generation system that renders Microsoft Word documents from JSON plans. DocGen2 implements a "React for Docs" paradigm, treating document creation as a rendering process.

## Quick Start

### Building the CLI
```bash
go build -o docgen ./cmd/renderer/
```

### Basic Usage
```bash
./docgen -shell assets/shell/template_shell.docx \
         -components assets/components/ \
         -plan assets/plans/test_plan_01.json \
         -output output/generated_document.docx
```

## Development

### Running Tests
```bash
# Run unit tests
go test ./internal/docgen/...

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
go run ./cmd/renderer/main.go \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentCategoryTitle.json \
    -output ./output/cli_smoke_test_DocumentCategoryTitle.docx

# DocumentTitle smoke test
go run ./cmd/renderer/main.go \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentTitle.json \
    -output ./output/cli_smoke_test_DocumentTitle.docx

# DocumentSubject smoke test
go run ./cmd/renderer/main.go \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_DocumentSubject.json \
    -output ./output/cli_smoke_test_DocumentSubject.docx

# TestBlock smoke test
go run ./cmd/renderer/main.go \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_TestBlock.json \
    -output ./output/cli_smoke_test_TestBlock.docx

# AuthorBlock smoke test
go run ./cmd/renderer/main.go \
    -shell ./assets/shell/template_shell.docx \
    -components ./assets/components/ \
    -plan ./assets/plans/smoke_test_AuthorBlock.json \
    -output ./output/cli_smoke_test_AuthorBlock.docx
```

#### 3. Full Integration Test
Test all components working together in a complete document:

```bash
go run ./cmd/renderer/main.go \
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
