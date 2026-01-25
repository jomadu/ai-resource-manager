package registry

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

func TestGitRegistry_GetPackage_WithArchiveAndPatterns(t *testing.T) {
	// Create a tar.gz archive with test files
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	files := map[string]string{
		"security/rule1.yml":     "security content 1",
		"security/rule2.yml":     "security content 2",
		"general/rule3.yml":      "general content 3",
		"experimental/rule4.yml": "experimental content 4",
	}

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatalf("failed to write tar header: %v", err)
		}
		if _, err := tarWriter.Write([]byte(content)); err != nil {
			t.Fatalf("failed to write tar content: %v", err)
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatalf("failed to close tar writer: %v", err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}

	// Test extraction and filtering
	archiveFile := &core.File{
		Path:    "test-ruleset/rules.tar.gz",
		Content: buf.Bytes(),
		Size:    int64(buf.Len()),
	}

	// Extract
	extractor := core.NewExtractor()
	extracted, err := extractor.ExtractAndMerge([]*core.File{archiveFile})
	if err != nil {
		t.Fatalf("extraction failed: %v", err)
	}

	t.Logf("Extracted %d files:", len(extracted))
	for _, f := range extracted {
		t.Logf("  - %s", f.Path)
	}

	// Apply patterns
	include := []string{"security/**/*.yml"}
	exclude := []string{"**/experimental/**"}

	var filtered []*core.File
	for _, file := range extracted {
		// Use the same matching logic as GitRegistry
		if len(include) == 0 && len(exclude) == 0 {
			filtered = append(filtered, file)
			continue
		}

		// Check exclude patterns first
		excluded := false
		for _, pattern := range exclude {
			if core.MatchPattern(pattern, file.Path) {
				excluded = true
				break
			}
		}
		if excluded {
			t.Logf("  %s: EXCLUDED by %v", file.Path, exclude)
			continue
		}

		// If no include patterns, file is included (not excluded)
		if len(include) == 0 {
			filtered = append(filtered, file)
			t.Logf("  %s: INCLUDED (no include patterns)", file.Path)
			continue
		}

		// Check include patterns
		included := false
		for _, pattern := range include {
			if core.MatchPattern(pattern, file.Path) {
				included = true
				break
			}
		}
		if included {
			filtered = append(filtered, file)
			t.Logf("  %s: INCLUDED by %v", file.Path, include)
		} else {
			t.Logf("  %s: NOT INCLUDED", file.Path)
		}
	}

	t.Logf("Filtered to %d files:", len(filtered))
	for _, f := range filtered {
		t.Logf("  - %s", f.Path)
	}

	// Verify results
	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered files, got %d", len(filtered))
	}

	expectedPaths := map[string]bool{
		"security/rule1.yml": false,
		"security/rule2.yml": false,
	}

	for _, f := range filtered {
		if _, ok := expectedPaths[f.Path]; ok {
			expectedPaths[f.Path] = true
		} else {
			t.Errorf("unexpected file in filtered results: %s", f.Path)
		}
	}

	for path, found := range expectedPaths {
		if !found {
			t.Errorf("expected file not found in filtered results: %s", path)
		}
	}
}
