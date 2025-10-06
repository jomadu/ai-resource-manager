package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed resources",
		Long:  "List all installed resources, or use subcommands for specific resource types.",
		RunE:  runListAll,
	}

	// Add subcommands
	cmd.AddCommand(newListRulesetCmd())
	cmd.AddCommand(newListPromptsetCmd())

	return cmd
}

func newListRulesetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset",
		Short: "List installed rulesets",
		RunE:  runListRuleset,
	}

	cmd.Flags().Bool("sort-priority", false, "Sort by priority (highest first) instead of alphanumeric")

	return cmd
}

func newListPromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset",
		Short: "List installed promptsets",
		RunE:  runListPromptset,
	}
}

func runListAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: Implement unified list when service interface is updated
	return armService.ShowList(ctx, false) // Temporary fallback to rulesets only
}

func runListRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sortPriority, _ := cmd.Flags().GetBool("sort-priority")

	return armService.ShowList(ctx, sortPriority)
}

func runListPromptset(cmd *cobra.Command, args []string) error {
	// TODO: Implement promptset list when service interface is updated
	return fmt.Errorf("promptset list not yet implemented - service interface needs to be updated first")
}
