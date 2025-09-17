## Engineering Specification: DocGen Service (MVP)

**Document ID:** ENG-DG-001-GO-MVP
**Version:** 0.1 (MVP)
**Target Audience:** Senior Software Engineer / Development Team

### 1.0 Project Mandate & MVP Goal

**The goal of this MVP is to build a functional, single-endpoint document generation service that can be run locally.**

This service will prove the core architectural concept: rendering a document from a JSON plan using a component library. It will be built with cloud deployment in mind (containerization, environment-based config) but will not include any specific cloud integrations or non-essential API endpoints. The primary deliverable is a containerized application that accepts a JSON plan and returns a `.docx` file.

### 2.0 MVP Scope: "What's In" vs. "What's Out"

| Feature | IN for MVP | OUT for MVP (Post-MVP Enhancements) |
| :--- | :--- | :--- |
| **Core Logic** | ✅ **Document Assembly:** Renders a plan using a shell and components. | |
| **API Endpoints**| ✅ **Single `/generate` endpoint:** The core functionality. | ❌ **`/validate-plan` endpoint:** Validation will happen implicitly within the `/generate` call. The Planner service can infer failure from a 400 error. |
| | | ❌ **`/list-components` endpoint:** The component schema will be manually shared with the Planner team (e.g., via a static JSON file in a shared repo). |
| **Validation** | ✅ **Implicit CUE Validation:** The `/generate` endpoint MUST validate the plan before assembly. | ❌ **Detailed, structured error responses:** A simple `400 Bad Request` with the raw CUE error string in the body is sufficient. |
| **Configuration** | ✅ **File paths** (shell, components, schema) are loaded from environment variables. | ❌ Complex configuration providers or secret management. |
| **Error Handling**| ✅ Basic logging of errors to `stdout`. | ❌ Structured JSON logging, metrics, or tracing. |
| **Deployment** | ✅ **Dockerfile:** A simple, multi-stage `Dockerfile` is required. | ❌ **CI/CD Pipeline:** Manual `docker build` and `docker run` is acceptable for the MVP. No `cloudbuild.yaml` needed. |
| **Dependencies** | ✅ Minimal Go modules (`etree` for XML is recommended). | ❌ Anything beyond the essentials. |

### 3.0 MVP API Endpoint Specification

The service must expose only **one** endpoint.

#### `POST /generate`

*   **Description:** Accepts a Document Plan, validates it, assembles it, and returns the final `.docx` file. This is the sole entry point for the MVP.
*   **Request Body:** The full `plan.json` object.
*   **Workflow:**
    1.  Parse the JSON request body.
    2.  Validate the parsed JSON against the CUE schema loaded from disk.
        *   **If validation fails,** log the error and immediately return `41. **Goal:** Build the core document assembly engine. This is the heart of the MVP.
    2.  **Tasks:**
        *   Implement the `shell.go` module to load a `.docx` into an in-memory representation.
        *   Implement the `component.go` module to load the library of `.component.xml` files.
        *   Implement the `assembler.go` logic. This is the highest-priority task. It must:
            *   Take an in-memory shell.
            *   Take a plan and iterate through its components.
            *   Render props into each component template (`strings.NewReplacer`).
            *   Use an XML library (`etree`) to parse and append the rendered XML nodes to the shell's `document.xml`.
            *   Package the result back into a `.docx` byte slice.
    *   **Crucial:** Create a comprehensive **table-driven unit test** (`TestAssembler`) that feeds the assembler a simple but complete plan and asserts that the output is a non-empty, valid byte slice. This test is your primary development tool.
3.  **Definition of Done:** The `assembler` package can be called from a test to successfully produce a multi-component `.docx` file.

#### Phase 3: HTTP Server & Containerization (1 week)

1.  **Goal:** Wrap the engine in an HTTP server and create the Docker image.
2.  **Tasks:**
    *   Implement the `main.go` file and the `/generate` handler.
    *   Wire the handler to call the Validator first, then the Assembler.
    *   Add basic `log.Printf` statements for request start/end and errors.
    *   Write the multi-stage `Dockerfile`.
    *   Write a `README.md` with simple instructions: how to set environment variables, build the Docker image, and run the container, including a sample `curl` command to test the endpoint.
3.  **Definition of Done:** A developer can clone the repo, run `docker build` and `docker run`, and successfully generate a document by sending a JSON plan to the running container's endpoint.

### 5.0 Definition of MVP Success

The MVP is successful when a developer can run the application as a Docker container on their local machine and successfully generate a valid, multi-component `.docx` document by POSTing a valid JSON plan to the `/generate` endpoint. The service must also correctly reject an invalid plan with a `400 Bad Request`.