package installer

import (
	"github.com/jomadu/ai-rules-manager/internal/manifest"
)

// NewInstaller creates the appropriate installer based on sink layout configuration.
func NewInstaller(sink *manifest.SinkConfig) Installer {
	switch sink.Layout {
	case "flat":
		return NewFlatInstaller(sink.Directory, sink.CompileTarget)
	case "hierarchical", "":
		return NewHierarchicalInstaller(sink.Directory, sink.CompileTarget)
	default:
		// Default to hierarchical for unknown layouts
		return NewHierarchicalInstaller(sink.Directory, sink.CompileTarget)
	}
}
