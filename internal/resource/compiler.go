package resource

import (
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// DefaultCompiler compiles resource files using generators
type DefaultCompiler struct {
	ruleGen     RuleGenerator
	promptGen   PromptGenerator
	filenameGen FilenameGenerator
}

// NewCompiler creates a new compiler for the specified target
func NewCompiler(target CompileTarget) (Compiler, error) {
	ruleFactory := NewRuleGeneratorFactory()
	ruleGen, err := ruleFactory.NewRuleGenerator(target)
	if err != nil {
		return nil, err
	}

	promptFactory := NewPromptGeneratorFactory()
	promptGen, err := promptFactory.NewPromptGenerator(target)
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
		promptGen:   promptGen,
		filenameGen: filenameGen,
	}, nil
}

// CompileRuleset compiles a ruleset to the target format
func (c *DefaultCompiler) CompileRuleset(namespace string, ruleset *Ruleset) ([]*types.File, error) {
	var compiledFiles []*types.File
	for ruleID, rule := range ruleset.Spec.Rules {
		filename := c.filenameGen.GenerateFilename(ruleset.Metadata.ID, ruleID)
		content := c.ruleGen.GenerateRule(namespace, ruleset, ruleID, &rule)
		compiledFiles = append(compiledFiles, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}

	return compiledFiles, nil
}

// CompilePromptset compiles a promptset to the target format
func (c *DefaultCompiler) CompilePromptset(namespace string, promptset *Promptset) ([]*types.File, error) {
	var compiledFiles []*types.File
	for promptID, prompt := range promptset.Spec.Prompts {
		filename := c.filenameGen.GenerateFilename(promptset.Metadata.ID, promptID)
		content := c.promptGen.GeneratePrompt(namespace, promptset, promptID, &prompt)
		compiledFiles = append(compiledFiles, &types.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}

	return compiledFiles, nil
}
