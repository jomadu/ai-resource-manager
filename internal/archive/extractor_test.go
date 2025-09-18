package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestExtractAndMerge(t *testing.T) {
	extractor := NewExtractor()

	// Create test tar.gz content
	tarGzContent := createTestTarGz(t, map[string]string{
		"ruleset-zip-1.yml":                     "version: 1.0\nrules: []",
		"build/cursor/ruleset-zip-1_rule-1.mdc": "# Rule 1\nContent from archive",
	})

	// Create test zip content
	zipContent := createTestZip(t, map[string]string{
		"ruleset-zip-2.yml":                     "version: 1.0\nrules: []",
		"build/amazonq/ruleset-zip-2_rule-1.md": "# Rule 2\nContent from zip",
	})

	files := []types.File{
		{
			Path:    "ruleset-zip-1.tar.gz",
			Content: tarGzContent,
			Size:    int64(len(tarGzContent)),
		},
		{
			Path:    "ruleset-zip-2.zip",
			Content: zipContent,
			Size:    int64(len(zipContent)),
		},
		{
			Path:    "ruleset-a.yml",
			Content: []byte("version: 1.0\nrules: [loose-file]"),
			Size:    30,
		},
		{
			Path:    "build/cursor/existing.mdc",
			Content: []byte("# Existing\nLoose file content"),
			Size:    25,
		},
	}

	mergedFiles, err := extractor.ExtractAndMerge(files)
	if err != nil {
		t.Fatalf("ExtractAndMerge failed: %v", err)
	}

	// Verify expected files are present
	expectedFiles := map[string]bool{
		"ruleset-zip-1.yml":                     false,
		"build/cursor/ruleset-zip-1_rule-1.mdc": false,
		"ruleset-zip-2.yml":                     false,
		"build/amazonq/ruleset-zip-2_rule-1.md": false,
		"ruleset-a.yml":                         false,
		"build/cursor/existing.mdc":             false,
	}

	for _, file := range mergedFiles {
		if _, exists := expectedFiles[file.Path]; exists {
			expectedFiles[file.Path] = true
		} else {
			t.Errorf("Unexpected file: %s", file.Path)
		}
	}

	// Check all expected files were found
	for path, found := range expectedFiles {
		if !found {
			t.Errorf("Expected file not found: %s", path)
		}
	}

	// Verify archive wins over loose files (if there were conflicts)
	// In this test, there are no conflicts, so just verify content
	for _, file := range mergedFiles {
		if file.Path == "ruleset-a.yml" {
			content := string(file.Content)
			if content != "version: 1.0\nrules: [loose-file]" {
				t.Errorf("Loose file content incorrect: %s", content)
			}
		}
	}
}

func TestIsArchive(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		path     string
		expected bool
	}{
		{"file.tar.gz", true},
		{"file.zip", true},
		{"file.yml", false},
		{"file.md", false},
		{"file.tar", false},
		{"file.gz", false},
	}

	for _, test := range tests {
		result := extractor.isArchive(test.path)
		if result != test.expected {
			t.Errorf("isArchive(%s) = %v, want %v", test.path, result, test.expected)
		}
	}
}

func createTestTarGz(t *testing.T, files map[string]string) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Size: int64(len(content)),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			t.Fatal(err)
		}
		if _, err := tarWriter.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	if err := tarWriter.Close(); err != nil {
		t.Fatal(err)
	}
	if err := gzWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func createTestZip(t *testing.T, files map[string]string) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for path, content := range files {
		writer, err := zipWriter.Create(path)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := writer.Write([]byte(content)); err != nil {
			t.Fatal(err)
		}
	}

	if err := zipWriter.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}
