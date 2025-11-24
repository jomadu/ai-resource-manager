package main

import (
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  "List registries, sinks, rulesets, and promptsets",
	Run: func(cmd *cobra.Command, args []string) {
		listAll()
	},
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

func init() {
	// Add subcommands
	listCmd.AddCommand(listRegistryCmd)
	listCmd.AddCommand(listSinkCmd)
	listCmd.AddCommand(listRulesetCmd)
	listCmd.AddCommand(listPromptsetCmd)
}

func listRegistries() {
	if err := armService.ShowRegistryList(ctx); err != nil {
		handleCommandError(err)
	}
}

func listSinks() {
	if err := armService.ShowSinkList(ctx); err != nil {
		handleCommandError(err)
	}
}

func listRulesets() {
	if err := armService.ShowRulesetList(ctx, false); err != nil {
		handleCommandError(err)
	}
}

func listPromptsets() {
	if err := armService.ShowPromptsetList(ctx); err != nil {
		handleCommandError(err)
	}
}

func listAll() {
	if err := armService.ShowAllList(ctx, false); err != nil {
		handleCommandError(err)
	}
}
