package main

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "compile [paths...]",
		Short:        "Compile resource files to target format",
		Long:         `Compile resource files (rulesets and promptsets) to specific AI tool formats. Paths can be individual files or directories containing resource files.`,
		RunE:         runCompile,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().StringSliceP("target", "t", []string{}, "Target format (cursor, amazonq, markdown, copilot) - supports multiple values")
	cmd.Flags().StringP("output", "o", ".", "Output directory")
	cmd.Flags().StringP("namespace", "n", "", "Optional namespace for compiled rules")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	cmd.Flags().BoolP("recursive", "r", false, "Recursively find resource files in directories")
	cmd.Flags().BoolP("verbose", "v", false, "Show detailed compilation information")
	cmd.Flags().Bool("validate-only", false, "Validate resource syntax without compilation")
	cmd.Flags().StringSlice("include", nil, "Include patterns for file filtering")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns for file filtering")
	cmd.Flags().Bool("fail-fast", false, "Stop compilation on first error")

	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func runCompile(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse flags
	targets, _ := cmd.Flags().GetStringSlice("target")
	outputDir, _ := cmd.Flags().GetString("output")
	namespace, _ := cmd.Flags().GetString("namespace")
	force, _ := cmd.Flags().GetBool("force")
	recursive, _ := cmd.Flags().GetBool("recursive")
	verbose, _ := cmd.Flags().GetBool("verbose")
	validateOnly, _ := cmd.Flags().GetBool("validate-only")
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	failFast, _ := cmd.Flags().GetBool("fail-fast")

	// Create compile request
	req := &arm.CompileRequest{
		Paths:        args,
		Targets:      targets,
		OutputDir:    outputDir,
		Namespace:    namespace,
		Force:        force,
		Recursive:    recursive,
		Verbose:      verbose,
		ValidateOnly: validateOnly,
		Include:      include,
		Exclude:      exclude,
		FailFast:     failFast,
	}

	// Execute compilation
	return armService.CompileFiles(ctx, req)
}
