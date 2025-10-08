package main

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add resources",
	Long:  "Add registries and sinks to the ARM configuration",
}

var addRegistryCmd = &cobra.Command{
	Use:   "registry [--type <git|gitlab|cloudsmith>] [--branches BRANCH...] [--group-id ID] [--project-id ID] [--api-version VERSION] [--owner OWNER] [--repo REPO] [--force] NAME URL",
	Short: "Add a new registry",
	Long: `Add a new registry to the ARM configuration. This command supports different registry types (git, gitlab, cloudsmith)
and allows specifying additional parameters like GitLab group and project IDs, or Cloudsmith owner and repository for more precise targeting.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addRegistry(cmd, args[0], args[1])
	},
}

var addSinkCmd = &cobra.Command{
	Use:   "sink [--type <cursor|copilot|amazonq>] [--layout <hierarchical|flat>] [--compile-to <md|cursor|amazonq|copilot>] [--force] NAME PATH",
	Short: "Add a new sink",
	Long:  `Add a new sink to the ARM configuration. A sink defines where and how compiled rulesets and promptsets should be output.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		addSink(cmd, args[0], args[1])
	},
}

func init() {
	// Add subcommands
	addCmd.AddCommand(addRegistryCmd)
	addCmd.AddCommand(addSinkCmd)

	// Add registry flags
	addRegistryCmd.Flags().String("type", "git", "Registry type (git, gitlab, cloudsmith)")
	addRegistryCmd.Flags().StringSlice("branches", []string{}, "Git branches to track")
	addRegistryCmd.Flags().String("group-id", "", "GitLab group ID")
	addRegistryCmd.Flags().String("project-id", "", "GitLab project ID")
	addRegistryCmd.Flags().String("api-version", "", "API version")
	addRegistryCmd.Flags().String("owner", "", "Cloudsmith owner")
	addRegistryCmd.Flags().String("repo", "", "Cloudsmith repository")
	addRegistryCmd.Flags().Bool("force", false, "Overwrite existing registry")

	// Add sink flags
	addSinkCmd.Flags().String("type", "", "Sink type (cursor, copilot, amazonq)")
	addSinkCmd.Flags().String("layout", "hierarchical", "Layout type (hierarchical, flat)")
	addSinkCmd.Flags().String("compile-to", "cursor", "Compile target (md, cursor, amazonq, copilot)")
	addSinkCmd.Flags().Bool("force", false, "Overwrite existing sink")
}

func addRegistry(cmd *cobra.Command, name, url string) {
	registryType, _ := cmd.Flags().GetString("type")
	branches, _ := cmd.Flags().GetStringSlice("branches")
	groupID, _ := cmd.Flags().GetString("group-id")
	projectID, _ := cmd.Flags().GetString("project-id")
	apiVersion, _ := cmd.Flags().GetString("api-version")
	owner, _ := cmd.Flags().GetString("owner")
	repo, _ := cmd.Flags().GetString("repo")
	force, _ := cmd.Flags().GetBool("force")

	// Build options map
	options := make(map[string]interface{})
	if len(branches) > 0 {
		options["branches"] = branches
	}
	if groupID != "" {
		options["group_id"] = groupID
	}
	if projectID != "" {
		options["project_id"] = projectID
	}
	if apiVersion != "" {
		options["api_version"] = apiVersion
	}
	if owner != "" {
		options["owner"] = owner
	}
	if repo != "" {
		options["repository"] = repo
	}

	if err := armService.AddRegistry(ctx, name, url, registryType, options, force); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func addSink(cmd *cobra.Command, name, path string) {
	sinkType, _ := cmd.Flags().GetString("type")
	layout, _ := cmd.Flags().GetString("layout")
	compileTo, _ := cmd.Flags().GetString("compile-to")
	force, _ := cmd.Flags().GetBool("force")

	// Handle type shortcuts
	if sinkType != "" {
		switch sinkType {
		case "cursor":
			layout = "hierarchical"
			compileTo = "cursor"
		case "copilot":
			layout = "flat"
			compileTo = "copilot"
		case "amazonq":
			layout = "hierarchical"
			compileTo = "markdown"
		}
	}

	if err := armService.AddSink(ctx, name, path, layout, compileTo, force); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}
