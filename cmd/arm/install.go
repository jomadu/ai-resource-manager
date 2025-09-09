package main

import (
	"context"
	"fmt"

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
	include = GetDefaultIncludePatterns(include)

	// Install each ruleset
	for _, ruleset := range rulesets {
		err := armService.InstallRuleset(ctx, ruleset.Registry, ruleset.Name, ruleset.Version, include, exclude)
		if err != nil {
			return fmt.Errorf("failed to install %s/%s: %w", ruleset.Registry, ruleset.Name, err)
		}
	}

	return nil
}
