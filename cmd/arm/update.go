package main

import (
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update resources",
	Long:  "Update rulesets and promptsets to their latest available versions",
	Run: func(cmd *cobra.Command, args []string) {
		updateAll()
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
	updateCmd.AddCommand(updateRulesetCmd)
	updateCmd.AddCommand(updatePromptsetCmd)
}

func updateAll() {
	if err := armService.UpdateAll(ctx); err != nil {
		handleCommandError(err)
	}
}

func updateRulesets(names []string) {
	if len(names) == 0 {
		// Update all rulesets
		if err := armService.UpdateAllRulesets(ctx); err != nil {
			handleCommandError(err)
		}
	} else {
		// Update specific rulesets
		for _, name := range names {
			registry, err := parseRegistry(name)
			if err != nil {
				handleCommandError(err)
			}

			ruleset, err := parsePackage(name)
			if err != nil {
				handleCommandError(err)
			}

			if err := armService.UpdateRuleset(ctx, registry, ruleset); err != nil {
				handleCommandError(err)
			}
		}
	}
}

func updatePromptsets(names []string) {
	if len(names) == 0 {
		// Update all promptsets
		if err := armService.UpdateAllPromptsets(ctx); err != nil {
			handleCommandError(err)
		}
	} else {
		// Update specific promptsets
		for _, name := range names {
			registry, err := parseRegistry(name)
			if err != nil {
				handleCommandError(err)
			}

			promptset, err := parsePackage(name)
			if err != nil {
				handleCommandError(err)
			}

			if err := armService.UpdatePromptset(ctx, registry, promptset); err != nil {
				handleCommandError(err)
			}
		}
	}
}
