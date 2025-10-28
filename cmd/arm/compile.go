package main

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

var compileCmd = &cobra.Command{
	Use:   "compile [--target <md|cursor|amazonq|copilot>] [--force] [--recursive] [--validate-only] [--include GLOB...] [--exclude GLOB...] [--fail-fast] INPUT_PATH... [OUTPUT_PATH]",
	Short: "Compile resources",
	Long: `Compile rulesets and promptsets from source files. This command compiles source ruleset and promptset files to platform-specific formats.

It supports different target platforms (md, cursor, amazonq, copilot), recursive directory processing, validation-only mode, and various filtering and output options. This is useful for development and testing of rulesets and promptsets before publishing to registries.

When using --validate-only, OUTPUT_PATH is optional and will be ignored if provided.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		compileFiles(cmd, args)
	},
}

func init() {
	// Add compile flags
	compileCmd.Flags().String("target", "cursor", "Target platform (md, cursor, amazonq, copilot)")
	compileCmd.Flags().Bool("force", false, "Force overwrite existing files")
	compileCmd.Flags().Bool("recursive", false, "Process directories recursively")
	compileCmd.Flags().Bool("validate-only", false, "Validate only (no output files)")
	compileCmd.Flags().StringSlice("include", []string{"**/*.yml", "**/*.yaml"}, "Include patterns")
	compileCmd.Flags().StringSlice("exclude", []string{}, "Exclude patterns")
	compileCmd.Flags().Bool("fail-fast", false, "Stop on first error")
	compileCmd.Flags().String("namespace", "", "Namespace for compiled files")
	compileCmd.Flags().Bool("verbose", false, "Verbose output")
}

func compileFiles(cmd *cobra.Command, args []string) {
	target, _ := cmd.Flags().GetString("target")
	force, _ := cmd.Flags().GetBool("force")
	recursive, _ := cmd.Flags().GetBool("recursive")
	validateOnly, _ := cmd.Flags().GetBool("validate-only")
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	failFast, _ := cmd.Flags().GetBool("fail-fast")
	namespace, _ := cmd.Flags().GetString("namespace")
	verbose, _ := cmd.Flags().GetBool("verbose")

	var inputPaths []string
	var outputPath string

	// Handle arguments based on validate-only mode
	if validateOnly {
		// In validate-only mode, all args are input paths
		inputPaths = args
		outputPath = "" // Will be ignored
	} else {
		// In normal mode, require at least 2 args (input + output)
		if len(args) < 2 {
			handleCommandError(fmt.Errorf("compile requires at least 2 arguments: INPUT_PATH... OUTPUT_PATH"))
			return
		}
		inputPaths = args[:len(args)-1]
		outputPath = args[len(args)-1]
	}

	req := &arm.CompileRequest{
		Paths:        inputPaths,
		Targets:      []string{target},
		OutputDir:    outputPath,
		Namespace:    namespace,
		Force:        force,
		Recursive:    recursive,
		Verbose:      verbose,
		ValidateOnly: validateOnly,
		Include:      include,
		Exclude:      exclude,
		FailFast:     failFast,
	}

	if err := armService.CompileFiles(ctx, req); err != nil {
		handleCommandError(err)
		return
	}
}
