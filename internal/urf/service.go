package urf

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Service orchestrates URF processing operations
type Service struct {
	parser          Parser
	compilerFactory CompilerFactory
}

// NewService creates a new URF service
func NewService() *Service {
	return &Service{
		parser:          NewParser(),
		compilerFactory: NewDefaultCompilerFactory(),
	}
}

// ProcessFiles processes files for URF content
func (s *Service) ProcessFiles(files []*types.File) ([]*Ruleset, error) {
	var parsedFiles []*Ruleset
	for _, file := range files {
		if !s.parser.IsURF(file) {
			continue
		}

		ruleset, err := s.parser.Parse(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URF file %s: %w", file.Path, err)
		}

		parsedFiles = append(parsedFiles, ruleset)
	}

	return parsedFiles, nil
}

// CompileFiles compiles URF files to target format
func (s *Service) CompileFiles(rulesets []*Ruleset, target CompileTarget, namespace string) ([]*types.File, error) {
	compiler, err := s.compilerFactory.GetCompiler(target)
	if err != nil {
		return nil, err
	}

	var allFiles []*types.File
	for _, ruleset := range rulesets {
		files, err := compiler.Compile(namespace, ruleset)
		if err != nil {
			return nil, fmt.Errorf("failed to compile URF %s: %w", ruleset.Metadata.ID, err)
		}
		allFiles = append(allFiles, files...)
	}
	return allFiles, nil
}

// GetSupportedTargets returns supported compilation targets
func (s *Service) GetSupportedTargets() []CompileTarget {
	return s.compilerFactory.SupportedTargets()
}
