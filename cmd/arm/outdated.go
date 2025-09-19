package main

import (
	"context"

	"github.com/spf13/cobra"
)

func newOutdatedCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "outdated",
		Short: "Show outdated rulesets",
		RunE:  runOutdated,
	}

	cmd.Flags().StringP("output", "o", "table", "Output format (table or json)")
	cmd.Flags().Bool("no-spinner", false, "Disable spinner for machine-readable output")

	return cmd
}

func runOutdated(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	outputFormat, _ := cmd.Flags().GetString("output")
	noSpinner, _ := cmd.Flags().GetBool("no-spinner")
	return armService.ShowOutdated(ctx, outputFormat, noSpinner)
}
