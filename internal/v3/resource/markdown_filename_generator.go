package resource

import "fmt"

// MarkdownFilenameGenerator generates markdown filenames
type MarkdownFilenameGenerator struct{}

// GenerateFilename generates a markdown filename
func (g *MarkdownFilenameGenerator) GenerateFilename(rulesetID, ruleID string) string {
	return fmt.Sprintf("%s_%s.md", rulesetID, ruleID)
}
