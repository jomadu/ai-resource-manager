package arm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/index"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/registry"
	"github.com/jomadu/ai-rules-manager/internal/resource"
)

func (a *ArmService) CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error {

	cacheDir := cache.GetRegistriesDir()
	entries, err := os.ReadDir(cacheDir)
	if err != nil {
		if os.IsNotExist(err) {
			a.ui.Success("No cache directory found - nothing to clean")
			return nil
		}
		return fmt.Errorf("failed to read cache directory: %w", err)
	}

	cleanedCount := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		baseDir := filepath.Join(cacheDir, entry.Name())

		// Try to read registry metadata from existing index
		registryKeyObj, err := a.readRegistryMetadata(baseDir)
		if err != nil {
			// Remove corrupted registry and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			cleanedCount++
			continue
		}

		packageCache, err := cache.NewRegistryPackageCache(registryKeyObj)
		if err != nil {
			// Remove registry with invalid metadata and its lock
			registryKey := entry.Name()
			_ = os.RemoveAll(baseDir)
			_ = os.Remove(filepath.Join(cacheDir, ".locks", registryKey+".lock"))
			cleanedCount++
			continue
		}

		if err := packageCache.Cleanup(maxAge); err != nil {
			return fmt.Errorf("failed to clean cache for registry %s: %w", entry.Name(), err)
		}
		cleanedCount++
	}

	a.ui.Success(fmt.Sprintf("Cache cleaned: processed %d registries, removed versions older than %v", cleanedCount, maxAge))
	return nil
}

func (a *ArmService) NukeCache(ctx context.Context) error {
	cacheDir := cache.GetCacheDir()
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}
	a.ui.Success("Cache directory removed successfully")
	return nil
}

// readRegistryMetadata reads registry metadata from the index file
func (a *ArmService) readRegistryMetadata(baseDir string) (interface{}, error) {
	indexPath := filepath.Join(baseDir, "index.json")

	data, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, err
	}

	var index cache.RegistryIndex
	if err := json.Unmarshal(data, &index); err != nil {
		return nil, err
	}

	if index.RegistryMetadata == nil {
		return nil, fmt.Errorf("missing registry metadata")
	}

	return index.RegistryMetadata, nil
}

// CleanSinks performs selective cleanup of sink directories based on ARM index
func (a *ArmService) CleanSinks(ctx context.Context) error {
	// Get all configured sinks
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	cleanedCount := 0
	for sinkName, sinkConfig := range sinks {
		cleaned, err := a.cleanSinkDirectory(ctx, sinkConfig)
		if err != nil {
			return fmt.Errorf("failed to clean sink %s: %w", sinkName, err)
		}
		cleanedCount += cleaned
	}

	a.ui.Success(fmt.Sprintf("Sink cleanup completed: processed %d sinks, removed %d orphaned files", len(sinks), cleanedCount))
	return nil
}

// NukeSinks removes all sink directories and their contents
func (a *ArmService) NukeSinks(ctx context.Context) error {
	// Get all configured sinks
	sinks, err := a.manifestManager.GetSinks(ctx)
	if err != nil {
		return fmt.Errorf("failed to get sinks: %w", err)
	}

	removedCount := 0
	for sinkName, sinkConfig := range sinks {
		// Remove the entire sink directory
		err := os.RemoveAll(sinkConfig.Directory)
		if err != nil {
			return fmt.Errorf("failed to remove sink directory %s: %w", sinkName, err)
		}
		removedCount++
	}

	a.ui.Success(fmt.Sprintf("Sink nuke completed: removed %d sink directories", removedCount))
	return nil
}

// cleanSinkDirectory performs selective cleanup of a sink directory based on the ARM index
func (a *ArmService) cleanSinkDirectory(ctx context.Context, sinkConfig manifest.SinkConfig) (int, error) {
	sinkDir := sinkConfig.Directory

	// Check if sink directory exists
	if _, err := os.Stat(sinkDir); os.IsNotExist(err) {
		return 0, nil // Sink directory doesn't exist, nothing to clean
	}

	// Load the ARM index for this sink
	indexManager := index.NewIndexManager(sinkDir, sinkConfig.GetLayout(), sinkConfig.CompileTarget)

	// Get all files that should exist according to the index
	allFiles, err := indexManager.GetAllInstalledFiles(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get expected files from index: %w", err)
	}

	// Convert to map for easy lookup
	expectedFiles := make(map[string]bool)
	for _, filePath := range allFiles {
		expectedFiles[filePath] = true
	}

	// Find and remove orphaned files
	var orphanedFiles []string
	err = filepath.Walk(sinkDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if this file is expected
		if !expectedFiles[path] {
			orphanedFiles = append(orphanedFiles, path)
		}

		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to find orphaned files: %w", err)
	}

	// Remove orphaned files
	removedCount := 0
	for _, filePath := range orphanedFiles {
		if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
			return removedCount, fmt.Errorf("failed to remove orphaned file %s: %w", filePath, err)
		}
		removedCount++
	}

	// Clean up empty directories
	if err := a.cleanupEmptyDirectories(sinkDir, sinkConfig.GetLayout()); err != nil {
		return removedCount, fmt.Errorf("failed to cleanup empty directories: %w", err)
	}

	return removedCount, nil
}

// cleanupEmptyDirectories removes empty directories after file cleanup
func (a *ArmService) cleanupEmptyDirectories(sinkDir, layout string) error {
	if layout == "flat" {
		// For flat layout, we don't need to clean up directories since files are in the root
		return nil
	}

	// For hierarchical layout, clean up empty arm/ subdirectories
	armDir := filepath.Join(sinkDir, "arm")
	if _, err := os.Stat(armDir); os.IsNotExist(err) {
		return nil // arm directory doesn't exist
	}

	// Walk up from the deepest level and remove empty directories
	return filepath.Walk(armDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != armDir {
			// Check if directory is empty
			entries, err := os.ReadDir(path)
			if err != nil {
				return err
			}

			if len(entries) == 0 {
				// Directory is empty, remove it
				if err := os.Remove(path); err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// CompileFiles compiles resource files to target formats
func (a *ArmService) CompileFiles(ctx context.Context, req *CompileRequest) error {
	// 1. Use targets directly from request
	targets := req.Targets

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

// Shared helper methods for version resolution and registry management

// getLatestVersions gets the latest and wanted versions for a package from the registry
func (a *ArmService) getLatestVersions(ctx context.Context, registry, packageName string) (latest, wanted string, err error) {
	// Get registry instance
	registryInstance, err := a.getRegistryInstance(ctx, registry)
	if err != nil {
		return "", "", fmt.Errorf("failed to get registry %s: %w", registry, err)
	}

	// Get all available versions
	versions, err := registryInstance.ListVersions(ctx, packageName)
	if err != nil {
		return "", "", fmt.Errorf("failed to list versions for %s/%s: %w", registry, packageName, err)
	}

	if len(versions) == 0 {
		return "", "", fmt.Errorf("no versions found for %s/%s", registry, packageName)
	}

	// Get the latest version (first in sorted list)
	latest = versions[0].Version

	// For wanted version, we need to resolve the constraint from manifest
	// Try to get manifest entry (could be ruleset or promptset)
	var constraint string

	// Try ruleset first
	rulesetManifest, err := a.manifestManager.GetRuleset(ctx, registry, packageName)
	if err == nil {
		constraint = rulesetManifest.Version
	} else {
		// Try promptset
		promptsetManifest, err := a.manifestManager.GetPromptset(ctx, registry, packageName)
		if err == nil {
			constraint = promptsetManifest.Version
		} else {
			// Fail fast: package must be in manifest to determine wanted version
			return "", "", fmt.Errorf("package %s/%s not found in manifest - cannot determine wanted version", registry, packageName)
		}
	}

	// Resolve the constraint to get wanted version
	resolved, err := registryInstance.ResolveVersion(ctx, packageName, constraint)
	if err != nil {
		// Fail fast: constraint resolution failure indicates configuration error
		return "", "", fmt.Errorf("failed to resolve constraint %s for %s/%s: %w", constraint, registry, packageName, err)
	}

	wanted = resolved.Version.Version
	return latest, wanted, nil
}

// getRegistryInstance creates a registry instance for the given registry name
func (a *ArmService) getRegistryInstance(ctx context.Context, registryName string) (registry.Registry, error) {
	// Get registry config from manifest
	registries, err := a.manifestManager.GetRegistries(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get registries: %w", err)
	}

	registryConfig, exists := registries[registryName]
	if !exists {
		return nil, fmt.Errorf("registry %s not found in manifest", registryName)
	}

	// Create registry instance
	registryInstance, err := registry.NewRegistry(registryName, registryConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create registry %s: %w", registryName, err)
	}

	return registryInstance, nil
}

// expandVersionShorthand expands version shorthand notation to full semantic version constraints
func expandVersionShorthand(constraint string) string {
	// Match pure major version (e.g., "1")
	if matched, _ := regexp.MatchString(`^\d+$`, constraint); matched {
		return "^" + constraint + ".0.0"
	}
	// Match major.minor version (e.g., "1.0")
	if matched, _ := regexp.MatchString(`^\d+\.\d+$`, constraint); matched {
		return "~" + constraint + ".0"
	}
	return constraint
}
