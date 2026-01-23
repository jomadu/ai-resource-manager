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

	parts = append(parts,
		"---",
		fmt.Sprintf("namespace: %s", namespace),
		"ruleset:",
		fmt.Sprintf("  id: %s", ruleset.Metadata.ID),
		fmt.Sprintf("  name: %s", ruleset.Metadata.Name),
		"  rules:",
	)

	// Add all rule IDs in sorted order
	ruleIDs := make([]string, 0, len(ruleset.Spec.Rules))
	for id := range ruleset.Spec.Rules {
		ruleIDs = append(ruleIDs, id)
	}
	sort.Strings(ruleIDs)

	for _, id := range ruleIDs {
		parts = append(parts, fmt.Sprintf("    - %s", id))
	}

	parts = append(parts,
		"rule:",
		fmt.Sprintf("  id: %s", ruleID),
		fmt.Sprintf("  name: %s", rule.Name),
		fmt.Sprintf("  enforcement: %s", strings.ToUpper(rule.Enforcement)),
	)

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
						parts[len(parts)-1] += fmt.Sprintf(`%q`, file)
					} else {
						parts[len(parts)-1] += fmt.Sprintf(`, %q`, file)
					}
				}
				parts[len(parts)-1] += "]"
			}
		}
	}

	parts = append(parts, "---")

	return strings.Join(parts, "\n")
}
