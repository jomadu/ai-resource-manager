package main

import (
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add registries and sinks",
	Long:  "Add registries and sinks to the ARM configuration",
}

var addRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "Add a new registry",
	Long:  `Add a new registry to the ARM configuration. Use subcommands for different registry types (git, gitlab, cloudsmith).`,
}

var addRegistryGitCmd = &cobra.Command{
	Use:   "git --url URL [--branches BRANCH...] [--force] NAME",
	Short: "Add a new Git registry",
	Long:  `Add a new Git registry to the ARM configuration. Git registries use Git repositories for storing and versioning rulesets and promptsets.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addGitRegistry(cmd, args[0])
	},
}

var addRegistryGitLabCmd = &cobra.Command{
	Use:   "gitlab [--url URL] [--group-id ID] [--project-id ID] [--api-version VERSION] [--force] NAME",
	Short: "Add a new GitLab registry",
	Long:  `Add a new GitLab registry to the ARM configuration. URL defaults to https://gitlab.com if not specified.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addGitLabRegistry(cmd, args[0])
	},
}

var addRegistryCloudsmithCmd = &cobra.Command{
	Use:   "cloudsmith [--url URL] [--owner OWNER] [--repo REPO] [--force] NAME",
	Short: "Add a new Cloudsmith registry",
	Long:  `Add a new Cloudsmith registry to the ARM configuration. URL defaults to https://api.cloudsmith.io if not specified.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addCloudsmithRegistry(cmd, args[0])
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

	// Add registry type subcommands
	addRegistryCmd.AddCommand(addRegistryGitCmd)
	addRegistryCmd.AddCommand(addRegistryGitLabCmd)
	addRegistryCmd.AddCommand(addRegistryCloudsmithCmd)

	// Git registry flags
	addRegistryGitCmd.Flags().String("url", "", "Git repository URL (required)")
	addRegistryGitCmd.MarkFlagRequired("url")
	addRegistryGitCmd.Flags().StringSlice("branches", []string{}, "Git branches to track")
	addRegistryGitCmd.Flags().Bool("force", false, "Overwrite existing registry")

	// GitLab registry flags
	addRegistryGitLabCmd.Flags().String("url", "https://gitlab.com", "GitLab instance URL")
	addRegistryGitLabCmd.Flags().String("group-id", "", "GitLab group ID")
	addRegistryGitLabCmd.Flags().String("project-id", "", "GitLab project ID")
	addRegistryGitLabCmd.Flags().String("api-version", "", "API version (defaults to v4)")
	addRegistryGitLabCmd.Flags().Bool("force", false, "Overwrite existing registry")

	// Cloudsmith registry flags
	addRegistryCloudsmithCmd.Flags().String("url", "https://api.cloudsmith.io", "Cloudsmith API URL")
	addRegistryCloudsmithCmd.Flags().String("owner", "", "Cloudsmith owner (required)")
	addRegistryCloudsmithCmd.MarkFlagRequired("owner")
	addRegistryCloudsmithCmd.Flags().String("repo", "", "Cloudsmith repository (required)")
	addRegistryCloudsmithCmd.MarkFlagRequired("repo")
	addRegistryCloudsmithCmd.Flags().Bool("force", false, "Overwrite existing registry")

	// Add sink flags
	addSinkCmd.Flags().String("type", "", "Sink type (cursor, copilot, amazonq)")
	addSinkCmd.Flags().String("layout", "hierarchical", "Layout type (hierarchical, flat)")
	addSinkCmd.Flags().String("compile-to", "cursor", "Compile target (md, cursor, amazonq, copilot)")
	addSinkCmd.Flags().Bool("force", false, "Overwrite existing sink")
}

func addGitRegistry(cmd *cobra.Command, name string) {
	url, _ := cmd.Flags().GetString("url")
	branches, _ := cmd.Flags().GetStringSlice("branches")
	force, _ := cmd.Flags().GetBool("force")

	options := make(map[string]interface{})
	if len(branches) > 0 {
		options["branches"] = branches
	}

	if err := armService.AddRegistry(ctx, name, url, "git", options, force); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func addGitLabRegistry(cmd *cobra.Command, name string) {
	url, _ := cmd.Flags().GetString("url")
	groupID, _ := cmd.Flags().GetString("group-id")
	projectID, _ := cmd.Flags().GetString("project-id")
	apiVersion, _ := cmd.Flags().GetString("api-version")
	force, _ := cmd.Flags().GetBool("force")

	options := make(map[string]interface{})
	if groupID != "" {
		options["group_id"] = groupID
	}
	if projectID != "" {
		options["project_id"] = projectID
	}
	if apiVersion != "" {
		options["api_version"] = apiVersion
	}

	if err := armService.AddRegistry(ctx, name, url, "gitlab", options, force); err != nil {
		cmd.PrintErrln("Error:", err)
		return
	}
}

func addCloudsmithRegistry(cmd *cobra.Command, name string) {
	url, _ := cmd.Flags().GetString("url")
	owner, _ := cmd.Flags().GetString("owner")
	repo, _ := cmd.Flags().GetString("repo")
	force, _ := cmd.Flags().GetBool("force")

	options := make(map[string]interface{})
	if owner != "" {
		options["owner"] = owner
	}
	if repo != "" {
		options["repository"] = repo
	}

	if err := armService.AddRegistry(ctx, name, url, "cloudsmith", options, force); err != nil {
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
