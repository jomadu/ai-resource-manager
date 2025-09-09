package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed rulesets",
		RunE:  runList,
	}
}

func runList(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	installed, err := armService.ListInstalledRulesets(ctx)
	if err != nil {
		return fmt.Errorf("failed to list installed rulesets: %w", err)
	}

	FormatInstalledRulesets(installed)
	return nil
}
