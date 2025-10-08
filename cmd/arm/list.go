package main

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  "List registries, sinks, rulesets, promptsets, and packages",
}

var listRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "List all registries",
	Long:  "List all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		listRegistries()
	},
}

var listSinkCmd = &cobra.Command{
	Use:   "sink",
	Short: "List all sinks",
	Long:  "List all configured sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		listSinks()
	},
}

var listRulesetCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "List all rulesets",
	Long:  "List all installed rulesets.",
	Run: func(cmd *cobra.Command, args []string) {
		listRulesets()
	},
}

var listPromptsetCmd = &cobra.Command{
	Use:   "promptset",
	Short: "List all promptsets",
	Long:  "List all installed promptsets.",
	Run: func(cmd *cobra.Command, args []string) {
		listPromptsets()
	},
}

var listPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "List all packages",
	Long:  "List all installed packages across all sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		listPackages()
	},
}

func init() {
	// Add subcommands
	listCmd.AddCommand(listRegistryCmd)
	listCmd.AddCommand(listSinkCmd)
	listCmd.AddCommand(listRulesetCmd)
	listCmd.AddCommand(listPromptsetCmd)
	listCmd.AddCommand(listPackageCmd)
}

func listRegistries() {
	if err := armService.ListRegistries(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func listSinks() {
	if err := armService.ListSinks(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func listRulesets() {
	if err := armService.ShowRulesetList(ctx, false); err != nil {
		// TODO: Handle error properly
		return
	}
}

func listPromptsets() {
	if err := armService.ShowPromptsetList(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func listPackages() {
	if err := armService.ShowAllList(ctx, false); err != nil {
		// TODO: Handle error properly
		return
	}
}
