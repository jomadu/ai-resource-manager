package compiler

import (
	"fmt"
	"sort"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

// CompileRuleset compiles a ruleset to the tool format
func CompileRuleset(tool Tool, namespace string, ruleset *resource.RulesetResource) ([]*core.File, error) {
	var compiledFiles []*core.File

	// Sort rule IDs for consistent ordering
	ruleIDs := make([]string, 0, len(ruleset.Spec.Rules))
	for ruleID := range ruleset.Spec.Rules {
		ruleIDs = append(ruleIDs, ruleID)
	}
	sort.Strings(ruleIDs)

	// Create rule filename generator
	ruleFilenameFactory := NewRuleFilenameGeneratorFactory()
	ruleFilenameGen, err := ruleFilenameFactory.NewRuleFilenameGenerator(tool)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule filename generator: %w", err)
	}

	// Create rule generator
	ruleFactory := NewRuleGeneratorFactory()
	ruleGen, err := ruleFactory.NewRuleGenerator(tool)
	if err != nil {
		return nil, fmt.Errorf("failed to create rule generator: %w", err)
	}

	for _, ruleID := range ruleIDs {
		// Generate filename
		filename, err := ruleFilenameGen.GenerateRuleFilename(ruleset.Metadata.ID, ruleID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate filename for rule %s: %w", ruleID, err)
		}
		
		// Generate content
		content, err := ruleGen.GenerateRule(namespace, ruleset, ruleID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate content for rule %s: %w", ruleID, err)
		}
		
		compiledFiles = append(compiledFiles, &core.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}

	return compiledFiles, nil
}

// CompilePromptset compiles a promptset to the tool format
func CompilePromptset(tool Tool, namespace string, promptset *resource.PromptsetResource) ([]*core.File, error) {
	var compiledFiles []*core.File

	// Sort prompt IDs for consistent ordering
	promptIDs := make([]string, 0, len(promptset.Spec.Prompts))
	for promptID := range promptset.Spec.Prompts {
		promptIDs = append(promptIDs, promptID)
	}
	sort.Strings(promptIDs)

	// Create prompt filename generator
	promptFilenameFactory := NewPromptFilenameGeneratorFactory()
	promptFilenameGen, err := promptFilenameFactory.NewPromptFilenameGenerator(tool)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt filename generator: %w", err)
	}

	// Create prompt generator
	promptFactory := NewPromptGeneratorFactory()
	promptGen, err := promptFactory.NewPromptGenerator(tool)
	if err != nil {
		return nil, fmt.Errorf("failed to create prompt generator: %w", err)
	}

	for _, promptID := range promptIDs {
		// Generate filename
		filename, err := promptFilenameGen.GeneratePromptFilename(promptset.Metadata.ID, promptID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate filename for prompt %s: %w", promptID, err)
		}
		
		// Generate content
		content, err := promptGen.GeneratePrompt(namespace, promptset, promptID)
		if err != nil {
			return nil, fmt.Errorf("failed to generate content for prompt %s: %w", promptID, err)
		}
		
		compiledFiles = append(compiledFiles, &core.File{
			Path:    filename,
			Content: []byte(content),
			Size:    int64(len(content)),
		})
	}

	return compiledFiles, nil
}