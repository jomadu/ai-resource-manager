package lockfile

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

func TestGenerateChecksum(t *testing.T) {
	files := []types.File{
		{Path: "file1.md", Content: []byte("content1")},
		{Path: "file2.md", Content: []byte("content2")},
	}

	checksum := GenerateChecksum(files)

	if !IsValidChecksum(checksum) {
		t.Errorf("Generated checksum is not valid: %s", checksum)
	}

	// Should be deterministic
	checksum2 := GenerateChecksum(files)
	if checksum != checksum2 {
		t.Errorf("Checksum should be deterministic, got %s and %s", checksum, checksum2)
	}
}

func TestVerifyChecksum(t *testing.T) {
	files := []types.File{
		{Path: "file1.md", Content: []byte("content1")},
		{Path: "file2.md", Content: []byte("content2")},
	}

	checksum := GenerateChecksum(files)

	if !VerifyChecksum(files, checksum) {
		t.Error("Checksum verification should pass for same files")
	}

	// Different content should fail
	differentFiles := []types.File{
		{Path: "file1.md", Content: []byte("different")},
		{Path: "file2.md", Content: []byte("content2")},
	}

	if VerifyChecksum(differentFiles, checksum) {
		t.Error("Checksum verification should fail for different content")
	}
}

func TestIsValidChecksum(t *testing.T) {
	validChecksum := "sha256:abcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"
	if !IsValidChecksum(validChecksum) {
		t.Error("Valid checksum should pass validation")
	}

	invalidChecksums := []string{
		"md5:abc123",
		"sha256:short",
		"invalid",
		"",
	}

	for _, invalid := range invalidChecksums {
		if IsValidChecksum(invalid) {
			t.Errorf("Invalid checksum should fail validation: %s", invalid)
		}
	}
}
