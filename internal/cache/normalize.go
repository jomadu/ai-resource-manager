package cache

import (
	"crypto/sha256"
	"fmt"
	"sort"
	"strings"
)

func normalizeURL(url string) string {
	url = strings.TrimSuffix(url, ".git")
	url = strings.TrimSuffix(url, "/")
	return strings.ToLower(url)
}

func normalizePatterns(patterns []string) []string {
	if len(patterns) == 0 {
		return patterns
	}
	normalized := make([]string, len(patterns))
	copy(normalized, patterns)
	sort.Strings(normalized)
	return normalized
}

func sha256Hash(input string) string {
	hash := sha256.Sum256([]byte(input))
	return fmt.Sprintf("%x", hash)
}
