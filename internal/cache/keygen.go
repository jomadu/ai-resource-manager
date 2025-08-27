package cache

import "github.com/jomadu/ai-rules-manager/internal/arm"

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
	return "" // TODO: implement SHA256 hash
}

func (d *GitKeyGen) RulesetKey(selector arm.ContentSelector) string {
	return "" // TODO: implement SHA256 hash
}