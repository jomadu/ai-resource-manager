package installer

import "github.com/jomadu/ai-rules-manager/internal/manifest"

// NewInstaller creates the appropriate installer based on sink layout configuration.
func NewInstaller(sink *manifest.SinkConfig) Installer {
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
