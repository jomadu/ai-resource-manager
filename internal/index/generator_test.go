package index

import (
	"strings"
	"testing"
)

func TestDefaultIndexGenerator_CreateRuleset(t *testing.T) {
	generator := &DefaultIndexGenerator{}
	data := &IndexData{
		Rulesets: map[string]map[string]RulesetInfo{
			"ai-rules": {
				"security": {
					Version:   "1.0.0",
					Priority:  100,
					FilePaths: []string{"./security/rule1.md"},
				},
			},
		},
		Files: map[string]FileInfo{
			"./security/rule1.md": {Registry: "ai-rules", Ruleset: "security"},
		},
	}

	ruleset := generator.CreateRuleset(data)

	if ruleset.Metadata.ID != "arm" {
		t.Errorf("expected ID 'arm', got %s", ruleset.Metadata.ID)
	}
	if len(ruleset.Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(ruleset.Rules))
	}
}

func TestDefaultIndexGenerator_GenerateBody(t *testing.T) {
	generator := &DefaultIndexGenerator{}
	data := &IndexData{
		Rulesets: map[string]map[string]RulesetInfo{
			"ai-rules": {
				"security": {
					Version:   "1.0.0",
					Priority:  100,
					FilePaths: []string{"./security/rule1.md", "./security/rule2.md"},
				},
				"style": {
					Version:   "2.0.0",
					Priority:  50,
					FilePaths: []string{"./style/rule1.md"},
				},
			},
		},
	}

	body := generator.GenerateBody(data)

	if !strings.Contains(body, "# ARM Rulesets") {
		t.Error("expected body to contain header")
	}
	if !strings.Contains(body, "ai-rules/security@1.0.0") {
		t.Error("expected body to contain security ruleset")
	}
	if !strings.Contains(body, "ai-rules/style@2.0.0") {
		t.Error("expected body to contain style ruleset")
	}
	if !strings.Contains(body, "**Priority:** 100") {
		t.Error("expected body to contain priority 100")
	}
	if !strings.Contains(body, "**Priority:** 50") {
		t.Error("expected body to contain priority 50")
	}
	if !strings.Contains(body, "./security/rule1.md") {
		t.Error("expected body to contain file path")
	}
}

func TestDefaultIndexGenerator_GenerateBody_Empty(t *testing.T) {
	generator := &DefaultIndexGenerator{}
	data := &IndexData{
		Rulesets: make(map[string]map[string]RulesetInfo),
		Files:    make(map[string]FileInfo),
	}

	body := generator.GenerateBody(data)

	if !strings.Contains(body, "# ARM Rulesets") {
		t.Error("expected body to contain header even when empty")
	}
	if !strings.Contains(body, "## Installed Rulesets") {
		t.Error("expected body to contain section header")
	}
}
