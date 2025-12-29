package storage

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
)

// DIRECTORY STRUCTURE:
// ~/.arm/
// ├── .armrc                                # ARM configuration (optional)
// └── storage/
//     └── registries/
//         └── {registry-key}/
//             ├── metadata.json             # Registry metadata + timestamps
//             ├── repo/                     # Git clone (git registries only)
//             │   ├── .git/
//             │   └── source-files...
//             └── packages/
//                 └── {package-key}/
//                     ├── metadata.json     # Package metadata + timestamps
//                     └── {version}/
//                         ├── metadata.json # Version metadata + timestamps
//                         └── files/
//                             └── extracted-files...
//
// CONCURRENCY PROTECTION:
// Cross-process file locking prevents concurrent access issues:
// - Registry: locks {registry-dir}.lock for metadata operations
// - Repo: locks {repo-dir}.lock for git operations
// - PackageCache: locks {package-dir}.lock per package for file operations
//
// Storage provides three main components:
// - Registry: manages registry directory structure and metadata
// - Repo: manages git repository operations with locking
// - PackageCache: manages package file storage with per-package locking

// GenerateKey creates a hash from any object for use as directory name
func GenerateKey(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}