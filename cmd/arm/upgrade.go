package main

import (
	"github.com/spf13/cobra"
)

var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade resources",
	Long:  "Upgrade rulesets and promptsets to their latest available versions, ignoring version constraints",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeAll()
	},
}

var upgradeRulesetCmd = &cobra.Command{
	Use:   "ruleset [REGISTRY_NAME/RULESET_NAME...]",
	Short: "Upgrade rulesets",
	Long:  "Upgrade one or more rulesets to their latest available versions, ignoring version constraints.",
	Run: func(cmd *cobra.Command, args []string) {
		upgradeRulesets(args)
	},
}

var upgradePromptsetCmd = &cobra.Command{
	Use:   "promptset [REGISTRY_NAME/PROMPTSET_NAME...]",
	Short: "Upgrade promptsets",
	Long:  "Upgrade one or more promptsets to their latest available versions, ignoring version constraints.",
	Run: func(cmd *cobra.Command, args []string) {
		upgradePromptsets(args)
	},
}

func init() {
	// Add subcommands
	upgradeCmd.AddCommand(upgradeRulesetCmd)
	upgradeCmd.AddCommand(upgradePromptsetCmd)
}

func upgradeAll() {
	if err := armService.UpgradeAll(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func upgradeRulesets(names []string) {
	if len(names) == 0 {
		// Upgrade all rulesets
		if err := armService.UpdateAllRulesets(ctx); err != nil {
			// TODO: Handle error properly
			return
		}
	} else {
		// Upgrade specific rulesets
		for _, name := range names {
			registry, ruleset := parseRegistryPackage(name)
			if err := armService.UpgradeRuleset(ctx, registry, ruleset); err != nil {
				// TODO: Handle error properly
				return
			}
		}
	}
}

func upgradePromptsets(names []string) {
	if len(names) == 0 {
		// Upgrade all promptsets
		if err := armService.UpdateAllPromptsets(ctx); err != nil {
			// TODO: Handle error properly
			return
		}
	} else {
		// Upgrade specific promptsets
		for _, name := range names {
			registry, promptset := parseRegistryPackage(name)
			if err := armService.UpgradePromptset(ctx, registry, promptset); err != nil {
				// TODO: Handle error properly
				return
			}
		}
	}
}
