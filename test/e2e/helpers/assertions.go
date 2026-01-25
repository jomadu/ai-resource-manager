package helpers

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// FileExists checks if a file exists
func FileExists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	return err == nil
}

// AssertFileExists asserts that a file exists
func AssertFileExists(t *testing.T, path string) {
	t.Helper()
	if !FileExists(t, path) {
		t.Errorf("expected file to exist: %s", path)
	}
}

// AssertFileNotExists asserts that a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	t.Helper()
	if FileExists(t, path) {
		t.Errorf("expected file to not exist: %s", path)
	}
}

// AssertFileContains asserts that a file contains the given content
func AssertFileContains(t *testing.T, path, content string) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", path, err)
	}
	if !contains(string(data), content) {
		t.Errorf("file %s does not contain expected content:\nExpected substring: %s\nActual content: %s", path, content, string(data))
	}
}

// AssertDirExists asserts that a directory exists
func AssertDirExists(t *testing.T, path string) {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Errorf("expected directory to exist: %s (error: %v)", path, err)
		return
	}
	if !info.IsDir() {
		t.Errorf("expected %s to be a directory", path)
	}
}

// AssertDirNotExists asserts that a directory does not exist
func AssertDirNotExists(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err == nil {
		t.Errorf("expected directory to not exist: %s", path)
	}
}

// AssertJSONField asserts that a JSON file contains a specific field with the expected value
func AssertJSONField(t *testing.T, path, field string, expected interface{}) {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read JSON file %s: %v", path, err)
	}
	
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("failed to parse JSON file %s: %v", path, err)
	}
	
	actual, ok := obj[field]
	if !ok {
		t.Errorf("JSON file %s does not contain field %s", path, field)
		return
	}
	
	if actual != expected {
		t.Errorf("JSON field %s in %s: expected %v, got %v", field, path, expected, actual)
	}
}

// ReadJSON reads and parses a JSON file
func ReadJSON(t *testing.T, path string) map[string]interface{} {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read JSON file %s: %v", path, err)
	}
	
	var obj map[string]interface{}
	if err := json.Unmarshal(data, &obj); err != nil {
		t.Fatalf("failed to parse JSON file %s: %v", path, err)
	}
	
	return obj
}

// CountFiles counts the number of files in a directory (non-recursive)
func CountFiles(t *testing.T, dir string) int {
	t.Helper()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("failed to read directory %s: %v", dir, err)
	}
	
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			count++
		}
	}
	return count
}

// CountFilesRecursive counts the number of files in a directory recursively
func CountFilesRecursive(t *testing.T, dir string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("failed to walk directory %s: %v", dir, err)
	}
	return count
}

// CountFilesWithExtension counts files with a specific extension in a directory recursively
func CountFilesWithExtension(t *testing.T, dir string, ext string) int {
	t.Helper()
	count := 0
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ext {
			count++
		}
		return nil
	})
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatalf("failed to walk directory %s: %v", dir, err)
	}
	return count
}

// DirExists checks if a directory exists
func DirExists(dir string) bool {
	info, err := os.Stat(dir)
	if err != nil {
		return false
	}
	return info.IsDir()
}

// contains checks if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 || 
		(len(s) > 0 && len(substr) > 0 && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
