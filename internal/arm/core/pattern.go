package core

import "strings"

// MatchPattern matches a glob pattern against a file path.
// Supports *, **, and literal paths.
// Pattern and path separators are normalized to forward slashes.
func MatchPattern(pattern, path string) bool {
	// Normalize path separators
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	path = strings.ReplaceAll(path, "\\", "/")

	// Handle literal paths (no wildcards)
	if !strings.Contains(pattern, "*") {
		return pattern == path
	}

	// Handle patterns with ** (recursive directory matching)
	if strings.Contains(pattern, "**") {
		return matchDoublestar(pattern, path)
	}

	// Handle simple wildcard patterns (no **)
	return matchSimpleWildcard(pattern, path)
}

// matchDoublestar handles patterns containing **
func matchDoublestar(pattern, path string) bool {
	// Split pattern by /**/ to get segments
	parts := strings.Split(pattern, "/**/")

	if len(parts) == 1 {
		// Pattern is **/ prefix or /** suffix
		if strings.HasPrefix(pattern, "**/") {
			suffix := pattern[3:]
			// Match if path ends with suffix or contains /suffix
			if matchSimpleWildcard(suffix, path) {
				return true
			}
			// Try matching against any path segment
			pathParts := strings.Split(path, "/")
			for i := range pathParts {
				subPath := strings.Join(pathParts[i:], "/")
				if matchSimpleWildcard(suffix, subPath) {
					return true
				}
			}
			return false
		}

		if strings.HasSuffix(pattern, "/**") {
			prefix := pattern[:len(pattern)-3]
			return strings.HasPrefix(path, prefix+"/") || path == prefix
		}

		return false
	}

	// Pattern has /** in the middle (e.g., "security/**/*.yml")
	// Match prefix, then try to match suffix from any point in remaining path
	prefix := parts[0]
	suffix := parts[len(parts)-1]

	// Path must start with prefix (if prefix is not empty)
	if prefix != "" && !strings.HasPrefix(path, prefix+"/") && path != prefix {
		return false
	}

	// Get the remaining path after prefix
	remaining := path
	if prefix != "" {
		if path == prefix {
			// Path exactly matches prefix, check if suffix can match empty
			return suffix == "" || suffix == "*"
		}
		remaining = path[len(prefix)+1:]
	}

	// Try to match suffix from any point in remaining path
	if suffix == "" {
		return true
	}

	// Try matching suffix against remaining path and all its sub-paths
	pathParts := strings.Split(remaining, "/")
	for i := range pathParts {
		subPath := strings.Join(pathParts[i:], "/")
		if matchSimpleWildcard(suffix, subPath) {
			return true
		}
	}

	return false
}

// matchSimpleWildcard handles patterns with * but not **
func matchSimpleWildcard(pattern, path string) bool {
	if !strings.Contains(pattern, "*") {
		return pattern == path
	}

	parts := strings.Split(pattern, "*")
	pos := 0

	for i, part := range parts {
		if part == "" {
			continue
		}

		idx := strings.Index(path[pos:], part)
		if idx == -1 {
			return false
		}

		// First part must match at the beginning
		if i == 0 && idx != 0 {
			return false
		}

		pos += idx + len(part)
	}

	// Last part must match at the end (if not empty)
	if len(parts) > 0 && parts[len(parts)-1] != "" && !strings.HasSuffix(path, parts[len(parts)-1]) {
		return false
	}

	return true
}
