package resource

import "fmt"

// CursorFilenameGenerator generates Cursor-specific filenames
type CursorFilenameGenerator struct{}

// GenerateFilename generates a Cursor filename
func (g *CursorFilenameGenerator) GenerateFilename(rulesetID, ruleID string) string {
	return fmt.Sprintf("%s_%s.mdc", rulesetID, ruleID)
}
