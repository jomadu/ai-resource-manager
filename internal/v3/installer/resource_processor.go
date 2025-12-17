package installer

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// compileResourceFiles detects resource files, compiles them, and replaces them with compiled output
func compileResourceFiles(files []types.File, registry, resourceName, version string, compiler resource.Compiler) ([]types.File, error) {
	parser := resource.NewParser()

	namespace := fmt.Sprintf("%s/%s@%s", registry, resourceName, version)
	var processedFiles []types.File

	// Process each file
	for _, file := range files {
		switch {
		case parser.IsRuleset(&file):
			// Parse and compile as ruleset
			ruleset, err := parser.ParseRuleset(&file)
			if err != nil {
				return nil, fmt.Errorf("failed to parse ruleset file %s: %w", file.Path, err)
			}

			compiledFiles, err := compiler.CompileRuleset(namespace, ruleset)
			if err != nil {
				return nil, fmt.Errorf("failed to compile ruleset file %s: %w", file.Path, err)
			}

			// Add compiled files
			for _, compiledFile := range compiledFiles {
				processedFiles = append(processedFiles, *compiledFile)
			}
		case parser.IsPromptset(&file):
			// Parse and compile as promptset
			promptset, err := parser.ParsePromptset(&file)
			if err != nil {
				return nil, fmt.Errorf("failed to parse promptset file %s: %w", file.Path, err)
			}

			compiledFiles, err := compiler.CompilePromptset(namespace, promptset)
			if err != nil {
				return nil, fmt.Errorf("failed to compile promptset file %s: %w", file.Path, err)
			}

			// Add compiled files
			for _, compiledFile := range compiledFiles {
				processedFiles = append(processedFiles, *compiledFile)
			}
		default:
			// Keep non-resource files as-is
			processedFiles = append(processedFiles, file)
		}
	}

	return processedFiles, nil
}
