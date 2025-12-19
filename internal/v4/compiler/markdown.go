package compiler

import (
	"fmt"

	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

// MarkdownRuleGenerator generates markdown rule content
type MarkdownRuleGenerator struct{}

func (g *MarkdownRuleGenerator) GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) (string, error) {
	rule, exists := ruleset.Spec.Rules[ruleID]
	if !exists {
		return "", fmt.Errorf("rule %s not found in ruleset", ruleID)
	}

	// Markdown format: metadata + content
	metadata := GenerateRuleMetadata(namespace, ruleset, ruleID, &rule)
	return metadata + "\n\n" + rule.Body, nil
}

// MarkdownPromptGenerator generates markdown prompt content
type MarkdownPromptGenerator struct{}

func (g *MarkdownPromptGenerator) GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) (string, error) {
	prompt, exists := promptset.Spec.Prompts[promptID]
	if !exists {
		return "", fmt.Errorf("prompt %s not found in promptset", promptID)
	}

	// Prompts are just content, no metadata
	return prompt.Body, nil
}

// MarkdownRuleFilenameGenerator generates markdown rule filenames
type MarkdownRuleFilenameGenerator struct{}

func (g *MarkdownRuleFilenameGenerator) GenerateRuleFilename(rulesetID, ruleID string) (string, error) {
	if rulesetID == "" {
		return "", fmt.Errorf("rulesetID cannot be empty")
	}
	if ruleID == "" {
		return "", fmt.Errorf("ruleID cannot be empty")
	}

	return rulesetID + "_" + ruleID + ".md", nil
}

// MarkdownPromptFilenameGenerator generates markdown prompt filenames
type MarkdownPromptFilenameGenerator struct{}

func (g *MarkdownPromptFilenameGenerator) GeneratePromptFilename(promptsetID, promptID string) (string, error) {
	if promptsetID == "" {
		return "", fmt.Errorf("promptsetID cannot be empty")
	}
	if promptID == "" {
		return "", fmt.Errorf("promptID cannot be empty")
	}

	return promptsetID + "_" + promptID + ".md", nil
}