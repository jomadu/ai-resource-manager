package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/urf"
	"github.com/spf13/cobra"
)

func newCompileCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compile [file...]",
		Short: "Compile URF files to target formats",
		Long: `Compile Universal Rule Format (URF) files to specific AI tool formats.

Supports compilation to various targets:
- cursor: Cursor-compatible .mdc files with YAML frontmatter
- amazonq: Amazon Q compatible .md files
- copilot: GitHub Copilot .instructions.md files
- markdown: Generic markdown files

Examples:
  arm compile rules.yaml --target cursor
  arm compile rules.yaml --target cursor,amazonq,copilot --output ./compiled
  arm compile ./rules/ --target cursor --recursive --output .cursor/rules
  arm compile rules.yaml --target cursor --dry-run --verbose
  arm compile *.yaml --validate-only`,
		RunE: runCompile,
		Args: cobra.MinimumNArgs(1),
	}

	// Core flags
	cmd.Flags().StringP("target", "t", "", "Target format(s) - comma-separated (cursor, amazonq, markdown, copilot) [REQUIRED]")
	cmd.Flags().StringP("output", "o", ".", "Output directory (defaults to current directory)")
	cmd.Flags().StringP("namespace", "n", "", "Namespace for compiled rules (defaults to filename)")
	cmd.Flags().BoolP("force", "f", false, "Overwrite existing files")
	cmd.Flags().BoolP("recursive", "r", false, "Recursively find URF files in directories")

	// Processing flags
	cmd.Flags().Bool("dry-run", false, "Show what would be compiled without writing files")
	cmd.Flags().BoolP("verbose", "v", false, "Show detailed compilation information")
	cmd.Flags().Bool("validate-only", false, "Validate URF syntax without compilation")
	cmd.Flags().Bool("fail-fast", false, "Stop compilation on first error")

	// Filtering flags (reuse existing ARM patterns)
	cmd.Flags().StringSlice("include", nil, "Include patterns for file filtering")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns for file filtering")

	// Mark target as required
	_ = cmd.MarkFlagRequired("target")

	return cmd
}

func runCompile(cmd *cobra.Command, args []string) error {
	// Parse and validate flags
	targetStr, _ := cmd.Flags().GetString("target")
	outputDir, _ := cmd.Flags().GetString("output")
	namespace, _ := cmd.Flags().GetString("namespace")
	force, _ := cmd.Flags().GetBool("force")
	recursive, _ := cmd.Flags().GetBool("recursive")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	verbose, _ := cmd.Flags().GetBool("verbose")
	validateOnly, _ := cmd.Flags().GetBool("validate-only")
	failFast, _ := cmd.Flags().GetBool("fail-fast")
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")

	// Parse and validate targets (comma-separated)
	targets, err := parseTargets(targetStr)
	if err != nil {
		return fmt.Errorf("invalid target specification: %w", err)
	}

	// Apply default include patterns for YAML files
	include = GetDefaultIncludePatterns(include)

	// Validate conflicting flags
	if validateOnly && (dryRun || force) {
		return fmt.Errorf("--validate-only cannot be used with --dry-run or --force")
	}

	// Create service and compile request
	service := arm.NewArmService()

	request := &arm.CompileRequest{
		Files:        args,
		Targets:      targets,
		OutputDir:    outputDir,
		Namespace:    namespace,
		Force:        force,
		Recursive:    recursive,
		DryRun:       dryRun,
		Verbose:      verbose,
		ValidateOnly: validateOnly,
		FailFast:     failFast,
		Include:      include,
		Exclude:      exclude,
	}

	// Execute compilation
	ctx := context.Background()
	result, err := service.CompileFiles(ctx, request)
	if err != nil {
		return fmt.Errorf("compilation failed: %w", err)
	}

	// Display results
	return displayCompileResults(result, verbose, dryRun)
}

// parseTargets parses comma-separated target string and validates each target
func parseTargets(targetStr string) ([]urf.CompileTarget, error) {
	if targetStr == "" {
		return nil, fmt.Errorf("target is required")
	}

	targetParts := strings.Split(targetStr, ",")
	targets := make([]urf.CompileTarget, 0, len(targetParts))
	seen := make(map[string]bool)

	for _, part := range targetParts {
		target := strings.TrimSpace(part)
		if target == "" {
			continue
		}

		// Check for duplicates
		if seen[target] {
			return nil, fmt.Errorf("duplicate target: %s", target)
		}
		seen[target] = true

		// Add target (validation will happen when creating compiler)
		compileTarget := urf.CompileTarget(target)
		targets = append(targets, compileTarget)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid targets specified")
	}

	return targets, nil
}

// displayCompileResults formats and displays compilation results
func displayCompileResults(result *arm.CompileResult, verbose, dryRun bool) error {
	if result == nil {
		return fmt.Errorf("no compilation result to display")
	}

	// Determine exit code based on errors
	exitCode := 0
	if result.Stats.Errors > 0 {
		if result.Stats.FilesCompiled == 0 {
			exitCode = 2 // All failures
		} else {
			exitCode = 1 // Some failures
		}
	}

	// Display summary statistics
	fmt.Printf("Compilation Summary:\n")
	fmt.Printf("  Files processed: %d\n", result.Stats.FilesProcessed)
	fmt.Printf("  Files compiled:  %d\n", result.Stats.FilesCompiled)
	if result.Stats.FilesSkipped > 0 {
		fmt.Printf("  Files skipped:   %d\n", result.Stats.FilesSkipped)
	}
	if result.Stats.Errors > 0 {
		fmt.Printf("  Errors:          %d\n", result.Stats.Errors)
	}
	if result.Stats.RulesGenerated > 0 {
		fmt.Printf("  Rules generated: %d\n", result.Stats.RulesGenerated)
	}

	// Show per-target statistics if multiple targets
	if len(result.Stats.TargetStats) > 1 {
		fmt.Printf("\nPer-target compilation:\n")
		for target, count := range result.Stats.TargetStats {
			fmt.Printf("  %s: %d files\n", target, count)
		}
	}

	// Show compiled files in verbose mode or for dry-run
	if verbose || dryRun {
		if len(result.CompiledFiles) > 0 {
			fmt.Printf("\nCompiled files:\n")
			for _, file := range result.CompiledFiles {
				if dryRun {
					fmt.Printf("  [DRY RUN] %s -> %s (%s)\n", file.SourcePath, file.TargetPath, file.Target)
				} else {
					fmt.Printf("  %s -> %s (%s)\n", file.SourcePath, file.TargetPath, file.Target)
				}
			}
		}
	}

	// Show skipped files if any
	if len(result.Skipped) > 0 && verbose {
		fmt.Printf("\nSkipped files:\n")
		for _, skipped := range result.Skipped {
			fmt.Printf("  %s (%s)\n", skipped.Path, skipped.Reason)
		}
	}

	// Show errors if any
	if len(result.Errors) > 0 {
		fmt.Printf("\nErrors:\n")
		for _, err := range result.Errors {
			if err.Target != "" {
				fmt.Printf("  %s [%s]: %s\n", err.FilePath, err.Target, err.Error)
			} else {
				fmt.Printf("  %s: %s\n", err.FilePath, err.Error)
			}
		}
	}

	// Return appropriate exit code by using os.Exit indirectly
	if exitCode != 0 {
		return fmt.Errorf("compilation completed with %d errors", result.Stats.Errors)
	}

	return nil
}
