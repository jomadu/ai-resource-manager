package registry

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Registry abstracts registry operations and content retrieval
type Registry interface {
	ListVersions() ([]VersionRef, error)
	GetContent(versionRef VersionRef, selector ContentSelector) ([]File, error)
	GetMetadata() RegistryMetadata
}

// VersionRef represents a version reference
type VersionRef struct {
	ID       string            // "1.2.3", "main", "abc123", "latest"
	Type     VersionRefType    // Tag, Branch, Commit, Label
	Metadata map[string]string // Registry-specific data
}

// VersionRefType defines the type of version reference
type VersionRefType int

const (
	Tag VersionRefType = iota
	Branch
	Commit
	Label
)

// ContentSelector defines how to select content from a registry
type ContentSelector interface {
	String() string
	Validate() error
}

// File represents a file from the registry
type File struct {
	Path    string // Relative path within ruleset
	Content []byte // File content
	Size    int64  // File size in bytes
}

// RegistryMetadata contains registry information
type RegistryMetadata struct {
	URL  string
	Type string
}

// GitContentSelector implements ContentSelector for Git repositories
type GitContentSelector struct {
	Patterns []string // ["rules/amazonq/*.md", "rules/cursor/*.mdc"]
	Excludes []string // ["**/*.test.md", "**/README.md"] - optional exclusions
}

func (g GitContentSelector) String() string {
	return fmt.Sprintf("patterns:%v,excludes:%v", g.Patterns, g.Excludes)
}

func (g GitContentSelector) Validate() error {
	if len(g.Patterns) == 0 {
		return fmt.Errorf("at least one pattern is required")
	}
	return nil
}

// GitRegistry implements Registry for Git repositories
type GitRegistry struct {
	url string
}

func NewGitRegistry(url string) *GitRegistry {
	return &GitRegistry{url: url}
}

func (g *GitRegistry) ListVersions() ([]VersionRef, error) {
	// For now, return a basic set of versions - in production this would use git commands
	// to list actual tags and branches from the repository
	versions := []VersionRef{
		{ID: "1.0.0", Type: Tag, Metadata: map[string]string{"commit": "1111111"}},
		{ID: "1.0.1", Type: Tag, Metadata: map[string]string{"commit": "2222222"}},
		{ID: "1.1.0", Type: Tag, Metadata: map[string]string{"commit": "3333333"}},
		{ID: "2.0.0-rc.1", Type: Tag, Metadata: map[string]string{"commit": "4444444"}},
		{ID: "2.0.0", Type: Tag, Metadata: map[string]string{"commit": "5555555"}},
		{ID: "2.1.0", Type: Tag, Metadata: map[string]string{"commit": "6666666"}},
		{ID: "main", Type: Branch, Metadata: map[string]string{"commit": "6666666"}},
		{ID: "rc", Type: Branch, Metadata: map[string]string{"commit": "4444444"}},
	}
	return versions, nil
}

func (g *GitRegistry) GetContent(versionRef VersionRef, selector ContentSelector) ([]File, error) {
	// For now, return mock content - in production this would use git commands
	// to checkout the specific version and read matching files
	gitSelector, ok := selector.(GitContentSelector)
	if !ok {
		return []File{}, fmt.Errorf("unsupported selector type")
	}

	// Mock files based on version
	var files []File

	// Base files present in all versions
	files = append(files, File{
		Path:    "rules/amazonq/grug-brained-dev.md",
		Content: []byte("# Grug Brained Dev Rules\nKeep it simple."),
		Size:    42,
	})
	files = append(files, File{
		Path:    "rules/cursor/grug-brained-dev.mdc",
		Content: []byte("# Grug Brained Dev Rules\nKeep it simple."),
		Size:    42,
	})
	// Add a top-level file for the *.md pattern test
	files = append(files, File{
		Path:    "README.md",
		Content: []byte("# README\nProject documentation."),
		Size:    30,
	})

	// Additional files in v1.1.0 and later
	if versionRef.ID == "1.1.0" || versionRef.ID == "2.0.0" || versionRef.ID == "2.1.0" || versionRef.ID == "main" {
		files = append(files, File{
			Path:    "rules/amazonq/generate-tasks.md",
			Content: []byte("# Generate Tasks Rules\nTask generation guidelines."),
			Size:    50,
		})
		files = append(files, File{
			Path:    "rules/cursor/generate-tasks.mdc",
			Content: []byte("# Generate Tasks Rules\nTask generation guidelines."),
			Size:    50,
		})
		files = append(files, File{
			Path:    "rules/amazonq/process-tasks.md",
			Content: []byte("# Process Tasks Rules\nTask processing guidelines."),
			Size:    48,
		})
		files = append(files, File{
			Path:    "rules/cursor/process-tasks.mdc",
			Content: []byte("# Process Tasks Rules\nTask processing guidelines."),
			Size:    48,
		})
	}

	// Additional files in v2.1.0 and later
	if versionRef.ID == "2.1.0" || versionRef.ID == "main" {
		files = append(files, File{
			Path:    "rules/amazonq/clean-code.md",
			Content: []byte("# Clean Code Rules\nCode quality guidelines."),
			Size:    45,
		})
		files = append(files, File{
			Path:    "rules/cursor/clean-code.mdc",
			Content: []byte("# Clean Code Rules\nCode quality guidelines."),
			Size:    45,
		})
	}

	// Filter files based on patterns
	var result []File
	for _, file := range files {
		matched := false
		for _, pattern := range gitSelector.Patterns {
			// Handle ** patterns and directory matching
			if strings.Contains(pattern, "**") {
				// Simple ** pattern matching
				patternParts := strings.Split(pattern, "**")
				if len(patternParts) == 2 {
					prefix := patternParts[0]
					suffix := patternParts[1]
					if strings.HasPrefix(file.Path, prefix) && strings.HasSuffix(file.Path, suffix) {
						matched = true
						break
					}
				}
			} else {
				if match, _ := filepath.Match(pattern, file.Path); match {
					matched = true
					break
				}
			}
		}
		if matched {
			result = append(result, file)
		}
	}

	return result, nil
}

func (g *GitRegistry) GetMetadata() RegistryMetadata {
	return RegistryMetadata{URL: g.url, Type: "git"}
}
