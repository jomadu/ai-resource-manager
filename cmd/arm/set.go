package main

import (
	"github.com/spf13/cobra"
)

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration",
	Long:  "Set configuration values for registries, sinks, rulesets, and promptsets",
}

var setRegistryCmd = &cobra.Command{
	Use:   "registry NAME KEY VALUE",
	Short: "Set registry configuration",
	Long:  "Set configuration values for a specific registry.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		setRegistry(cmd, args[0], args[1], args[2])
	},
}

var setSinkCmd = &cobra.Command{
	Use:   "sink NAME KEY VALUE",
	Short: "Set sink configuration",
	Long:  "Set configuration values for a specific sink.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		setSink(cmd, args[0], args[1], args[2])
	},
}

var setRulesetCmd = &cobra.Command{
	Use:   "ruleset REGISTRY_NAME/RULESET_NAME KEY VALUE",
	Short: "Set ruleset configuration",
	Long:  "Set configuration values for a specific ruleset.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		setRuleset(cmd, args[0], args[1], args[2])
	},
}

var setPromptsetCmd = &cobra.Command{
	Use:   "promptset REGISTRY_NAME/PROMPTSET_NAME KEY VALUE",
	Short: "Set promptset configuration",
	Long:  "Set configuration values for a specific promptset.",
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		setPromptset(cmd, args[0], args[1], args[2])
	},
}

func init() {
	// Add subcommands
	setCmd.AddCommand(setRegistryCmd)
	setCmd.AddCommand(setSinkCmd)
	setCmd.AddCommand(setRulesetCmd)
	setCmd.AddCommand(setPromptsetCmd)
}

func setRegistry(cmd *cobra.Command, name, key, value string) {
	if err := armService.SetRegistryConfig(ctx, name, key, value); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func setSink(cmd *cobra.Command, name, key, value string) {
	if err := armService.SetSinkConfig(ctx, name, key, value); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func setRuleset(cmd *cobra.Command, packageName, key, value string) {
	// Parse registry/ruleset from packageName
	registry, err := parseRegistry(packageName)
	if err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
	
	ruleset, err := parsePackage(packageName)
	if err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}

	if err := armService.SetRulesetConfig(ctx, registry, ruleset, key, value); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func setPromptset(cmd *cobra.Command, packageName, key, value string) {
	// Parse registry/promptset from packageName
	registry, err := parseRegistry(packageName)
	if err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
	
	promptset, err := parsePackage(packageName)
	if err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}

	if err := armService.SetPromptsetConfig(ctx, registry, promptset, key, value); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}
