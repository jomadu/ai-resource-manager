package main

import (
	"github.com/spf13/cobra"
)

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Check for outdated resources",
	Long:  "Check for outdated rulesets and promptsets across configured registries",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedAll(cmd)
	},
}

var outdatedRulesetCmd = &cobra.Command{
	Use:   "ruleset [--output <table|json|list>]",
	Short: "Check for outdated rulesets",
	Long:  "Check for outdated rulesets across all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedRulesets(cmd)
	},
}

var outdatedPromptsetCmd = &cobra.Command{
	Use:   "promptset [--output <table|json|list>]",
	Short: "Check for outdated promptsets",
	Long:  "Check for outdated promptsets across all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedPromptsets(cmd)
	},
}

func init() {
	// Add subcommands
	outdatedCmd.AddCommand(outdatedRulesetCmd)
	outdatedCmd.AddCommand(outdatedPromptsetCmd)

	// Add output format flags
	outdatedRulesetCmd.Flags().String("output", "table", "Output format (table, json, list)")
	outdatedPromptsetCmd.Flags().String("output", "table", "Output format (table, json, list)")

	// Add output format flag to main command
	outdatedCmd.Flags().String("output", "table", "Output format (table, json, list)")
}

func checkOutdatedAll(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	if err := armService.ShowAllOutdated(ctx, output, false); err != nil {
		// TODO: Handle error properly
		return
	}
}

func checkOutdatedRulesets(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	if err := armService.ShowRulesetOutdated(ctx, output, false); err != nil {
		// TODO: Handle error properly
		return
	}
}

func checkOutdatedPromptsets(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	if err := armService.ShowPromptsetOutdated(ctx, output, false); err != nil {
		// TODO: Handle error properly
		return
	}
}
