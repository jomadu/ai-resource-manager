package compiler

import (
	"fmt"
	"sort"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/arm/resource"
)

// GenerateRuleMetadata generates rule metadata (simple function since consistent across tools)
func GenerateRuleMetadata(namespace string, ruleset *resource.RulesetResource, ruleID string, rule *resource.Rule) string {
	var parts []string
	
	parts = append(parts, "---")
	parts = append(parts, fmt.Sprintf("namespace: %s", namespace))
	parts = append(parts, "ruleset:")
	parts = append(parts, fmt.Sprintf("  id: %s", ruleset.Metadata.ID))
	parts = append(parts, fmt.Sprintf("  name: %s", ruleset.Metadata.Name))
	parts = append(parts, "  rules:")
	
	// Add all rule IDs in sorted order
	ruleIDs := make([]string, 0, len(ruleset.Spec.Rules))
	for id := range ruleset.Spec.Rules {
		ruleIDs = append(ruleIDs, id)
	}
	sort.Strings(ruleIDs)
	
	for _, id := range ruleIDs {
		parts = append(parts, fmt.Sprintf("    - %s", id))
	}
	
	parts = append(parts, "rule:")
	parts = append(parts, fmt.Sprintf("  id: %s", ruleID))
	parts = append(parts, fmt.Sprintf("  name: %s", rule.Name))
	parts = append(parts, fmt.Sprintf("  enforcement: %s", strings.ToUpper(rule.Enforcement)))
	
	if rule.Priority > 0 {
		parts = append(parts, fmt.Sprintf("  priority: %d", rule.Priority))
	}
	
	if len(rule.Scope) > 0 {
		parts = append(parts, "  scope:")
		for _, scope := range rule.Scope {
			if len(scope.Files) > 0 {
				parts = append(parts, "    - files: [")
				for i, file := range scope.Files {
					if i == 0 {
						parts[len(parts)-1] += fmt.Sprintf(`"%s"`, file)
					} else {
						parts[len(parts)-1] += fmt.Sprintf(`, "%s"`, file)
					}
				}
				parts[len(parts)-1] += "]"
			}
		}
	}
	
	parts = append(parts, "---")
	
	return strings.Join(parts, "\n")
}