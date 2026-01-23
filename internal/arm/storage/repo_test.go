package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRepo(t *testing.T) {
	repo := NewRepo("/tmp/test-repo")
	assert.NotNil(t, repo)
}

func TestGetTags_RepoNotCloned(t *testing.T) {
	// Create source repo with tags
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit").
		Tag("v1.0.0").
		AddFile("feature.txt", "new feature").
		Commit("Add feature").
		Tag("v1.1.0")
	err := builder.Build()
	require.NoError(t, err)

	// Create target repo for cloning
	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	// Test GetTags clones and returns tags
	ctx := context.Background()
	tags, err := repo.GetTags(ctx, sourceDir)

	// TODO: implement and change assertion
	assert.NoError(t, err)
	assert.Equal(t, []string{"v1.0.0", "v1.1.0"}, tags)
}

func TestGetTags_NoTags(t *testing.T) {
	// Create source repo with no tags
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	tags, err := repo.GetTags(ctx, sourceDir)

	assert.NoError(t, err)
	assert.Empty(t, tags)
}

func TestGetBranches_RepoNotCloned(t *testing.T) {
	// Create source repo with branches
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit").
		Branch("feature").
		AddFile("feature.txt", "feature content").
		Commit("Add feature").
		Checkout("main")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	branches, err := repo.GetBranches(ctx, sourceDir)

	assert.NoError(t, err)
	assert.Contains(t, branches, "main")
	assert.Contains(t, branches, "feature")
}

func TestGetBranchHeadCommitHash_ValidBranch(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	hash, err := repo.GetBranchHeadCommitHash(ctx, sourceDir, "main")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 40) // git commit hash is 40 chars
}

func TestGetBranchHeadCommitHash_InvalidBranch(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	hash, err := repo.GetBranchHeadCommitHash(ctx, sourceDir, "nonexistent")

	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestGetTagCommitHash_ValidTag(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit").
		Tag("v1.0.0")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	hash, err := repo.GetTagCommitHash(ctx, sourceDir, "v1.0.0")

	assert.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Len(t, hash, 40)
}

func TestGetTagCommitHash_InvalidTag(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	hash, err := repo.GetTagCommitHash(ctx, sourceDir, "nonexistent")

	assert.Error(t, err)
	assert.Empty(t, hash)
}

func TestGetFilesFromCommit_ValidCommit(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		AddFile("src/main.go", "package main").
		Commit("Initial commit").
		Tag("v1.0.0")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	// Get commit hash first
	ctx := context.Background()
	hash, err := repo.GetTagCommitHash(ctx, sourceDir, "v1.0.0")
	require.NoError(t, err)

	// Get files from commit
	files, err := repo.GetFilesFromCommit(ctx, sourceDir, hash)

	assert.NoError(t, err)
	assert.Len(t, files, 2)

	// Check files exist
	fileNames := make([]string, len(files))
	for i, f := range files {
		fileNames[i] = f.Path
	}
	assert.Contains(t, fileNames, "README.md")
	assert.Contains(t, fileNames, "src/main.go")
}

func TestGetFilesFromCommit_InvalidCommit(t *testing.T) {
	sourceDir := t.TempDir()
	testRepo := NewTestRepo(t, sourceDir)

	builder := testRepo.Builder().
		Init().
		AddFile("README.md", "# Test Repo").
		Commit("Initial commit")
	err := builder.Build()
	require.NoError(t, err)

	targetDir := t.TempDir()
	repo := NewRepo(targetDir)

	ctx := context.Background()
	files, err := repo.GetFilesFromCommit(ctx, sourceDir, "invalidhash")

	assert.Error(t, err)
	assert.Nil(t, files)
}
