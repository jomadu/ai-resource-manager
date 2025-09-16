package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Manage ARM configuration including registries and sinks.",
	}

	cmd.AddCommand(configRegistryCmd)
	cmd.AddCommand(configSinkCmd)
	cmd.AddCommand(configRulesetCmd)
	cmd.AddCommand(configListCmd)

	return cmd
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

Sinks define where installed rules should be placed in your local filesystem. Each sink targets a single directory for rule installation. Rulesets are explicitly assigned to sinks during installation.

Sinks support two layout modes:
- Hierarchical Layout (default): Preserves directory structure from rulesets
- Flat Layout: Places all files in a single directory with hash-prefixed names

Available commands:
  add     Add a new sink
  remove  Remove an existing sink
  update  Update sink configuration

Examples:
  arm config sink add cursor .cursor/rules --layout hierarchical
  arm config sink add q .amazonq/rules --layout hierarchical
  arm config sink add github .github/instructions --layout flat`,
}

var configRulesetCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "Manage ruleset configuration",
	Long: `Manage ruleset configuration for ARM.

Rulesets are collections of AI rules that can be configured with priorities, version constraints, and sink assignments.

Available commands:
  update  Update ruleset configuration (triggers reinstall)

Examples:
  arm config ruleset update ai-rules/ruleset priority 150
  arm config ruleset update ai-rules/ruleset version 1.1.0
  arm config ruleset update ai-rules/ruleset sinks cursor,q`,
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
		force, _ := cmd.Flags().GetBool("force")

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
			return manifestManager.AddGitRegistry(context.Background(), name, gitConfig, force)
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
	Use:   "add <name> <directory>",
	Short: "Add a new sink",
	Long: `Add a sink to the configuration.

Sinks define where installed rules should be placed in your local filesystem. Each sink targets a single directory for rule installation.

Arguments:
  name       Sink name (used to reference the sink)
  directory  Target directory for rule installation

Flags:
  --type        Sink type with defaults (cursor, copilot, amazonq) - REQUIRED unless --compile-to is specified
  --layout      Layout mode: hierarchical or flat (overrides type default)
  --compile-to  Target format for compilation (cursor, amazonq, markdown, copilot) - REQUIRED unless --type is specified

Type Defaults:
- cursor: hierarchical layout, cursor compile target
- copilot: flat layout, copilot compile target
- amazonq: hierarchical layout, amazonq compile target

Examples:
  arm config sink add cursor .cursor/rules --type cursor
  arm config sink add copilot .github/copilot --type copilot
  arm config sink add q .amazonq/rules --type amazonq
  arm config sink add custom .custom/rules --layout flat --compile-to markdown`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		directory := args[1]
		layout, _ := cmd.Flags().GetString("layout")
		compileToStr, _ := cmd.Flags().GetString("compile-to")
		typeStr, _ := cmd.Flags().GetString("type")
		force, _ := cmd.Flags().GetBool("force")

		return armService.AddSink(context.Background(), name, directory, typeStr, layout, compileToStr, force)
	},
}

var sinkRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an existing sink",
	Long: `Remove a sink from the configuration.

This will remove the sink configuration and automatically clean it from all ruleset configurations. Files will be removed from the sink's directory.

Arguments:
  name  Sink name to remove

Examples:
  arm config sink remove q
  arm config sink remove cursor`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		manifestManager := manifest.NewFileManager()
		// Get sink before removal
		sink, err := manifestManager.GetSink(context.Background(), name)
		if err != nil {
			return err
		}
		// Remove from manifest and clean from all rulesets
		err = manifestManager.RemoveSink(context.Background(), name)
		if err != nil {
			return err
		}
		// Clean files from sink directory
		return armService.SyncRemovedSink(context.Background(), sink)
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		manifestManager := manifest.NewFileManager()
		registries, err := manifestManager.GetRawRegistries(context.Background())
		if err == nil {
			fmt.Println("Registries:")
			for name, reg := range registries {
				url, _ := reg["url"].(string)
				regType, _ := reg["type"].(string)
				fmt.Printf("  %s: %s (%s)\n", name, url, regType)
			}
		}

		sinks, err := manifestManager.GetSinks(context.Background())
		if err == nil {
			fmt.Println("Sinks:")
			for name, sink := range sinks {
				fmt.Printf("  %s:\n", name)
				fmt.Printf("    directory: %s\n", sink.Directory)
				layout := sink.Layout
				if layout == "" {
					layout = "hierarchical"
				}
				fmt.Printf("    layout: %s\n", layout)
				fmt.Printf("    compileTarget: %s\n", sink.CompileTarget)
			}
		} else {
			fmt.Println("Sinks: (none configured)")
		}
		return nil
	},
}

var registryUpdateCmd = &cobra.Command{
	Use:   "update <name> <field> <value>",
	Short: "Update registry field",
	Long: `Update a specific field in an existing registry configuration.

Arguments:
  name   Registry name
  field  Field to update (url, type, branches)
  value  New field value (comma-separated for branches)

Examples:
  arm config registry update ai-rules url https://new-url
  arm config registry update ai-rules branches main,develop`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]
		manifestManager := manifest.NewFileManager()
		return manifestManager.UpdateGitRegistry(context.Background(), name, field, value)
	},
}

var sinkUpdateCmd = &cobra.Command{
	Use:   "update <name> <field> <value>",
	Short: "Update sink field",
	Long: `Update a specific field in an existing sink configuration.

Arguments:
  name   Sink name
  field  Field to update (directory, layout, compileTarget)
  value  New field value

Examples:
  arm config sink update q directory .amazonq/rules
  arm config sink update q layout flat
  arm config sink update q compileTarget amazonq`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]
		manifestManager := manifest.NewFileManager()
		// Update sink config
		err := manifestManager.UpdateSink(context.Background(), name, field, value)
		if err != nil {
			return err
		}
		return nil
	},
}

var rulesetUpdateCmd = &cobra.Command{
	Use:   "update <name> <field> <value>",
	Short: "Update ruleset configuration",
	Long: `Update a specific field in an existing ruleset configuration. This triggers a reinstall of the ruleset.

Arguments:
  name   Ruleset name (registry/ruleset)
  field  Field to update (priority, version, sinks, include, exclude)
  value  New field value

Examples:
  arm config ruleset update ai-rules/ruleset priority 150
  arm config ruleset update ai-rules/ruleset version 1.1.0
  arm config ruleset update ai-rules/ruleset sinks cursor,q
  arm config ruleset update ai-rules/ruleset include "**/*.py,**/*.js"
  arm config ruleset update ai-rules/ruleset exclude "**/test/**,**/node_modules/**"`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]

		// Parse ruleset name
		rulesets, err := ParseRulesetArgs([]string{name})
		if err != nil {
			return fmt.Errorf("failed to parse ruleset name: %w", err)
		}
		ruleset := rulesets[0]

		// Use service to update ruleset config
		return armService.UpdateRulesetConfig(context.Background(), ruleset.Registry, ruleset.Name, field, value)
	},
}

func init() {
	configRegistryCmd.AddCommand(registryAddCmd)
	configRegistryCmd.AddCommand(registryRemoveCmd)
	configRegistryCmd.AddCommand(registryUpdateCmd)
	configSinkCmd.AddCommand(sinkAddCmd)
	configSinkCmd.AddCommand(sinkRemoveCmd)
	configSinkCmd.AddCommand(sinkUpdateCmd)
	configRulesetCmd.AddCommand(rulesetUpdateCmd)

	registryAddCmd.Flags().String("type", "git", "Registry type (git, http)")
	registryAddCmd.Flags().StringSlice("branches", nil, "Git branches to track (default: main,master)")
	registryAddCmd.Flags().Bool("force", false, "Overwrite existing registry")
	sinkAddCmd.Flags().String("type", "", "Sink type with defaults (cursor, copilot, amazonq)")
	sinkAddCmd.Flags().String("layout", "", "Layout mode (hierarchical, flat)")
	sinkAddCmd.Flags().String("compile-to", "", "Target format for compilation (cursor, amazonq, markdown, copilot)")
	sinkAddCmd.Flags().Bool("force", false, "Overwrite existing sink")
}
