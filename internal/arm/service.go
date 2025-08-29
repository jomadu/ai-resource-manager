package arm

import (
	"context"
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/installer"
	"github.com/jomadu/ai-rules-manager/internal/lockfile"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
)

// Service provides the main ARM functionality for managing AI rule rulesets.
type Service interface {
	Install(ctx context.Context, ruleset, version string, include, exclude []string) error
	InstallFromManifest(ctx context.Context) error
	Uninstall(ctx context.Context, ruleset string) error
	Update(ctx context.Context, ruleset string) error
	Outdated(ctx context.Context) ([]OutdatedRuleset, error)
	List(ctx context.Context) ([]InstalledRuleset, error)
	Info(ctx context.Context, ruleset string) (*RulesetInfo, error)
	Version() VersionInfo
}

// ArmService orchestrates all ARM operations.
type ArmService struct {
	configManager   config.Manager
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
	installer       installer.Installer
}

// NewArmService creates a new ARM service instance.
func NewArmService(
	configManager config.Manager,
	manifestManager manifest.Manager,
	lockFileManager lockfile.Manager,
	installer installer.Installer,
) *ArmService {
	return &ArmService{
		configManager:   configManager,
		manifestManager: manifestManager,
		lockFileManager: lockFileManager,
		installer:       installer,
	}
}

func (a *ArmService) Install(ctx context.Context, ruleset, version string, include, exclude []string) error {
	return errors.New("not implemented")
}

func (a *ArmService) InstallFromManifest(ctx context.Context) error {
	return errors.New("not implemented")
}

func (a *ArmService) Uninstall(ctx context.Context, ruleset string) error {
	return errors.New("not implemented")
}

func (a *ArmService) Update(ctx context.Context, ruleset string) error {
	return errors.New("not implemented")
}

func (a *ArmService) Outdated(ctx context.Context) ([]OutdatedRuleset, error) {
	return nil, errors.New("not implemented")
}

func (a *ArmService) List(ctx context.Context) ([]InstalledRuleset, error) {
	return nil, errors.New("not implemented")
}

func (a *ArmService) Info(ctx context.Context, ruleset string) (*RulesetInfo, error) {
	return nil, errors.New("not implemented")
}

func (a *ArmService) Version() VersionInfo {
	return VersionInfo{}
}
