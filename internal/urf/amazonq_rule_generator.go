package urf

import (
	"fmt"
	"strings"
)

// AmazonQRuleGenerator generates Amazon Q-compatible rule files
type AmazonQRuleGenerator struct {
	metadataGen RuleMetadataGenerator
}

// GenerateRule generates an Amazon Q rule file
func (g *AmazonQRuleGenerator) GenerateRule(namespace string, ruleset *Ruleset, rule *Rule) string {
	var content strings.Builder

	// URF metadata block (Amazon Q doesn't use tool-specific frontmatter)
	content.WriteString(g.metadataGen.GenerateRuleMetadata(namespace, ruleset, rule))

	// Rule title and body
	enforcement := strings.ToUpper(rule.Enforcement)
	content.WriteString(fmt.Sprintf("# %s (%s)\n\n", rule.Name, enforcement))
	content.WriteString(rule.Body)

	return content.String()
}
