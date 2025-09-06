package cache

import (
	"context"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GitRepoCache provides git repository operations.
type GitRepoCache interface {
	GetTags(ctx context.Context) ([]string, error)
	GetBranches(ctx context.Context) ([]string, error)
	GetBranchHeadCommitHash(ctx context.Context, branch string) (string, error)
	GetTagCommitHash(ctx context.Context, tag string) (string, error)
	GetFilesFromCommit(ctx context.Context, commit string, selector types.ContentSelector) ([]types.File, error)
}
