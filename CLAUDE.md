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
  - `example-component-extraction.md`: AI-assisted component authoring workflow
  - `asset-generation-procedure.md`: Manual component creation workflow
- `/assets/`: Document generation assets
  - `/components/`: Parameterized OpenXML component files (`.component.xml`)

## Development Status

This is currently a **specification and design phase project**. The actual Go implementation has not yet been started. All Go code, CUE schemas, and the HTTP API are still to be implemented according to the MVP specification.

## Component Authoring Workflow

Components are created by:
1. **Isolating**: Extract single visual elements from master templates into scaffold documents
2. **Extracting**: Unzip scaffold `.docx` files to access clean `document.xml`
3. **Parameterizing**: Replace hard-coded text with `{{prop_name}}` placeholders
4. **Saving**: Store as `.component.xml` files in `/assets/components/`

## MVP Implementation Plan

The Go microservice should implement:
- Single `POST /generate` endpoint accepting JSON document plans
- CUE validation of incoming plans
- Document assembly engine combining shell + components
- Docker containerization
- Basic error handling and logging

Components use `strings.NewReplacer` for prop substitution, and the `etree` library is recommended for XML manipulation.

## Key Files to Read First

1. `docs/docgen-mvp-spec.md` - Complete implementation specification
2. `docs/docgen-vision.md` - Architecture philosophy and design rationale
3. `docs/example-component-extraction.md` - AI-assisted component creation workflow

## Development Notes

- No package.json or dependencies exist yet - this is a Go project to be created
- The system is designed for cloud deployment (GCP Cloud Run) but starts with local Docker
- Components must be parameterized with semantic styles, not direct formatting
- All business logic validation happens in CUE schemas, not Go code