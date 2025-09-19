package main

import (
	"context"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed rulesets",
		RunE:  runList,
	}

	cmd.Flags().Bool("sort-priority", false, "Sort by priority (highest first) instead of alphanumeric")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	sortPriority, _ := cmd.Flags().GetBool("sort-priority")

	return armService.ShowList(ctx, sortPriority)
}
