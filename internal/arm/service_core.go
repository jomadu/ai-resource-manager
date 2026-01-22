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
func NewArmService() *ArmService {
	return &ArmService{
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
		ui:              ui.New(false), // Set debug to false by default
	}
}

// ShowVersion displays the current version, build information, and build datetime
func (a *ArmService) ShowVersion() error {
	a.ui.VersionInfo(a.Version())
	return nil
}

func (a *ArmService) Version() version.VersionInfo {
	return version.GetVersionInfo()
}
