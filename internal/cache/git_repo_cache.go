package cache

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GitRepoCache provides git repository operations.
type GitRepoCache interface {
	GetTags(ctx context.Context) ([]string, error)
	GetBranches(ctx context.Context) ([]string, error)
	GetCommitHash(ctx context.Context, ref string) (string, error)
	GetFiles(ctx context.Context, ref string, selector types.ContentSelector) ([]types.File, error)
}
