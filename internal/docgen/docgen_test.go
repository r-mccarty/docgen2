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