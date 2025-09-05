package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Registry provides version-controlled access to ruleset repositories.
type Registry interface {
	ListVersions(ctx context.Context) ([]types.Version, error)
	ResolveVersion(ctx context.Context, constraint string) (*resolver.ResolvedVersion, error)
	GetContent(ctx context.Context, version types.Version, selector types.ContentSelector) ([]types.File, error)
}
