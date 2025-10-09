package main

import (
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed information",
	Long:  "Show detailed information about registries, sinks, rulesets, and promptsets",
	Run: func(cmd *cobra.Command, args []string) {
		infoAll()
	},
}

var infoRegistryCmd = &cobra.Command{
	Use:   "registry [NAME]...",
	Short: "Show registry information",
	Long:  "Display detailed information about one or more registries.",
	Run: func(cmd *cobra.Command, args []string) {
		infoRegistries(args)
	},
}

var infoSinkCmd = &cobra.Command{
	Use:   "sink [NAME]...",
	Short: "Show sink information",
	Long:  "Display detailed information about one or more sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		infoSinks(args)
	},
}

var infoRulesetCmd = &cobra.Command{
	Use:   "ruleset [REGISTRY_NAME/RULESET_NAME...]",
	Short: "Show ruleset information",
	Long:  "Display detailed information about one or more rulesets.",
	Run: func(cmd *cobra.Command, args []string) {
		infoRulesets(args)
	},
}

var infoPromptsetCmd = &cobra.Command{
	Use:   "promptset [REGISTRY_NAME/PROMPTSET_NAME...]",
	Short: "Show promptset information",
	Long:  "Display detailed information about one or more promptsets.",
	Run: func(cmd *cobra.Command, args []string) {
		infoPromptsets(args)
	},
}

func init() {
	// Add subcommands
	infoCmd.AddCommand(infoRegistryCmd)
	infoCmd.AddCommand(infoSinkCmd)
	infoCmd.AddCommand(infoRulesetCmd)
	infoCmd.AddCommand(infoPromptsetCmd)
}

func infoRegistries(names []string) {
	if err := armService.ShowRegistryInfo(ctx, names); err != nil {
		// TODO: Handle error properly
		return
	}
}

func infoSinks(names []string) {
	if err := armService.ShowSinkInfo(ctx, names); err != nil {
		// TODO: Handle error properly
		return
	}
}

func infoRulesets(names []string) {
	if err := armService.ShowRulesetInfo(ctx, names); err != nil {
		// TODO: Handle error properly
		return
	}
}

func infoPromptsets(names []string) {
	if err := armService.ShowPromptsetInfo(ctx, names); err != nil {
		// TODO: Handle error properly
		return
	}
}

func infoAll() {
	if err := armService.ShowAllInfo(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}
