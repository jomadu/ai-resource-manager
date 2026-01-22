package resource

import (
	"fmt"
	"strings"
)

// MarkdownRuleGenerator generates markdown rule files
type MarkdownRuleGenerator struct {
	metadataGen RuleMetadataGenerator
}

// GenerateRule generates a markdown rule file
func (g *MarkdownRuleGenerator) GenerateRule(namespace string, ruleset *Ruleset, ruleID string, rule *Rule) string {
	var content strings.Builder

	// Resource metadata block (Amazon Q doesn't use tool-specific frontmatter)
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
