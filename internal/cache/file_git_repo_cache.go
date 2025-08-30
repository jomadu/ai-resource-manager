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
	repoDir     string
	url         string
	initialized bool
}

// NewGitRepoCache creates a new git repository cache.
func NewGitRepoCache(registryKey, repoName, url string) *FileGitRepoCache {
	homeDir, _ := os.UserHomeDir()
	repoDir := filepath.Join(homeDir, ".arm", "cache", "registries", registryKey, "repository", repoName)

	return &FileGitRepoCache{
		repoDir: repoDir,
		url:     url,
	}
}

func (g *FileGitRepoCache) ensureInitialized(ctx context.Context) error {
	if g.initialized {
		return nil
	}

	if _, err := os.Stat(filepath.Join(g.repoDir, ".git")); os.IsNotExist(err) {
		// Clone if repo doesn't exist
		if err := os.MkdirAll(filepath.Dir(g.repoDir), 0o755); err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, "git", "clone", g.url, g.repoDir)
		if err := cmd.Run(); err != nil {
			return err
		}
	} else {
		// Fetch if repo exists
		cmd := exec.CommandContext(ctx, "git", "fetch", "--all", "--tags")
		cmd.Dir = g.repoDir
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	g.initialized = true
	return nil
}

func (g *FileGitRepoCache) GetTags(ctx context.Context) ([]string, error) {
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "git", "tag", "-l")
	cmd.Dir = g.repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	tags := strings.Fields(string(output))
	return tags, nil
}

func (g *FileGitRepoCache) GetBranches(ctx context.Context) ([]string, error) {
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}
	cmd := exec.CommandContext(ctx, "git", "branch", "-r", "--format=%(refname:short)")
	cmd.Dir = g.repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	branches := strings.Fields(string(output))
	for i, branch := range branches {
		branches[i] = strings.TrimPrefix(branch, "origin/")
	}
	return branches, nil
}

func (g *FileGitRepoCache) GetFiles(ctx context.Context, ref string, selector types.ContentSelector) ([]types.File, error) {
	if err := g.ensureInitialized(ctx); err != nil {
		return nil, err
	}

	// Checkout the specified ref
	cmd := exec.CommandContext(ctx, "git", "checkout", ref)
	cmd.Dir = g.repoDir
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	var files []types.File
	err := filepath.WalkDir(g.repoDir, func(path string, d fs.DirEntry, err error) error {
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

	return files, err
}
