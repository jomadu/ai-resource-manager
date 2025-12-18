package sink

import (
	"context"
	"core"
)

type ResourceType string

const (
	ResourceTypeRuleset   ResourceType = "ruleset"
	ResourceTypePromptset ResourceType = "promptset"
)

type PackageInstallation struct {
	Metadata     core.PackageMetadata
	ResourceType ResourceType
	FilePaths    []string
}

// Manager handles installation and management of packages for a specific sink directory.
// Each Manager instance is bound to a single sink directory and handles all operations
// for that directory (compilation, layout, file placement, etc.).
type Manager interface {
	// InstallRuleset installs a ruleset package to this sink.
	// Priority is required for rulesets to determine rule ordering when multiple rulesets are installed.
	InstallRuleset(ctx context.Context, ruleset *core.Package, priority int) error

	// InstallPromptset installs a promptset package to this sink.
	// Promptsets do not use priority as they are not ordered.
	InstallPromptset(ctx context.Context, promptset *core.Package) error

	UninstallPackage(ctx context.Context, packageMetaData *core.PackageMetadata) error
	ListInstalledPackages(ctx context.Context) ([]*PackageInstallation, error)
	IsPackageInstalled(ctx context.Context, packageMetaData *core.PackageMetadata) (bool, error)
}
