# DocumentCategoryTitle Component

## Purpose

Renders a category header with decorative underline, typically used to identify the document type or category (e.g., "TEST PROCEDURE", "DESIGN SPECIFICATION", "USER MANUAL"). This component provides visual separation and hierarchical structure at the top of document sections.

## Visual Description

- Right-aligned category title text using Word's "Title" paragraph style
- Decorative horizontal line element below the text
- Line spans approximately 5.94 inches (5943600 EMUs) with 3pt stroke weight
- Modern drawing element with VML fallback for compatibility
- Black solid line with rounded caps

## Props

| Prop Name | Type | Required | Description |
|-----------|------|----------|-------------|
| `category_title` | string | Yes | The category or document type text to display (e.g., "TEST PROCEDURE", "DESIGN SPECIFICATION") |

## Usage Example

```json
{
  "component": "DocumentCategoryTitle",
  "props": {
    "category_title": "TEST PROCEDURE"
  }
}
```

## Styling Notes

- Uses Word's built-in "Title" paragraph style for semantic consistency
- Right-aligned text positioning
- Bold formatting with 28pt font size for the decorative line context
- Complex drawing element with both modern (`<w:drawing>`) and legacy (`<w:pict>`) fallbacks
- Strategic namespace declarations only where needed (mc, wp, v namespaces)

## Technical Details

- Contains `mc:AlternateContent` for modern/legacy compatibility
- Modern drawing uses WordprocessingML drawing namespace
- VML fallback ensures compatibility with older Word versions
- Optimized XML structure with revision metadata removed (~60% size reduction)
- Essential paragraph properties preserved for styling integrity

## Standard Page Layout Position

In a typical company document layout, this component appears in the following vertical order:
1. Header (to be created)
2. **DocumentCategoryTitle** ← This component
3. DocumentTitle
4. DocumentSubject
5. TestBlock
6. AuthorBlock