# DocumentTitle Component

## Purpose

Renders a structured document title block commonly used for technical documents, reports, and specifications. The component displays the main document title within a structured data template (SDT) that integrates with Word's document properties system.

## Visual Description

- Right-aligned title text
- Uses Word's "Title" paragraph style
- Styled with "IntenseReference" character style
- 32pt font size with Times New Roman/Arial font family
- Integrated with Word's document metadata binding

## Props

| Prop Name | Type | Required | Description |
|-----------|------|----------|-------------|
| `document_title` | string | Yes | The main title text to display. Can be multi-part (e.g., "PCA-1153-01/02 (12V Supervisor Board) Safe To Mate") |

## Usage Example

```json
{
  "component": "DocumentTitle",
  "props": {
    "document_title": "Engineering Design Verification Test - CFC-400XS Extended DVT Procedure"
  }
}
```

## Styling Notes

- The component preserves Word's built-in "Title" and "IntenseReference" styles
- Text is automatically right-aligned
- The SDT structure allows Word to treat this as document metadata
- Font characteristics are defined through semantic styles, not direct formatting

## Technical Details

- Contains structured document tag (SDT) with metadata binding
- Preserves document property integration for Word compatibility
- Optimized XML with revision tracking metadata removed
- Essential paragraph and run properties maintained for styling integrity