package arm

import (
	"context"
	"time"
)

// Service provides the main ARM functionality for managing AI rule rulesets.
type Service interface {
	// Ruleset operations
	InstallRuleset(ctx context.Context, req *InstallRequest) error
	InstallManifest(ctx context.Context) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	UpdateAllRulesets(ctx context.Context) error
	UpdateRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error

	// Configuration operations
	ShowConfig(ctx context.Context) error
	AddRegistry(ctx context.Context, name, url, regType string, options map[string]interface{}) error
	RemoveRegistry(ctx context.Context, name string) error
	AddSink(ctx context.Context, name, directory, sinkType, layout, compileTarget string, force bool) error
	RemoveSink(ctx context.Context, name string) error
	UpdateRegistry(ctx context.Context, name, field, value string) error
	UpdateSink(ctx context.Context, name, field, value string) error

	// Cache operations
	CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error
	NukeCache(ctx context.Context) error

	// Info operations
	ShowVersion() error
	ShowRulesetInfo(ctx context.Context, rulesets []string) error
	ShowOutdated(ctx context.Context, outputFormat string, noSpinner bool) error
	ShowList(ctx context.Context, sortByPriority bool) error
}
