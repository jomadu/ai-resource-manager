package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// Registry provides version-controlled access to ruleset repositories.
type Registry interface {
	ListVersions(ctx context.Context) ([]arm.VersionRef, error)
	GetContent(ctx context.Context, version arm.VersionRef, selector arm.ContentSelector) ([]arm.File, error)
}

// Repository handles low-level Git operations for registry access.
type Repository interface {
	Clone(ctx context.Context, url string) error
	Fetch(ctx context.Context) error
	Pull(ctx context.Context) error
	GetTags(ctx context.Context) ([]string, error)
	GetBranches(ctx context.Context) ([]string, error)
	Checkout(ctx context.Context, ref string) error
	GetFiles(ctx context.Context, selector arm.ContentSelector) ([]arm.File, error)
}
