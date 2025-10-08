package main

import (
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install packages, rulesets, and promptsets",
	Long:  "Install packages, rulesets, and promptsets to their assigned sinks",
}

var installPackageCmd = &cobra.Command{
	Use:   "package",
	Short: "Install all configured packages",
	Long:  "Install all configured packages (rulesets and promptsets) to their assigned sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		installPackages()
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
	installCmd.AddCommand(installPackageCmd)
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

func installPackages() {
	if err := armService.InstallAll(ctx); err != nil {
		// TODO: Handle error properly
		return
	}
}

func installRuleset(cmd *cobra.Command, packageName string, sinks []string) {
	priority, _ := cmd.Flags().GetInt("priority")
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	// Parse registry/ruleset from packageName
	registry, ruleset, version := parsePackageName(packageName)

	req := &arm.InstallRulesetRequest{
		Registry: registry,
		Ruleset:  ruleset,
		Version:  version,
		Priority: priority,
		Include:  include,
		Exclude:  exclude,
		Sinks:    sinks,
	}

	if err := armService.InstallRuleset(ctx, req); err != nil {
		// TODO: Handle error properly
		return
	}
}

func installPromptset(cmd *cobra.Command, packageName string, sinks []string) {
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	// Parse registry/promptset from packageName
	registry, promptset, version := parsePackageName(packageName)

	req := &arm.InstallPromptsetRequest{
		Registry:  registry,
		Promptset: promptset,
		Version:   version,
		Include:   include,
		Exclude:   exclude,
		Sinks:     sinks,
	}

	if err := armService.InstallPromptset(ctx, req); err != nil {
		// TODO: Handle error properly
		return
	}
}

// parsePackageName parses a package name like "registry/package@version" or "registry/package"
func parsePackageName(packageName string) (registry, pkgName, version string) {
	parts := strings.Split(packageName, "/")
	if len(parts) != 2 {
		// TODO: Handle error
		return "", "", ""
	}

	registry = parts[0]
	packageWithVersion := parts[1]

	// Check for version
	if strings.Contains(packageWithVersion, "@") {
		versionParts := strings.Split(packageWithVersion, "@")
		pkgName = versionParts[0]
		version = versionParts[1]
	} else {
		pkgName = packageWithVersion
		version = ""
	}

	return registry, pkgName, version
}
