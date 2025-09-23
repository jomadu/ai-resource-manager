package arm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

// CompileFiles compiles URF files to target formats
func (a *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
	// 1. Discover files
	files, err := a.discoverFiles(req.Files, req.Recursive, req.Include, req.Exclude)
	if err != nil {
		return fmt.Errorf("failed to discover files: %w", err)
	}

	if len(files) == 0 {
		a.ui.Warning("No URF files found matching the criteria")
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
			if err := a.validateURFFile(filePath); err != nil {
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

// discoverFiles finds URF files based on the input patterns
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
		} else if a.isURFFile(input) {
			files = append(files, input)
		}
	}

	return files, nil
}

// discoverInDirectory finds URF files in a directory
func (a *ArmService) discoverInDirectory(dir string, recursive bool, _, _ []string) ([]string, error) {
	var files []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && a.isURFFile(path) {
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
				if a.isURFFile(path) {
					files = append(files, path)
				}
			}
		}
	}

	return files, nil
}

// isURFFile checks if a file is a URF file by extension
func (a *ArmService) isURFFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	return ext == ".yaml" || ext == ".yml"
}

// validateURFFile validates a URF file
func (a *ArmService) validateURFFile(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	file := &types.File{
		Path:    filePath,
		Content: content,
		Size:    int64(len(content)),
	}

	parser := urf.NewParser()
	_, err = parser.Parse(file)
	return err
}

// compileFile compiles a single URF file to multiple targets
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
		compiler, err := urf.NewCompiler(urf.CompileTarget(target))
		if err != nil {
			return 0, fmt.Errorf("failed to create compiler for target %s: %w", target, err)
		}

		compiledFiles, err := compiler.Compile(namespace, file)
		if err != nil {
			return 0, fmt.Errorf("failed to compile for target %s: %w", target, err)
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
