package main

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update packages, rulesets, and promptsets",
	Long:  "Update packages, rulesets, and promptsets to their latest available versions",
}

var updatePackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Update all packages",
	Long:  "Update all installed packages to their latest available versions.",
	Run: func(cmd *cobra.Command, args []string) {
		updatePackages()
	},
}

var updateRulesetCmd = &cobra.Command{
	Use:   "ruleset [REGISTRY_NAME/RULESET_NAME...]",
	Short: "Update rulesets",
	Long:  "Update one or more rulesets to their latest available versions.",
	Run: func(cmd *cobra.Command, args []string) {
		updateRulesets(args)
	},
}

var updatePromptsetCmd = &cobra.Command{
	Use:   "promptset [REGISTRY_NAME/PROMPTSET_NAME...]",
	Short: "Update promptsets",
	Long:  "Update one or more promptsets to their latest available versions.",
	Run: func(cmd *cobra.Command, args []string) {
		updatePromptsets(args)
	},
}

func init() {
	// Add subcommands
	updateCmd.AddCommand(updatePackageCmd)
	updateCmd.AddCommand(updateRulesetCmd)
	updateCmd.AddCommand(updatePromptsetCmd)
}

func updatePackages() {
	if err := armService.UpdateAll(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func updateRulesets(names []string) {
	if len(names) == 0 {
		// Update all rulesets
		if err := armService.UpdateAllRulesets(ctx); err != nil {
			// TODO: Handle error properly
			return
		}
	} else {
		// Update specific rulesets
		for _, name := range names {
			registry, ruleset := parseRegistryPackage(name)
			if err := armService.UpdateRuleset(ctx, registry, ruleset); err != nil {
				// TODO: Handle error properly
				return
			}
		}
	}
}

func updatePromptsets(names []string) {
	if len(names) == 0 {
		// Update all promptsets
		if err := armService.UpdateAllPromptsets(ctx); err != nil {
			// TODO: Handle error properly
			return
		}
	} else {
		// Update specific promptsets
		for _, name := range names {
			registry, promptset := parseRegistryPackage(name)
			if err := armService.UpdatePromptset(ctx, registry, promptset); err != nil {
				// TODO: Handle error properly
				return
			}
		}
	}
}
