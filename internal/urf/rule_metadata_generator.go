package urf

import (
	"fmt"
	"strings"
)

// DefaultRuleMetadataGenerator implements shared metadata generation
type DefaultRuleMetadataGenerator struct{}

// NewRuleMetadataGenerator creates a new metadata generator
func NewRuleMetadataGenerator() RuleMetadataGenerator {
	return &DefaultRuleMetadataGenerator{}
}

// GenerateRuleMetadata generates the shared metadata block
func (g *DefaultRuleMetadataGenerator) GenerateRuleMetadata(namespace string, ruleset *Ruleset, ruleID string, rule *Rule) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("namespace: %s\n", namespace))
	content.WriteString("ruleset:\n")
	content.WriteString(fmt.Sprintf("  id: %s\n", ruleset.Metadata.ID))
	content.WriteString(fmt.Sprintf("  name: %s\n", ruleset.Metadata.Name))
	if ruleset.Metadata.Version != "" {
		content.WriteString(fmt.Sprintf("  version: %s\n", ruleset.Metadata.Version))
	}
	content.WriteString("  rules:\n")
	for id := range ruleset.Rules {
		content.WriteString(fmt.Sprintf("    - %s\n", id))
	}
	content.WriteString("rule:\n")
	content.WriteString(fmt.Sprintf("  id: %s\n", ruleID))
	content.WriteString(fmt.Sprintf("  name: %s\n", rule.Name))
	if rule.Enforcement != "" {
		content.WriteString(fmt.Sprintf("  enforcement: %s\n", strings.ToUpper(rule.Enforcement)))
	}
	if rule.Priority != 0 {
		content.WriteString(fmt.Sprintf("  priority: %d\n", rule.Priority))
	}
	if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
		content.WriteString("  scope:\n")
		content.WriteString("    - files: [")
		for i, file := range rule.Scope[0].Files {
			if i > 0 {
				content.WriteString(", ")
			}
			content.WriteString(fmt.Sprintf("%q", file))
		}
		content.WriteString("]\n")
	}
	content.WriteString("---\n\n")

	return content.String()
}
