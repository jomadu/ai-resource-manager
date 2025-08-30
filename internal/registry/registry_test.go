package registry

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// Mock implementations for testing
type mockCache struct {
	versions map[string][]string
	content  map[string][]arm.File
	setErr   error
	getErr   error
}

func (m *mockCache) ListVersions(ctx context.Context, registryKey, rulesetKey string) ([]string, error) {
	key := registryKey + "/" + rulesetKey
	return m.versions[key], nil
}

func (m *mockCache) Get(ctx context.Context, registryKey, rulesetKey, version string) ([]arm.File, error) {
	if m.getErr != nil {
		return nil, m.getErr
	}
	key := registryKey + "/" + rulesetKey + "/" + version
	return m.content[key], nil
}

func (m *mockCache) Set(ctx context.Context, registryKey, rulesetKey, version string, files []arm.File) error {
	if m.setErr != nil {
		return m.setErr
	}
	key := registryKey + "/" + rulesetKey + "/" + version
	if m.content == nil {
		m.content = make(map[string][]arm.File)
	}
	m.content[key] = files
	return nil
}

func (m *mockCache) InvalidateRegistry(ctx context.Context, registryKey string) error { return nil }
func (m *mockCache) InvalidateRuleset(ctx context.Context, registryKey, rulesetKey string) error {
	return nil
}
func (m *mockCache) InvalidateVersion(ctx context.Context, registryKey, rulesetKey, version string) error {
	return nil
}

type mockRepository struct {
	tags     []string
	branches []string
	files    []arm.File
	err      error
}

func (m *mockRepository) Clone(ctx context.Context, url string) error    { return m.err }
func (m *mockRepository) Fetch(ctx context.Context) error                { return m.err }
func (m *mockRepository) Pull(ctx context.Context) error                 { return m.err }
func (m *mockRepository) Checkout(ctx context.Context, ref string) error { return m.err }

func (m *mockRepository) GetTags(ctx context.Context) ([]string, error) {
	return m.tags, m.err
}

func (m *mockRepository) GetBranches(ctx context.Context) ([]string, error) {
	return m.branches, m.err
}

func (m *mockRepository) GetFiles(ctx context.Context, selector arm.ContentSelector) ([]arm.File, error) {
	return m.files, m.err
}

type mockKeyGenerator struct {
	registryKey string
	rulesetKey  string
}

func (m *mockKeyGenerator) RegistryKey(url, registryType string) string {
	return m.registryKey
}

func (m *mockKeyGenerator) RulesetKey(selector arm.ContentSelector) string {
	return m.rulesetKey
}

func TestGitRegistry_ListVersions(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		branches []string
		repoErr  error
		want     []arm.VersionRef
		wantErr  bool
	}{
		{
			name:     "returns tags and branches as version refs",
			tags:     []string{"v1.0.0", "v1.1.0"},
			branches: []string{"main", "develop"},
			want: []arm.VersionRef{
				{ID: "v1.0.0", Type: arm.Tag},
				{ID: "v1.1.0", Type: arm.Tag},
				{ID: "main", Type: arm.Branch},
				{ID: "develop", Type: arm.Branch},
			},
		},
		{
			name:    "handles repository error",
			repoErr: errors.New("git error"),
			wantErr: true,
		},
		{
			name: "handles empty repository",
			want: []arm.VersionRef{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &mockCache{}
			repo := &mockRepository{
				tags:     tt.tags,
				branches: tt.branches,
				err:      tt.repoErr,
			}
			keyGen := &mockKeyGenerator{registryKey: "test-reg", rulesetKey: "test-ruleset"}

			g := NewGitRegistry(cache, repo, keyGen, "https://github.com/test/repo.git", "git")
			g.initialized = true
			got, err := g.ListVersions(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("ListVersions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalVersionRefs(got, tt.want) {
				t.Errorf("ListVersions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRegistry_GetContent(t *testing.T) {
	testFiles := []arm.File{
		{Path: "rules/test.md", Content: []byte("# Test Rule"), Size: 11},
		{Path: "config.json", Content: []byte(`{"test": true}`), Size: 15},
	}

	tests := []struct {
		name        string
		version     arm.VersionRef
		selector    arm.ContentSelector
		cachedFiles []arm.File
		repoFiles   []arm.File
		cacheErr    error
		repoErr     error
		want        []arm.File
		wantErr     bool
		expectCache bool
	}{
		{
			name:        "returns cached content when available",
			version:     arm.VersionRef{ID: "v1.0.0", Type: arm.Tag},
			selector:    arm.ContentSelector{Include: []string{"*.md"}},
			cachedFiles: testFiles,
			want:        testFiles,
			expectCache: false,
		},
		{
			name:        "fetches from repo when not cached",
			version:     arm.VersionRef{ID: "v1.0.0", Type: arm.Tag},
			selector:    arm.ContentSelector{Include: []string{"*.md"}},
			cacheErr:    errors.New("not found"),
			repoFiles:   testFiles,
			want:        testFiles,
			expectCache: true,
		},
		{
			name:     "handles repository error",
			version:  arm.VersionRef{ID: "v1.0.0", Type: arm.Tag},
			selector: arm.ContentSelector{Include: []string{"*.md"}},
			cacheErr: errors.New("not found"),
			repoErr:  errors.New("git error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &mockCache{
				content: make(map[string][]arm.File),
				getErr:  tt.cacheErr,
			}
			if tt.cachedFiles != nil {
				cache.content["test-reg/test-ruleset/v1.0.0"] = tt.cachedFiles
				cache.getErr = nil
			}

			repo := &mockRepository{
				files: tt.repoFiles,
				err:   tt.repoErr,
			}
			keyGen := &mockKeyGenerator{registryKey: "test-reg", rulesetKey: "test-ruleset"}

			g := NewGitRegistry(cache, repo, keyGen, "https://github.com/test/repo.git", "git")
			g.initialized = true
			got, err := g.GetContent(context.Background(), tt.version, tt.selector)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetContent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalFiles(got, tt.want) {
				t.Errorf("GetContent() = %v, want %v", got, tt.want)
			}

			if tt.expectCache {
				cached := cache.content["test-reg/test-ruleset/v1.0.0"]
				if !equalFiles(cached, tt.want) {
					t.Errorf("Expected files to be cached, got %v, want %v", cached, tt.want)
				}
			}
		})
	}
}

func equalVersionRefs(a, b []arm.VersionRef) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].ID != b[i].ID || a[i].Type != b[i].Type {
			return false
		}
	}
	return true
}

func TestGitRegistry_GetTags(t *testing.T) {
	tests := []struct {
		name    string
		tags    []string
		repoErr error
		want    []string
		wantErr bool
	}{
		{
			name: "returns tags from repository",
			tags: []string{"v1.0.0", "v1.1.0", "v2.0.0"},
			want: []string{"v1.0.0", "v1.1.0", "v2.0.0"},
		},
		{
			name:    "handles repository error",
			repoErr: errors.New("git error"),
			wantErr: true,
		},
		{
			name: "handles empty repository",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &mockCache{}
			repo := &mockRepository{
				tags: tt.tags,
				err:  tt.repoErr,
			}
			keyGen := &mockKeyGenerator{registryKey: "test-reg", rulesetKey: "test-ruleset"}

			g := NewGitRegistry(cache, repo, keyGen, "https://github.com/test/repo.git", "git")
			g.initialized = true
			got, err := g.GetTags(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetTags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalStringSlices(got, tt.want) {
				t.Errorf("GetTags() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGitRegistry_GetBranches(t *testing.T) {
	tests := []struct {
		name     string
		branches []string
		repoErr  error
		want     []string
		wantErr  bool
	}{
		{
			name:     "returns branches from repository",
			branches: []string{"main", "develop", "feature/test"},
			want:     []string{"main", "develop", "feature/test"},
		},
		{
			name:    "handles repository error",
			repoErr: errors.New("git error"),
			wantErr: true,
		},
		{
			name: "handles empty repository",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache := &mockCache{}
			repo := &mockRepository{
				branches: tt.branches,
				err:      tt.repoErr,
			}
			keyGen := &mockKeyGenerator{registryKey: "test-reg", rulesetKey: "test-ruleset"}

			g := NewGitRegistry(cache, repo, keyGen, "https://github.com/test/repo.git", "git")
			g.initialized = true
			got, err := g.GetBranches(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("GetBranches() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !equalStringSlices(got, tt.want) {
				t.Errorf("GetBranches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func equalFiles(a, b []arm.File) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i].Path != b[i].Path || !bytes.Equal(a[i].Content, b[i].Content) || a[i].Size != b[i].Size {
			return false
		}
	}
	return true
}
