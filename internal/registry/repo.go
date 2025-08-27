package registry

import (
	"context"
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// GitRepo implements Git repository operations.
type GitRepo struct{}

// NewGitRepo creates a new Git repository handler.
func NewGitRepo() *GitRepo {
	return &GitRepo{}
}

func (r *GitRepo) Clone(ctx context.Context, url string) error {
	return errors.New("not implemented")
}

func (r *GitRepo) Fetch(ctx context.Context) error {
	return errors.New("not implemented")
}

func (r *GitRepo) Pull(ctx context.Context) error {
	return errors.New("not implemented")
}

func (r *GitRepo) GetTags(ctx context.Context) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (r *GitRepo) GetBranches(ctx context.Context) ([]string, error) {
	return nil, errors.New("not implemented")
}

func (r *GitRepo) Checkout(ctx context.Context, ref string) error {
	return errors.New("not implemented")
}

func (r *GitRepo) GetFiles(ctx context.Context, selector arm.ContentSelector) ([]arm.File, error) {
	return nil, errors.New("not implemented")
}