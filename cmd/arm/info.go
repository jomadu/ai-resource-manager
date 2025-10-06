package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info",
		Short: "Show resource information",
		Long:  "Show information about all installed resources, or use subcommands for specific resource types.",
		RunE:  runInfoAll,
	}

	// Add subcommands
	cmd.AddCommand(newInfoRulesetCmd())
	cmd.AddCommand(newInfoPromptsetCmd())

	return cmd
}

func newInfoRulesetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ruleset [registry/ruleset...]",
		Short: "Show ruleset information",
		Long:  "Show information about specific rulesets or all installed rulesets.",
		RunE:  runInfoRuleset,
	}
}

func newInfoPromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset [registry/promptset...]",
		Short: "Show promptset information",
		Long:  "Show information about specific promptsets or all installed promptsets.",
		RunE:  runInfoPromptset,
	}
}

func runInfoAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// TODO: Implement unified info when service interface is updated
	return armService.ShowRulesetInfo(ctx, []string{}) // Temporary fallback to all rulesets
}

func runInfoRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Convert args to ruleset strings for service
	var rulesetStrings []string
	rulesetStrings = append(rulesetStrings, args...)

	return armService.ShowRulesetInfo(ctx, rulesetStrings)
}

func runInfoPromptset(cmd *cobra.Command, args []string) error {
	// TODO: Implement promptset info when service interface is updated
	return fmt.Errorf("promptset info not yet implemented - service interface needs to be updated first")
}
