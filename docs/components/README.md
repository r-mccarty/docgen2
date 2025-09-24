# DocGen Component Library

This directory contains documentation for all available components in the DocGen document generation system.

## Component Overview

Components are reusable, parameterizable OpenXML snippets that render specific visual elements in Word documents. Each component:

- Is defined as a `.component.xml` file in `/assets/components/`
- Accepts specific props via `{{ prop_name }}` placeholders
- Maintains semantic styling through Word's built-in styles
- Can be composed together in document plans to create complete documents

## Available Components

### Layout & Structure
- [DocumentCategoryTitle](./DocumentCategoryTitle.md) - Category header with decorative underline
- [DocumentTitle](./DocumentTitle.md) - Main document title with structured data display
- [DocumentSubject](./DocumentSubject.md) - Document subject/revision line

### Section & Content Components
- HeadingOneSection - Level 1 heading with content paragraph
- HeadingTwoSection - Level 2 heading with content paragraph
- [TestBlock](./TestBlock.md) - Test form with tester, date, serial number, result, and additional info fields
- [AuthorBlock](./AuthorBlock.md) - Author contact information block with company details

### Document Metadata
- RevisionBlock - Document revision history table with revision letter, date, changes, and authority

### Tables & Data
- ReferenceTable - Reference documents table supporting up to 4 document entries
- AcronymList - Symbols, abbreviations, and acronyms table supporting up to 5 entries

### Figures & Media
- ImageFigure - Figure component with image reference and auto-numbered caption

### Navigation Components
- TableOfContents - Auto-generated table of contents supporting up to 3 sections
- TableList - List of tables supporting up to 3 table references
- FigureList - List of figures supporting up to 3 figure references

### Document Structure
- StandardHeader - Document header with classification and company information
- StandardFooter - Document footer with restrictions text and classification

### Legal & Compliance
- DisclaimerStatement - Multi-part disclaimer block with technical rights, export control, handling, and disclaimer statements
- DisclaimerPage - Full page of general instructions and warnings for technical documents

## Standard Company Document Layout

### Title Page Layout
For typical company documents, components should be arranged in this vertical order on the first page:

1. **StandardHeader** - Document classification, logos, and company information
2. **DocumentCategoryTitle** - Document type/category identifier
3. **DocumentTitle** - Main document title
4. **DocumentSubject** - Document number/revision
5. **TestBlock** - Test execution details (for test documents)
6. **AuthorBlock** - Author and company contact information
7. **StandardFooter** - Restrictions text and classification

### Front Matter Pages
For comprehensive documents, include these pages after the title page:

1. **RevisionBlock** - Document revision history (typically on separate page)
2. **DisclaimerPage** or **DisclaimerStatement** - Legal and compliance information
3. **TableOfContents** - Document navigation
4. **TableList** - List of tables (if applicable)
5. **FigureList** - List of figures (if applicable)
6. **AcronymList** - Symbols, abbreviations, and acronyms

### Body Content
For document body sections, use:

- **HeadingOneSection** - Major section headings with introductory content
- **HeadingTwoSection** - Subsection headings with content
- **ReferenceTable** - Tables of reference documents
- **ImageFigure** - Figures and diagrams with captions

## Usage in Document Plans

Components are referenced in JSON document plans using this structure:

```json
{
  "component": "ComponentName",
  "props": {
    "prop_name": "value",
    "another_prop": "another value"
  }
}
```

## Component Creation Guidelines

When creating new components, follow the lessons learned in `CLAUDE.md`:

1. **XML Simplification**: Remove all revision tracking metadata and redundant namespaces
2. **Strategic Namespaces**: Add namespaces only where first needed on specific elements
3. **Parameterization**: Replace hard-coded text with `{{ prop_name }}` placeholders
4. **Styling Preservation**: Maintain essential paragraph and run properties for visual consistency

For detailed component creation workflows, see:
- `docs/example-component-extraction.md` - AI-assisted component authoring
- `docs/asset-generation-procedure.md` - Manual component creation process