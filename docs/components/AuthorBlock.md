# AuthorBlock Component

## Purpose

Renders a complete author contact information block with company details, typically used at the bottom of title pages or document sections. This component provides comprehensive contact information in a professional, right-aligned format.

## Visual Description

- Bold "Prepared by" header
- Right-aligned contact information layout
- Author name integrated with Word's document metadata
- Complete company address block
- Phone, fax, and website contact details
- Hyperlinked website URL with Word's "Hyperlink" style

## Props

| Prop Name | Type | Required | Description |
|-----------|------|----------|-------------|
| `author_name` | string | Yes | Name of the document author/preparer |
| `company_name` | string | Yes | Company or organization name |
| `address_line1` | string | Yes | First line of company address |
| `address_line2` | string | Yes | Second line of company address (suite, unit, etc.) |
| `city_state_zip` | string | Yes | City, state, and ZIP code |
| `phone` | string | Yes | Phone number (format: "(XXX) XXX-XXXX") |
| `fax` | string | Yes | Fax number (format: "(XXX) XXX-XXXX") |
| `website` | string | Yes | Company website URL |

## Usage Example

```json
{
  "component": "AuthorBlock",
  "props": {
    "author_name": "Ryan McCarty",
    "company_name": "Innoflight",
    "address_line1": "9985 Pacific Heights Blvd.",
    "address_line2": "Suite 250",
    "city_state_zip": "San Diego, CA 92121",
    "phone": "(858) 638-1580",
    "fax": "(858) 638-1581",
    "website": "https://www.innoflight.com"
  }
}
```

## Styling Notes

- Right-aligned paragraph layout for all contact information
- Single line spacing (240 units) with automatic line rule
- Bold formatting for "Prepared by" header
- Hyperlink style applied to website URL
- Author name integrated with Word's document creator metadata binding
- Consistent spacing with no after-paragraph spacing (w:after="0")

## Technical Details

- Contains structured document tag (SDT) with author metadata binding
- Hyperlink element with relationship ID for website
- Bookmark elements for Word navigation compatibility
- Optimized XML structure with revision metadata removed
- Essential paragraph properties preserved for formatting consistency

## Standard Page Layout Position

In a typical company document layout, this component appears in the following vertical order:
1. Header (to be created)
2. DocumentCategoryTitle
3. DocumentTitle
4. DocumentSubject
5. TestBlock
6. **AuthorBlock** ← This component (bottom of page)