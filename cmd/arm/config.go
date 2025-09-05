package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
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
	Long: `Manage registry configuration for ARM.

Registries are remote sources where rulesets are stored and versioned, similar to npm registries for JavaScript packages. ARM supports Git-based registries that can point to GitHub repositories, GitLab projects, or any Git remote containing rule collections.

Available commands:
  add     Add a new registry
  remove  Remove an existing registry

Examples:
  arm config registry add ai-rules https://github.com/user/rules-repo --type git
  arm config registry remove ai-rules`,
}

var configSinkCmd = &cobra.Command{
	Use:   "sink",
	Short: "Manage sink configuration",
	Long: `Manage sink configuration for ARM.

Sinks define where installed rules should be placed in your local filesystem and which AI tools should receive them. Each sink targets specific directories and can filter rulesets using include/exclude patterns.

Sinks support two layout modes:
- Hierarchical Layout (default): Preserves directory structure from rulesets
- Flat Layout: Places all files in a single directory with hash-prefixed names

Available commands:
  add     Add a new sink
  remove  Remove an existing sink

Examples:
  arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"
  arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
  arm config sink add github --directories .github/instructions --layout flat`,
}

var registryAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new registry",
	Long: `Add a registry to the manifest.

Registries are remote sources where rulesets are stored and versioned, similar to npm registries for JavaScript packages. When you configure a registry, you're creating a named connection to a repository that contains multiple rulesets with proper semantic versioning.

Arguments:
  name  Registry name (used to reference the registry)
  url   Registry URL (Git repository URL)

Flags:
  --type  Registry type (default: git)

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
		switch registryType {
		case "git":
			branches, _ := cmd.Flags().GetStringSlice("branches")
			if len(branches) == 0 {
				branches = []string{"main", "master"}
			}
			gitConfig := registry.GitRegistryConfig{
				RegistryConfig: registry.RegistryConfig{URL: url, Type: registryType},
				Branches:       branches,
			}
			return manifestManager.AddGitRegistry(context.Background(), name, gitConfig)
		default:
			return fmt.Errorf("registry type %s is not implemented", registryType)
		}
	},
}

var registryRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an existing registry",
	Long: `Remove a registry from the manifest.

This will remove the registry configuration but will not affect any already installed rulesets from that registry. To remove installed rulesets, use 'arm uninstall'.

Arguments:
  name  Registry name to remove

Examples:
  arm config registry remove ai-rules
  arm config registry remove company-rules`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		manifestManager := manifest.NewFileManager()
		return manifestManager.RemoveRegistry(context.Background(), name)
	},
}

var sinkAddCmd = &cobra.Command{
	Use:   "add <name>",
	Short: "Add a new sink",
	Long: `Add a sink to the configuration.

Sinks define where installed rules should be placed in your local filesystem and which AI tools should receive them. Each sink targets specific directories and can filter rulesets using include/exclude patterns.

Arguments:
  name  Sink name (used to reference the sink)

Flags:
  --directories   Target directories for rule installation (required)
  --include       Include patterns to filter rulesets
  --exclude       Exclude patterns to filter rulesets
  --layout        Layout mode: hierarchical (default) or flat

Layout Modes:
- Hierarchical: Preserves directory structure from rulesets
- Flat: Places all files in single directory with hash-prefixed names

Examples:
  arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"
  arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
  arm config sink add github --directories .github/instructions --include "ai-rules/*" --layout flat`,
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
	Short: "Remove an existing sink",
	Long: `Remove a sink from the configuration.

This will remove the sink configuration but will not delete any files that were previously installed to the sink's directories. To clean up installed files, manually delete them or reinstall rulesets after reconfiguring sinks.

Arguments:
  name  Sink name to remove

Examples:
  arm config sink remove q
  arm config sink remove cursor`,
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
	registryAddCmd.Flags().StringSlice("branches", nil, "Git branches to track (default: main,master)")
	sinkAddCmd.Flags().StringSlice("directories", nil, "Sink directories")
	sinkAddCmd.Flags().StringSlice("include", nil, "Sink include patterns")
	sinkAddCmd.Flags().StringSlice("exclude", nil, "Sink exclude patterns")
	sinkAddCmd.Flags().String("layout", "hierarchical", "Layout mode (hierarchical, flat)")
}
