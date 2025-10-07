package index

import (
	"fmt"
	"sort"

	"github.com/jomadu/ai-rules-manager/internal/resource"
)

type IndexGenerator interface {
	CreateRuleset(data *IndexData) *resource.Ruleset
	GenerateBody(data *IndexData) string
}

type DefaultIndexGenerator struct{}

func (g *DefaultIndexGenerator) CreateRuleset(data *IndexData) *resource.Ruleset {
	return &resource.Ruleset{
		APIVersion: "v1",
		Kind:       "Ruleset",
		Metadata: resource.Metadata{
			ID:   "arm",
			Name: "ARM Rulesets Index",
		},
		Spec: resource.RulesetSpec{
			Rules: map[string]resource.Rule{
				"index": {
					Name:        "ARM Rulesets Index",
					Enforcement: "must",
					Priority:    1000,
					Body:        g.GenerateBody(data),
				},
			},
		},
	}
}

func (g *DefaultIndexGenerator) GenerateBody(data *IndexData) string {
	body := "# ARM Rulesets\n\n"
	body += "This file defines the installation priorities for rulesets managed by ARM.\n\n"
	body += "## Priority Rules\n\n"
	body += "**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:\n\n"
	body += "1. **Higher priority numbers take precedence** over lower priority numbers\n"
	body += "2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets\n"
	body += "3. **Always consult this index** to resolve any ambiguity about which rules to follow\n\n"
	body += "## Installed Rulesets\n\n"

	// Collect all rulesets with their priorities
	type rulesetEntry struct {
		registry string
		name     string
		info     RulesetInfo
	}

	var entries []rulesetEntry
	for registry, rulesets := range data.Rulesets {
		for name, info := range rulesets {
			entries = append(entries, rulesetEntry{
				registry: registry,
				name:     name,
				info:     info,
			})
		}
	}

	// Sort by priority (high to low)
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].info.Priority > entries[j].info.Priority
	})

	// Generate output in priority order
	for _, entry := range entries {
		body += fmt.Sprintf("### %s/%s@%s\n", entry.registry, entry.name, entry.info.Version)
		body += fmt.Sprintf("- **Priority:** %d\n", entry.info.Priority)
		body += "- **Rules:**\n"
		for _, file := range entry.info.FilePaths {
			body += fmt.Sprintf("  - %s\n", file)
		}
		body += "\n"
	}

	return body
}
