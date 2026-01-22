package storage

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

// RepoInterface defines git repository operations
type RepoInterface interface {
	GetTags(ctx context.Context, url string) ([]string, error)
	GetBranches(ctx context.Context, url string) ([]string, error)
	GetBranchHeadCommitHash(ctx context.Context, url, branch string) (string, error)
	GetTagCommitHash(ctx context.Context, url, tag string) (string, error)
	GetFilesFromCommit(ctx context.Context, url, commit string) ([]*core.File, error)
}

// Repo implements git operations using system git commands with cross-process locking
type Repo struct {
	repoDir string
	lock    *FileLock // Protects git operations
}

// NewRepo creates new repo instance
func NewRepo(repoDir string) RepoInterface {
	return &Repo{
		repoDir: repoDir,
		lock:    NewFileLock(repoDir),
	}
}

// GetTags returns all git tags
func (r *Repo) GetTags(ctx context.Context, url string) ([]string, error) {
	if err := r.lock.Lock(ctx); err != nil {
		return nil, err
	}
	defer r.lock.Unlock()
	
	if err := r.ensureCloned(ctx, url); err != nil {
		return nil, err
	}
	
	cmd := exec.Command("git", "tag", "-l")
	cmd.Dir = r.repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	return lines, nil
}

// GetBranches returns all git branches
func (r *Repo) GetBranches(ctx context.Context, url string) ([]string, error) {
	if err := r.lock.Lock(ctx); err != nil {
		return nil, err
	}
	defer r.lock.Unlock()
	
	if err := r.ensureCloned(ctx, url); err != nil {
		return nil, err
	}
	
	cmd := exec.Command("git", "branch", "-r", "--format=%(refname:short)")
	cmd.Dir = r.repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var branches []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip origin/HEAD pointer
		if line != "" && !strings.HasPrefix(line, "origin/HEAD") {
			// Remove origin/ prefix
			branch := strings.TrimPrefix(line, "origin/")
			branches = append(branches, branch)
		}
	}
	return branches, nil
}

// GetBranchHeadCommitHash returns commit hash for branch head
func (r *Repo) GetBranchHeadCommitHash(ctx context.Context, url, branch string) (string, error) {
	if err := r.lock.Lock(ctx); err != nil {
		return "", err
	}
	defer r.lock.Unlock()
	
	if err := r.ensureCloned(ctx, url); err != nil {
		return "", err
	}
	
	cmd := exec.Command("git", "rev-parse", "origin/"+branch)
	cmd.Dir = r.repoDir
	output, err := cmd.Output()
	if err != nil {
		// Try without origin/ prefix for local branches
		cmd = exec.Command("git", "rev-parse", branch)
		cmd.Dir = r.repoDir
		output, err = cmd.Output()
		if err != nil {
			return "", err
		}
	}
	
	return strings.TrimSpace(string(output)), nil
}

// GetTagCommitHash returns commit hash for tag
func (r *Repo) GetTagCommitHash(ctx context.Context, url, tag string) (string, error) {
	if err := r.lock.Lock(ctx); err != nil {
		return "", err
	}
	defer r.lock.Unlock()
	
	if err := r.ensureCloned(ctx, url); err != nil {
		return "", err
	}
	
	cmd := exec.Command("git", "rev-parse", tag)
	cmd.Dir = r.repoDir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	
	return strings.TrimSpace(string(output)), nil
}

// GetFilesFromCommit returns all files from specific commit
func (r *Repo) GetFilesFromCommit(ctx context.Context, url, commit string) ([]*core.File, error) {
	if err := r.lock.Lock(ctx); err != nil {
		return nil, err
	}
	defer r.lock.Unlock()
	
	if err := r.ensureCloned(ctx, url); err != nil {
		return nil, err
	}
	
	// Get list of files in commit
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", commit)
	cmd.Dir = r.repoDir
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	
	filePaths := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(filePaths) == 1 && filePaths[0] == "" {
		return []*core.File{}, nil
	}
	
	var files []*core.File
	for _, path := range filePaths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}
		
		// Get file content from commit
		cmd := exec.Command("git", "show", commit+":"+path)
		cmd.Dir = r.repoDir
		content, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		
		files = append(files, &core.File{
			Path:    path,
			Content: content,
			Size:    int64(len(content)),
		})
	}
	
	return files, nil
}

// ensureCloned clones repo if not exists, fetches if exists
func (r *Repo) ensureCloned(ctx context.Context, url string) error {
	// Check if repo directory exists and has .git
	if _, err := os.Stat(filepath.Join(r.repoDir, ".git")); err == nil {
		// Repo exists, fetch updates
		cmd := exec.Command("git", "fetch", "--all", "--tags")
		cmd.Dir = r.repoDir
		return cmd.Run()
	}
	
	// Repo doesn't exist, clone it
	cmd := exec.Command("git", "clone", url, r.repoDir)
	return cmd.Run()
}