package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed rulesets",
		RunE:  runList,
	}

	cmd.Flags().Bool("show-priority", false, "Display ruleset installation priorities")
	cmd.Flags().Bool("sort-priority", false, "Sort by priority (highest first) instead of alphanumeric")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	showPriority, _ := cmd.Flags().GetBool("show-priority")
	sortPriority, _ := cmd.Flags().GetBool("sort-priority")

	installed, err := armService.ListInstalledRulesets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list installed rulesets: %w", err)
	}

	FormatInstalledRulesets(installed, showPriority, sortPriority)
	return nil
}
