package arm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/resource"
)

// CompileFiles compiles resource files to target formats
func (a *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
	// 1. Parse targets
	targets := strings.Split(req.Target, ",")
	for i, target := range targets {
		targets[i] = strings.TrimSpace(target)
	}

	// 2. Use parser to discover and parse all resources
	parser := resource.NewParser()

	// Parse rulesets from all paths
	rulesets, err := parser.ParseRulesets(req.Paths, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return fmt.Errorf("failed to parse rulesets: %w", err)
	}

	// Parse promptsets from all paths
	promptsets, err := parser.ParsePromptsets(req.Paths, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return fmt.Errorf("failed to parse promptsets: %w", err)
	}

	if len(rulesets) == 0 && len(promptsets) == 0 {
		a.ui.Warning("No resource files found matching the criteria")
		return nil
	}

	// 3. Process resources
	var errors []error
	stats := CompileStats{}

	// Process rulesets
	for _, ruleset := range rulesets {
		stats.FilesProcessed++

		if req.Verbose {
			a.ui.CompileStep(fmt.Sprintf("Processing ruleset %s", ruleset.Metadata.ID))
		}

		if req.ValidateOnly {
			if req.Verbose {
				a.ui.Success(fmt.Sprintf("✓ ruleset %s validated", ruleset.Metadata.ID))
			}
			continue
		}

		compiled, err := a.compileRuleset(ruleset, targets, req)
		if err != nil {
			a.ui.Error(fmt.Errorf("compilation failed for ruleset %s: %w", ruleset.Metadata.ID, err))
			errors = append(errors, err)
			stats.Errors++
			if req.FailFast {
				return err
			}
			continue
		}
		stats.FilesCompiled++
		stats.RulesGenerated += compiled
	}

	// Process promptsets
	for _, promptset := range promptsets {
		stats.FilesProcessed++

		if req.Verbose {
			a.ui.CompileStep(fmt.Sprintf("Processing promptset %s", promptset.Metadata.ID))
		}

		if req.ValidateOnly {
			if req.Verbose {
				a.ui.Success(fmt.Sprintf("✓ promptset %s validated", promptset.Metadata.ID))
			}
			continue
		}

		compiled, err := a.compilePromptset(promptset, targets, req)
		if err != nil {
			a.ui.Error(fmt.Errorf("compilation failed for promptset %s: %w", promptset.Metadata.ID, err))
			errors = append(errors, err)
			stats.Errors++
			if req.FailFast {
				return err
			}
			continue
		}
		stats.FilesCompiled++
		stats.RulesGenerated += compiled
	}

	// 4. Display results
	a.ui.CompileComplete(stats, req.ValidateOnly)

	if len(errors) > 0 {
		return fmt.Errorf("compilation completed with %d errors", len(errors))
	}

	return nil
}

// compileRuleset compiles a single ruleset to multiple targets
func (a *ArmService) compileRuleset(ruleset *resource.Ruleset, targets []string, req *CompileRequest) (int, error) {
	// Determine namespace
	namespace := req.Namespace
	if namespace == "" {
		namespace = ruleset.Metadata.ID
	}

	totalRules := 0

	for _, target := range targets {
		compiler, err := resource.NewCompiler(resource.CompileTarget(target))
		if err != nil {
			return 0, fmt.Errorf("failed to create compiler for target %s: %w", target, err)
		}

		// Compile the ruleset
		compiledFiles, err := compiler.CompileRuleset(namespace, ruleset)
		if err != nil {
			return 0, fmt.Errorf("failed to compile ruleset for target %s: %w", target, err)
		}

		// Write compiled files
		outputDir := req.OutputDir
		if len(targets) > 1 {
			outputDir = filepath.Join(req.OutputDir, target)
		}

		// Create output directory
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return 0, fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}

		for _, compiledFile := range compiledFiles {
			outputPath := filepath.Join(outputDir, compiledFile.Path)

			if !req.Force {
				if _, err := os.Stat(outputPath); err == nil {
					return 0, fmt.Errorf("output file %s already exists (use --force to overwrite)", outputPath)
				}
			}

			if err := os.WriteFile(outputPath, compiledFile.Content, 0o644); err != nil {
				return 0, fmt.Errorf("failed to write output file %s: %w", outputPath, err)
			}
			if req.Verbose {
				a.ui.CompileStep(fmt.Sprintf("Wrote %s", outputPath))
			}
		}

		totalRules += len(compiledFiles)
	}

	return totalRules, nil
}

// compilePromptset compiles a single promptset to multiple targets
func (a *ArmService) compilePromptset(promptset *resource.Promptset, targets []string, req *CompileRequest) (int, error) {
	// Determine namespace
	namespace := req.Namespace
	if namespace == "" {
		namespace = promptset.Metadata.ID
	}

	totalPrompts := 0

	for _, target := range targets {
		compiler, err := resource.NewCompiler(resource.CompileTarget(target))
		if err != nil {
			return 0, fmt.Errorf("failed to create compiler for target %s: %w", target, err)
		}

		// Compile the promptset
		compiledFiles, err := compiler.CompilePromptset(namespace, promptset)
		if err != nil {
			return 0, fmt.Errorf("failed to compile promptset for target %s: %w", target, err)
		}

		// Write compiled files
		outputDir := req.OutputDir
		if len(targets) > 1 {
			outputDir = filepath.Join(req.OutputDir, target)
		}

		// Create output directory
		if err := os.MkdirAll(outputDir, 0o755); err != nil {
			return 0, fmt.Errorf("failed to create output directory %s: %w", outputDir, err)
		}

		for _, compiledFile := range compiledFiles {
			outputPath := filepath.Join(outputDir, compiledFile.Path)

			if !req.Force {
				if _, err := os.Stat(outputPath); err == nil {
					return 0, fmt.Errorf("output file %s already exists (use --force to overwrite)", outputPath)
				}
			}

			if err := os.WriteFile(outputPath, compiledFile.Content, 0o644); err != nil {
				return 0, fmt.Errorf("failed to write output file %s: %w", outputPath, err)
			}
			if req.Verbose {
				a.ui.CompileStep(fmt.Sprintf("Wrote %s", outputPath))
			}
		}

		totalPrompts += len(compiledFiles)
	}

	return totalPrompts, nil
}
