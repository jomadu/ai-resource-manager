package arm

import (
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/ui"
	"github.com/jomadu/ai-rules-manager/internal/version"
)

// ArmService orchestrates all ARM operations.
type ArmService struct {
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
	ui              ui.Interface
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService(uiInterface ui.Interface) *ArmService {
	return &ArmService{
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
		ui:              uiInterface,
	}
}

func (a *ArmService) Version() version.VersionInfo {
	return version.GetVersionInfo()
}
