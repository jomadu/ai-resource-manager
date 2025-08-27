package cache

import (
	"crypto/sha256"
	"fmt"

	"github.com/jomadu/ai-rules-manager/pkg/registry"
)

// CacheKeyGenerator generates cache keys for different registry types
type CacheKeyGenerator interface {
	RegistryKey(url string) string
	RulesetKey(rulesetName string, selector registry.ContentSelector) string
	VersionKey(versionRef registry.VersionRef) string
}

// GitCacheKeyGenerator implements CacheKeyGenerator for Git registries
type GitCacheKeyGenerator struct{}

func NewGitCacheKeyGenerator() *GitCacheKeyGenerator {
	return &GitCacheKeyGenerator{}
}

func (g *GitCacheKeyGenerator) RegistryKey(url string) string {
	hash := sha256.Sum256([]byte(url + "git"))
	return fmt.Sprintf("%x", hash)
}

func (g *GitCacheKeyGenerator) RulesetKey(rulesetName string, selector registry.ContentSelector) string {
	hash := sha256.Sum256([]byte(rulesetName + selector.String()))
	return fmt.Sprintf("%x", hash)
}

func (g *GitCacheKeyGenerator) VersionKey(versionRef registry.VersionRef) string {
	return versionRef.ID // Use commit hash directly for Git
}