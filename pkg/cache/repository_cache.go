package cache

import (
	"path/filepath"

	"github.com/jomadu/ai-rules-manager/pkg/registry"
)

// RepositoryCache manages cached Git repositories and operations
type RepositoryCache interface {
	// EnsureRepository clones or updates the repository, returns repo path
	EnsureRepository(url, registryKey string) (repoPath string, err error)
	// ListVersions returns available tags and branches from cached repo
	ListVersions(repoPath string) ([]registry.VersionRef, error)
	// GetFilesAtCommit extracts files at specific commit with pattern matching
	GetFilesAtCommit(repoPath, commit string, selector registry.ContentSelector) ([]registry.File, error)
}

// GitRepositoryCache implements RepositoryCache for Git repositories
type GitRepositoryCache struct {
	basePath string
}

func NewGitRepositoryCache(basePath string) *GitRepositoryCache {
	return &GitRepositoryCache{basePath: basePath}
}

func (g *GitRepositoryCache) EnsureRepository(url, registryKey string) (string, error) {
	repoPath := filepath.Join(g.basePath, "registries", registryKey, "repository")
	// TODO: implement git clone/fetch operations
	return repoPath, nil
}

func (g *GitRepositoryCache) ListVersions(repoPath string) ([]registry.VersionRef, error) {
	// TODO: implement git tag/branch listing
	return nil, nil
}

func (g *GitRepositoryCache) GetFilesAtCommit(repoPath, commit string, selector registry.ContentSelector) ([]registry.File, error) {
	// TODO: implement git checkout + file extraction
	return nil, nil
}