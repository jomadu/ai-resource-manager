package urf

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// CursorCompiler compiles URF to Cursor format
type CursorCompiler struct {
	metadataGen MetadataGenerator
}

// NewCursorCompiler creates a new Cursor compiler
func NewCursorCompiler() Compiler {
	return &CursorCompiler{
		metadataGen: NewMetadataGenerator(),
	}
}

// Compile compiles URF to Cursor format
func (c *CursorCompiler) Compile(urf *URFFile) ([]*types.File, error) {
	var files []*types.File
	for _, rule := range urf.Rules {
		filename := fmt.Sprintf("%s_%s.mdc", urf.Metadata.ID, rule.ID)
		content := c.generateContent(urf, &rule)
		files = append(files, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}
	return files, nil
}

// generateContent generates Cursor-specific content
func (c *CursorCompiler) generateContent(urf *URFFile, rule *Rule) string {
	var content strings.Builder

	// Cursor frontmatter
	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("description: %q\n", rule.Description))
	if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
		content.WriteString(fmt.Sprintf("globs: %v\n", rule.Scope[0].Files))
	}
	if rule.Enforcement == "must" {
		content.WriteString("alwaysApply: true\n")
	}
	content.WriteString("---\n\n")

	// Namespace metadata
	content.WriteString(c.metadataGen.GenerateMetadata(urf, rule, urf.Metadata.ID))

	// Rule content
	enforcement := strings.ToUpper(rule.Enforcement)
	content.WriteString(fmt.Sprintf("# %s (%s)\n\n", rule.Name, enforcement))
	content.WriteString(rule.Body)

	return content.String()
}
