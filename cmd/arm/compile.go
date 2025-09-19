package main

import (
	"fmt"
	"strings"

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
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	validateOnly, _ := cmd.Flags().GetBool("validate-only")
	force, _ := cmd.Flags().GetBool("force")
	include, _ := cmd.Flags().GetStringSlice("include")

	// Parse and validate targets (comma-separated)
	targets, err := parseTargets(targetStr)
	if err != nil {
		return fmt.Errorf("invalid target specification: %w", err)
	}

	// Apply default include patterns for YAML files
	_ = GetDefaultIncludePatterns(include)

	// Validate conflicting flags
	if validateOnly && (dryRun || force) {
		return fmt.Errorf("--validate-only cannot be used with --dry-run or --force")
	}

	// TODO: Create CompileRequest and call service
	// TODO: Display results

	fmt.Printf("Parsed targets: %v, output: %s, files: %v\n", targets, outputDir, args)
	return nil
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

		// Validate target
		compileTarget := urf.CompileTarget(target)
		if !isValidTarget(compileTarget) {
			return nil, fmt.Errorf("unsupported target: %s (supported: cursor, amazonq, markdown, copilot)", target)
		}

		targets = append(targets, compileTarget)
	}

	if len(targets) == 0 {
		return nil, fmt.Errorf("no valid targets specified")
	}

	return targets, nil
}

// isValidTarget checks if the target is supported
func isValidTarget(target urf.CompileTarget) bool {
	switch target {
	case urf.TargetCursor, urf.TargetAmazonQ, urf.TargetMarkdown, urf.TargetCopilot:
		return true
	default:
		return false
	}
}
