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
	Install(ctx context.Context, registry, ruleset, version string, include, exclude []string) error
	InstallFromManifest(ctx context.Context) error
	Uninstall(ctx context.Context, registry, ruleset string) error
	Update(ctx context.Context, registry, ruleset string) error
	UpdateFromManifest(ctx context.Context) error
	Outdated(ctx context.Context) ([]OutdatedRuleset, error)
	List(ctx context.Context) ([]InstalledRuleset, error)
	Info(ctx context.Context, registry, ruleset string) (*RulesetInfo, error)
	InfoAll(ctx context.Context) ([]*RulesetInfo, error)
	Version() VersionInfo
}

// ArmService orchestrates all ARM operations.
type ArmService struct {
	configManager   config.Manager
	manifestManager manifest.Manager
	lockFileManager lockfile.Manager
	installer       installer.Installer
}

// NewArmService creates a new ARM service instance with all dependencies.
func NewArmService() *ArmService {
	return &ArmService{
		configManager:   config.NewFileManager(),
		manifestManager: manifest.NewFileManager(),
		lockFileManager: lockfile.NewFileManager(),
		installer:       installer.NewFileInstaller(),
	}
}

func (a *ArmService) Install(ctx context.Context, registry, ruleset, version string, include, exclude []string) error {
	// 1. Validate registry exists in config
	// 2. Parse version constraint (^1.0.0, ~1.0.0, =1.0.0, branch name, or latest)
	// 3. Resolve version from registry (git tags/branches)
	// 4. Download ruleset files matching include patterns
	// 5. Remove existing version directories if updating
	// 6. Update manifest with new ruleset entry
	// 7. Update lockfile with resolved version and metadata
	// 8. Install files to configured sink directories
	return errors.New("not implemented")
}

func (a *ArmService) InstallFromManifest(ctx context.Context) error {
	// 1. Check if arm.json exists, if not check arm.lock
	// 2. If only arm.lock exists, generate arm.json from lockfile
	// 3. For each ruleset in manifest, call Install()
	return errors.New("not implemented")
}

func (a *ArmService) Uninstall(ctx context.Context, registry, ruleset string) error {
	// 1. Remove ruleset entry from manifest
	// 2. Remove ruleset entry from lockfile
	// 3. Remove installed files from sink directories
	// 4. Clean up empty ARM directories if no rulesets remain
	return errors.New("not implemented")
}

func (a *ArmService) Update(ctx context.Context, registry, ruleset string) error {
	// 1. Get current constraint and include/exclude from manifest
	// 2. Resolve latest version within constraint from registry
	// 3. If newer version available, call Install() with new version
	return errors.New("not implemented")
}

func (a *ArmService) Outdated(ctx context.Context) ([]OutdatedRuleset, error) {
	// 1. For each ruleset in lockfile:
	//    - Get current resolved version
	//    - Calculate wanted version within constraint
	//    - Get latest available version from registry
	// 2. Return table data with Current/Wanted/Latest columns
	// 3. Show "All rulesets are up to date!" if none outdated
	return nil, errors.New("not implemented")
}

func (a *ArmService) List(ctx context.Context) ([]InstalledRuleset, error) {
	// 1. Read lockfile to get all installed rulesets
	// 2. Return list in format: registry/ruleset@version
	// 3. Sort by registry then ruleset name
	return nil, errors.New("not implemented")
}

func (a *ArmService) Info(ctx context.Context, registry, ruleset string) (*RulesetInfo, error) {
	// 1. Get ruleset details from lockfile and manifest
	// 2. Get registry URL and type from config
	// 3. Calculate version information (constraint/resolved/wanted/latest)
	// 4. Find matching sinks based on include/exclude patterns
	// 5. List installation directories
	// 6. Return formatted info structure
	return nil, errors.New("not implemented")
}

func (a *ArmService) UpdateFromManifest(ctx context.Context) error {
	// 1. For each ruleset in manifest, call Update()
	return errors.New("not implemented")
}

func (a *ArmService) InfoAll(ctx context.Context) ([]*RulesetInfo, error) {
	// 1. Get all installed rulesets from List()
	// 2. For each ruleset, call Info() and collect results
	return nil, errors.New("not implemented")
}

func (a *ArmService) Version() VersionInfo {
	// 1. Return build-time version info (version, commit, arch)
	// 2. Format: "arm 1.2.3\ncommit: a1b2c3d4\narch: darwin/arm64"
	return VersionInfo{}
}
