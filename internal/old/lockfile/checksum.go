package lockfile

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GenerateChecksum creates a SHA-256 checksum of the resolved files.
func GenerateChecksum(files []types.File) string {
	// Sort files by path for deterministic ordering
	sort.Slice(files, func(i, j int) bool {
		return files[i].Path < files[j].Path
	})

	hasher := sha256.New()

	for _, file := range files {
		// Include file path and content in checksum
		hasher.Write([]byte(file.Path))
		hasher.Write([]byte{0}) // separator
		hasher.Write(file.Content)
		hasher.Write([]byte{0}) // separator
	}

	return fmt.Sprintf("sha256:%x", hasher.Sum(nil))
}

// VerifyChecksum verifies that the files match the expected checksum.
func VerifyChecksum(files []types.File, expectedChecksum string) bool {
	actualChecksum := GenerateChecksum(files)
	return actualChecksum == expectedChecksum
}

// IsValidChecksum checks if a checksum has the correct format.
func IsValidChecksum(checksum string) bool {
	return strings.HasPrefix(checksum, "sha256:") && len(checksum) == 71 // "sha256:" + 64 hex chars
}
