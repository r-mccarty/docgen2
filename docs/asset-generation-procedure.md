## Guide: Authoring DocGen Assets (The Manual Workflow)

### Part 1: Your Mission & Mindset

Your goal is to deconstruct a master document template into a set of reusable, styled "Lego bricks" (Components) and a foundational "baseplate" (the Shell).

**The Golden Rule:** From this point forward, stop thinking about formatting (bold, centered, etc.). Instead, think in terms of **semantic styles**. You are not making the word "Tester:" bold; you are applying the `ComponentLabel` style to it. This distinction is everything.

**Prerequisites:**
1.  **Microsoft Word:** You will need the desktop application.
2.  **Enable the Developer Tab in Word:** Go to `File > Options > Customize Ribbon` and check the box for "Developer". This isn't strictly necessary for this workflow but is invaluable for advanced Word tasks.
3.  **A Text Editor:** A good text editor like VS Code (with an XML formatter) will make parameterizing the components much easier.
4.  **A File Archiver:** A tool like 7-Zip or the built-in Windows/Mac archiver to handle `.zip` files.

---

### Phase I: The Foundation — Creating the Shell & Styles

The Shell Document is the heart of your document's look and feel. It contains **NO content**, only the definitions for every style your components will ever need.

#### Step 1: Start with a "Golden" Master Template
Open the existing `.docx` document that represents the ideal final output. This is your source of truth. Save it under a new name, like `template_shell_WIP.docx`.

#### Step 2: Define and Refine Your Styles
This is the most important step. You need to create a consistent set of styles that will be used by all your components.

1.  In Word, open the **Styles Pane** (on the Home tab, click the small arrow in the corner of the Styles section).
2.  Go through your document and identify every distinct type of text. For each one, create a style.
    *   **Example:** For the "Test Details" section:
        *   The labels like "**Tester:**" and "**Serial Number:**" are conceptually the same. Create a new style called `ComponentLabel`. Set its font to be bold, size 10pt, etc.
        *   The values like "Jonathan Schallert" and "INF-0656" are also the same. Create a new style called `ComponentValue`. Set its font to be regular, size 10pt, etc.
        *   The main document title ("Engineering Design Verification Test") should have a style like `DocMainTitle`.
3.  **Apply these new styles** to the text in your `template_shell_WIP.docx`. This makes it easy to see them all in one place.
4.  **Pro Tip:** Use the "Modify..." option in the Styles pane to fine-tune spacing, fonts, and other properties. Be meticulous here, as this defines the look of your final documents.

#### Step 3: Strip the Content to Create the Shell
Once all your styles are defined and applied, it's time to create the pure shell.

1.  **Select all content** in the document (`Ctrl+A`).
2.  **Press Delete.** The document body should now be completely blank.
3.  **Crucially, the styles you created are still saved in the document!** You can verify this in the Styles pane.
4.  Go into the Header & Footer and delete all content there as well. The header and footer sections themselves should remain, but they should be empty.
5.  **Save this empty document** as `template_shell.docx`.

**Deliverable:** You now have a `template_shell.docx` file. It looks empty, but it's packed with all the style definitions your engine will need. This file goes into the `/assets/` directory of the Go project.

---

### Phase II: The Building Blocks — Creating the Component Library

Now you will create the individual "Lego bricks."

**The Component Workflow Rule:** For every component, you will create a temporary, single-purpose Word document. **Never work inside the giant XML of your original master template.**

#### Step 4: Create a Component (Example: `TestDetails`)

Let's build the `TestDetails` component from our example.

**A. Create the Scaffold Document:**
   - Create a brand new, blank Word document. Save it as `TestDetails_scaffold.docx`.

**B. Isolate the Visual Element:**
   - Go back to your original, fully-contented master document.
   - Select and copy *only the table* containing the Test Details.
   - Paste this table into `TestDetails_scaffold.docx`.

**C. Apply Your Styles:**
   - Ensure the content inside this scaffold document uses the styles you defined in Phase I. For example, make sure "Tester:" has the `ComponentLabel` style applied, and "Jonathan Schallert" has the `ComponentValue` style. This links the component to the shell.

**D. Extract the Clean XML:**
   1. Close Word.
   2. Rename `TestDetails_scaffold.docx` to `TestDetails_scaffold.zip`.
   3. Unzip the file. You will get a folder with the OpenXML structure inside.
   4. Navigate to the `word/` subfolder.
   5. Open `document.xml` in your text editor (like VS Code). It will be small and clean!

**E. Identify and Copy the Core Snippet:**
   - The content you pasted was a table. In the XML, find the block that starts with `<w:tbl>` and ends with `</w:tbl>`.
   - This block is your component. Copy the entire thing to your clipboard.

**F. Parameterize the Snippet:**
   1. Create a new, blank file in your text editor.
   2. Paste the `<w:tbl>...</w:tbl>` snippet into it.
   3. Now, you will replace the hard-coded text with `{{prop_name}}` placeholders. Search for the `<w:t>` (text) tags.
      *   Find `<w:t>Jonathan Schallert</w:t>` and change it to `<w:t>{{ tester_name }}</w:t>`.
      *   Find `<w:t>6/20/2024</w:t>` and change it to `<w:t>{{ test_date }}</w:t>`.
      *   Find `<w:t>INF-0656</w:t>` and change it to `<w:t>{{ serial_number }}</w:t>`.
      *   ...and so on for `test_result` and `completed_by`.

**G. Save the Component File:**
   - Save this new file as `TestDetails.component.xml`. The name is important: the part before `.component.xml` (`TestDetails`) is the name the Planner service will use in the JSON plan.
   - Place this file in the `/assets/components/` directory of the Go project.

#### Step 5: Rinse and Repeat

You have now created your first component! Follow the **exact same process (Steps A-G)** for every other reusable part of your document:

*   **`DocumentTitle.component.xml`**: The title, subtitle, doc number block.
*   **`Paragraph.component.xml`**: A single paragraph. Its prop will be `{{ text }}`.
*   **`Section.component.xml`**: This one is special. It will contain the XML for a heading (e.g., using the "Heading 1" style) and a placeholder where child components will be injected. Your Go developer will need to implement a special way to handle injecting children here. For now, just create the heading part with a prop for `{{ title }}`.
*   **`ResultsTable.component.xml`**: A table with a header row. The props for this will be more complex (`{{ headers }}` and `{{ rows }})`, and the Go developer will need to write logic to loop through the data and generate the `<w:tr>` (table row) elements dynamically. For the template, just create a table with a header and one placeholder data row.

**Deliverables:** A populated `/assets/components/` directory, full of parameterized `.component.xml` files.

---

### Phase III: The Rule Book — Creating the CUE Schema

While the Go engineer will write the CUE code, you are responsible for defining the business rules.

#### Step 6: Document the Component "API"
For each component you created, write down a simple list of the props you defined.

*   **Component: `TestDetails`**
    *   `tester_name` (string)
    *   `test_date` (string)
    *   `serial_number` (string)
    *   `test_result` (string, must be "PASS" or "FAIL")
    *   `completed_by` (string)

This list is the specification that the Go engineer will translate into the `rules.cue` file, and it's the documentation the Planner team will use to create valid JSON.

By following these phases, you will systematically and reliably create the complete set of assets needed for the DocGen engine to function.