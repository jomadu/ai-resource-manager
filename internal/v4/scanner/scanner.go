package scanner

import (
	"os"
	"path/filepath"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/filetype"
	"github.com/jomadu/ai-resource-manager/internal/v4/parser"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

// Config holds directory scanning parameters
type Config struct {
	Dirs      []string
	Recursive bool
	Include   []string
	Exclude   []string
}

// ScanResult holds the results of scanning directories for resources
type ScanResult struct {
	Rulesets   []*resource.RulesetResource
	Promptsets []*resource.PromptsetResource
}

// ScanResources finds and parses all resource files in directories
func ScanResources(config Config) (*ScanResult, error) {
	result := &ScanResult{
		Rulesets:   []*resource.RulesetResource{},
		Promptsets: []*resource.PromptsetResource{},
	}

	for _, dir := range config.Dirs {
		err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
			// Skip if directory doesn't exist or permission denied
			if err != nil {
				return nil
			}

			// Skip directories
			if d.IsDir() {
				// If not recursive, skip subdirectories
				if !config.Recursive && path != dir {
					return filepath.SkipDir
				}
				return nil
			}

			// Apply include/exclude patterns
			if !matchesPatterns(path, config.Include, config.Exclude) {
				return nil
			}

			// Try to parse as ruleset
			if filetype.IsRulesetFile(path) {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil // Skip files we can't read
				}

				file := &core.File{Path: path, Content: content}
				ruleset, err := parser.ParseRuleset(file)
				if err != nil {
					return nil // Skip invalid files
				}

				result.Rulesets = append(result.Rulesets, ruleset)
				return nil
			}

			// Try to parse as promptset
			if filetype.IsPromptsetFile(path) {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil // Skip files we can't read
				}

				file := &core.File{Path: path, Content: content}
				promptset, err := parser.ParsePromptset(file)
				if err != nil {
					return nil // Skip invalid files
				}

				result.Promptsets = append(result.Promptsets, promptset)
			}

			return nil
		})

		// Ignore walk errors (directory doesn't exist, etc)
		if err != nil {
			continue
		}
	}

	return result, nil
}

// matchesPatterns checks if file path matches include/exclude patterns
func matchesPatterns(path string, include, exclude []string) bool {
	// If include patterns specified, file must match at least one
	if len(include) > 0 {
		matched := false
		for _, pattern := range include {
			if m, _ := filepath.Match(pattern, filepath.Base(path)); m {
				matched = true
				break
			}
		}
		if !matched {
			return false
		}
	}

	// If exclude patterns specified, file must not match any
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
	}

	return true
}
