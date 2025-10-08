package main

import (
	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall packages, rulesets, and promptsets",
	Long:  "Uninstall packages, rulesets, and promptsets from their assigned sinks",
}

var uninstallPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Uninstall all packages",
	Long:  "Uninstall all configured packages from their assigned sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		uninstallPackages()
	},
}

var uninstallRulesetCmd = &cobra.Command{
	Use:   "ruleset REGISTRY_NAME/RULESET_NAME",
	Short: "Uninstall a ruleset",
	Long:  "Uninstall a specific ruleset from all sinks.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uninstallRuleset(args[0])
	},
}

var uninstallPromptsetCmd = &cobra.Command{
	Use:   "promptset REGISTRY_NAME/PROMPTSET_NAME",
	Short: "Uninstall a promptset",
	Long:  "Uninstall a specific promptset from all sinks.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		uninstallPromptset(args[0])
	},
}

func init() {
	// Add subcommands
	uninstallCmd.AddCommand(uninstallPackageCmd)
	uninstallCmd.AddCommand(uninstallRulesetCmd)
	uninstallCmd.AddCommand(uninstallPromptsetCmd)
}

func uninstallPackages() {
	if err := armService.UninstallAll(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func uninstallRuleset(packageName string) {
	registry, ruleset := parseRegistryPackage(packageName)

	if err := armService.UninstallRuleset(ctx, registry, ruleset); err != nil {
		// TODO: Handle error properly
		return
	}
}

func uninstallPromptset(packageName string) {
	registry, promptset := parseRegistryPackage(packageName)

	if err := armService.UninstallPromptset(ctx, registry, promptset); err != nil {
		// TODO: Handle error properly
		return
	}
}
