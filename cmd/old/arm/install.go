package main

import (
	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install new or configured resources",
	Long:  "Install new or already configured rulesets and promptsets to their assigned sinks",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		installAll()
	},
}

var installRulesetCmd = &cobra.Command{
	Use:   "ruleset [--priority PRIORITY] [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/RULESET_NAME[@VERSION] SINK_NAME...",
	Short: "Install a ruleset",
	Long:  "Install a specific ruleset from a registry to one or more sinks.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		installRuleset(cmd, args[0], args[1:])
	},
}

var installPromptsetCmd = &cobra.Command{
	Use:   "promptset [--include GLOB...] [--exclude GLOB...] REGISTRY_NAME/PROMPTSET[@VERSION] SINK_NAME...",
	Short: "Install a promptset",
	Long:  "Install a specific promptset from a registry to one or more sinks.",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		installPromptset(cmd, args[0], args[1:])
	},
}

func init() {
	// Add subcommands
	installCmd.AddCommand(installRulesetCmd)
	installCmd.AddCommand(installPromptsetCmd)

	// Add ruleset flags
	installRulesetCmd.Flags().Int("priority", 100, "Ruleset priority")
	installRulesetCmd.Flags().StringSlice("include", []string{"**/*.yml", "**/*.yaml"}, "Include patterns")
	installRulesetCmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")

	// Add promptset flags
	installPromptsetCmd.Flags().StringSlice("include", []string{"**/*.yml", "**/*.yaml"}, "Include patterns")
	installPromptsetCmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")
}

func installAll() {
	if err := armService.InstallAll(ctx); err != nil {
		handleCommandError(err)
	}
}

func installRuleset(cmd *cobra.Command, packageName string, sinks []string) {
	priority, _ := cmd.Flags().GetInt("priority")
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	// Parse registry/ruleset from packageName
	registry, err := parseRegistry(packageName)
	if err != nil {
		handleCommandError(err)
	}

	ruleset, err := parsePackage(packageName)
	if err != nil {
		handleCommandError(err)
	}

	version, err := parseVersion(packageName)
	if err != nil {
		handleCommandError(err)
	}

	// Use constructor and fluent API
	req := arm.NewInstallRulesetRequest(registry, ruleset, version, sinks).
		WithPriority(priority).
		WithInclude(include).
		WithExclude(exclude)

	if err := armService.InstallRuleset(ctx, req); err != nil {
		handleCommandError(err)
	}
}

func installPromptset(cmd *cobra.Command, packageName string, sinks []string) {
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	// Parse registry/promptset from packageName
	registry, err := parseRegistry(packageName)
	if err != nil {
		handleCommandError(err)
	}

	promptset, err := parsePackage(packageName)
	if err != nil {
		handleCommandError(err)
	}

	version, err := parseVersion(packageName)
	if err != nil {
		handleCommandError(err)
	}

	req := &arm.InstallPromptsetRequest{
		Registry:  registry,
		Promptset: promptset,
		Version:   version,
		Include:   include,
		Exclude:   exclude,
		Sinks:     sinks,
	}

	if err := armService.InstallPromptset(ctx, req); err != nil {
		handleCommandError(err)
	}
}
