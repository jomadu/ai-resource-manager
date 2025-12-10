package arm

import (
	"context"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/manifest"
)

// Service provides the main ARM functionality for managing AI resources (rulesets and promptsets).
type Service interface {
	// Registry Management
	AddGitRegistry(ctx context.Context, name, url string, branches []string, force bool) error
	AddGitLabRegistry(ctx context.Context, name, url, projectID, groupID, apiVersion string, force bool) error
	AddCloudsmithRegistry(ctx context.Context, name, url, owner, repository string, force bool) error
	RemoveRegistry(ctx context.Context, name string) error
	SetRegistryConfig(ctx context.Context, name, field, value string) error
	GetRegistries(ctx context.Context) (map[string]map[string]interface{}, error)

	// Sink Management
	AddSink(ctx context.Context, name, directory, layout, compileTarget string, force bool) error
	RemoveSink(ctx context.Context, name string) error
	SetSinkConfig(ctx context.Context, name, field, value string) error
	GetSinks(ctx context.Context) (map[string]manifest.SinkConfig, error)

	// Package Management (unified operations)
	InstallAll(ctx context.Context) error
	UpdateAll(ctx context.Context) error
	UpgradeAll(ctx context.Context) error
	UninstallAll(ctx context.Context) error
	GetOutdatedPackages(ctx context.Context) ([]*OutdatedPackage, error)

	// Ruleset Management
	InstallRuleset(ctx context.Context, req *InstallRulesetRequest) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	UpdateAllRulesets(ctx context.Context) error
	UpgradeRuleset(ctx context.Context, registry, ruleset string) error
	SetRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error
	GetRulesets(ctx context.Context) ([]*RulesetInfo, error)

	// Promptset Management
	InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error
	UninstallPromptset(ctx context.Context, registry, promptset string) error
	UpdatePromptset(ctx context.Context, registry, promptset string) error
	UpdateAllPromptsets(ctx context.Context) error
	UpgradePromptset(ctx context.Context, registry, promptset string) error
	SetPromptsetConfig(ctx context.Context, registry, promptset, field, value string) error
	GetPromptsets(ctx context.Context) ([]*PromptsetInfo, error)

	// Utilities
	CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error
	NukeCache(ctx context.Context) error
	CleanSinks(ctx context.Context) error
	NukeSinks(ctx context.Context) error
	CompileFiles(ctx context.Context, req *CompileRequest) error
}
