package storage

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestRepoBuilder helps construct test git repositories
type TestRepoBuilder interface {
	Init() TestRepoBuilder
	Branch(name string) TestRepoBuilder
	Checkout(branch string) TestRepoBuilder
	AddFile(path, content string) TestRepoBuilder
	RemoveFile(path string) TestRepoBuilder
	Commit(message string) TestRepoBuilder
	Tag(name string) TestRepoBuilder
	Merge(branch string) TestRepoBuilder
	Build() error
}

// TestRepo creates and manages test git repository
type TestRepo struct {
	repoDir string
	t       *testing.T
}

// NewTestRepo creates new test repo helper
func NewTestRepo(t *testing.T, repoDir string) *TestRepo {
	return &TestRepo{
		repoDir: repoDir,
		t:       t,
	}
}

// Builder returns builder for constructing git repository
func (tr *TestRepo) Builder() TestRepoBuilder {
	return &testRepoBuilder{
		repoDir: tr.repoDir,
		t:       tr.t,
	}
}

// GetRepoDir returns the repository directory path
func (tr *TestRepo) GetRepoDir() string {
	return tr.repoDir
}

// testRepoBuilder implements the builder pattern for git operations
type testRepoBuilder struct {
	repoDir string
	t       *testing.T
}

func (b *testRepoBuilder) Init() TestRepoBuilder {
	err := os.MkdirAll(b.repoDir, 0o755)
	require.NoError(b.t, err)

	cmd := exec.Command("git", "init", "--initial-branch=main")
	cmd.Dir = b.repoDir
	err = cmd.Run()
	require.NoError(b.t, err)

	// Set test user config
	b.runGitCmd("config", "user.name", "Test User")
	b.runGitCmd("config", "user.email", "test@example.com")

	return b
}

func (b *testRepoBuilder) Branch(name string) TestRepoBuilder {
	b.runGitCmd("checkout", "-b", name)
	return b
}

func (b *testRepoBuilder) Checkout(branch string) TestRepoBuilder {
	b.runGitCmd("checkout", branch)
	return b
}

func (b *testRepoBuilder) AddFile(path, content string) TestRepoBuilder {
	fullPath := filepath.Join(b.repoDir, path)
	err := os.MkdirAll(filepath.Dir(fullPath), 0o755)
	require.NoError(b.t, err)

	err = os.WriteFile(fullPath, []byte(content), 0o644)
	require.NoError(b.t, err)

	b.runGitCmd("add", path)
	return b
}

func (b *testRepoBuilder) RemoveFile(path string) TestRepoBuilder {
	b.runGitCmd("rm", path)
	return b
}

func (b *testRepoBuilder) Commit(message string) TestRepoBuilder {
	b.runGitCmd("commit", "-m", message)
	return b
}

func (b *testRepoBuilder) Tag(name string) TestRepoBuilder {
	b.runGitCmd("tag", name)
	return b
}

func (b *testRepoBuilder) Merge(branch string) TestRepoBuilder {
	b.runGitCmd("merge", branch)
	return b
}

func (b *testRepoBuilder) Build() error {
	// Nothing to do - all operations are executed immediately
	return nil
}

func (b *testRepoBuilder) runGitCmd(args ...string) {
	cmd := exec.Command("git", args...)
	cmd.Dir = b.repoDir
	err := cmd.Run()
	require.NoError(b.t, err)
}
