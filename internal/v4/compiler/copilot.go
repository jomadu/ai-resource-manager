package compiler

import (
	"fmt"

	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

// CopilotRuleGenerator generates copilot rule content
type CopilotRuleGenerator struct{}

func (g *CopilotRuleGenerator) GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) (string, error) {
	rule, exists := ruleset.Spec.Rules[ruleID]
	if !exists {
		return "", fmt.Errorf("rule %s not found in ruleset", ruleID)
	}

	// Copilot format: metadata + content (same as markdown)
	metadata := GenerateRuleMetadata(namespace, ruleset, ruleID, &rule)
	return metadata + "\n\n" + rule.Body, nil
}

// CopilotPromptGenerator generates copilot prompt content
type CopilotPromptGenerator struct{}

func (g *CopilotPromptGenerator) GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) (string, error) {
	prompt, exists := promptset.Spec.Prompts[promptID]
	if !exists {
		return "", fmt.Errorf("prompt %s not found in promptset", promptID)
	}

	// Prompts are just content, no metadata
	return prompt.Body, nil
}

// CopilotRuleFilenameGenerator generates copilot rule filenames
type CopilotRuleFilenameGenerator struct{}

func (g *CopilotRuleFilenameGenerator) GenerateRuleFilename(rulesetID, ruleID string) (string, error) {
	if rulesetID == "" {
		return "", fmt.Errorf("rulesetID cannot be empty")
	}
	if ruleID == "" {
		return "", fmt.Errorf("ruleID cannot be empty")
	}

	return rulesetID + "_" + ruleID + ".instructions.md", nil
}

// CopilotPromptFilenameGenerator generates copilot prompt filenames
type CopilotPromptFilenameGenerator struct{}

func (g *CopilotPromptFilenameGenerator) GeneratePromptFilename(promptsetID, promptID string) (string, error) {
	if promptsetID == "" {
		return "", fmt.Errorf("promptsetID cannot be empty")
	}
	if promptID == "" {
		return "", fmt.Errorf("promptID cannot be empty")
	}

	return promptsetID + "_" + promptID + ".md", nil
}