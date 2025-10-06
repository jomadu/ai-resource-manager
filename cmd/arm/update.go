package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
		Long:  "Update all configured resources, or use subcommands for specific resource types.",
		RunE:  runUpdateAll,
	}

	// Add subcommands
	cmd.AddCommand(newUpdateRulesetCmd())
	cmd.AddCommand(newUpdatePromptsetCmd())

	return cmd
}

func newUpdateRulesetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ruleset [registry/ruleset...]",
		Short: "Update rulesets",
		Long:  "Update specific rulesets or all rulesets if none specified.",
		RunE:  runUpdateRuleset,
	}
}

func newUpdatePromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset [registry/promptset...]",
		Short: "Update promptsets",
		Long:  "Update specific promptsets or all promptsets if none specified.",
		RunE:  runUpdatePromptset,
	}
}

func runUpdateAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	// TODO: Implement unified update when service interface is updated
	return armService.UpdateAllRulesets(ctx) // Temporary fallback to rulesets only
}

func runUpdateRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// If no arguments, update all rulesets
	if len(args) == 0 {
		return armService.UpdateAllRulesets(ctx)
	}

	// Parse arguments
	rulesets, err := ParsePackageArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Update each ruleset
	for _, ruleset := range rulesets {
		err := armService.UpdateRuleset(ctx, ruleset.Registry, ruleset.Name)
		if err != nil {
			return fmt.Errorf("failed to update %s/%s: %w", ruleset.Registry, ruleset.Name, err)
		}
	}

	return nil
}

func runUpdatePromptset(cmd *cobra.Command, args []string) error {
	// TODO: Implement promptset update when service interface is updated
	return fmt.Errorf("promptset update not yet implemented - service interface needs to be updated first")
}
