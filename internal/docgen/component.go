package docgen

import (
	"fmt"
	"html"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// LoadComponents loads all .component.xml files from the specified directory
func LoadComponents(componentsDir string) (map[string]string, error) {
	components := make(map[string]string)

	err := filepath.Walk(componentsDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".component.xml") {
			// Extract component name from filename (remove .component.xml extension)
			componentName := strings.TrimSuffix(info.Name(), ".component.xml")

			content, err := readComponentFile(path)
			if err != nil {
				return fmt.Errorf("failed to read component %s: %w", componentName, err)
			}

			components[componentName] = content
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load components from %s: %w", componentsDir, err)
	}

	return components, nil
}

// readComponentFile reads the content of a component file
func readComponentFile(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// RenderComponent renders a component template with the given props
func RenderComponent(template string, props map[string]interface{}) (string, error) {
	if len(props) == 0 {
		return template, nil
	}

	// Build replacement pairs for strings.NewReplacer
	var replacements []string
	for key, value := range props {
		placeholder := fmt.Sprintf("{{ %s }}", key)
		// Convert value to string and XML-escape it
		valueStr := fmt.Sprintf("%v", value)
		escapedValue := html.EscapeString(valueStr)

		replacements = append(replacements, placeholder, escapedValue)
	}

	replacer := strings.NewReplacer(replacements...)
	return replacer.Replace(template), nil
}

// GetComponent retrieves a component template by name
func (e *Engine) GetComponent(componentName string) (string, error) {
	template, exists := e.components[componentName]
	if !exists {
		return "", &ComponentNotFoundError{ComponentName: componentName}
	}
	return template, nil
}