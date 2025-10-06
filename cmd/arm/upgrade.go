package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newUpgradeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade resources",
		Long:  "Upgrade all configured resources to latest versions, or use subcommands for specific resource types.",
		RunE:  runUpgradeAll,
	}

	// Add subcommands
	cmd.AddCommand(newUpgradeRulesetCmd())
	cmd.AddCommand(newUpgradePromptsetCmd())

	return cmd
}

func newUpgradeRulesetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ruleset [registry/ruleset...]",
		Short: "Upgrade rulesets",
		Long:  "Upgrade specific rulesets to latest versions or all rulesets if none specified.",
		RunE:  runUpgradeRuleset,
	}
}

func newUpgradePromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset [registry/promptset...]",
		Short: "Upgrade promptsets",
		Long:  "Upgrade specific promptsets to latest versions or all promptsets if none specified.",
		RunE:  runUpgradePromptset,
	}
}

func runUpgradeAll(cmd *cobra.Command, args []string) error {
	// TODO: Implement unified upgrade when service interface is updated
	return fmt.Errorf("unified upgrade not yet implemented - service interface needs to be updated first")
}

func runUpgradeRuleset(cmd *cobra.Command, args []string) error {
	// If no arguments, upgrade all rulesets
	if len(args) == 0 {
		// TODO: Implement upgrade all rulesets when service interface is updated
		return fmt.Errorf("upgrade all rulesets not yet implemented - service interface needs to be updated first")
	}

	// Parse arguments
	_, err := ParsePackageArgs(args)
	if err != nil {
		return fmt.Errorf("failed to parse arguments: %w", err)
	}

	// TODO: Implement individual ruleset upgrade when service interface is updated
	return fmt.Errorf("individual ruleset upgrade not yet implemented - service interface needs to be updated first")
}

func runUpgradePromptset(cmd *cobra.Command, args []string) error {
	// TODO: Implement promptset upgrade when service interface is updated
	return fmt.Errorf("promptset upgrade not yet implemented - service interface needs to be updated first")
}
