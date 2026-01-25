package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// GitRepo represents a test Git repository
type GitRepo struct {
	Path string
	t    *testing.T
}

// NewGitRepo creates a new test Git repository
func NewGitRepo(t *testing.T, dir string) *GitRepo {
	t.Helper()
	
	repo := &GitRepo{
		Path: dir,
		t:    t,
	}
	
	// Initialize Git repo
	repo.run("git", "init")
	repo.run("git", "config", "user.email", "test@example.com")
	repo.run("git", "config", "user.name", "Test User")
	
	return repo
}

// WriteFile writes a file to the repository
func (r *GitRepo) WriteFile(path, content string) {
	r.t.Helper()
	fullPath := filepath.Join(r.Path, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		r.t.Fatalf("failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		r.t.Fatalf("failed to write file %s: %v", fullPath, err)
	}
}

// Commit commits all changes with the given message
func (r *GitRepo) Commit(message string) {
	r.t.Helper()
	r.run("git", "add", ".")
	r.run("git", "commit", "-m", message)
}

// Tag creates a tag at the current commit
func (r *GitRepo) Tag(tag string) {
	r.t.Helper()
	r.run("git", "tag", tag)
}

// Branch creates a new branch
func (r *GitRepo) Branch(name string) {
	r.t.Helper()
	r.run("git", "branch", name)
}

// Checkout checks out a branch
func (r *GitRepo) Checkout(ref string) {
	r.t.Helper()
	r.run("git", "checkout", ref)
}

// run executes a command in the repository directory
func (r *GitRepo) run(name string, args ...string) {
	r.t.Helper()
	cmd := exec.Command(name, args...)
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("command failed: %s %v\nOutput: %s\nError: %v", name, args, output, err)
	}
}
