package urf

import "fmt"

// AmazonQFilenameGenerator generates Amazon Q-specific filenames
type AmazonQFilenameGenerator struct{}

// GenerateFilename generates an Amazon Q filename
func (g *AmazonQFilenameGenerator) GenerateFilename(rulesetID, ruleID string) string {
	return fmt.Sprintf("%s_%s.md", rulesetID, ruleID)
}
