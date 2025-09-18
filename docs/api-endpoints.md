# DocGen2 API Endpoints

This document provides comprehensive documentation for the DocGen2 HTTP microservice API endpoints, including specifications, examples, and current development status.

## Service Overview

DocGen2 is a declarative document generation microservice that renders Microsoft Word documents from JSON plans. The service implements a "React for Docs" paradigm with reusable, parameterizable components.

**Base URL**: `http://localhost:8080` (default)
**Content Type**: `application/json` (for requests), `application/vnd.openxmlformats-officedocument.wordprocessingml.document` (for responses)

## Configuration

The service supports the following environment variables:

- `PORT`: Server port (default: `8080`)
- `DOCGEN_SHELL_PATH`: Path to shell document (default: `./assets/shell/template_shell.docx`)
- `DOCGEN_COMPONENTS_DIR`: Components directory (default: `./assets/components/`)

## API Endpoints

### 1. POST /generate

**Status**: ✅ **Production Ready** (Phase 2 Complete)

Generates a Microsoft Word document from a JSON document plan.

#### Request

- **Method**: `POST`
- **URL**: `/generate`
- **Content-Type**: `application/json`
- **Body**: JSON document plan following the [document plan specification](/docs/document-plan-spec.md)

#### Request Body Schema

```json
{
  "doc_props": {
    "filename": "string (optional, defaults to 'generated_document.docx')"
  },
  "body": [
    {
      "component": "string (component name)",
      "props": {
        "prop_name": "string (component-specific properties)"
      }
    }
  ]
}
```

#### Response

- **Success Status**: `200 OK`
- **Content-Type**: `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
- **Headers**:
  - `Content-Disposition`: `attachment; filename="[filename].docx"`
  - `Content-Length`: Document size in bytes
- **Body**: Binary DOCX file data

#### Error Responses

| Status Code | Description | Response Body |
|-------------|-------------|---------------|
| `400 Bad Request` | Invalid JSON format | `"Invalid JSON format"` |
| `400 Bad Request` | Empty document plan body | `"Document plan body cannot be empty"` |
| `405 Method Not Allowed` | Non-POST request | `"Method not allowed"` |
| `500 Internal Server Error` | Document generation failed | `"Failed to generate document"` |

#### Example Request

```bash
curl -X POST http://localhost:8080/generate \
  -H "Content-Type: application/json" \
  -d '{
    "doc_props": {
      "filename": "MyFirstRenderedDocument.docx"
    },
    "body": [
      {
        "component": "DocumentCategoryTitle",
        "props": {
          "category_title": "This Title Was Rendered by DocGen!"
        }
      }
    ]
  }' \
  --output generated.docx
```

#### Example Response Headers

```
HTTP/1.1 200 OK
Content-Type: application/vnd.openxmlformats-officedocument.wordprocessingml.document
Content-Disposition: attachment; filename="MyFirstRenderedDocument.docx"
Content-Length: 15432
```

---

### 2. GET /health

**Status**: ✅ **Production Ready** (Phase 2 Complete)

Returns service health status and component information.

#### Request

- **Method**: `GET`
- **URL**: `/health`

#### Response

- **Success Status**: `200 OK`
- **Content-Type**: `application/json`

#### Response Schema

```json
{
  "status": "healthy",
  "service": "docgen-service",
  "components_loaded": "number",
  "available_components": ["array of component names"]
}
```

#### Error Responses

| Status Code | Description | Response Body |
|-------------|-------------|---------------|
| `405 Method Not Allowed` | Non-GET request | `"Method not allowed"` |
| `500 Internal Server Error` | Response encoding failed | `"Failed to encode response"` |

#### Example Request

```bash
curl http://localhost:8080/health
```

#### Example Response

```json
{
  "status": "healthy",
  "service": "docgen-service",
  "components_loaded": 5,
  "available_components": [
    "AuthorBlock",
    "DocumentCategoryTitle",
    "DocumentSubject",
    "DocumentTitle",
    "TestBlock"
  ]
}
```

---

### 3. GET /components

**Status**: ✅ **Production Ready** (Phase 2 Complete)

Returns detailed information about available components.

#### Request

- **Method**: `GET`
- **URL**: `/components`

#### Response

- **Success Status**: `200 OK`
- **Content-Type**: `application/json`

#### Response Schema

```json
{
  "components": ["array of component names"],
  "count": "number",
  "note": "string (reference to detailed documentation)"
}
```

#### Error Responses

| Status Code | Description | Response Body |
|-------------|-------------|---------------|
| `405 Method Not Allowed` | Non-GET request | `"Method not allowed"` |
| `500 Internal Server Error` | Response encoding failed | `"Failed to encode response"` |

#### Example Request

```bash
curl http://localhost:8080/components
```

#### Example Response

```json
{
  "components": [
    "AuthorBlock",
    "DocumentCategoryTitle",
    "DocumentSubject",
    "DocumentTitle",
    "TestBlock"
  ],
  "count": 5,
  "note": "Detailed component specifications available in /docs/components/"
}
```

## Development Status

### Phase 2 Complete ✅

The HTTP microservice implementation is **production-ready** with the following features:

- **Core Endpoints**: All 3 endpoints (`/generate`, `/health`, `/components`) fully implemented
- **Document Generation**: Complete integration with DocGen engine for DOCX assembly
- **Error Handling**: Comprehensive error handling with appropriate HTTP status codes
- **Logging**: Request/response logging with performance metrics
- **Configuration**: Environment-based configuration for deployment flexibility
- **Testing**: Comprehensive test suite including unit tests and HTTP integration tests
- **Docker Ready**: Multi-stage Dockerfile with distroless base image
- **Production Features**:
  - Graceful shutdown handling
  - Request timeouts (15s read, 30s write, 60s idle)
  - Proper MIME types and headers
  - Content-Length headers for efficient streaming

### Available Components ✅

The component library includes 5 production-ready components:

| Component | Description | Documentation |
|-----------|-------------|---------------|
| `DocumentCategoryTitle` | Category header with decorative underline | [docs/components/DocumentCategoryTitle.md](/docs/components/DocumentCategoryTitle.md) |
| `DocumentTitle` | Main document title with metadata integration | [docs/components/DocumentTitle.md](/docs/components/DocumentTitle.md) |
| `DocumentSubject` | Document subject/revision line | [docs/components/DocumentSubject.md](/docs/components/DocumentSubject.md) |
| `TestBlock` | Test form with 5 input fields | [docs/components/TestBlock.md](/docs/components/TestBlock.md) |
| `AuthorBlock` | Author contact information block | [docs/components/AuthorBlock.md](/docs/components/AuthorBlock.md) |

### Dual-Mode Operation ✅

The service supports both HTTP server and CLI modes:

```bash
# HTTP Server Mode
go run ./cmd/server -server

# CLI Mode
go run ./cmd/server \
  -shell assets/shell/template_shell.docx \
  -components assets/components/ \
  -plan assets/plans/test_plan_01.json \
  -output generated.docx
```

### Future Enhancements (Post-MVP)

- **CUE Validation**: Schema validation layer for business rules enforcement
- **Additional Components**: Expanding the component library based on requirements
- **Metrics/Observability**: Prometheus metrics and distributed tracing
- **Rate Limiting**: Request throttling for production deployment

## Testing

The API endpoints are thoroughly tested with:

- **Unit Tests**: Individual handler testing with mocked dependencies
- **Integration Tests**: Full HTTP request/response cycle testing
- **End-to-End Tests**: Complete document generation workflow validation

Run tests with:

```bash
go test ./...
```

## Deployment

### Local Development

```bash
# Start server
go run ./cmd/server -server

# Test endpoints
curl http://localhost:8080/health
```

### Docker Deployment

```bash
# Build container
docker build -t docgen-service .

# Run container
docker run -p 8080:8080 docgen-service
```

### Cloud Deployment

The service is designed for deployment on cloud platforms like GCP Cloud Run with:

- Stateless design
- Environment-based configuration
- Graceful shutdown handling
- Proper HTTP status codes and headers

## Related Documentation

- [DocGen MVP Specification](/docs/docgen-mvp-spec.md) - Complete implementation specification
- [Document Plan Specification](/docs/document-plan-spec.md) - JSON plan format and examples
- [Component Library Documentation](/docs/components/) - Individual component specifications
- [DocGen Architecture Vision](/docs/docgen-vision.md) - Overall design philosophy