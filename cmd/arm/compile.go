package main

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "compile [file...]",
		Short:        "Compile URF files to target format",
		Long:         `Compile Universal Rule Format (URF) files to specific AI tool formats.`,
		RunE:         runCompile,
		Args:         cobra.MinimumNArgs(1),
		SilenceUsage: true,
	}

	cmd.Flags().StringP("target", "t", "", "Target format (cursor, amazonq, markdown, copilot) - supports comma-separated")
	cmd.Flags().StringP("output", "o", ".", "Output directory")
	cmd.Flags().StringP("namespace", "n", "", "Optional namespace for compiled rules")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	cmd.Flags().BoolP("recursive", "r", false, "Recursively find URF files in directories")
	cmd.Flags().BoolP("verbose", "v", false, "Show detailed compilation information")
	cmd.Flags().Bool("validate-only", false, "Validate URF syntax without compilation")
	cmd.Flags().StringSlice("include", nil, "Include patterns for file filtering")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns for file filtering")
	cmd.Flags().Bool("fail-fast", false, "Stop compilation on first error")

	_ = cmd.MarkFlagRequired("target")
	return cmd
}

func runCompile(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse flags
	target, _ := cmd.Flags().GetString("target")
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
		Files:        args,
		Target:       target,
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
