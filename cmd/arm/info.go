package main

import (
	"context"
	"fmt"

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

	// If no arguments, show info for all installed rulesets
	if len(args) == 0 {
		infos, err := armService.GetAllRulesetInfo(ctx)
		if err != nil {
			return fmt.Errorf("failed to get ruleset information: %w", err)
		}

		for _, info := range infos {
			FormatRulesetInfo(info, false)
			fmt.Println()
		}
		return nil
	}

	// Parse arguments
	rulesets, err := ParseRulesetArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Show info for each specified ruleset
	for i, ruleset := range rulesets {
		info, err := armService.GetRulesetInfo(ctx, ruleset.Registry, ruleset.Name)
		if err != nil {
			return fmt.Errorf("failed to get info for %s/%s: %w", ruleset.Registry, ruleset.Name, err)
		}

		detailed := len(rulesets) == 1
		FormatRulesetInfo(info, detailed)
		if i < len(rulesets)-1 {
			fmt.Println()
		}
	}

	return nil
}
