package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newOutdatedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "Show outdated resources",
		Long:  "Show all outdated resources, or use subcommands for specific resource types.",
		RunE:  runOutdatedAll,
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table, json, or list)")
	cmd.Flags().Bool("no-spinner", false, "Disable spinner for machine-readable output")

	// Add subcommands
	cmd.AddCommand(newOutdatedRulesetCmd())
	cmd.AddCommand(newOutdatedPromptsetCmd())

	return cmd
}

func newOutdatedRulesetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset",
		Short: "Show outdated rulesets",
		RunE:  runOutdatedRuleset,
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table, json, or list)")
	cmd.Flags().Bool("no-spinner", false, "Disable spinner for machine-readable output")

	return cmd
}

func newOutdatedPromptsetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promptset",
		Short: "Show outdated promptsets",
		RunE:  runOutdatedPromptset,
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table, json, or list)")
	cmd.Flags().Bool("no-spinner", false, "Disable spinner for machine-readable output")

	return cmd
}

func runOutdatedAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	outputFormat, _ := cmd.Flags().GetString("output")
	noSpinner, _ := cmd.Flags().GetBool("no-spinner")
	// TODO: Implement unified outdated when service interface is updated
	return armService.ShowOutdated(ctx, outputFormat, noSpinner) // Temporary fallback to rulesets only
}

func runOutdatedRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	outputFormat, _ := cmd.Flags().GetString("output")
	noSpinner, _ := cmd.Flags().GetBool("no-spinner")
	return armService.ShowOutdated(ctx, outputFormat, noSpinner)
}

func runOutdatedPromptset(cmd *cobra.Command, args []string) error {
	// TODO: Implement promptset outdated when service interface is updated
	return fmt.Errorf("promptset outdated not yet implemented - service interface needs to be updated first")
}
