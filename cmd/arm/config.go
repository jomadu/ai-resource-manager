package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage configuration",
		Long:  "Manage ARM configuration including registries, sinks, and resources.",
	}

	cmd.AddCommand(configRegistryCmd)
	cmd.AddCommand(configSinkCmd)
	cmd.AddCommand(configRulesetCmd)
	cmd.AddCommand(configPromptsetCmd)

	return cmd
}

var configRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Manage registry configuration",
	Long: `Manage registry configuration for ARM.

Registries are remote sources where rulesets are stored and versioned, similar to npm registries for JavaScript packages. ARM supports Git, GitLab, and Cloudsmith registries for storing and distributing resources.

Available commands:
  add     Add a new registry
  remove  Remove an existing registry
  set     Set registry configuration values

Examples:
  arm config registry add my-org https://github.com/my-org/arm-registry --type git
  arm config registry add my-gitlab https://gitlab.com --type gitlab --gitlab-group-id 123
  arm config registry add sample-registry https://app.cloudsmith.com/sample-org/arm-registry --type cloudsmith
  arm config registry remove my-org`,
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
  set     Set sink configuration values

Examples:
  arm config sink add cursor-rules .cursor/rules --type cursor
  arm config sink add q-rules .amazonq/rules --type amazonq
  arm config sink add copilot-rules .github/copilot --type copilot`,
}

var configRulesetCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "Manage ruleset configuration",
	Long: `Manage ruleset configuration for ARM.

Rulesets are collections of AI rules that can be configured with priorities, version constraints, and sink assignments.

Available commands:
  set     Set ruleset configuration values (triggers reinstall)

Examples:
  arm config ruleset update ai-rules/ruleset priority 150
  arm config ruleset update ai-rules/ruleset version 1.1.0
  arm config ruleset update ai-rules/ruleset sinks cursor,q`,
}

var configPromptsetCmd = &cobra.Command{
	Use:   "promptset",
	Short: "Manage promptset configuration",
	Long: `Manage promptset configuration for ARM.

Promptsets are collections of AI prompts that can be configured with version constraints and sink assignments. Unlike rulesets, promptsets do not have priorities.

Available commands:
  set     Set promptset configuration values (triggers reinstall)

Examples:
  arm config promptset update ai-rules/promptset version 1.1.0
  arm config promptset update ai-rules/promptset sinks cursor,q`,
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
  --type                Registry type (git, gitlab, cloudsmith)
  --git-branches        Git branches to track (for git type)
  --gitlab-group-id     GitLab group ID (for gitlab type)
  --gitlab-project-id   GitLab project ID (for gitlab type)
  --gitlab-api-version  GitLab API version (for gitlab type)

Examples:
  arm config registry add my-org https://github.com/my-org/arm-registry --type git
  arm config registry add my-gitlab https://gitlab.com --type gitlab --gitlab-group-id 123
  arm config registry add my-gitlab-project https://gitlab.com --type gitlab --gitlab-project-id 456
  arm config registry add sample-registry https://app.cloudsmith.com/sample-org/arm-registry --type cloudsmith`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		url := args[1]
		registryType, _ := cmd.Flags().GetString("type")
		if registryType == "" {
			registryType = "git"
		}
		options := make(map[string]interface{})
		switch registryType {
		case "git":
			branches, _ := cmd.Flags().GetStringSlice("git-branches")
			options["branches"] = branches
		case "gitlab":
			gitlabProjectID, _ := cmd.Flags().GetString("gitlab-project-id")
			gitlabGroupID, _ := cmd.Flags().GetString("gitlab-group-id")
			apiVersion, _ := cmd.Flags().GetString("gitlab-api-version")
			if apiVersion == "" {
				apiVersion = "v4"
			}
			if gitlabProjectID == "" && gitlabGroupID == "" {
				return fmt.Errorf("either --gitlab-project-id or --gitlab-group-id must be specified for GitLab registries")
			}
			options["project_id"] = gitlabProjectID
			options["group_id"] = gitlabGroupID
			options["api_version"] = apiVersion
		case "cloudsmith":
			// Parse URL to extract owner and repository
			owner, repository, err := registry.ParseCloudsmithURL(url)
			if err != nil {
				return fmt.Errorf("failed to parse Cloudsmith URL: %w", err)
			}
			options["owner"] = owner
			options["repository"] = repository
		default:
			return fmt.Errorf("registry type %s is not implemented", registryType)
		}
		return armService.AddRegistry(context.Background(), name, url, registryType, options)
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
  arm config registry remove my-org
  arm config registry remove my-gitlab`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return armService.RemoveRegistry(context.Background(), name)
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
  arm config sink add cursor-rules .cursor/rules --type cursor
  arm config sink add q-rules .amazonq/rules --type amazonq
  arm config sink add copilot-rules .github/copilot --type copilot
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
  arm config sink remove cursor-rules
  arm config sink remove q-rules`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		return armService.RemoveSink(context.Background(), name)
	},
}

var registrySetCmd = &cobra.Command{
	Use:   "set <name> <field> <value>",
	Short: "Set registry field",
	Long: `Update a specific field in an existing registry configuration.

Arguments:
  name   Registry name
  field  Field to update (url, type, git_branches)
  value  New field value (comma-separated for branches)

Examples:
  arm config registry set my-org url https://github.com/my-org/new-arm-registry
  arm config registry set my-gitlab gitlab_group_id 789`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]
		return armService.SetRegistryConfig(context.Background(), name, field, value)
	},
}

var sinkSetCmd = &cobra.Command{
	Use:   "set <name> <field> <value>",
	Short: "Set sink field",
	Long: `Update a specific field in an existing sink configuration.

Arguments:
  name   Sink name
  field  Field to update (directory, layout, compileTarget)
  value  New field value

Examples:
  arm config sink set cursor-rules layout flat
  arm config sink set cursor-rules directory .cursor/new-rules
  arm config sink set cursor-rules compile_target md`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]
		return armService.SetSinkConfig(context.Background(), name, field, value)
	},
}

var rulesetSetCmd = &cobra.Command{
	Use:   "set <name> <field> <value>",
	Short: "Set ruleset configuration",
	Long: `Update a specific field in an existing ruleset configuration. This triggers a reinstall of the ruleset.

Arguments:
  name   Ruleset name (registry/ruleset)
  field  Field to update (priority, version, sinks, include, exclude)
  value  New field value

Examples:
  arm config ruleset set my-org/clean-code-ruleset priority 150
  arm config ruleset set my-org/clean-code-ruleset version ^2.0.0
  arm config ruleset set my-org/clean-code-ruleset sinks cursor-rules,q-rules`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, value := args[0], args[1], args[2]

		// Parse ruleset name
		ruleset, err := ParsePackageArg(name)
		if err != nil {
			return fmt.Errorf("failed to parse ruleset name: %w", err)
		}

		// Use service to update ruleset config
		return armService.SetRulesetConfig(context.Background(), ruleset.Registry, ruleset.Name, field, value)
	},
}

var promptsetSetCmd = &cobra.Command{
	Use:   "set <name> <field> <value>",
	Short: "Set promptset configuration",
	Long: `Update a specific field in an existing promptset configuration. This triggers a reinstall of the promptset.

Arguments:
  name   Promptset name (registry/promptset)
  field  Field to update (version, sinks, include, exclude) - priority not supported for promptsets
  value  New field value

Examples:
  arm config promptset set my-org/code-review-promptset version ^2.0.0
  arm config promptset set my-org/code-review-promptset sinks cursor-prompts,q-prompts`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, field, _ := args[0], args[1], args[2]

		// Validate that priority is not being set for promptsets
		if field == "priority" {
			return fmt.Errorf("priority is not supported for promptsets")
		}

		// Parse promptset name
		_, err := ParsePackageArg(name)
		if err != nil {
			return fmt.Errorf("failed to parse promptset name: %w", err)
		}

		// TODO: Implement promptset config update when service interface is updated
		return fmt.Errorf("promptset config update not yet implemented - service interface needs to be updated first")
	},
}

func init() {
	configRegistryCmd.AddCommand(registryAddCmd)
	configRegistryCmd.AddCommand(registryRemoveCmd)
	configRegistryCmd.AddCommand(registrySetCmd)
	configSinkCmd.AddCommand(sinkAddCmd)
	configSinkCmd.AddCommand(sinkRemoveCmd)
	configSinkCmd.AddCommand(sinkSetCmd)
	configRulesetCmd.AddCommand(rulesetSetCmd)
	configPromptsetCmd.AddCommand(promptsetSetCmd)

	registryAddCmd.Flags().String("type", "git", "Registry type (git, gitlab, cloudsmith)")
	registryAddCmd.Flags().StringSlice("git-branches", nil, "Git branches to track (for git type)")

	registryAddCmd.Flags().String("gitlab-project-id", "", "GitLab project ID (for gitlab type)")
	registryAddCmd.Flags().String("gitlab-group-id", "", "GitLab group ID (for gitlab type)")
	registryAddCmd.Flags().String("gitlab-api-version", "v4", "GitLab API version (for gitlab type)")
	registryAddCmd.Flags().Bool("force", false, "Overwrite existing registry")
	sinkAddCmd.Flags().String("type", "", "Sink type with defaults (cursor, copilot, amazonq)")
	sinkAddCmd.Flags().String("layout", "", "Layout mode (hierarchical, flat)")
	sinkAddCmd.Flags().String("compile-to", "", "Target format for compilation (cursor, amazonq, markdown, copilot)")
	sinkAddCmd.Flags().Bool("force", false, "Overwrite existing sink")
}
