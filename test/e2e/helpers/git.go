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
	repo.run("init")
	repo.run("config", "user.email", "test@example.com")
	repo.run("config", "user.name", "Test User")

	return repo
}

// WriteFile writes a file to the repository
func (r *GitRepo) WriteFile(path, content string) {
	r.t.Helper()
	fullPath := filepath.Join(r.Path, path)
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		r.t.Fatalf("failed to create directory %s: %v", dir, err)
	}
	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		r.t.Fatalf("failed to write file %s: %v", fullPath, err)
	}
}

// Commit commits all changes with the given message
func (r *GitRepo) Commit(message string) {
	r.t.Helper()
	r.run("add", ".")
	r.run("commit", "-m", message)
}

// Tag creates a tag at the current commit
func (r *GitRepo) Tag(tag string) {
	r.t.Helper()
	r.run("tag", tag)
}

// Branch creates a new branch
func (r *GitRepo) Branch(name string) {
	r.t.Helper()
	r.run("branch", name)
}

// Checkout checks out a branch
func (r *GitRepo) Checkout(ref string) {
	r.t.Helper()
	r.run("checkout", ref)
}

// run executes a git command in the repository directory
func (r *GitRepo) run(args ...string) {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Path
	if output, err := cmd.CombinedOutput(); err != nil {
		r.t.Fatalf("git command failed: %v\nOutput: %s\nError: %v", args, output, err)
	}
}
