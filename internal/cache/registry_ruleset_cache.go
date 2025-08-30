package cache

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// RegistryRulesetCache provides registry-scoped storage for cached rulesets.
type RegistryRulesetCache interface {
	ListVersions(ctx context.Context, rulesetKey string) ([]string, error)
	GetRulesetVersion(ctx context.Context, rulesetKey, version string) ([]types.File, error)
	SetRulesetVersion(ctx context.Context, rulesetKey, version string, files []types.File) error
	InvalidateRuleset(ctx context.Context, rulesetKey string) error
	InvalidateVersion(ctx context.Context, rulesetKey, version string) error
}
