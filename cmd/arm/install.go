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
	include = GetDefaultIncludePatterns(include)

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
