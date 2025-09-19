package main

import (
	"context"

	"github.com/spf13/cobra"
)

func newInfoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "info [ruleset...]",
		Short: "Show ruleset information",
		Long:  "Show information about specific rulesets or all installed rulesets.",
		RunE:  runInfo,
	}
}

func runInfo(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Convert args to ruleset strings for service
	var rulesetStrings []string
	rulesetStrings = append(rulesetStrings, args...)

	return armService.ShowRulesetInfo(ctx, rulesetStrings)
}
