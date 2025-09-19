package main

import (
	"fmt"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/version"
)

// CompileOutputFormatter handles formatting and display of compilation results
type CompileOutputFormatter struct {
	verbose bool
	dryRun  bool
}

// NewCompileOutputFormatter creates a new output formatter
func NewCompileOutputFormatter(verbose, dryRun bool) *CompileOutputFormatter {
	return &CompileOutputFormatter{
		verbose: verbose,
		dryRun:  dryRun,
	}
}

// DisplayResults formats and displays compilation results with appropriate exit handling
func (f *CompileOutputFormatter) DisplayResults(result *arm.CompileResult) error {
	if result == nil {
		return fmt.Errorf("no compilation result to display")
	}

	// Display summary statistics
	f.displaySummary(result)

	// Display per-target statistics for multi-target builds
	if len(result.Stats.TargetStats) > 1 {
		f.displayTargetStats(result)
	}

	// Display file details in verbose mode or for dry-run
	if f.verbose || f.dryRun {
		f.displayFileDetails(result)
	}

	// Display skipped files if any (verbose mode only)
	if len(result.Skipped) > 0 && f.verbose {
		f.displaySkippedFiles(result)
	}

	// Display errors if any
	if len(result.Errors) > 0 {
		f.displayErrors(result)
	}

	// Return error for exit code handling
	return f.getExitError(result)
}

// displaySummary shows compilation summary statistics
func (f *CompileOutputFormatter) displaySummary(result *arm.CompileResult) {
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
}

// displayTargetStats shows per-target compilation statistics
func (f *CompileOutputFormatter) displayTargetStats(result *arm.CompileResult) {
	fmt.Printf("\nPer-target compilation:\n")
	for target, count := range result.Stats.TargetStats {
		fmt.Printf("  %s: %d files\n", target, count)
	}
}

// displayFileDetails shows detailed file compilation information
func (f *CompileOutputFormatter) displayFileDetails(result *arm.CompileResult) {
	if len(result.CompiledFiles) > 0 {
		fmt.Printf("\nCompiled files:\n")
		for _, file := range result.CompiledFiles {
			if f.dryRun {
				fmt.Printf("  [DRY RUN] %s -> %s (%s)\n",
					file.SourcePath, file.TargetPath, file.Target)
			} else {
				fmt.Printf("  %s -> %s (%s)\n",
					file.SourcePath, file.TargetPath, file.Target)
			}
		}
	}
}

// displaySkippedFiles shows files that were skipped during compilation
func (f *CompileOutputFormatter) displaySkippedFiles(result *arm.CompileResult) {
	fmt.Printf("\nSkipped files:\n")
	for _, skipped := range result.Skipped {
		fmt.Printf("  %s (%s)\n", skipped.Path, skipped.Reason)
	}
}

// displayErrors shows compilation errors
func (f *CompileOutputFormatter) displayErrors(result *arm.CompileResult) {
	fmt.Printf("\nErrors:\n")
	for _, err := range result.Errors {
		if err.Target != "" {
			fmt.Printf("  %s [%s]: %s\n", err.FilePath, err.Target, err.Error)
		} else {
			fmt.Printf("  %s: %s\n", err.FilePath, err.Error)
		}
	}
}

// getExitError determines the appropriate exit code and error message
func (f *CompileOutputFormatter) getExitError(result *arm.CompileResult) error {
	if result.Stats.Errors == 0 {
		return nil // Success
	}

	// Determine exit code based on compilation results
	if result.Stats.FilesCompiled == 0 {
		// Total failure - no files compiled successfully
		os.Exit(2)
	} else {
		// Partial failure - some files compiled successfully
		return fmt.Errorf("compilation completed with %d errors", result.Stats.Errors)
	}

	return nil // This line should never be reached due to os.Exit above
}

// DisplayValidationResults shows URF validation-only results
func (f *CompileOutputFormatter) DisplayValidationResults(result *arm.CompileResult) {
	fmt.Printf("URF Validation Summary:\n")
	fmt.Printf("  Files processed: %d\n", result.Stats.FilesProcessed)

	if result.Stats.Errors > 0 {
		fmt.Printf("  Validation errors: %d\n", result.Stats.Errors)
		f.displayErrors(result)
	} else {
		fmt.Printf("  All files are valid URF format\n")
	}

	if len(result.Skipped) > 0 && f.verbose {
		f.displaySkippedFiles(result)
	}
}

// DisplayDryRunPlan shows what would be compiled in dry-run mode
func (f *CompileOutputFormatter) DisplayDryRunPlan(result *arm.CompileResult) {
	fmt.Printf("Dry Run - Compilation Plan:\n")
	fmt.Printf("  Would process: %d files\n", result.Stats.FilesProcessed)
	fmt.Printf("  Would compile: %d files\n", result.Stats.FilesCompiled)

	if len(result.Stats.TargetStats) > 1 {
		fmt.Printf("\nTarget breakdown:\n")
		for target, count := range result.Stats.TargetStats {
			fmt.Printf("  %s: %d files\n", target, count)
		}
	}

	f.displayFileDetails(result)

	if len(result.Skipped) > 0 {
		f.displaySkippedFiles(result)
	}
}

// DisplayCompileProgress shows real-time compilation progress (if implemented)
func (f *CompileOutputFormatter) DisplayCompileProgress(sourceFile, targetFile, target string) {
	if f.verbose {
		fmt.Printf("Compiling: %s -> %s (%s)\n", sourceFile, targetFile, target)
	}
}

// Global output functions used by other ARM commands

// WriteError writes an error message to stderr
func WriteError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// FormatVersionInfo formats version information for display
func FormatVersionInfo(versionInfo version.VersionInfo) {
	fmt.Printf("ARM %s\n", versionInfo.Version)
	if versionInfo.Commit != "" {
		fmt.Printf("Commit: %s\n", versionInfo.Commit)
	}
	if versionInfo.Timestamp != "" {
		fmt.Printf("Built: %s\n", versionInfo.Timestamp)
	}
	fmt.Printf("Architecture: %s\n", versionInfo.Arch)
}

// FormatRulesetInfo formats a single ruleset's information for display
func FormatRulesetInfo(info *arm.RulesetInfo, verbose bool) {
	fmt.Printf("Registry: %s\n", info.Registry)
	fmt.Printf("Ruleset: %s\n", info.Name)

	if info.Installation.Version != "" {
		fmt.Printf("Version: %s\n", info.Installation.Version)
	}

	if verbose {
		if info.Manifest.Constraint != "" {
			fmt.Printf("Constraint: %s\n", info.Manifest.Constraint)
		}
		if info.Manifest.Priority > 0 {
			fmt.Printf("Priority: %d\n", info.Manifest.Priority)
		}
		if len(info.Manifest.Sinks) > 0 {
			fmt.Printf("Sinks: %v\n", info.Manifest.Sinks)
		}
	}
}

// FormatInstalledRulesets formats a list of installed rulesets for display
func FormatInstalledRulesets(rulesets []*arm.RulesetInfo, verbose, summary bool) {
	if len(rulesets) == 0 {
		fmt.Println("No rulesets installed")
		return
	}

	if summary {
		fmt.Printf("Total installed rulesets: %d\n", len(rulesets))
		return
	}

	fmt.Printf("Installed rulesets (%d):\n", len(rulesets))
	for _, ruleset := range rulesets {
		fmt.Printf("  %s/%s", ruleset.Registry, ruleset.Name)
		if ruleset.Installation.Version != "" {
			fmt.Printf("@%s", ruleset.Installation.Version)
		}
		fmt.Println()

		if verbose && len(ruleset.Manifest.Sinks) > 0 {
			fmt.Printf("    Sinks: %v\n", ruleset.Manifest.Sinks)
		}
	}
}

// FormatOutdatedRulesets formats a list of outdated rulesets for display and returns an error for command chains
func FormatOutdatedRulesets(outdated []arm.OutdatedRuleset, _ string) error {
	if len(outdated) == 0 {
		fmt.Println("All rulesets are up to date")
		return nil
	}

	fmt.Printf("Outdated rulesets (%d):\n", len(outdated))
	for _, ruleset := range outdated {
		fmt.Printf("  %s/%s: %s -> %s\n",
			ruleset.RulesetInfo.Registry,
			ruleset.RulesetInfo.Name,
			ruleset.RulesetInfo.Installation.Version,
			ruleset.Latest)
	}

	return nil
}
