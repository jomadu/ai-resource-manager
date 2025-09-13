package urf

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// AmazonQCompiler compiles URF to Amazon Q format
type AmazonQCompiler struct {
	metadataGen MetadataGenerator
}

// NewAmazonQCompiler creates a new Amazon Q compiler
func NewAmazonQCompiler() Compiler {
	return &AmazonQCompiler{
		metadataGen: NewMetadataGenerator(),
	}
}

// Compile compiles URF to Amazon Q format
func (c *AmazonQCompiler) Compile(urf *URFFile) ([]*types.File, error) {
	var files []*types.File
	for _, rule := range urf.Rules {
		filename := fmt.Sprintf("%s_%s.md", urf.Metadata.ID, rule.ID)
		content := c.generateContent(urf, &rule)
		files = append(files, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}
	return files, nil
}

// generateContent generates Amazon Q-specific content
func (c *AmazonQCompiler) generateContent(urf *URFFile, rule *Rule) string {
	var content strings.Builder

	// Amazon Q frontmatter
	content.WriteString(c.metadataGen.GenerateMetadata(urf, rule, urf.Metadata.ID))

	// Rule content
	enforcement := strings.ToUpper(rule.Enforcement)
	content.WriteString(fmt.Sprintf("# %s (%s)\n\n", rule.Name, enforcement))
	content.WriteString(rule.Body)

	return content.String()
}
