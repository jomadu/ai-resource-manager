package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [ruleset...]",
		Short: "Install rulesets",
		Long:  "Install rulesets from a registry. If no ruleset is specified, installs from manifest.",
		RunE:  runInstall,
	}

	cmd.Flags().StringSlice("include", nil, "Include patterns")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns")
	cmd.Flags().StringSlice("sinks", nil, "Target sinks")
	cmd.Flags().Int("priority", 100, "Ruleset installation priority (1-1000+)")

	return cmd
}

func runInstall(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// If no arguments, install from manifest
	if len(args) == 0 {
		return armService.InstallManifest(ctx)
	}

	// Parse arguments
	rulesets, err := ParseRulesetArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Get flags
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	sinks, _ := cmd.Flags().GetStringSlice("sinks")
	priority, _ := cmd.Flags().GetInt("priority")
	include = GetDefaultIncludePatterns(include)

	// Validate priority
	if priority < 1 {
		return fmt.Errorf("priority must be a positive integer (got %d)", priority)
	}

	// Require sinks for new installations
	if len(sinks) == 0 {
		return fmt.Errorf("--sinks is required for installing rulesets")
	}

	// Install each ruleset
	for _, ruleset := range rulesets {
		req := &arm.InstallRequest{
			Registry: ruleset.Registry,
			Ruleset:  ruleset.Name,
			Version:  ruleset.Version,
			Priority: priority,
			Include:  include,
			Exclude:  exclude,
			Sinks:    sinks,
		}
		err := armService.InstallRuleset(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to install %s/%s: %w", ruleset.Registry, ruleset.Name, err)
		}
	}

	return nil
}
