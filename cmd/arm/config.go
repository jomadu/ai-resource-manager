package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configRegistryCmd)
	configCmd.AddCommand(configSinkCmd)
	configCmd.AddCommand(configListCmd)
}

var configRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage registry configuration",
}

var configSinkCmd = &cobra.Command{
	Use:   "sink",
	Short: "Manage sink configuration",
}

var registryAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add registry configuration",
	Long: `Add a registry to the manifest.

Registries are remote sources where rulesets are stored and versioned.

Examples:
  arm config registry add ai-rules https://github.com/user/rules-repo --type git
  arm config registry add company-rules https://gitlab.com/company/rules --type git`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]
		registryType, _ := cmd.Flags().GetString("type")
		if registryType == "" {
			registryType = "git"
		}

		manifestManager := manifest.NewFileManager()
		return manifestManager.AddRegistry(context.Background(), name, url, registryType)
	},
}

var registryRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove registry configuration",
	Long: `Remove a registry from the manifest.

Examples:
  arm config registry remove ai-rules`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		manifestManager := manifest.NewFileManager()
		return manifestManager.RemoveRegistry(context.Background(), name)
	},
}

var sinkAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add sink configuration",
	Long: `Add a sink to the configuration.

Sinks define where installed rules should be placed in your filesystem.

Examples:
  arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"
  arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
  arm config sink add copilot --directories .copilot/rules --include "ai-rules/*" --layout flat`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		directories, _ := cmd.Flags().GetStringSlice("directories")
		include, _ := cmd.Flags().GetStringSlice("include")
		exclude, _ := cmd.Flags().GetStringSlice("exclude")
		layout, _ := cmd.Flags().GetString("layout")

		configManager := config.NewFileManager()
		return configManager.AddSinkWithLayout(context.Background(), name, directories, include, exclude, layout)
	},
}

var sinkRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove sink configuration",
	Long: `Remove a sink from the configuration.

Examples:
  arm config sink remove q`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		configManager := config.NewFileManager()
		return configManager.RemoveSink(context.Background(), name)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestManager := manifest.NewFileManager()
		registries, err := manifestManager.GetRegistries(context.Background())
		if err == nil {
			fmt.Println("Registries:")
			for name, reg := range registries {
				fmt.Printf("  %s: %s (%s)\n", name, reg.URL, reg.Type)
			}
		}

		configManager := config.NewFileManager()
		sinks, err := configManager.GetSinks(context.Background())
		if err == nil {
			fmt.Println("Sinks:")
			for name, sink := range sinks {
				fmt.Printf("  %s:\n", name)
				fmt.Printf("    directories: %v\n", sink.Directories)
				fmt.Printf("    include: %v\n", sink.Include)
				fmt.Printf("    exclude: %v\n", sink.Exclude)
				layout := sink.Layout
				if layout == "" {
					layout = "hierarchical"
				}
				fmt.Printf("    layout: %s\n", layout)
			}
		} else {
			fmt.Println("Sinks: (none configured)")
		}
		return nil
	},
}

func init() {
	configRegistryCmd.AddCommand(registryAddCmd)
	configRegistryCmd.AddCommand(registryRemoveCmd)
	configSinkCmd.AddCommand(sinkAddCmd)
	configSinkCmd.AddCommand(sinkRemoveCmd)

	registryAddCmd.Flags().String("type", "git", "Registry type (git, http)")
	sinkAddCmd.Flags().StringSlice("directories", nil, "Sink directories")
	sinkAddCmd.Flags().StringSlice("include", nil, "Sink include patterns")
	sinkAddCmd.Flags().StringSlice("exclude", nil, "Sink exclude patterns")
	sinkAddCmd.Flags().String("layout", "hierarchical", "Layout mode (hierarchical, flat)")
}
