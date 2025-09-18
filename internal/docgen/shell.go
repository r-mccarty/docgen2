package docgen

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
)

// LoadShell loads a DOCX shell document into memory
func LoadShell(shellPath string) (InMemoryDocx, error) {
	zipReader, err := zip.OpenReader(shellPath)
	if err != nil {
		return nil, &ShellLoadError{Path: shellPath, Err: err}
	}
	defer zipReader.Close()

	shell := make(InMemoryDocx)

	for _, file := range zipReader.File {
		reader, err := file.Open()
		if err != nil {
			return nil, &ShellLoadError{Path: shellPath, Err: fmt.Errorf("failed to open file %s: %w", file.Name, err)}
		}

		content, err := io.ReadAll(reader)
		reader.Close()
		if err != nil {
			return nil, &ShellLoadError{Path: shellPath, Err: fmt.Errorf("failed to read file %s: %w", file.Name, err)}
		}

		shell[file.Name] = content
	}

	return shell, nil
}

// Clone creates a deep copy of the shell document for safe concurrent use
func (shell InMemoryDocx) Clone() InMemoryDocx {
	clone := make(InMemoryDocx)
	for path, content := range shell {
		// Create a new slice and copy the content
		contentCopy := make([]byte, len(content))
		copy(contentCopy, content)
		clone[path] = contentCopy
	}
	return clone
}

// ToBytes serializes the in-memory DOCX back to a byte slice
func (shell InMemoryDocx) ToBytes() ([]byte, error) {
	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	for path, content := range shell {
		fileWriter, err := zipWriter.Create(path)
		if err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to create zip entry for %s: %w", path, err)
		}

		if _, err := fileWriter.Write(content); err != nil {
			zipWriter.Close()
			return nil, fmt.Errorf("failed to write content for %s: %w", path, err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		return nil, fmt.Errorf("failed to close zip writer: %w", err)
	}

	return buf.Bytes(), nil
}