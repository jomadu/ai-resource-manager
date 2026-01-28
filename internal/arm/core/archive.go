package core

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Extractor handles archive extraction and merging
type Extractor struct{}

// NewExtractor creates a new archive extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// Extract extracts archives from files to subdirectories (no merge)
func (e *Extractor) Extract(files []*File) ([]*File, error) {
	var result []*File

	for _, file := range files {
		if e.isArchive(file.Path) {
			extractedFiles, err := e.extractArchive(file)
			if err != nil {
				return nil, fmt.Errorf("failed to extract %s: %w", file.Path, err)
			}
			result = append(result, extractedFiles...)
		} else {
			result = append(result, file)
		}
	}

	return result, nil
}

// isArchive checks if file is a supported archive format
func (e *Extractor) isArchive(path string) bool {
	return strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".zip")
}

// extractArchive extracts a single archive file to subdirectory
func (e *Extractor) extractArchive(file *File) ([]*File, error) {
	subdirName := e.getSubdirName(file.Path)
	
	if strings.HasSuffix(file.Path, ".tar.gz") {
		return e.extractTarGz(file, subdirName)
	}
	if strings.HasSuffix(file.Path, ".zip") {
		return e.extractZip(file, subdirName)
	}
	return nil, fmt.Errorf("unsupported archive format: %s", file.Path)
}

// getSubdirName returns subdirectory name from archive filename
func (e *Extractor) getSubdirName(archivePath string) string {
	base := filepath.Base(archivePath)
	base = strings.TrimSuffix(base, ".tar.gz")
	base = strings.TrimSuffix(base, ".zip")
	return base
}

// extractTarGz extracts a tar.gz archive to subdirectory
func (e *Extractor) extractTarGz(file *File, subdirName string) ([]*File, error) {
	// Create temp file for streaming extraction
	tempFile, err := os.CreateTemp("", "arm-extract-*.tar.gz")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()
	defer func() { _ = tempFile.Close() }()

	// Write content to temp file
	if _, err := tempFile.Write(file.Content); err != nil {
		return nil, err
	}
	if _, err := tempFile.Seek(0, 0); err != nil {
		return nil, err
	}

	// Open gzip reader
	gzReader, err := gzip.NewReader(tempFile)
	if err != nil {
		return nil, err
	}
	defer func() { _ = gzReader.Close() }()

	// Open tar reader
	tarReader := tar.NewReader(gzReader)

	var extractedFiles []*File
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Skip directories
		if header.Typeflag == tar.TypeDir {
			continue
		}

		// Sanitize path to prevent directory traversal
		cleanName := filepath.Clean(header.Name)
		if cleanName == "." || filepath.IsAbs(header.Name) || strings.Contains(cleanName, "..") {
			continue
		}

		// Read file content
		content, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, &File{
			Path:    filepath.Join(subdirName, cleanName),
			Content: content,
			Size:    header.Size,
		})
	}

	return extractedFiles, nil
}

// extractZip extracts a zip archive to subdirectory
func (e *Extractor) extractZip(file *File, subdirName string) ([]*File, error) {
	// Create temp file for streaming extraction
	tempFile, err := os.CreateTemp("", "arm-extract-*.zip")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.Remove(tempFile.Name()) }()

	// Write content to temp file
	if _, err := tempFile.Write(file.Content); err != nil {
		return nil, err
	}
	if err := tempFile.Close(); err != nil {
		return nil, err
	}

	// Open zip reader
	zipReader, err := zip.OpenReader(tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer func() { _ = zipReader.Close() }()

	var extractedFiles []*File
	for _, zipFile := range zipReader.File {
		// Skip directories
		if zipFile.FileInfo().IsDir() {
			continue
		}

		// Sanitize path to prevent directory traversal
		cleanName := filepath.Clean(zipFile.Name)
		if cleanName == "." || filepath.IsAbs(zipFile.Name) || strings.Contains(cleanName, "..") {
			continue
		}

		// Open file in zip
		rc, err := zipFile.Open()
		if err != nil {
			return nil, err
		}

		// Read content
		content, err := io.ReadAll(rc)
		if closeErr := rc.Close(); closeErr != nil {
			return nil, closeErr
		}
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, &File{
			Path:    filepath.Join(subdirName, cleanName),
			Content: content,
			Size:    int64(zipFile.UncompressedSize64),
		})
	}

	return extractedFiles, nil
}
