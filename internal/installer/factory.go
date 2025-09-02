package installer

import "github.com/jomadu/ai-rules-manager/internal/config"

// NewInstaller creates the appropriate installer based on sink layout configuration.
func NewInstaller(sink *config.SinkConfig) Installer {
	switch sink.Layout {
	case "flat":
		return NewFlatInstaller()
	case "hierarchical", "":
		return NewHierarchicalInstaller()
	default:
		// Default to hierarchical for unknown layouts
		return NewHierarchicalInstaller()
	}
}
