package arm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// CompileFiles compiles resource files to target formats
func (a *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
	// 1. Discover files
	files, err := a.discoverFiles(req.Files, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return fmt.Errorf("failed to discover files: %w", err)
	}

	if len(files) == 0 {
		a.ui.Warning("No resource files found matching the criteria")
		return nil
	}

	// 2. Parse targets
	targets := strings.Split(req.Target, ",")
	for i, target := range targets {
		targets[i] = strings.TrimSpace(target)
	}

	// 3. Process each file
	var errors []error
	stats := CompileStats{}

	for _, filePath := range files {
		stats.FilesProcessed++

		if req.Verbose {
			a.ui.CompileStep(fmt.Sprintf("Processing %s", filePath))
		}

		if req.ValidateOnly {
			if err := a.validateResourceFile(filePath); err != nil {
				a.ui.Error(fmt.Errorf("validation failed for %s: %w", filePath, err))
				errors = append(errors, err)
				stats.Errors++
				if req.FailFast {
					return err
				}
				continue
			}
			if req.Verbose {
				a.ui.Success(fmt.Sprintf("✓ %s validated", filePath))
			}
			continue
		}

		compiled, err := a.compileFile(filePath, targets, req)
		if err != nil {
			a.ui.Error(fmt.Errorf("compilation failed for %s: %w", filePath, err))
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

// discoverFiles finds resource files based on the input patterns
func (a *ArmService) discoverFiles(inputs []string, recursive bool, include, exclude []string) ([]string, error) {
	var files []string

	for _, input := range inputs {
		stat, err := os.Stat(input)
		if err != nil {
			return nil, fmt.Errorf("failed to stat %s: %w", input, err)
		}

		if stat.IsDir() {
			dirFiles, err := a.discoverInDirectory(input, recursive, include, exclude)
			if err != nil {
				return nil, err
			}
			files = append(files, dirFiles...)
		} else if a.isResourceFile(input) {
			files = append(files, input)
		}
	}

	return files, nil
}

// discoverInDirectory finds resource files in a directory
func (a *ArmService) discoverInDirectory(dir string, recursive bool, _, _ []string) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && a.isResourceFile(path) {
				files = append(files, path)
			}
			return nil
		})
		if err != nil {
			return nil, fmt.Errorf("failed to walk directory %s: %w", dir, err)
		}
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, fmt.Errorf("failed to read directory %s: %w", dir, err)
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				path := filepath.Join(dir, entry.Name())
				if a.isResourceFile(path) {
					files = append(files, path)
				}
			}
		}
	}

	return files, nil
}

// isResourceFile checks if a file is a resource file by extension
// TODO: this is dumb. it should be checking if the file is both the yml extension, and can be parsed as any resource kind.
func (a *ArmService) isResourceFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

// validateResourceFile validates a resource file
func (a *ArmService) validateResourceFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	file := &types.File{
		Path:    filePath,
		Content: content,
		Size:    int64(len(content)),
	}

	parser := resource.NewParser()
	// Try parsing as either ruleset or promptset
	if _, err := parser.ParseRuleset(file); err != nil {
		if _, err := parser.ParsePromptset(file); err != nil {
			return fmt.Errorf("file is not a valid resource (neither ruleset nor promptset)")
		}
	}
	return err
}

// compileFile compiles a single resource file to multiple targets
func (a *ArmService) compileFile(filePath string, targets []string, req *CompileRequest) (int, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return 0, fmt.Errorf("failed to read file: %w", err)
	}

	file := &types.File{
		Path:    filePath,
		Content: content,
		Size:    int64(len(content)),
	}

	// Determine namespace
	namespace := req.Namespace
	if namespace == "" {
		namespace = strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))
	}

	totalRules := 0

	for _, target := range targets {
		compiler, err := resource.NewCompiler(resource.CompileTarget(target))
		if err != nil {
			return 0, fmt.Errorf("failed to create compiler for target %s: %w", target, err)
		}

		// Try to compile as ruleset first, then promptset
		parser := resource.NewParser()
		var compiledFiles []*types.File
		var compileErr error
		if _, err := parser.ParseRuleset(file); err == nil {
			compiledFiles, compileErr = compiler.CompileRuleset(namespace, file)
		} else if _, err := parser.ParsePromptset(file); err == nil {
			compiledFiles, compileErr = compiler.CompilePromptset(namespace, file)
		} else {
			return 0, fmt.Errorf("file is not a valid resource (neither ruleset nor promptset)")
		}
		if compileErr != nil {
			return 0, fmt.Errorf("failed to compile for target %s: %w", target, compileErr)
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
