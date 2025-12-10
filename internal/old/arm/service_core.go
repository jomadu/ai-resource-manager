package arm

import (
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
)

// ArmService orchestrates all ARM operations.
type ArmService struct {
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService() *ArmService {
	return &ArmService{
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
	}
}
