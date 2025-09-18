package docgen

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestEngineInitialization(t *testing.T) {
	// Test with the actual assets
	shellPath := "../../assets/shell/template_shell.docx"
	componentsDir := "../../assets/components"

	engine, err := NewEngine(shellPath, componentsDir)
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	if engine == nil {
		t.Fatal("Engine is nil")
	}

	// Check that components were loaded
	components := engine.GetLoadedComponents()
	if len(components) == 0 {
		t.Fatal("No components were loaded")
	}

	t.Logf("Loaded components: %v", components)

	// Verify that DocumentCategoryTitle component exists
	found := false
	for _, comp := range components {
		if comp == "DocumentCategoryTitle" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("DocumentCategoryTitle component not found")
	}
}

func TestComponentRendering(t *testing.T) {
	template := "Hello {{ name }}, welcome to {{ place }}!"
	props := map[string]interface{}{
		"name":  "World",
		"place": "DocGen",
	}

	result, err := RenderComponent(template, props)
	if err != nil {
		t.Fatalf("Failed to render component: %v", err)
	}

	expected := "Hello World, welcome to DocGen!"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestComponentRenderingWithXMLEscaping(t *testing.T) {
	template := "Title: {{ title }}"
	props := map[string]interface{}{
		"title": "Test & Demo <Example>",
	}

	result, err := RenderComponent(template, props)
	if err != nil {
		t.Fatalf("Failed to render component: %v", err)
	}

	expected := "Title: Test &amp; Demo &lt;Example&gt;"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestDocumentAssembly(t *testing.T) {
	// Test with the actual assets
	shellPath := "../../assets/shell/template_shell.docx"
	componentsDir := "../../assets/components"
	planPath := "../../assets/plans/test_plan_01.json"

	// Initialize engine
	engine, err := NewEngine(shellPath, componentsDir)
	if err != nil {
		t.Fatalf("Failed to initialize engine: %v", err)
	}

	// Read test plan
	planData, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("Failed to read test plan: %v", err)
	}

	var plan DocumentPlan
	if err := json.Unmarshal(planData, &plan); err != nil {
		t.Fatalf("Failed to parse test plan: %v", err)
	}

	// Assemble document
	result, err := engine.Assemble(plan)
	if err != nil {
		t.Fatalf("Failed to assemble document: %v", err)
	}

	if len(result) == 0 {
		t.Fatal("Generated document is empty")
	}

	t.Logf("Generated document size: %d bytes", len(result))

	// Write test output for manual verification
	outputPath := "../../output/test_output.docx"
	os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err := os.WriteFile(outputPath, result, 0644); err != nil {
		t.Logf("Warning: Could not write test output file: %v", err)
	} else {
		t.Logf("Test output written to: %s", outputPath)
	}
}

func TestShellCloning(t *testing.T) {
	shell := InMemoryDocx{
		"test1.xml": []byte("content1"),
		"test2.xml": []byte("content2"),
	}

	clone := shell.Clone()

	// Verify the clone has the same content
	if len(clone) != len(shell) {
		t.Errorf("Clone has different number of files: %d vs %d", len(clone), len(shell))
	}

	for path, content := range shell {
		cloneContent, exists := clone[path]
		if !exists {
			t.Errorf("File %s missing from clone", path)
			continue
		}

		if string(content) != string(cloneContent) {
			t.Errorf("Content differs for file %s", path)
		}
	}

	// Verify modifying clone doesn't affect original
	clone["test1.xml"] = []byte("modified")
	if string(shell["test1.xml"]) == "modified" {
		t.Error("Modifying clone affected original")
	}
}

// setupTestEngine creates a test engine for use in unit tests
func setupTestEngine(t *testing.T) *Engine {
	engine, err := NewEngine("../../assets/shell/template_shell.docx", "../../assets/components/")
	if err != nil {
		t.Fatalf("Failed to create test engine: %v", err)
	}
	return engine
}

// writeTestOutput writes test output to a structured directory for manual inspection
func writeTestOutput(t *testing.T, name string, data []byte) {
	outputDir := "../../test_output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		t.Logf("Warning: Could not create test output directory: %v", err)
		return
	}

	outputPath := filepath.Join(outputDir, name+".docx")
	if err := os.WriteFile(outputPath, data, 0644); err != nil {
		t.Logf("Warning: Could not write test output file %s: %v", outputPath, err)
	} else {
		t.Logf("Test output written to: %s", outputPath)
	}
}

// TestAssembleComponents provides comprehensive unit testing of the component library
func TestAssembleComponents(t *testing.T) {
	engine := setupTestEngine(t)

	testCases := []struct {
		name string
		plan DocumentPlan
	}{
		{
			name: "DocumentCategoryTitle_Component",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_DocumentCategoryTitle.docx"},
				Body: []ComponentInstance{
					{
						Component: "DocumentCategoryTitle",
						Props: map[string]interface{}{
							"category_title": "TEST PROCEDURE",
						},
					},
				},
			},
		},
		{
			name: "DocumentTitle_Component",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_DocumentTitle.docx"},
				Body: []ComponentInstance{
					{
						Component: "DocumentTitle",
						Props: map[string]interface{}{
							"document_title": "Engineering Design Verification Test - CFC-400XS Extended DVT Procedure",
						},
					},
				},
			},
		},
		{
			name: "DocumentSubject_Component",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_DocumentSubject.docx"},
				Body: []ComponentInstance{
					{
						Component: "DocumentSubject",
						Props: map[string]interface{}{
							"document_subject": "DOC-2791, Rev A",
						},
					},
				},
			},
		},
		{
			name: "TestBlock_Component",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_TestBlock.docx"},
				Body: []ComponentInstance{
					{
						Component: "TestBlock",
						Props: map[string]interface{}{
							"tester_name":     "Ryan McCarty",
							"test_date":       "6/20/2024",
							"serial_number":   "INF-0656",
							"test_result":     "PASS",
							"additional_info": "Completed all test iterations successfully",
						},
					},
				},
			},
		},
		{
			name: "AuthorBlock_Component",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_AuthorBlock.docx"},
				Body: []ComponentInstance{
					{
						Component: "AuthorBlock",
						Props: map[string]interface{}{
							"author_name":    "Ryan McCarty",
							"company_name":   "Innoflight",
							"address_line1":  "9985 Pacific Heights Blvd.",
							"address_line2":  "Suite 250",
							"city_state_zip": "San Diego, CA 92121",
							"phone":          "(858) 638-1580",
							"fax":            "(858) 638-1581",
							"website":        "https://www.innoflight.com",
						},
					},
				},
			},
		},
		{
			name: "Full_Integration_Test",
			plan: DocumentPlan{
				DocProps: DocProps{Filename: "unit_test_full_integration.docx"},
				Body: []ComponentInstance{
					{
						Component: "DocumentCategoryTitle",
						Props: map[string]interface{}{
							"category_title": "DESIGN VERIFICATION PROCEDURE",
						},
					},
					{
						Component: "DocumentTitle",
						Props: map[string]interface{}{
							"document_title": "PCA-1153-01/02 (12V Supervisor Board) Safe To Mate",
						},
					},
					{
						Component: "DocumentSubject",
						Props: map[string]interface{}{
							"document_subject": "DOC-3421, Rev B",
						},
					},
					{
						Component: "TestBlock",
						Props: map[string]interface{}{
							"tester_name":     "Sarah Chen",
							"test_date":       "9/18/2024",
							"serial_number":   "PCA-1153-SN-001",
							"test_result":     "PASS",
							"additional_info": "All electrical and mechanical specifications verified",
						},
					},
					{
						Component: "AuthorBlock",
						Props: map[string]interface{}{
							"author_name":    "Sarah Chen",
							"company_name":   "Innoflight",
							"address_line1":  "9985 Pacific Heights Blvd.",
							"address_line2":  "Suite 250",
							"city_state_zip": "San Diego, CA 92121",
							"phone":          "(858) 638-1580",
							"fax":            "(858) 638-1581",
							"website":        "https://www.innoflight.com",
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Assemble the document
			result, err := engine.Assemble(tc.plan)
			if err != nil {
				t.Fatalf("Failed to assemble document for %s: %v", tc.name, err)
			}

			if len(result) == 0 {
				t.Fatalf("Generated document is empty for %s", tc.name)
			}

			t.Logf("%s: Generated document size: %d bytes", tc.name, len(result))

			// Write test output for manual verification
			writeTestOutput(t, tc.name, result)
		})
	}
}