package cache

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestGitKeyGen_RegistryKey(t *testing.T) {
	keyGen := NewGitKeyGen()

	tests := []struct {
		name         string
		url          string
		registryType string
		want         string
	}{
		{
			name:         "github_git_registry",
			url:          "https://github.com/my-user/ai-rules",
			registryType: "git",
			want:         "sha256(\"https://github.com/my-user/ai-rules\" + \"git\")",
		},
		{
			name:         "normalized_url_with_trailing_slash",
			url:          "https://github.com/my-user/ai-rules/",
			registryType: "git",
			want:         "sha256(\"https://github.com/my-user/ai-rules\" + \"git\")",
		},
		{
			name:         "different_registry_type",
			url:          "https://github.com/my-user/ai-rules",
			registryType: "http",
			want:         "sha256(\"https://github.com/my-user/ai-rules\" + \"http\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyGen.RegistryKey(tt.url, tt.registryType)
			if got == "" {
				t.Error("RegistryKey() returned empty string")
			}
			// TODO: Verify actual SHA256 hash when implementation is complete
		})
	}
}

func TestGitKeyGen_RulesetKey(t *testing.T) {
	keyGen := NewGitKeyGen()

	tests := []struct {
		name     string
		selector types.ContentSelector
		want     string
	}{
		{
			name: "amazonq_rules_selector",
			selector: types.ContentSelector{
				Include: []string{"rules/amazonq/*.md"},
			},
			want: "sha256(\"rules/amazonq/*.md\")",
		},
		{
			name: "cursor_rules_selector",
			selector: types.ContentSelector{
				Include: []string{"rules/cursor/*.mdc"},
			},
			want: "sha256(\"rules/cursor/*.mdc\")",
		},
		{
			name: "multiple_includes",
			selector: types.ContentSelector{
				Include: []string{"rules/amazonq/*.md", "rules/cursor/*.mdc"},
			},
			want: "sha256(\"rules/amazonq/*.md,rules/cursor/*.mdc\")",
		},
		{
			name: "with_excludes",
			selector: types.ContentSelector{
				Include: []string{"rules/**/*.md"},
				Exclude: []string{"rules/test/*.md"},
			},
			want: "sha256(\"rules/**/*.md\" + \"rules/test/*.md\")",
		},
		{
			name:     "empty_selector",
			selector: types.ContentSelector{},
			want:     "sha256(\"\")",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := keyGen.RulesetKey(tt.selector)
			if got == "" {
				t.Error("RulesetKey() returned empty string")
			}
			// TODO: Verify actual SHA256 hash when implementation is complete
		})
	}
}

func TestGitKeyGen_ConsistentKeys(t *testing.T) {
	keyGen := NewGitKeyGen()

	// Test that same inputs produce same keys
	url := "https://github.com/my-user/ai-rules"
	registryType := "git"

	key1 := keyGen.RegistryKey(url, registryType)
	key2 := keyGen.RegistryKey(url, registryType)

	if key1 != key2 {
		t.Errorf("RegistryKey() not consistent: %s != %s", key1, key2)
	}

	selector := types.ContentSelector{
		Include: []string{"rules/amazonq/*.md"},
	}

	rulesetKey1 := keyGen.RulesetKey(selector)
	rulesetKey2 := keyGen.RulesetKey(selector)

	if rulesetKey1 != rulesetKey2 {
		t.Errorf("RulesetKey() not consistent: %s != %s", rulesetKey1, rulesetKey2)
	}
}

func TestGitKeyGen_DifferentInputsDifferentKeys(t *testing.T) {
	keyGen := NewGitKeyGen()

	// Different URLs should produce different keys
	key1 := keyGen.RegistryKey("https://github.com/user1/repo", "git")
	key2 := keyGen.RegistryKey("https://github.com/user2/repo", "git")

	if key1 == key2 {
		t.Error("Different URLs produced same registry key")
	}

	// Different selectors should produce different keys
	selector1 := types.ContentSelector{Include: []string{"rules/amazonq/*.md"}}
	selector2 := types.ContentSelector{Include: []string{"rules/cursor/*.mdc"}}

	rulesetKey1 := keyGen.RulesetKey(selector1)
	rulesetKey2 := keyGen.RulesetKey(selector2)

	if rulesetKey1 == rulesetKey2 {
		t.Error("Different selectors produced same ruleset key")
	}
}

func TestGitKeyGen_Normalization(t *testing.T) {
	keyGen := NewGitKeyGen()

	// URL normalization - trailing slash
	key1 := keyGen.RegistryKey("https://github.com/user/repo", "git")
	key2 := keyGen.RegistryKey("https://github.com/user/repo/", "git")
	if key1 != key2 {
		t.Error("URLs with/without trailing slash should produce same key")
	}

	// URL normalization - case
	key3 := keyGen.RegistryKey("https://github.com/User/Repo", "git")
	key4 := keyGen.RegistryKey("https://github.com/user/repo", "git")
	if key3 != key4 {
		t.Error("URLs with different case should produce same key")
	}

	// Registry type normalization
	key5 := keyGen.RegistryKey("https://github.com/user/repo", "GIT")
	key6 := keyGen.RegistryKey("https://github.com/user/repo", "git")
	if key5 != key6 {
		t.Error("Registry type should be case-insensitive")
	}

	// Include order normalization
	selector1 := types.ContentSelector{Include: []string{"rules/amazonq/*.md", "rules/cursor/*.mdc"}}
	selector2 := types.ContentSelector{Include: []string{"rules/cursor/*.mdc", "rules/amazonq/*.md"}}
	rulesetKey1 := keyGen.RulesetKey(selector1)
	rulesetKey2 := keyGen.RulesetKey(selector2)
	if rulesetKey1 != rulesetKey2 {
		t.Error("Include patterns in different order should produce same key")
	}

	// Exclude order normalization
	selector3 := types.ContentSelector{Include: []string{"rules/*.md"}, Exclude: []string{"rules/test/*.md", "rules/tmp/*.md"}}
	selector4 := types.ContentSelector{Include: []string{"rules/*.md"}, Exclude: []string{"rules/tmp/*.md", "rules/test/*.md"}}
	rulesetKey3 := keyGen.RulesetKey(selector3)
	rulesetKey4 := keyGen.RulesetKey(selector4)
	if rulesetKey3 != rulesetKey4 {
		t.Error("Exclude patterns in different order should produce same key")
	}
}
