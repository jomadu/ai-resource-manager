package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update [ruleset...]",
		Short: "Update rulesets",
		Long:  "Update rulesets. If no ruleset is specified, updates all rulesets.",
		RunE:  runUpdate,
	}
}

func runUpdate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// If no arguments, update all
	if len(args) == 0 {
		return armService.UpdateAllRulesets(ctx)
	}

	// Parse arguments
	rulesets, err := ParseRulesetArgs(args)
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
