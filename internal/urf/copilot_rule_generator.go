package urf

import (
	"fmt"
	"strings"
)

// CopilotRuleGenerator generates GitHub Copilot-compatible rule files
type CopilotRuleGenerator struct {
	metadataGen RuleMetadataGenerator
}

// GenerateRule generates a GitHub Copilot rule file
func (g *CopilotRuleGenerator) GenerateRule(namespace string, ruleset *Ruleset, rule *Rule) string {
	var content strings.Builder

	// GitHub Copilot frontmatter with applyTo
	content.WriteString("---\n")
	if len(rule.Scope) > 0 {
		var patterns []string
		for _, scope := range rule.Scope {
			patterns = append(patterns, scope.Files...)
		}
		content.WriteString(fmt.Sprintf("applyTo: %q\n", strings.Join(patterns, ",")))
	} else {
		content.WriteString("applyTo: \"**\"\n")
	}
	content.WriteString("---\n\n")

	// URF metadata block
	content.WriteString(g.metadataGen.GenerateRuleMetadata(namespace, ruleset, rule))

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
