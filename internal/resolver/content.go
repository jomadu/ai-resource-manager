package resolver

import (
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// ContentResolver applies include/exclude patterns to filter ruleset files.
type ContentResolver interface {
	ResolveContent(selector arm.ContentSelector, files []arm.File) ([]arm.File, error)
}

// GitContentResolver implements glob pattern matching for file filtering.
type GitContentResolver struct{}

// NewGitContentResolver creates a new Git-based content resolver.
func NewGitContentResolver() *GitContentResolver {
	return &GitContentResolver{}
}

func (g *GitContentResolver) ResolveContent(selector arm.ContentSelector, files []arm.File) ([]arm.File, error) {
	return nil, errors.New("not implemented")
}