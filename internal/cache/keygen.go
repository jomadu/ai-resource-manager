package cache

import (
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// KeyGenerator creates consistent hash keys for cache storage.
type KeyGenerator interface {
	RegistryKey(url, registryType string) string
	RulesetKey(selector arm.ContentSelector) string
}

// GitKeyGen implements SHA256-based key generation.
type GitKeyGen struct{}

// NewGitKeyGen creates a new Git-based key generator.
func NewGitKeyGen() *GitKeyGen {
	return &GitKeyGen{}
}

func (d *GitKeyGen) RegistryKey(url, registryType string) string {
	normalizedURL := normalizeURL(url)
	normalizedType := strings.ToLower(registryType)
	input := normalizedURL + normalizedType
	return sha256Hash(input)
}

func (d *GitKeyGen) RulesetKey(selector arm.ContentSelector) string {
	normalizedIncludes := normalizePatterns(selector.Include)
	normalizedExcludes := normalizePatterns(selector.Exclude)
	input := strings.Join(normalizedIncludes, ",")
	if len(normalizedExcludes) > 0 {
		input += "|" + strings.Join(normalizedExcludes, ",")
	}
	return sha256Hash(input)
}
