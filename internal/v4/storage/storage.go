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
// Storage provides two main components:
// - Registry: manages registry directory structure and metadata
// - PackageCache: manages package file storage within registry

// GenerateKey creates a hash from any object for use as directory name
func GenerateKey(obj interface{}) (string, error) {
	data, err := json.Marshal(obj)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:]), nil
}