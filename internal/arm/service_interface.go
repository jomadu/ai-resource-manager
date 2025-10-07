package arm

import (
	"context"
	"time"
)

// Service provides the main ARM functionality for managing AI resources (rulesets and promptsets).
type Service interface {
	// Ruleset operations
	InstallRuleset(ctx context.Context, req *InstallRulesetRequest) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	SetRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error

	// Promptset operations
	InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error
	UninstallPromptset(ctx context.Context, registry, promptset string) error
	UpdatePromptset(ctx context.Context, registry, promptset string) error
	SetPromptsetConfig(ctx context.Context, registry, promptset, field, value string) error

	// Unified operations
	InstallAll(ctx context.Context) error
	UpdateAll(ctx context.Context) error
	UpgradeAll(ctx context.Context) error
	UninstallAll(ctx context.Context) error

	// Configuration operations
	AddRegistry(ctx context.Context, name, url, regType string, options map[string]interface{}) error
	RemoveRegistry(ctx context.Context, name string) error
	AddSink(ctx context.Context, name, directory, sinkType, layout, compileTarget string, force bool) error
	RemoveSink(ctx context.Context, name string) error
	SetRegistryConfig(ctx context.Context, name, field, value string) error
	SetSinkConfig(ctx context.Context, name, field, value string) error

	// Cache operations
	CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error
	NukeCache(ctx context.Context) error

	// Info operations
	ShowVersion() error
	ShowRulesetInfo(ctx context.Context, rulesets []string) error
	ShowPromptsetInfo(ctx context.Context, promptsets []string) error
	ShowOutdated(ctx context.Context, outputFormat string, noSpinner bool) error
	ShowList(ctx context.Context, sortByPriority bool) error

	// Compile operations
	CompileFiles(ctx context.Context, req *CompileRequest) error
}
