package compiler

import (
	"fmt"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/arm/resource"
)

// CursorRuleGenerator generates cursor rule content
type CursorRuleGenerator struct{}

func (g *CursorRuleGenerator) GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) (string, error) {
	rule, exists := ruleset.Spec.Rules[ruleID]
	if !exists {
		return "", fmt.Errorf("rule %s not found in ruleset", ruleID)
	}

	// Cursor format: frontmatter + metadata + content
	frontmatter := g.generateCursorFrontmatter(&rule)
	metadata := GenerateRuleMetadata(namespace, ruleset, ruleID, &rule)
	return frontmatter + "\n\n" + metadata + "\n\n" + rule.Body, nil
}

func (g *CursorRuleGenerator) generateCursorFrontmatter(rule *resource.Rule) string {
	var parts []string
	
	parts = append(parts, "---")
	
	if rule.Description != "" {
		parts = append(parts, fmt.Sprintf(`description: "%s"`, rule.Description))
	}
	
	// Add globs if scope is defined
	if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
		globs := strings.Join(rule.Scope[0].Files, ", ")
		parts = append(parts, fmt.Sprintf("globs: %s", globs))
	}
	
	// Add alwaysApply for must enforcement
	if rule.Enforcement == "must" {
		parts = append(parts, "alwaysApply: true")
	}
	
	parts = append(parts, "---")
	
	return strings.Join(parts, "\n")
}

// CursorPromptGenerator generates cursor prompt content
type CursorPromptGenerator struct{}

func (g *CursorPromptGenerator) GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) (string, error) {
	prompt, exists := promptset.Spec.Prompts[promptID]
	if !exists {
		return "", fmt.Errorf("prompt %s not found in promptset", promptID)
	}

	// Prompts are just content, no metadata
	return prompt.Body, nil
}

// CursorRuleFilenameGenerator generates cursor rule filenames
type CursorRuleFilenameGenerator struct{}

func (g *CursorRuleFilenameGenerator) GenerateRuleFilename(rulesetID, ruleID string) (string, error) {
	if rulesetID == "" {
		return "", fmt.Errorf("rulesetID cannot be empty")
	}
	if ruleID == "" {
		return "", fmt.Errorf("ruleID cannot be empty")
	}

	return rulesetID + "_" + ruleID + ".mdc", nil
}

// CursorPromptFilenameGenerator generates cursor prompt filenames
type CursorPromptFilenameGenerator struct{}

func (g *CursorPromptFilenameGenerator) GeneratePromptFilename(promptsetID, promptID string) (string, error) {
	if promptsetID == "" {
		return "", fmt.Errorf("promptsetID cannot be empty")
	}
	if promptID == "" {
		return "", fmt.Errorf("promptID cannot be empty")
	}

	return promptsetID + "_" + promptID + ".md", nil
}