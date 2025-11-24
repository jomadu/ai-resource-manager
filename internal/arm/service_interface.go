package arm

import (
	"context"
	"time"
)

// Service provides the main ARM functionality for managing AI resources (rulesets and promptsets).
type Service interface {
	// Core operations
	ShowVersion() error

	// Registry Management
	AddRegistry(ctx context.Context, name, url, regType string, options map[string]interface{}, force bool) error
	RemoveRegistry(ctx context.Context, name string) error
	SetRegistryConfig(ctx context.Context, name, field, value string) error
	ShowRegistryList(ctx context.Context) error
	ShowRegistryInfo(ctx context.Context, registries []string) error

	// Sink Management
	AddSink(ctx context.Context, name, directory, layout, compileTarget string, force bool) error
	RemoveSink(ctx context.Context, name string) error
	SetSinkConfig(ctx context.Context, name, field, value string) error
	ShowSinkList(ctx context.Context) error
	ShowSinkInfo(ctx context.Context, sinks []string) error

	// Package Management (unified operations)
	InstallAll(ctx context.Context) error
	UpdateAll(ctx context.Context) error
	UpgradeAll(ctx context.Context) error
	UninstallAll(ctx context.Context) error
	ShowAllOutdated(ctx context.Context, outputFormat string, noSpinner bool) error
	ShowAllList(ctx context.Context, sortByPriority bool) error
	ShowAllInfo(ctx context.Context) error

	// Ruleset Management
	InstallRuleset(ctx context.Context, req *InstallRulesetRequest) error
	UninstallRuleset(ctx context.Context, registry, ruleset string) error
	UpdateRuleset(ctx context.Context, registry, ruleset string) error
	UpdateAllRulesets(ctx context.Context) error
	UpgradeRuleset(ctx context.Context, registry, ruleset string) error
	SetRulesetConfig(ctx context.Context, registry, ruleset, field, value string) error
	ShowRulesetInfo(ctx context.Context, rulesets []string) error
	ShowRulesetList(ctx context.Context, sortByPriority bool) error
	ShowRulesetOutdated(ctx context.Context, outputFormat string, noSpinner bool) error

	// Promptset Management
	InstallPromptset(ctx context.Context, req *InstallPromptsetRequest) error
	UninstallPromptset(ctx context.Context, registry, promptset string) error
	UpdatePromptset(ctx context.Context, registry, promptset string) error
	UpdateAllPromptsets(ctx context.Context) error
	UpgradePromptset(ctx context.Context, registry, promptset string) error
	SetPromptsetConfig(ctx context.Context, registry, promptset, field, value string) error
	ShowPromptsetInfo(ctx context.Context, promptsets []string) error
	ShowPromptsetList(ctx context.Context) error
	ShowPromptsetOutdated(ctx context.Context, outputFormat string, noSpinner bool) error

	// Utilities
	CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error
	NukeCache(ctx context.Context) error
	CleanSinks(ctx context.Context) error
	NukeSinks(ctx context.Context) error
	CompileFiles(ctx context.Context, req *CompileRequest) error
}
