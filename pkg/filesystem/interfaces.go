package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileSystemManager handles atomic file operations for ruleset installation
type FileSystemManager interface {
	// Install files to configured sink directories atomically
	Install(sinkDir, registry, ruleset, version string, files []File) error

	// Uninstall removes ruleset files and cleans up empty directories
	Uninstall(sinkDir, registry, ruleset, version string) error

	// List returns installed files for a ruleset
	List(sinkDir, registry, ruleset, version string) ([]string, error)
}

// File represents a file to be installed
type File struct {
	Path    string
	Content []byte
}

// AtomicFileSystemManager implements FileSystemManager with atomic operations
type AtomicFileSystemManager struct {
	basePath string
}

func NewAtomicFileSystemManager(basePath string) *AtomicFileSystemManager {
	return &AtomicFileSystemManager{basePath: basePath}
}

func (a *AtomicFileSystemManager) Install(sinkDir, registry, ruleset, version string, files []File) error {
	if sinkDir == "" || registry == "" || ruleset == "" || version == "" {
		return fmt.Errorf("all parameters must be non-empty")
	}

	targetDir := filepath.Join(sinkDir, "arm", registry, ruleset, version)
	tempDir := targetDir + ".tmp"

	// Create temp directory
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}

	// Write files to temp directory
	for _, file := range files {
		filePath := filepath.Join(tempDir, file.Path)
		dir := filepath.Dir(filePath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			os.RemoveAll(tempDir)
			return err
		}
		if err := os.WriteFile(filePath, file.Content, 0644); err != nil {
			os.RemoveAll(tempDir)
			return err
		}
	}

	// Remove existing installation if it exists
	if _, err := os.Stat(targetDir); err == nil {
		if err := os.RemoveAll(targetDir); err != nil {
			os.RemoveAll(tempDir)
			return err
		}
	}

	// Atomic move
	if err := os.Rename(tempDir, targetDir); err != nil {
		os.RemoveAll(tempDir)
		return err
	}

	return nil
}

func (a *AtomicFileSystemManager) Uninstall(sinkDir, registry, ruleset, version string) error {
	targetDir := filepath.Join(sinkDir, "arm", registry, ruleset, version)

	// Remove the version directory
	if err := os.RemoveAll(targetDir); err != nil {
		return err
	}

	// Clean up empty parent directories
	rulesetDir := filepath.Join(sinkDir, "arm", registry, ruleset)
	if isEmpty, _ := isDirEmpty(rulesetDir); isEmpty {
		os.RemoveAll(rulesetDir)
	}

	registryDir := filepath.Join(sinkDir, "arm", registry)
	if isEmpty, _ := isDirEmpty(registryDir); isEmpty {
		os.RemoveAll(registryDir)
	}

	armDir := filepath.Join(sinkDir, "arm")
	if isEmpty, _ := isDirEmpty(armDir); isEmpty {
		os.RemoveAll(armDir)
	}

	return nil
}

func (a *AtomicFileSystemManager) List(sinkDir, registry, ruleset, version string) ([]string, error) {
	targetDir := filepath.Join(sinkDir, "arm", registry, ruleset, version)
	var files []string

	err := filepath.Walk(targetDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			relPath, _ := filepath.Rel(targetDir, path)
			files = append(files, relPath)
		}
		return nil
	})

	if os.IsNotExist(err) {
		return []string{}, nil
	}

	return files, err
}

func isDirEmpty(dir string) (bool, error) {
	f, err := os.Open(dir)
	if err != nil {
		return false, err
	}
	defer f.Close()

	_, err = f.Readdirnames(1)
	if err == nil {
		return false, nil
	}
	return true, nil
}
