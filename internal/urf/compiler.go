package urf

import (
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// DefaultCompiler compiles URF using a rule generator
type DefaultCompiler struct {
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
		ruleGen:     ruleGen,
		filenameGen: filenameGen,
	}, nil
}

// Compile compiles URF to the target format
func (c *DefaultCompiler) Compile(namespace string, ruleset *Ruleset) ([]*types.File, error) {
	var files []*types.File
	for _, rule := range ruleset.Rules {
		filename := c.filenameGen.GenerateFilename(ruleset.Metadata.ID, rule.ID)
		content := c.ruleGen.GenerateRule(namespace, ruleset, &rule)
		files = append(files, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}
	return files, nil
}
