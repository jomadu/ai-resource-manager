package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newOutdatedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "Show outdated rulesets",
		RunE:  runOutdated,
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")

	return cmd
}

func runOutdated(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	outdated, err := armService.GetOutdatedRulesets(ctx)
	if err != nil {
		return fmt.Errorf("failed to check for outdated rulesets: %w", err)
	}

	outputFormat, _ := cmd.Flags().GetString("output")
	return FormatOutdatedRulesets(outdated, outputFormat)
}
