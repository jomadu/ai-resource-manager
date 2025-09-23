package urf

import (
	"fmt"
	"strings"
)

// CursorRuleGenerator generates Cursor-compatible rule files
type CursorRuleGenerator struct {
	metadataGen RuleMetadataGenerator
}

// GenerateRule generates a Cursor rule file
func (g *CursorRuleGenerator) GenerateRule(namespace string, ruleset *Ruleset, ruleID string, rule *Rule) string {
	var content strings.Builder

	// Cursor-specific frontmatter
	content.WriteString("---\n")
	if rule.Description != "" {
		content.WriteString(fmt.Sprintf("description: %q\n", rule.Description))
	}

	if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
		content.WriteString("globs: [")
		for i, file := range rule.Scope[0].Files {
			if i > 0 {
				content.WriteString(", ")
			}
			content.WriteString(fmt.Sprintf("%q", file))
		}
		content.WriteString("]\n")
	} else {
		content.WriteString("globs: [\"**/*\"]\n")
	}

	if rule.Enforcement == "must" {
		content.WriteString("alwaysApply: true\n")
	}
	content.WriteString("---\n\n")

	// URF metadata block
	content.WriteString(g.metadataGen.GenerateRuleMetadata(namespace, ruleset, ruleID, rule))

	// Rule title and body
	if rule.Enforcement != "" {
		enforcement := strings.ToUpper(rule.Enforcement)
		content.WriteString(fmt.Sprintf("# %s (%s)\n\n", rule.Name, enforcement))
	} else {
		content.WriteString(fmt.Sprintf("# %s\n\n", rule.Name))
	}
	content.WriteString(rule.Body)

	return content.String()
}
