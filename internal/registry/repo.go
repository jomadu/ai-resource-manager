package registry

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// GitRepo implements Git repository operations.
type GitRepo struct {
	workDir string
}

// NewGitRepo creates a new Git repository handler.
func NewGitRepo(workDir string) *GitRepo {
	return &GitRepo{workDir: workDir}
}

func (r *GitRepo) Clone(ctx context.Context, url string) error {
	if err := os.MkdirAll(filepath.Dir(r.workDir), 0o755); err != nil {
		return err
	}
	cmd := exec.CommandContext(ctx, "git", "clone", url, r.workDir)
	return cmd.Run()
}

func (r *GitRepo) Fetch(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "fetch", "--all", "--tags")
	cmd.Dir = r.workDir
	return cmd.Run()
}

func (r *GitRepo) Pull(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, "git", "pull")
	cmd.Dir = r.workDir
	return cmd.Run()
}

func (r *GitRepo) GetTags(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "tag", "-l")
	cmd.Dir = r.workDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	tags := strings.Fields(string(output))
	return tags, nil
}

func (r *GitRepo) GetBranches(ctx context.Context) ([]string, error) {
	cmd := exec.CommandContext(ctx, "git", "branch", "-r", "--format=%(refname:short)")
	cmd.Dir = r.workDir
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

func (r *GitRepo) Checkout(ctx context.Context, ref string) error {
	cmd := exec.CommandContext(ctx, "git", "checkout", ref)
	cmd.Dir = r.workDir
	return cmd.Run()
}

func (r *GitRepo) GetFiles(ctx context.Context, selector arm.ContentSelector) ([]arm.File, error) {
	var files []arm.File

	err := filepath.WalkDir(r.workDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || strings.Contains(path, ".git") {
			return nil
		}

		relPath, err := filepath.Rel(r.workDir, path)
		if err != nil {
			return err
		}

		if r.matchesSelector(relPath, selector) {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			info, err := d.Info()
			if err != nil {
				return err
			}

			files = append(files, arm.File{
				Path:    relPath,
				Content: content,
				Size:    info.Size(),
			})
		}

		return nil
	})

	return files, err
}

func (r *GitRepo) matchesSelector(path string, selector arm.ContentSelector) bool {
	if len(selector.Include) == 0 {
		return true
	}

	for _, pattern := range selector.Include {
		if matched, _ := filepath.Match(pattern, path); matched {
			for _, exclude := range selector.Exclude {
				if matched, _ := filepath.Match(exclude, path); matched {
					return false
				}
			}
			return true
		}
	}

	return false
}
