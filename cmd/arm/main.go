package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var armService arm.Service

func main() {
	armService = arm.NewArmService()

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "arm",
	Short: "AI Rules Manager - Manage AI rule rulesets",
	Long:  "ARM helps you install, manage, and organize AI rule rulesets from various registries.",
}

func init() {
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(outdatedCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(cacheCmd)
	rootCmd.AddCommand(versionCmd)
}

var installCmd = &cobra.Command{
	Use:   "install [ruleset...]",
	Short: "Install rulesets",
	Long:  "Install rulesets from a registry. If no ruleset is specified, installs from manifest.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if len(args) == 0 {
			return armService.Install(ctx)
		}

		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")

		for _, arg := range args {
			registry, ruleset, version := parseRulesetArg(arg)
			if err := armService.InstallRuleset(ctx, registry, ruleset, version, include, exclude); err != nil {
				return err
			}
		}
		return nil
	},
}

var uninstallCmd = &cobra.Command{
	Use:   "uninstall <ruleset>",
	Short: "Uninstall a ruleset",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		registry, ruleset, _ := parseRulesetArg(args[0])
		return armService.Uninstall(ctx, registry, ruleset)
	},
}

var updateCmd = &cobra.Command{
	Use:   "update [ruleset...]",
	Short: "Update rulesets",
	Long:  "Update rulesets. If no ruleset is specified, updates all rulesets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if len(args) == 0 {
			return armService.Update(ctx)
		}

		for _, arg := range args {
			registry, ruleset, _ := parseRulesetArg(arg)
			if err := armService.UpdateRuleset(ctx, registry, ruleset); err != nil {
				return err
			}
		}
		return nil
	},
}

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Show outdated rulesets",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		outdated, err := armService.Outdated(ctx)
		if err != nil {
			return err
		}

		if len(outdated) == 0 {
			fmt.Println("All rulesets are up to date!")
			return nil
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Header("Registry", "Ruleset", "Current", "Wanted", "Latest")

		for _, r := range outdated {
			if err := table.Append(r.Registry, r.Name, r.Current, r.Wanted, r.Latest); err != nil {
				return err
			}
		}

		return table.Render()
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed rulesets",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		installed, err := armService.List(ctx)
		if err != nil {
			return err
		}

		for _, ruleset := range installed {
			fmt.Printf("%s/%s@%s\n", ruleset.Registry, ruleset.Name, ruleset.Version)
		}

		return nil
	},
}

var infoCmd = &cobra.Command{
	Use:   "info [ruleset...]",
	Short: "Show ruleset information",
	Long:  "Show information about specific rulesets or all installed rulesets.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		if len(args) == 0 {
			infos, err := armService.InfoAll(ctx)
			if err != nil {
				return err
			}

			for _, info := range infos {
				printRulesetInfo(info, false)
				fmt.Println()
			}
			return nil
		}

		for i, arg := range args {
			registry, ruleset, _ := parseRulesetArg(arg)
			info, err := armService.Info(ctx, registry, ruleset)
			if err != nil {
				return err
			}

			detailed := len(args) == 1
			printRulesetInfo(info, detailed)
			if i < len(args)-1 {
				fmt.Println()
			}
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  "Manage ARM configuration including registries and sinks.",
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version",
	RunE: func(cmd *cobra.Command, args []string) error {
		version := armService.Version()
		fmt.Printf("arm %s\n", version.Version)
		if version.Commit != "" {
			fmt.Printf("commit: %s\n", version.Commit)
		}
		if version.Arch != "" {
			fmt.Printf("arch: %s\n", version.Arch)
		}
		return nil
	},
}

func init() {
	installCmd.Flags().StringSlice("include", nil, "Include patterns")
	installCmd.Flags().StringSlice("exclude", nil, "Exclude patterns")
}

// parseRulesetArg parses registry/ruleset[@version] format
func parseRulesetArg(arg string) (registry, ruleset, version string) {
	// Simple parsing - in real implementation would be more robust
	parts := splitOnFirst(arg, "/")
	if len(parts) != 2 {
		return "", arg, ""
	}

	registry = parts[0]
	rulesetAndVersion := parts[1]

	versionParts := splitOnFirst(rulesetAndVersion, "@")
	ruleset = versionParts[0]
	if len(versionParts) > 1 {
		version = versionParts[1]
	}

	return registry, ruleset, version
}

func splitOnFirst(s, sep string) []string {
	if idx := findFirst(s, sep); idx != -1 {
		return []string{s[:idx], s[idx+len(sep):]}
	}
	return []string{s}
}

func findFirst(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func printRulesetInfo(info *arm.RulesetInfo, detailed bool) {
	if detailed {
		fmt.Printf("Ruleset: %s/%s\n", info.Registry, info.Name)
		fmt.Printf("Registry: %s (%s)\n", info.RegistryURL, info.RegistryType)
		fmt.Println("include:")
		for _, pattern := range info.Include {
			fmt.Printf("  - %s\n", pattern)
		}
		fmt.Println("Installed:")
		for _, path := range info.InstalledPaths {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Println("Sinks:")
		for _, sink := range info.Sinks {
			fmt.Printf("  - %s\n", sink)
		}
		fmt.Printf("Constraint: %s\n", info.Constraint)
		fmt.Printf("Resolved: %s\n", info.Resolved)
	} else {
		fmt.Printf("%s/%s\n", info.Registry, info.Name)
		fmt.Printf("  Registry: %s (%s)\n", info.RegistryURL, info.RegistryType)
		fmt.Println("  include:")
		for _, pattern := range info.Include {
			fmt.Printf("    - %s\n", pattern)
		}
		fmt.Println("  Installed:")
		for _, path := range info.InstalledPaths {
			fmt.Printf("    - %s\n", path)
		}
		fmt.Println("  Sinks:")
		for _, sink := range info.Sinks {
			fmt.Printf("    - %s\n", sink)
		}
		fmt.Printf("  Constraint: %s | Resolved: %s\n",
			info.Constraint, info.Resolved)
	}
}
