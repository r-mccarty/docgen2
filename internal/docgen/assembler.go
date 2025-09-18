package docgen

import (
	"fmt"

	"github.com/beevik/etree"
)

// AssembleDocument assembles components into the shell document according to the plan
func (e *Engine) AssembleDocument(plan DocumentPlan) ([]byte, error) {
	// Create a working copy of the shell document
	workingDoc := e.shell.Clone()

	// Get the document.xml content
	documentXML, exists := workingDoc["word/document.xml"]
	if !exists {
		return nil, NewDocGenError("assembly", fmt.Errorf("word/document.xml not found in shell document"))
	}

	// Parse the document XML
	doc := etree.NewDocument()
	if err := doc.ReadFromBytes(documentXML); err != nil {
		return nil, NewDocGenError("assembly", fmt.Errorf("failed to parse document.xml: %w", err))
	}

	// Find the document body
	body := doc.FindElement("//w:body")
	if body == nil {
		return nil, NewDocGenError("assembly", fmt.Errorf("w:body element not found in document.xml"))
	}

	// Process each component in the plan
	for _, componentInstance := range plan.Body {
		if err := e.addComponentToBody(body, componentInstance); err != nil {
			return nil, NewDocGenError("assembly", fmt.Errorf("failed to add component %s: %w", componentInstance.Component, err))
		}
	}

	// Serialize the modified document back to bytes
	doc.Indent(2)
	modifiedXML, err := doc.WriteToBytes()
	if err != nil {
		return nil, NewDocGenError("assembly", fmt.Errorf("failed to serialize modified document.xml: %w", err))
	}

	// Update the document.xml in our working copy
	workingDoc["word/document.xml"] = modifiedXML

	// Convert back to DOCX bytes
	result, err := workingDoc.ToBytes()
	if err != nil {
		return nil, NewDocGenError("assembly", fmt.Errorf("failed to create final DOCX: %w", err))
	}

	return result, nil
}

// addComponentToBody adds a rendered component to the document body
func (e *Engine) addComponentToBody(body *etree.Element, componentInstance ComponentInstance) error {
	// Get the component template
	template, err := e.GetComponent(componentInstance.Component)
	if err != nil {
		return err
	}

	// Render the component with props
	renderedXML, err := RenderComponent(template, componentInstance.Props)
	if err != nil {
		return fmt.Errorf("failed to render component: %w", err)
	}

	// Parse the rendered component XML
	componentDoc := etree.NewDocument()

	// Wrap the component XML in a temporary root to handle multiple top-level elements
	wrappedXML := fmt.Sprintf("<temp xmlns:w=\"http://schemas.openxmlformats.org/wordprocessingml/2006/main\" xmlns:mc=\"http://schemas.openxmlformats.org/markup-compatibility/2006\" xmlns:wp=\"http://schemas.openxmlformats.org/drawingml/2006/wordprocessingDrawing\" xmlns:wp14=\"http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing\" xmlns:a=\"http://schemas.openxmlformats.org/drawingml/2006/main\" xmlns:wps=\"http://schemas.microsoft.com/office/word/2010/wordprocessingShape\" xmlns:a14=\"http://schemas.microsoft.com/office/drawing/2010/main\" xmlns:v=\"urn:schemas-microsoft-com:vml\" xmlns:w10=\"urn:schemas-microsoft-com:office:word\" xmlns:o=\"urn:schemas-microsoft-com:office:office\">%s</temp>", renderedXML)

	if err := componentDoc.ReadFromString(wrappedXML); err != nil {
		return fmt.Errorf("failed to parse rendered component XML: %w", err)
	}

	// Find the temporary root and add all its children to the body
	tempRoot := componentDoc.Root()
	if tempRoot != nil {
		for _, child := range tempRoot.ChildElements() {
			// Clone the element to avoid modifying the original
			clonedChild := child.Copy()
			body.AddChild(clonedChild)
		}
	}

	return nil
}