package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade resources",
		Long:  "Upgrade all configured resources to latest versions, or use subcommands for specific resource types.",
		RunE:  runUpgradeAll,
	}

	// Add subcommands
	cmd.AddCommand(newUpgradeRulesetCmd())
	cmd.AddCommand(newUpgradePromptsetCmd())

	return cmd
}

func newUpgradeRulesetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ruleset [registry/ruleset...]",
		Short: "Upgrade rulesets",
		Long:  "Upgrade specific rulesets to latest versions or all rulesets if none specified.",
		RunE:  runUpgradeRuleset,
	}
}

func newUpgradePromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset [registry/promptset...]",
		Short: "Upgrade promptsets",
		Long:  "Upgrade specific promptsets to latest versions or all promptsets if none specified.",
		RunE:  runUpgradePromptset,
	}
}

func runUpgradeAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	return armService.UpgradeAll(ctx)
}

func runUpgradeRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// If no arguments, upgrade all rulesets
	if len(args) == 0 {
		// Use the unified upgrade which handles all rulesets
		return armService.UpgradeAll(ctx)
	}

	// Parse arguments and upgrade specific rulesets
	packages, err := ParsePackageArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	for _, pkg := range packages {
		err = armService.UpgradeRuleset(ctx, pkg.Registry, pkg.Name)
		if err != nil {
			return fmt.Errorf("failed to upgrade ruleset %s/%s: %w", pkg.Registry, pkg.Name, err)
		}
	}

	return nil
}

func runUpgradePromptset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// If no arguments, upgrade all promptsets
	if len(args) == 0 {
		// Use the unified upgrade which handles all promptsets
		return armService.UpgradeAll(ctx)
	}

	// Parse arguments and upgrade specific promptsets
	packages, err := ParsePackageArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	for _, pkg := range packages {
		err = armService.UpgradePromptset(ctx, pkg.Registry, pkg.Name)
		if err != nil {
			return fmt.Errorf("failed to upgrade promptset %s/%s: %w", pkg.Registry, pkg.Name, err)
		}
	}

	return nil
}
