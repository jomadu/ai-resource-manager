package registry

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Registry provides version-controlled access to resource repositories.
// It supports both rulesets and promptsets through generic package operations.
type Registry interface {
	// ListVersions returns all available versions for a given package name.
	// The package name can be either a ruleset or promptset name.
	ListVersions(ctx context.Context, packageName string) ([]types.Version, error)

	// ResolveVersion resolves a version constraint to a specific version.
	// Works for both rulesets and promptsets.
	ResolveVersion(ctx context.Context, packageName string, constraint string) (*resolver.ResolvedVersion, error)

	// GetContent retrieves the content of a specific package version.
	// The content can contain both rulesets and promptsets.
	GetContent(ctx context.Context, packageName string, version types.Version, selector types.ContentSelector) ([]types.File, error)
}
