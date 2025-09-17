## Engineering Specification: The "DocGen" Go Microservice

**Document ID:** ENG-DG-001-GO
**Version:** 1.0
**Target Audience:** Senior Software Engineer / Development Team

### 1.0 Project Mandate & Core Mission

This document outlines the implementation plan for the **DocGen Service**, a high-performance Go microservice.

**The mission of this service is singular: To efficiently and reliably render a valid Document Plan (JSON) into a compliant Microsoft Word (`.docx`) file.**

This service is the "rendering engine" of a larger document automation ecosystem. It is designed to be a stable, performant, and "unintelligent" backend. All complex business logic, planning, and content generation are explicitly handled by an upstream "Planner" service. DocGen's responsibility is solely to execute a well-formed plan.

### 2.0 Architectural Requirements & Constraints

*   **Language:** Go (latest stable version, currently 1.22+).
*   **Deployment Target:** GCP Cloud Run (containerized with Docker).
*   **Performance:** Must have low cold-start times and efficient memory usage.
*   **Dependencies:** Must have minimal external dependencies. Compilation should produce a single, static binary.
*   **API:** Must expose a RESTful HTTP API for all interactions.
*   **Validation:** Must be the single source of truth for Document Plan validation, using its native integration with the CUElang Go API.
*   **Concurrency:** The core engine must be thread-safe, allowing a single instantiated service to handle concurrent requests safely.
*   **Configuration:** All configurable paths (assets, schemas) should be manageable via environment variables for cloud-native operation.

### 3.0 Core API Endpoints

The service must expose the following three HTTP endpoints. All request and response bodies are `application/json` unless otherwise specified.

#### 3.1 `POST /generate`

*   **Description:** The primary endpoint. Accepts a Document Plan, validates it, assembles it, and returns the final `.docx` file.
*   **Request Body:** The full `plan.json` object.
*   **Workflow:**
    1.  Receive the request.
    2.  Perform a full, in-process CUE validation of the JSON body. If validation fails, immediately return `400 Bad Request` with a JSON body detailing the validation error.
    3.  If valid, proceed to the assembly process.
    4.  On successful assembly, package the document into a byte buffer.
    5.  Return `200 OK`.
*   **Success Response:**
    *   **Content-Type:** `application/vnd.openxmlformats-officedocument.wordprocessingml.document`
    *   **Content-Disposition:** `attachment; filename="<filename>"` (Filename should be extracted from a `doc_props.filename` field in the plan, with a sensible default).
    *   **Body:** The raw `[]byte` of the generated `.docx` file.
*   **Error Responses:**
    *   `400 Bad Request`: Plan validation failed or the JSON was malformed.
    *   `500 Internal Server Error`: An internal error occurred during assembly (e.g., a component was referenced but not found in the library).

#### 3.2 `POST /validate-plan`

*   **Description:** Exposes the CUE validation logic as a standalone endpoint for the Planner service's self-correction loop. This endpoint **does not** perform any document assembly.
*   **Request Body:** The full `plan.json` object.
*   **Workflow:**
    1.  Perform a full, in-process CUE validation of the JSON body.
    2.  If validation succeeds, return `200 OK` with a simple success message.
    3.  If validation fails, return `400 Bad Request` with a detailed JSON response containing the CUE validation error.
*   **Success Response (`200 OK`):**
    ```json
    { "status": "valid", "message": "Document plan is valid." }
    ```
*   **Error Response (`400 Bad Request`):**
    ```json
    {
      "status": "invalid",
      "message": "Document plan failed validation.",
      "errors": [
        {
          "path": "body[1].props.test_result",
          "message": "invalid value \"PENDING\" (out of bound \"PASS\"|\"FAIL\")"
        }
      ]
    }
    ```

#### 3.3 `GET /list-components`

*   **Description:** Provides discoverability for clients. It returns a JSON object describing the available components and their expected props.
*   **Implementation Strategy:** This schema should be derived directly from the CUE files. The service should parse the CUE schema definitions upon startup and generate a JSON representation of the component specifications. This ensures the API response is always in sync with the validation rules.
*   **Success Response (`200 OK`):** A JSON object detailing the component library schema.
    ```json
    {
      "components": {
        "DocumentTitle": {
          "description": "The main title block for the document.",
          "props": {
            "main_title": { "type": "string", "required": true },
            "subtitle": { "type": "string", "required": false },
            "doc_number": { "type": "string", "required": true }
          }
        },
        "TestDetails": { /* ... */ }
      }
    }
    ```

### 4.0 Recommended Go Package Structure

A clean, modular structure is recommended.

```
/docgen-service
├── /cmd/server/main.go        # Main application entry point, HTTP server setup
├── /internal
│   ├── /api/                  # HTTP handlers (generate, validate, list)
│   ├── /docgen/               # Core document generation engine package
│   │   ├── engine.go          # Main DocGenEngine struct and methods
│   │   ├── assembler.go       # Logic for assembling XML from components
│   │   ├── component.go       # Component library loading and management
│   │   ├── shell.go           # Logic for handling the shell .docx file
│   │   └── errors.go          # Custom error types
│   └── /validator/            # CUE validation logic
│       └── validator.go       # Wrapper for the CUE Go API
├── /assets/                   # Static assets, to be embedded or loaded
│   ├── /components/           # Component XML files (.component.xml)
│   ├── /schemas/              # CUE schema files (.cue)
│   └── template_shell.docx
├── go.mod
├── go.sum
└── Dockerfile
```

### 5.0 Core Technical Implementation Details

#### 5.1 Shell Document Handling
The service must treat the `template_shell.docx` as an immutable resource. On startup, the engine should unzip this file into an in-memory representation (a map of filenames to `[]byte`). The assembly process will then work on a *deep copy* of this in-memory structure for each request to ensure thread safety.

#### 5.2 XML Processing
*   The Go standard library's `encoding/xml` is sufficient for initial parsing. However, for the complex task of manipulating and appending XML node trees, a more robust third-party library is strongly recommended.
*   **Recommendation:** Use `github.com/beevik/etree` or a similar library that provides an ElementTree-like API for simple node manipulation. This will be critical for appending rendered component nodes into the shell's main `document.xml` body.

#### 5.3 Component Rendering
The prop injection is a simple text-templating problem.
*   Use the standard library's `strings.NewReplacer` for efficient, simultaneous replacement of all `{{prop_name}}` placeholders in a component's XML template. This is more performant than a series of single `strings.Replace` calls.
*   Ensure that prop values are properly XML-escaped to prevent injection issues (e.g., a prop value containing `<` or `&`).

#### 5.4 Dockerization
The `Dockerfile` should be a multi-stage build to produce a minimal final image.
*   **Stage 1 (Builder):** Use a `golang:alpine` image to build the static binary.
*   **Stage 2 (Final):** Start from a `scratch` or `distroless` image. Copy only the compiled binary and the required asset files (`assets/` directory) into the final image. This minimizes the attack surface and image size.

### 6.0 Phased Implementation Plan

This project can be broken down into three logical, testable milestones.

#### Milestone 1: The Core Validator (1-2 weeks)
1.  **Goal:** Implement the validation logic and the `/validate-plan` and `/list-components` endpoints.
2.  **Tasks:**
    *   Set up the project structure and basic HTTP server.
    *   Write the `validator` package to wrap the CUE Go API. It should load a schema from a file and have a function `Validate(plan map[string]interface{}) error`.
    *   Implement the `/validate-plan` handler.
    *   Implement the CUE-to-JSON schema logic for `/list-components`.
    *   Write comprehensive unit tests for the validator with valid and invalid plans.
3.  **Definition of Done:** The service is running and can successfully validate/reject plans and serve the component schema via its API.

#### Milestone 2: The Assembler Engine (2-3 weeks)
1.  **Goal:** Implement the core document assembly logic, without the HTTP wrapper.
2.  **Tasks:**
    *   Implement the `shell.go` module to load a `.docx` into an in-memory structure.
    *   Implement the `component.go` module to load the library of `.component.xml` files from a directory.
    *   Write the core `assembler.go` logic. This is the most complex part. It must:
        *   Take a copy of the in-memory shell.
        *   Take a component name and its props.
        *   Render the props into the component's XML template.
        *   Use an XML library (`etree`) to parse the rendered string and append the node to the shell's `document.xml` tree.
        *   Implement the final packaging logic to zip the modified in-memory structure back into a `[]byte`.
    *   Write extensive unit tests for the assembler, testing the rendering of individual components and a small sequence of components.
3.  **Definition of Done:** A set of unit tests can successfully call the assembler to produce a valid, multi-component `.docx` byte slice from a hardcoded plan.

#### Milestone 3: Full Integration & Deployment (1 week)
1.  **Goal:** Integrate the assembler with the HTTP server, build the Docker image, and prepare for deployment.
2.  **Tasks:**
    *   Implement the `/generate` handler, which wires together the validator and the assembler.
    *   Add robust logging and error handling to all API endpoints.
    *   Write the multi-stage `Dockerfile`.
    *   Create a simple `cloudbuild.yaml` or similar CI/CD pipeline for building and pushing the image to Google Artifact Registry.
    *   Perform end-to-end integration testing against all three API endpoints.
3.  **Definition of Done:** The service is fully containerized and can be deployed to Cloud Run, successfully serving requests and generating documents.

### 7.0 Definition of Success
The project will be considered successful when the DocGen service is deployed and reliably serves all three API endpoints according to this specification, enabling the upstream Planner service to successfully generate documents through it.