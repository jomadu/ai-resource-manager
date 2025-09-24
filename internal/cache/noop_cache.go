package cache

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// NoopRegistryRulesetCache is a no-op implementation of RegistryRulesetCache for testing
type NoopRegistryRulesetCache struct{}

// NewNoopRegistryRulesetCache creates a new no-op cache
func NewNoopRegistryRulesetCache() *NoopRegistryRulesetCache {
	return &NoopRegistryRulesetCache{}
}

func (n *NoopRegistryRulesetCache) ListVersions(ctx context.Context, keyObj interface{}) ([]string, error) {
	return []string{}, nil
}

func (n *NoopRegistryRulesetCache) GetRulesetVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error) {
	// Always return cache miss to force fresh fetches
	return nil, context.DeadlineExceeded // Standard cache miss error
}

func (n *NoopRegistryRulesetCache) SetRulesetVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error {
	// Do nothing - don't cache
	return nil
}

func (n *NoopRegistryRulesetCache) InvalidateRuleset(ctx context.Context, rulesetKey string) error {
	return nil
}

func (n *NoopRegistryRulesetCache) InvalidateVersion(ctx context.Context, rulesetKey, version string) error {
	return nil
}
