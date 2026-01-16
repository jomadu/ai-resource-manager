package registry

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

func calculateIntegrity(files []*core.File) string {
	h := sha256.New()

	var paths []string
	fileMap := make(map[string][]byte)
	for _, file := range files {
		paths = append(paths, file.Path)
		fileMap[file.Path] = file.Content
	}
	sort.Strings(paths)

	for _, path := range paths {
		h.Write([]byte(path))
		h.Write(fileMap[path])
	}

	return "sha256-" + hex.EncodeToString(h.Sum(nil))
}
