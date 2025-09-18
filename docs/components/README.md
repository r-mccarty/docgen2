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

### Content Blocks
- [TestBlock](./TestBlock.md) - Test form with tester, date, serial number, result, and additional info fields
- [AuthorBlock](./AuthorBlock.md) - Author contact information block with company details

## Standard Company Document Layout

For typical company documents, components should be arranged in this vertical order on the first page:

1. **Header** (to be created) - Document classification, logos, etc.
2. **DocumentCategoryTitle** - Document type/category identifier
3. **DocumentTitle** - Main document title
4. **DocumentSubject** - Document number/revision
5. **TestBlock** - Test execution details (for test documents)
6. **AuthorBlock** - Author and company contact information

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