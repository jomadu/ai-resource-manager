package core

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"fmt"
	"testing"
)

// Helper function to create a valid tar.gz archive
func createTarGz(files map[string][]byte) []byte {
	var buf bytes.Buffer
	gzWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzWriter)

	for path, content := range files {
		header := &tar.Header{
			Name: path,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		_ = tarWriter.WriteHeader(header)
		_, _ = tarWriter.Write(content)
	}

	_ = tarWriter.Close()
	_ = gzWriter.Close()
	return buf.Bytes()
}

// Helper function to create a valid zip archive
func createZip(files map[string][]byte) []byte {
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	for path, content := range files {
		writer, _ := zipWriter.Create(path)
		_, _ = writer.Write(content)
	}

	_ = zipWriter.Close()
	return buf.Bytes()
}

// TestIsArchive tests archive format detection
func TestIsArchive(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		// Normal cases
		{"tar.gz file", "archive.tar.gz", true},
		{"zip file", "archive.zip", true},
		{"yaml file", "config.yml", false},
		{"text file", "readme.txt", false},

		// Edge cases
		{"empty string", "", false},
		{"just .tar", "file.tar", false},
		{"just .gz", "file.gz", false},
		{"multiple extensions", "file.tar.gz.bak", false},
		{"uppercase TAR.GZ", "file.TAR.GZ", false},
		{"uppercase ZIP", "file.ZIP", false},

		// Extreme cases
		{"very long filename", "very-long-filename-that-goes-on-and-on-and-on.tar.gz", true},
		{"path with directories", "path/to/archive.tar.gz", true},
		{"hidden file", ".hidden.tar.gz", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.isArchive(tt.path)
			if result != tt.expected {
				t.Errorf("isArchive(%q) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

// TestExtractTarGz tests tar.gz extraction
func TestExtractTarGz(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name        string
		files       map[string][]byte
		expectError bool
		expectCount int
	}{
		// Normal cases
		{
			name: "single file",
			files: map[string][]byte{
				"file.txt": []byte("content"),
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "multiple files",
			files: map[string][]byte{
				"file1.txt": []byte("content1"),
				"file2.txt": []byte("content2"),
				"file3.txt": []byte("content3"),
			},
			expectError: false,
			expectCount: 3,
		},
		{
			name: "nested directories",
			files: map[string][]byte{
				"dir/file.txt":       []byte("content"),
				"dir/sub/file2.txt":  []byte("content2"),
				"dir/sub2/file3.txt": []byte("content3"),
			},
			expectError: false,
			expectCount: 3,
		},

		// Edge cases
		{
			name:        "empty archive",
			files:       map[string][]byte{},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "empty file",
			files: map[string][]byte{
				"empty.txt": []byte(""),
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "file with spaces in name",
			files: map[string][]byte{
				"file with spaces.txt": []byte("content"),
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "file with special characters",
			files: map[string][]byte{
				"file-name_123.txt": []byte("content"),
			},
			expectError: false,
			expectCount: 1,
		},

		// Extreme cases
		{
			name: "large file",
			files: map[string][]byte{
				"large.txt": bytes.Repeat([]byte("x"), 1024*1024), // 1MB
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "many files",
			files: func() map[string][]byte {
				files := make(map[string][]byte)
				for i := 0; i < 100; i++ {
					files[fmt.Sprintf("file%d.txt", i)] = []byte("content")
				}
				return files
			}(),
			expectError: false,
			expectCount: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := createTarGz(tt.files)
			file := &File{
				Path:    "test.tar.gz",
				Content: archive,
				Size:    int64(len(archive)),
			}

			extracted, err := extractor.extractTarGz(file)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(extracted) != tt.expectCount {
				t.Errorf("expected %d files, got %d", tt.expectCount, len(extracted))
			}

			// Verify content matches
			for _, extractedFile := range extracted {
				if expectedContent, ok := tt.files[extractedFile.Path]; ok {
					if !bytes.Equal(extractedFile.Content, expectedContent) {
						t.Errorf("content mismatch for %s", extractedFile.Path)
					}
				}
			}
		})
	}
}

// TestExtractTarGzSecurity tests security features
func TestExtractTarGzSecurity(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name        string
		files       map[string][]byte
		expectCount int
		description string
	}{
		{
			name: "directory traversal with ..",
			files: map[string][]byte{
				"../../../etc/passwd": []byte("malicious"),
				"normal.txt":          []byte("safe"),
			},
			expectCount: 1, // Only safe file
			description: "should skip files with .. in path",
		},
		{
			name: "absolute path",
			files: map[string][]byte{
				"/etc/passwd": []byte("malicious"),
				"normal.txt":  []byte("safe"),
			},
			expectCount: 1, // Only safe file
			description: "should skip absolute paths",
		},
		{
			name: "dot file",
			files: map[string][]byte{
				".": []byte("malicious"),
			},
			expectCount: 0,
			description: "should skip . path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := createTarGz(tt.files)
			file := &File{
				Path:    "test.tar.gz",
				Content: archive,
				Size:    int64(len(archive)),
			}

			extracted, err := extractor.extractTarGz(file)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(extracted) != tt.expectCount {
				t.Errorf("%s: expected %d files, got %d", tt.description, tt.expectCount, len(extracted))
			}
		})
	}
}

// TestExtractTarGzErrors tests error handling
func TestExtractTarGzErrors(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name    string
		content []byte
	}{
		{"invalid gzip", []byte("not a gzip file")},
		{"empty content", []byte{}},
		{"corrupted archive", []byte{0x1f, 0x8b, 0x08, 0x00, 0x00}}, // Partial gzip header
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{
				Path:    "test.tar.gz",
				Content: tt.content,
				Size:    int64(len(tt.content)),
			}

			_, err := extractor.extractTarGz(file)
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}

// TestExtractZip tests zip extraction
func TestExtractZip(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name        string
		files       map[string][]byte
		expectError bool
		expectCount int
	}{
		// Normal cases
		{
			name: "single file",
			files: map[string][]byte{
				"file.txt": []byte("content"),
			},
			expectError: false,
			expectCount: 1,
		},
		{
			name: "multiple files",
			files: map[string][]byte{
				"file1.txt": []byte("content1"),
				"file2.txt": []byte("content2"),
			},
			expectError: false,
			expectCount: 2,
		},

		// Edge cases
		{
			name:        "empty archive",
			files:       map[string][]byte{},
			expectError: false,
			expectCount: 0,
		},
		{
			name: "empty file",
			files: map[string][]byte{
				"empty.txt": []byte(""),
			},
			expectError: false,
			expectCount: 1,
		},

		// Extreme cases
		{
			name: "large file",
			files: map[string][]byte{
				"large.txt": bytes.Repeat([]byte("x"), 1024*1024), // 1MB
			},
			expectError: false,
			expectCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := createZip(tt.files)
			file := &File{
				Path:    "test.zip",
				Content: archive,
				Size:    int64(len(archive)),
			}

			extracted, err := extractor.extractZip(file)

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if len(extracted) != tt.expectCount {
				t.Errorf("expected %d files, got %d", tt.expectCount, len(extracted))
			}

			// Verify content matches
			for _, extractedFile := range extracted {
				if expectedContent, ok := tt.files[extractedFile.Path]; ok {
					if !bytes.Equal(extractedFile.Content, expectedContent) {
						t.Errorf("content mismatch for %s", extractedFile.Path)
					}
				}
			}
		})
	}
}

// TestExtractZipSecurity tests security features for zip
func TestExtractZipSecurity(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name        string
		files       map[string][]byte
		expectCount int
		description string
	}{
		{
			name: "directory traversal",
			files: map[string][]byte{
				"../../../etc/passwd": []byte("malicious"),
				"normal.txt":          []byte("safe"),
			},
			expectCount: 1,
			description: "should skip files with .. in path",
		},
		{
			name: "absolute path",
			files: map[string][]byte{
				"/etc/passwd": []byte("malicious"),
				"normal.txt":  []byte("safe"),
			},
			expectCount: 1,
			description: "should skip absolute paths",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			archive := createZip(tt.files)
			file := &File{
				Path:    "test.zip",
				Content: archive,
				Size:    int64(len(archive)),
			}

			extracted, err := extractor.extractZip(file)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(extracted) != tt.expectCount {
				t.Errorf("%s: expected %d files, got %d", tt.description, tt.expectCount, len(extracted))
			}
		})
	}
}

// TestExtractZipErrors tests error handling for zip
func TestExtractZipErrors(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name    string
		content []byte
	}{
		{"invalid zip", []byte("not a zip file")},
		{"empty content", []byte{}},
		{"corrupted archive", []byte{0x50, 0x4b, 0x03, 0x04}}, // Partial zip header
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := &File{
				Path:    "test.zip",
				Content: tt.content,
				Size:    int64(len(tt.content)),
			}

			_, err := extractor.extractZip(file)
			if err == nil {
				t.Error("expected error but got none")
			}
		})
	}
}

// TestExtractAndMerge tests the main extraction and merging logic
func TestExtractAndMerge(t *testing.T) {
	extractor := NewExtractor()

	tests := []struct {
		name        string
		setup       func() []*File
		expectCount int
		description string
	}{
		{
			name: "only loose files",
			setup: func() []*File {
				return []*File{
					{Path: "file1.txt", Content: []byte("content1")},
					{Path: "file2.txt", Content: []byte("content2")},
				}
			},
			expectCount: 2,
			description: "should return all loose files",
		},
		{
			name: "only archives",
			setup: func() []*File {
				return []*File{
					{
						Path:    "archive.tar.gz",
						Content: createTarGz(map[string][]byte{"file1.txt": []byte("content1")}),
					},
					{
						Path:    "archive.zip",
						Content: createZip(map[string][]byte{"file2.txt": []byte("content2")}),
					},
				}
			},
			expectCount: 2,
			description: "should extract all archives",
		},
		{
			name: "mixed loose and archives",
			setup: func() []*File {
				return []*File{
					{Path: "loose.txt", Content: []byte("loose")},
					{
						Path:    "archive.tar.gz",
						Content: createTarGz(map[string][]byte{"archived.txt": []byte("archived")}),
					},
				}
			},
			expectCount: 2,
			description: "should merge loose files and extracted archives",
		},
		{
			name: "archive overrides loose file",
			setup: func() []*File {
				return []*File{
					{Path: "file.txt", Content: []byte("loose")},
					{
						Path:    "archive.tar.gz",
						Content: createTarGz(map[string][]byte{"file.txt": []byte("archived")}),
					},
				}
			},
			expectCount: 1,
			description: "archive should override loose file with same path",
		},
		{
			name: "empty input",
			setup: func() []*File {
				return []*File{}
			},
			expectCount: 0,
			description: "should handle empty input",
		},
		{
			name: "nil input",
			setup: func() []*File {
				return nil
			},
			expectCount: 0,
			description: "should handle nil input",
		},
		{
			name: "multiple archives with same file",
			setup: func() []*File {
				return []*File{
					{
						Path:    "archive1.tar.gz",
						Content: createTarGz(map[string][]byte{"file.txt": []byte("first")}),
					},
					{
						Path:    "archive2.zip",
						Content: createZip(map[string][]byte{"file.txt": []byte("second")}),
					},
				}
			},
			expectCount: 1,
			description: "later archive should override earlier",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := tt.setup()
			merged, err := extractor.ExtractAndMerge(files)

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if len(merged) != tt.expectCount {
				t.Errorf("%s: expected %d files, got %d", tt.description, tt.expectCount, len(merged))
			}
		})
	}
}

// TestExtractAndMergeOverrideBehavior tests that archives override loose files
func TestExtractAndMergeOverrideBehavior(t *testing.T) {
	extractor := NewExtractor()

	files := []*File{
		{Path: "config.yml", Content: []byte("loose-content")},
		{
			Path:    "rules.tar.gz",
			Content: createTarGz(map[string][]byte{"config.yml": []byte("archive-content")}),
		},
	}

	merged, err := extractor.ExtractAndMerge(files)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(merged) != 1 {
		t.Fatalf("expected 1 file, got %d", len(merged))
	}

	if string(merged[0].Content) != "archive-content" {
		t.Errorf("expected archive content to override loose file, got: %s", string(merged[0].Content))
	}
}

// TestExtractAndMergeErrors tests error propagation
func TestExtractAndMergeErrors(t *testing.T) {
	extractor := NewExtractor()

	files := []*File{
		{Path: "good.txt", Content: []byte("content")},
		{Path: "bad.tar.gz", Content: []byte("invalid archive")},
	}

	_, err := extractor.ExtractAndMerge(files)
	if err == nil {
		t.Error("expected error from invalid archive")
	}
}
