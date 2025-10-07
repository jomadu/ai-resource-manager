package installer

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/resource"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// processResourceFiles detects resource files, compiles them, and replaces them with compiled output
func processResourceFiles(files []types.File, registry, ruleset, version string, compiler resource.Compiler) ([]types.File, error) {
	parser := resource.NewParser()

	namespace := fmt.Sprintf("%s/%s@%s", registry, ruleset, version)
	var processedFiles []types.File

	// Process each file
	for _, file := range files {
		if parser.IsResource(&file) {
			// Try to compile as ruleset first, then promptset
			var compiledFiles []*types.File
			var err error

			// Try ruleset compilation
			if _, parseErr := parser.ParseRuleset(&file); parseErr == nil {
				compiledFiles, err = compiler.CompileRuleset(namespace, &file)
			} else if _, parseErr := parser.ParsePromptset(&file); parseErr == nil {
				compiledFiles, err = compiler.CompilePromptset(namespace, &file)
			} else {
				return nil, fmt.Errorf("failed to parse resource file %s as either ruleset or promptset", file.Path)
			}

			if err != nil {
				return nil, fmt.Errorf("failed to compile resource file %s: %w", file.Path, err)
			}

			// Add compiled files
			for _, compiledFile := range compiledFiles {
				processedFiles = append(processedFiles, *compiledFile)
			}
		} else {
			// Keep non-resource files as-is
			processedFiles = append(processedFiles, file)
		}
	}

	return processedFiles, nil
}
