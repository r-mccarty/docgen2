# DocumentSubject Component

## Purpose

Renders a document subject or revision identifier line, typically used below the main title to display document numbers, revisions, or classification information.

## Visual Description

- Right-aligned subject text
- Bold formatting with 32pt font size
- Single line spacing with automatic rule
- Integrated with Word's document subject metadata

## Props

| Prop Name | Type | Required | Description |
|-----------|------|----------|-------------|
| `document_subject` | string | Yes | The subject text to display (e.g., "DOC-2145, Rev A", "UNCLASSIFIED", etc.) |

## Usage Example

```json
{
  "component": "DocumentSubject",
  "props": {
    "document_subject": "DOC-2791, Rev A"
  }
}
```

## Styling Notes

- Bold text formatting applied through run properties
- Right-aligned paragraph alignment
- 32pt font size for prominence
- Preserves Word's document subject metadata binding
- Single line spacing with automatic line rule

## Technical Details

- Contains structured document tag (SDT) with subject metadata binding
- Optimized XML structure with revision metadata removed
- Essential formatting properties preserved for consistent styling
- Compatible with Word's core document properties system