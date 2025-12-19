package compiler

import (
	"fmt"

	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

// AmazonQRuleGenerator generates amazonq rule content
type AmazonQRuleGenerator struct{}

func (g *AmazonQRuleGenerator) GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) (string, error) {
	rule, exists := ruleset.Spec.Rules[ruleID]
	if !exists {
		return "", fmt.Errorf("rule %s not found in ruleset", ruleID)
	}

	// AmazonQ format: metadata + content (same as markdown)
	metadata := GenerateRuleMetadata(namespace, ruleset, ruleID, &rule)
	return metadata + "\n\n" + rule.Body, nil
}

// AmazonQPromptGenerator generates amazonq prompt content
type AmazonQPromptGenerator struct{}

func (g *AmazonQPromptGenerator) GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) (string, error) {
	prompt, exists := promptset.Spec.Prompts[promptID]
	if !exists {
		return "", fmt.Errorf("prompt %s not found in promptset", promptID)
	}

	// Prompts are just content, no metadata
	return prompt.Body, nil
}

// AmazonQRuleFilenameGenerator generates amazonq rule filenames
type AmazonQRuleFilenameGenerator struct{}

func (g *AmazonQRuleFilenameGenerator) GenerateRuleFilename(rulesetID, ruleID string) (string, error) {
	if rulesetID == "" {
		return "", fmt.Errorf("rulesetID cannot be empty")
	}
	if ruleID == "" {
		return "", fmt.Errorf("ruleID cannot be empty")
	}

	return rulesetID + "_" + ruleID + ".md", nil
}

// AmazonQPromptFilenameGenerator generates amazonq prompt filenames
type AmazonQPromptFilenameGenerator struct{}

func (g *AmazonQPromptFilenameGenerator) GeneratePromptFilename(promptsetID, promptID string) (string, error) {
	if promptsetID == "" {
		return "", fmt.Errorf("promptsetID cannot be empty")
	}
	if promptID == "" {
		return "", fmt.Errorf("promptID cannot be empty")
	}

	return promptsetID + "_" + promptID + ".md", nil
}