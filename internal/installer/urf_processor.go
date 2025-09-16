package installer

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/types"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

// processURFFiles detects URF files, compiles them, and replaces them with compiled output
func processURFFiles(files []types.File, registry, ruleset, version string, compiler urf.Compiler) ([]types.File, error) {
	parser := urf.NewParser()

	namespace := fmt.Sprintf("%s/%s@%s", registry, ruleset, version)
	var processedFiles []types.File

	// Process each file
	for _, file := range files {
		if parser.IsURF(&file) {
			// Compile URF file
			compiledFiles, err := compiler.Compile(namespace, &file)
			if err != nil {
				return nil, fmt.Errorf("failed to compile URF file %s: %w", file.Path, err)
			}

			// Add compiled files
			for _, compiledFile := range compiledFiles {
				processedFiles = append(processedFiles, *compiledFile)
			}
		} else {
			// Keep non-URF files as-is
			processedFiles = append(processedFiles, file)
		}
	}

	return processedFiles, nil
}
