package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Registry provides version-controlled access to ruleset repositories.
type Registry interface {
	ListVersions(ctx context.Context, ruleset string) ([]types.Version, error)
	ResolveVersion(ctx context.Context, ruleset string, constraint string) (*resolver.ResolvedVersion, error)
	GetContent(ctx context.Context, ruleset string, version types.Version, selector types.ContentSelector) ([]types.File, error)
}
