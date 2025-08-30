package resolver

import (
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// ContentResolver applies include/exclude patterns to filter ruleset files.
type ContentResolver interface {
	ResolveContent(selector types.ContentSelector, files []types.File) ([]types.File, error)
}

// GitContentResolver implements glob pattern matching for file filtering.
type GitContentResolver struct{}

// NewGitContentResolver creates a new Git-based content resolver.
func NewGitContentResolver() *GitContentResolver {
	return &GitContentResolver{}
}

func (g *GitContentResolver) ResolveContent(selector types.ContentSelector, files []types.File) ([]types.File, error) {
	return nil, errors.New("not implemented")
}
