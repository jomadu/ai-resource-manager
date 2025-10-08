package main

import (
	"github.com/spf13/cobra"
)

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Check for outdated packages, rulesets, and promptsets",
	Long:  "Check for outdated packages, rulesets, and promptsets across all configured registries",
}

var outdatedPackageCmd = &cobra.Command{
	Use:   "package [--output <table|json|list>]",
	Short: "Check for outdated packages",
	Long:  "Check for outdated packages across all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedPackages(cmd)
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
	outdatedCmd.AddCommand(outdatedPackageCmd)
	outdatedCmd.AddCommand(outdatedRulesetCmd)
	outdatedCmd.AddCommand(outdatedPromptsetCmd)

	// Add output format flags
	outdatedPackageCmd.Flags().String("output", "table", "Output format (table, json, list)")
	outdatedRulesetCmd.Flags().String("output", "table", "Output format (table, json, list)")
	outdatedPromptsetCmd.Flags().String("output", "table", "Output format (table, json, list)")
}

func checkOutdatedPackages(cmd *cobra.Command) {
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
