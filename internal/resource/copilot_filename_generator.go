package resource

import "fmt"

// CopilotFilenameGenerator generates GitHub Copilot filenames
type CopilotFilenameGenerator struct{}

// GenerateFilename generates a GitHub Copilot filename
func (g *CopilotFilenameGenerator) GenerateFilename(rulesetID, ruleID string) string {
	return fmt.Sprintf("%s_%s.instructions.md", rulesetID, ruleID)
}
