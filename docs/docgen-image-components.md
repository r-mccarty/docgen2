## Guide: Authoring Image Components for DocGen

### Part 1: The Core Concept — Images and Relationships

In OpenXML, an image is **not** embedded directly into the `document.xml` file. Instead, the process works like this:

1.  **Storage:** The image file (e.g., `company_logo.png`) is stored in the `/word/media/` folder inside the `.docx` package.
2.  **Relationship:** A special "relationship" file, `word/_rels/document.xml.rels`, contains an entry that assigns a unique **Relationship ID (rId)** to that image file. It acts like a pointer, saying "`rId7` points to `../media/company_logo.png`".
3.  **Reference:** The main `document.xml` file doesn't mention the filename at all. It contains a drawing object that says "display the image identified by `rId7` here."

This means that to create a dynamic image component, our DocGen engine needs to do three things:
1.  Add the new image file to the media folder.
2.  Create a new, unique relationship ID for it in the `.rels` file.
3.  Use that new `rId` in the component's XML.

This guide will show you how to author a component that makes this possible.

---

### Phase I: Creating a Parameterizable Image Component

We will create a component that can display a different image based on a prop. The key is to create a template that the Go engine can easily find and modify.

#### Step 1: Create the Image Component Scaffold
1.  Create a brand new, blank Word document. Save it as `ImageBlock_scaffold.docx`.
2.  Insert a placeholder image into the document. **Any small image will do.** This image acts as a template for the size, position, and text wrapping properties we want our final image to have.
3.  **Crucially, set the image's "Alt Text"**.
    *   Right-click the image in Word and select "Edit Alt Text...".
    *   In the Alt Text box, enter a unique placeholder value that will act as our "prop name." A good convention is `{{prop:image_prop_name}}`. For this example, let's use:
      `{{prop:company_logo}}`
    *   This Alt Text is stored in the XML and will be our hook for the Go engine to identify this image reference.

#### Step 2: Extract the Image Component XML
1.  Save and close `ImageBlock_scaffold.docx`.
2.  Rename it to `.zip` and extract the contents.
3.  Navigate to `word/` and open `document.xml` in your text editor.

#### Step 3: Isolate and Prepare the Component Snippet
You will see a large block of XML for the image, likely starting with `<w:p>` and containing a `<w:drawing>` element. It will look complex, but you only need to focus on one part.

Inside the `<w:drawing>` element, find the `<a:blip>` tag. This tag contains the reference to the image. It will look something like this:

```xml
<a:blip r:embed="rId4" ... >
    ...
</a:blip>
```
The `r:embed="rId4"` is the reference we need to make dynamic.

**Your Task:**
1.  Copy the entire XML block for the image (starting from the parent `<w:p>` paragraph tag) into a new file named `ImageBlock.component.xml`.
2.  In this new file, find the `<a:blip>` tag.
3.  Change the `r:embed` attribute to use a placeholder that matches your Alt Text prop. The format should be `{{rId:prop_name}}`.
    *   Change `r:embed="rId4"` to `r:embed="{{rId:company_logo}}"`

The `{{prop:...}}` in the Alt Text tells the engine *which image component this is*, and the `{{rId:...}}` in the `r:embed` attribute tells the engine *where to inject the new Relationship ID it creates*.

Here is a simplified example of what your final `ImageBlock.component.xml` will look like:
```xml
<!-- ImageBlock.component.xml -->
<w:p>
  <w:pPr>...</w:pPr> <!-- Paragraph properties for alignment etc. -->
  <w:r>
    <w:drawing>
      <wp:inline ...>
        <!-- ... lots of other XML defining size and position ... -->
        <a:graphic>
          <a:graphicData uri="...">
            <pic:pic>
              <!-- ... Picture properties ... -->
              <pic:blipFill>
                <!-- THIS IS THE IMPORTANT PART -->
                <a:blip r:embed="{{rId:company_logo}}" /> 
                <a:stretch>
                  <a:fillRect/>
                </a:stretch>
              </pic:blipFill>
              <!-- ... Shape properties ... -->
            </pic:pic>
          </a:graphicData>
        </a:graphic>
      </wp:inline>
    </w:drawing>
  </w:r>
</w:p>
```

**Deliverable:** You now have an `ImageBlock.component.xml` in your `/assets/components/` directory.

---

### Phase II: The Document Plan (For the Planner Service)

To use this new component, the Planner service must provide the image data in the JSON plan. Raw image data is large, so the best practice is to use **Base64 encoding**.

The `props` for an image component in the JSON plan will look like this:

```json
{
  "component": "ImageBlock",
  "props": {
    "company_logo": {
      "filename": "innoflight_logo.png",
      "content_base64": "iVBORw0KGgoAAAANSUhEUgAAA...your...long...base64...string...here...=="
    }
  }
}
```

*   `company_logo`: This key **must match** the prop name you defined (`{{prop:company_logo}}`).
*   `filename`: The desired filename for the image inside the `.docx` package. This helps with debugging.
*   `content_base64`: The Base64-encoded string of the image file.

---

### Phase III: The Go Engine's Responsibility (For the Engineer)

This section documents the new logic the Go engineer must add to the DocGen service to handle this image component.

When the **Assembler** processes the plan, it needs a special workflow for images:

1.  **Detect Image Props:** Before rendering the main `document.xml`, the assembler must first iterate through the entire plan and find all props that are image objects (i.e., they have `filename` and `content_base64` keys).

2.  **Process Each Image:** For each unique image found in the plan:
    a. **Decode:** Decode the Base64 string back into raw image bytes (`[]byte`).
    b. **Add to Media:** Add a new entry to the in-memory representation of the `.docx` package at `word/media/the_filename.png` with the decoded bytes.
    c. **Generate New rId:** Create a new, unique Relationship ID that hasn't been used yet in `word/_rels/document.xml.rels` (e.g., if the last one is `rId10`, the new one is `rId11`).
    d. **Add Relationship:** Add a new `<Relationship>` XML node to the in-memory `.rels` file. This new node links the new `rId11` to the target file (`Target="../media/the_filename.png"`).
    e. **Store the Mapping:** Keep a map of the prop name to the new `rId`. For our example: `map["company_logo"] = "rId11"`.

3.  **Render Components:** Now, proceed with the normal component rendering process.
    a. When rendering `ImageBlock.component.xml`, the templating engine will encounter the `{{rId:company_logo}}` placeholder.
    b. The engine looks up `"company_logo"` in the map it created in step 2e and finds the value `"rId11"`.
    c. It replaces the placeholder, so the final XML in the document becomes `<a:blip r:embed="rId11" />`.

This process correctly adds the image file, creates the necessary relationship link, and injects the correct reference into the document body, allowing for fully dynamic, data-driven images.