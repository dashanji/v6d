package templates

import (
	"embed"
	"fmt"
	"path/filepath"
)

//go:embed vineyardd
var fs embed.FS

// EmbedTemplate is only used for implementing the interface
type EmbedTemplate struct{}

// NewEmbedTemplate returns a new EmbedTemplate
func NewEmbedTemplate() *EmbedTemplate {
	return &EmbedTemplate{}
}

// ReadFile reads a file from the embed.FS
func (e *EmbedTemplate) ReadFile(path string) ([]byte, error) {
	return fs.ReadFile(path)
}

// GetFilesRecursive returns all files in a directory
func (e *EmbedTemplate) GetFilesRecursive(dir string) ([]string, error) {
	path := filepath.Join(filepath.Dir(dir), dir)
	fd, err := fs.ReadDir(path)
	if err != nil {
		return []string{}, fmt.Errorf("ReadDir Error: %v", err)
	}
	files := []string{}
	for _, f := range fd {
		if !f.IsDir() {
			files = append(files, filepath.Join(path, f.Name()))
		}
	}
	return files, nil
}
