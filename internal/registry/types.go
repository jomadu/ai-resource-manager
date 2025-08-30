package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Registry provides version-controlled access to ruleset repositories.
type Registry interface {
	ListVersions(ctx context.Context) ([]types.VersionRef, error)
	ResolveVersion(ctx context.Context, constraint string) (*types.VersionRef, error)
	GetContent(ctx context.Context, version types.VersionRef, selector types.ContentSelector) ([]types.File, error)
}
