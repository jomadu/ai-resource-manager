package service

import (
	"context"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
	"github.com/jomadu/ai-resource-manager/internal/v4/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/v4/sink"
)

type PackageInfo struct {
	Installation sink.PackageInstallation
	LockInfo     packagelockfile.PackageLockInfo
	Config       map[string]interface{}
}

type OutdatedPackage struct {
	Current    core.PackageMetadata
	Constraint string
	Wanted     core.PackageMetadata
	Latest     core.PackageMetadata
}

// ARM
type ArmService interface {
	// ---------------------------------------------
	// Registry Management (Git, GitLab, Cloudsmith)
	// ---------------------------------------------
	// Add
	AddGitRegistry(ctx context.Context, name, url string, branches []string, force bool) error
	AddGitLabRegistry(ctx context.Context, name, url, projectID, groupID, apiVersion string, force bool) error
	AddCloudsmithRegistry(ctx context.Context, name, url, owner, repository string, force bool) error
	// Remove
	RemoveRegistry(ctx context.Context, name string) error
	// Get
	GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error)
	GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error)
	// Set
	SetRegistryName(ctx context.Context, name string, newName string) error
	SetRegistryURL(ctx context.Context, name string, url string) error
	SetGitRegistryBranches(ctx context.Context, name string, branches []string) error
	SetGitLabRegistryProjectID(ctx context.Context, name string, projectID string) error
	SetGitLabRegistryGroupID(ctx context.Context, name string, groupID string) error
	SetGitLabRegistryAPIVersion(ctx context.Context, name string, apiVersion string) error
	SetCloudsmithRegistryOwner(ctx context.Context, name string, owner string) error
	SetCloudsmithRegistryRepository(ctx context.Context, name string, repository string) error

	// ---------------
	// Sink Management
	// ---------------
	// Add
	AddSink(ctx context.Context, name, directory, layout, compileTarget string, force bool) error
	AddSinkOfType(ctx context.Context, name, directory string, sinkType string, force bool) error
	// Remove
	RemoveSink(ctx context.Context, name string) error
	// Get
	GetSinkConfig(ctx context.Context, name string) (*manifest.SinkConfig, error)
	GetAllSinkConfigs(ctx context.Context) (map[string]*manifest.SinkConfig, error)
	// Set
	SetSinkName(ctx context.Context, name string, newName string) error
	SetSinkType(ctx context.Context, name string, sinkType string) error
	SetSinkDirectory(ctx context.Context, name string, directory string) error
	SetSinkLayout(ctx context.Context, name string, layout string) error
	SetSinkCompileTarget(ctx context.Context, name string, compileTarget string) error

	// ------------------
	// Package Management
	// ------------------
	InstallAll(ctx context.Context) error
	UninstallAll(ctx context.Context) error
	UpdateAll(ctx context.Context) error
	UpgradeAll(ctx context.Context) error
	ListAll(ctx context.Context) ([]*PackageInfo, error)
	ListOutdated(ctx context.Context) ([]*OutdatedPackage, error)
	GetPackageInfo(ctx context.Context, registry, packageName string) (*PackageInfo, error)

	// ------------------
	// Ruleset Management
	// ------------------
	// Add
	InstallRuleset(ctx context.Context, registry, ruleset, version string, priority int, include []string, exclude []string, sinks []string) error
	// Remove
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	// Set
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	UpdateAllRulesets(ctx context.Context) error
	UpgradeRuleset(ctx context.Context, registry, ruleset string) error
	SetRulesetName(ctx context.Context, registry, ruleset, newName string) error
	SetRulesetVersion(ctx context.Context, registry, ruleset, version string) error
	SetRulesetPriority(ctx context.Context, registry, ruleset, priority int) error
	SetRulesetInclude(ctx context.Context, registry, ruleset string, include []string) error
	SetRulesetExclude(ctx context.Context, registry, ruleset string, exclude []string) error
	SetRulesetSinks(ctx context.Context, registry, ruleset string, sinks []string) error

	// --------------------
	// Promptset Management
	// --------------------
	// Add
	InstallPromptset(ctx context.Context, registry, promptset, version string, include []string, exclude []string, sinks []string) error
	// Remove
	UninstallPromptset(ctx context.Context, registry, ruleset string) error
	// Set
	UpdatePromptset(ctx context.Context, registry, ruleset string) error
	UpdateAllPromptsets(ctx context.Context) error
	UpgradePromptset(ctx context.Context, registry, ruleset string) error
	SetPromptsetName(ctx context.Context, registry, ruleset, newName string) error
	SetPromptsetVersion(ctx context.Context, registry, ruleset, version string) error
	SetPromptsetInclude(ctx context.Context, registry, ruleset string, include []string) error
	SetPromptsetExclude(ctx context.Context, registry, ruleset string, exclude []string) error
	SetPromptsetSinks(ctx context.Context, registry, ruleset string, sinks []string) error

	// --------
	// Cleaning
	// --------
	CleanCacheByAge(ctx context.Context, maxAge time.Duration) error
	CleanCacheByTimeSinceLastAccess(ctx context.Context, maxTimeSinceLastAccess time.Duration) error
	NukeCache(ctx context.Context) error
	CleanSinks(ctx context.Context) error
	NukeSinks(ctx context.Context) error

	// Compile
	CompileFiles(ctx context.Context, req *CompileRequest) error
}
