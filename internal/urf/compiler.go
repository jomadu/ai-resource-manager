package urf

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// DefaultCompiler compiles URF files using a rule generator
type DefaultCompiler struct {
	parser      Parser
	ruleGen     RuleGenerator
	filenameGen FilenameGenerator
}

// NewCompiler creates a new compiler for the specified target
func NewCompiler(target CompileTarget) (Compiler, error) {
	ruleFactory := NewRuleGeneratorFactory()
	ruleGen, err := ruleFactory.NewRuleGenerator(target)
	if err != nil {
		return nil, err
	}

	filenameFactory := NewFilenameGeneratorFactory()
	filenameGen, err := filenameFactory.NewFilenameGenerator(target)
	if err != nil {
		return nil, err
	}

	return &DefaultCompiler{
		parser:      NewParser(),
		ruleGen:     ruleGen,
		filenameGen: filenameGen,
	}, nil
}

// Compile compiles a single URF file to the target format
func (c *DefaultCompiler) Compile(namespace string, file *types.File) ([]*types.File, error) {
	if !c.parser.IsURF(file) {
		return nil, fmt.Errorf("file %s is not a valid URF file", file.Path)
	}

	ruleset, err := c.parser.Parse(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URF file %s: %w", file.Path, err)
	}

	var compiledFiles []*types.File
	for _, rule := range ruleset.Rules {
		filename := c.filenameGen.GenerateFilename(ruleset.Metadata.ID, rule.ID)
		content := c.ruleGen.GenerateRule(namespace, ruleset, &rule)
		compiledFiles = append(compiledFiles, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}

	return compiledFiles, nil
}
