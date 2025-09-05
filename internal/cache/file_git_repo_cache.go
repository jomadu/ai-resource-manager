package cache

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// FileGitRepoCache implements filesystem-based git repository operations.
type FileGitRepoCache struct {
	registryDir string
	repoDir     string
	url         string
}

// NewGitRepoCache creates a new git repository cache.
func NewGitRepoCache(keyObj interface{}, repoName, url string) (*FileGitRepoCache, error) {
	registryKey, err := GenerateKey(keyObj)
	if err != nil {
		return nil, err
	}

	registriesDir := GetRegistriesDir()
	registryDir := filepath.Join(registriesDir, registryKey)
	repoDir := filepath.Join(registryDir, "repository", repoName)

	return &FileGitRepoCache{
		registryDir: registryDir,
		repoDir:     repoDir,
		url:         url,
	}, nil
}

func (g *FileGitRepoCache) ensureInitialized(ctx context.Context) error {
	if _, err := os.Stat(filepath.Join(g.repoDir, ".git")); os.IsNotExist(err) {
		// Clone if repo doesn't exist
		if err := os.MkdirAll(filepath.Dir(g.repoDir), 0o755); err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, "git", "clone", g.url, g.repoDir)
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	return nil
}

// ensureUpToDate fetches the latest changes from the remote repository
func (g *FileGitRepoCache) ensureUpToDate(ctx context.Context) error {
	if err := g.ensureInitialized(ctx); err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, "git", "fetch", "--all", "--tags", "--prune")
	cmd.Dir = g.repoDir
	return cmd.Run()
}

func (g *FileGitRepoCache) GetTags(ctx context.Context) ([]string, error) {
	registryKey := filepath.Base(g.registryDir)
	var tags []string

	err := WithRegistryLock(registryKey, func() error {
		if err := g.ensureUpToDate(ctx); err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, "git", "tag", "-l")
		cmd.Dir = g.repoDir
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		tags = strings.Fields(string(output))
		return nil
	})

	return tags, err
}

func (g *FileGitRepoCache) GetBranches(ctx context.Context) ([]string, error) {
	registryKey := filepath.Base(g.registryDir)
	var branches []string

	err := WithRegistryLock(registryKey, func() error {
		if err := g.ensureUpToDate(ctx); err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, "git", "branch", "-r", "--format=%(refname:short)")
		cmd.Dir = g.repoDir
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		branches = strings.Fields(string(output))
		for i, branch := range branches {
			branches[i] = strings.TrimPrefix(branch, "origin/")
		}
		return nil
	})

	return branches, err
}

// GetCommitHash returns the commit hash for a given ref (branch or tag).
func (g *FileGitRepoCache) GetCommitHash(ctx context.Context, ref string) (string, error) {
	registryKey := filepath.Base(g.registryDir)
	var hash string

	err := WithRegistryLock(registryKey, func() error {
		if err := g.ensureUpToDate(ctx); err != nil {
			return err
		}

		cmd := exec.CommandContext(ctx, "git", "rev-parse", ref)
		cmd.Dir = g.repoDir
		output, err := cmd.Output()
		if err != nil {
			return err
		}
		hash = strings.TrimSpace(string(output))
		return nil
	})

	return hash, err
}

func (g *FileGitRepoCache) GetFiles(ctx context.Context, ref string, selector types.ContentSelector) ([]types.File, error) {
	registryKey := filepath.Base(g.registryDir)
	var files []types.File
	var walkErr error

	err := WithRegistryLock(registryKey, func() error {
		if err := g.ensureUpToDate(ctx); err != nil {
			return err
		}

		// Checkout the specified ref
		cmd := exec.CommandContext(ctx, "git", "checkout", ref)
		cmd.Dir = g.repoDir
		if err := cmd.Run(); err != nil {
			return err
		}

		walkErr = filepath.WalkDir(g.repoDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() || strings.Contains(path, ".git") {
				return nil
			}

			relPath, err := filepath.Rel(g.repoDir, path)
			if err != nil {
				return err
			}

			if selector.Matches(relPath) {
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}

				info, err := d.Info()
				if err != nil {
					return err
				}

				files = append(files, types.File{
					Path:    relPath,
					Content: content,
					Size:    info.Size(),
				})
			}

			return nil
		})

		return walkErr
	})

	if err != nil {
		return nil, err
	}
	return files, walkErr
}
