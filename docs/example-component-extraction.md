## Playbook: AI-Assisted Component Authoring

**Your Role:** The Architect/QA. You find and isolate the raw material and verify the final product.
**AI's Role:** The XML Technician. It performs the detailed, error-prone XML extraction and parameterization based on your instructions.

### The Workflow Cycle (Repeat for Each Component)

This cycle consists of five distinct steps: **1. Isolate, 2. Contextualize, 3. Instruct, 4. Implement (AI), 5. Verify.**

#### Step 1: ISOLATE (Your Task)
Your job is to create the clean, minimal XML snippet that the AI will work on. This is the most important step, as it provides high-quality input.

1.  **Create the Scaffold:** In Word, create a new, blank document (e.g., `TestDetails_scaffold.docx`).
2.  **Copy the Element:** From your master template, copy *only the single, complete visual element* you want to componentize (e.g., the "Test Details" table, the full title block, the disclaimer box).
3.  **Paste into Scaffold:** Paste the element into the scaffold document.
4.  **Extract the XML:** Rename the scaffold to `.zip`, unzip it, and get the `document.xml` file.

**Result:** You have a small, focused `document.xml` file that contains only the raw material for one component.

#### Step 2: CONTEXTUALIZE (Your Task)
Now, you prepare the prompt for your AI assistant. You need to give it three pieces of information: the raw material, the goal, and the specific details.

1.  **The Raw Material:** Copy the *entire content* of the small `document.xml` file you just extracted.
2.  **The Goal:** Define what you want to achieve.
3.  **The Details:** List the specific props you want to create and what text they should replace.

#### Step 3: INSTRUCT (Your Prompt to the AI)
This is where you combine the context into a clear, structured prompt. Use a template like the one below.

**Prompt Template:**

```text
Hello! I need your help creating a DocGen component from a piece of OpenXML. Please follow these instructions precisely.

**ROLE:**
You are an expert OpenXML technician. Your task is to extract an XML snippet from a `document.xml` file, clean it up, and parameterize it by replacing hard-coded text with placeholder props.

**CONTEXT:**
The component I want to create is called "**[Name of Component, e.g., DocumentCategoryTitle]**". It represents the main category title and its decorative line.

Here is the full content of the `document.xml` from a scaffold Word document containing only this element:
```xml
[PASTE THE ENTIRE CONTENTS OF YOUR document.xml HERE]
```

**INSTRUCTIONS:**
1.  **Identify the Core XML:** Analyze the XML I provided. The component consists of two paragraph (`<w:p>`) blocks. The first contains the title text, and the second contains the `<w:drawing>` element for a horizontal line.
2.  **Clean the XML:**
    *   For the first paragraph block containing the text "Acceptance Test Procedure", I want you to **remove the `<w:sdt>` (Content Control) wrapper**.
    *   Preserve the paragraph properties (`<w:pPr>`) from the original paragraph.
    *   Preserve the run (`<w:r>`) and text (`<w:t>`) elements from *inside* the `<w:sdtContent>` block.
3.  **Parameterize the Text:** Replace the hard-coded text with a placeholder prop.
    *   Replace `Acceptance Test Procedure` with `{{ category_title }}`.
4.  **Combine the Parts:** Create a final XML output that contains the cleaned, parameterized title paragraph followed immediately by the second, unmodified paragraph containing the `<w:drawing>` element.
5.  **Final Output:** Provide ONLY the final, combined XML snippet. Do not wrap it in a `<w:document>` or `<w:body>` tag. The output should start with the first `<w:p>` tag and end with the last `</w:p>` tag.

**Summary of Props to Create:**
*   `category_title` should replace "Acceptance Test Procedure".

Please generate the final component XML.
```

#### Step 4: IMPLEMENT (AI's Task)
The AI assistant will now process your request. Given the detailed instructions and the clean input, it should be able to produce the final component XML with high accuracy.

**Expected AI Output:**
```xml
<w:p>
  <w:pPr>
    <w:pStyle w:val="Title"/>
    <w:jc w:val="right"/>
  </w:pPr>
  <w:r>
    <w:t>{{ category_title }}</w:t>
  </w:r>
</w:p>
<w:p w14:paraId="00A0F9E5" w14:textId="77777777" w:rsidR="009A3938" w:rsidRDefault="009A3938" w:rsidP="00C20629">
    <w:pPr>
        <w:spacing w:after="0"/>
        <w:jc w:val="right"/>
        <w:rPr>
            <w:b/>
            <w:sz w:val="28"/>
        </w:rPr>
    </w:pPr>
    <w:r w:rsidRPr="0041538E">
        <w:rPr>
            <w:noProof/>
        </w:rPr>
        <mc:AlternateContent>
            <mc:Choice Requires="wps">
                <w:drawing>
                <!-- ... FULL DRAWING ELEMENT ... -->
                </w:drawing>
            </mc:Choice>
            <mc:Fallback>
                <!-- ... FALLBACK ELEMENT ... -->
            </mc:Fallback>
        </mc:AlternateContent>
    </w:r>
</w:p>
```

#### Step 5: VERIFY (Your Task)
This is your final quality assurance check.

1.  **Copy the AI's Output:** Take the XML the AI generated.
2.  **Create the Component File:** Paste the XML into a new file and save it as `[ComponentName].component.xml` (e.g., `DocumentCategoryTitle.component.xml`).
3.  **Review the Output:**
    *   Does it look right? Are both paragraph blocks present?
    *   Is the `<w:sdt>` wrapper gone?
    *   Is the `{{ category_title }}` prop correctly placed?
    *   **Quick Sanity Check:** Compare the prop names in the file with the list you created in the prompt. Make sure there are no typos.
4.  **Commit:** Once verified, commit the new component file to your project's `/assets/components/` directory.

By following this playbook, you create a powerful human-AI partnership. You handle the high-level strategy and quality control, while the AI handles the complex and repetitive low-level manipulation, dramatically accelerating the creation of your component library.