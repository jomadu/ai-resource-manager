package archive

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// Extractor handles archive extraction and merging
type Extractor struct{}

// NewExtractor creates a new archive extractor
func NewExtractor() *Extractor {
	return &Extractor{}
}

// ExtractAndMerge extracts archives from files and merges with loose files
func (e *Extractor) ExtractAndMerge(files []types.File) ([]types.File, error) {
	var mergedFiles []types.File
	fileMap := make(map[string]types.File)

	// First, add all loose files
	for _, file := range files {
		if !e.isArchive(file.Path) {
			fileMap[file.Path] = file
		}
	}

	// Then extract archives (archives win on conflicts)
	for _, file := range files {
		if e.isArchive(file.Path) {
			extractedFiles, err := e.extractArchive(file)
			if err != nil {
				return nil, fmt.Errorf("failed to extract %s: %w", file.Path, err)
			}

			// Add extracted files (overwriting loose files)
			for _, extracted := range extractedFiles {
				fileMap[extracted.Path] = extracted
			}
		}
	}

	// Convert map back to slice
	for _, file := range fileMap {
		mergedFiles = append(mergedFiles, file)
	}

	return mergedFiles, nil
}

// isArchive checks if file is a supported archive format
func (e *Extractor) isArchive(path string) bool {
	return strings.HasSuffix(path, ".tar.gz") || strings.HasSuffix(path, ".zip")
}

// extractArchive extracts a single archive file
func (e *Extractor) extractArchive(file types.File) ([]types.File, error) {
	if strings.HasSuffix(file.Path, ".tar.gz") {
		return e.extractTarGz(file)
	}
	if strings.HasSuffix(file.Path, ".zip") {
		return e.extractZip(file)
	}
	return nil, fmt.Errorf("unsupported archive format: %s", file.Path)
}

// extractTarGz extracts a tar.gz archive
func (e *Extractor) extractTarGz(file types.File) ([]types.File, error) {
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

	var extractedFiles []types.File
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

		// Read file content
		content, err := io.ReadAll(tarReader)
		if err != nil {
			return nil, err
		}

		extractedFiles = append(extractedFiles, types.File{
			Path:    header.Name,
			Content: content,
			Size:    header.Size,
		})
	}

	return extractedFiles, nil
}

// extractZip extracts a zip archive
func (e *Extractor) extractZip(file types.File) ([]types.File, error) {
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

	var extractedFiles []types.File
	for _, zipFile := range zipReader.File {
		// Skip directories
		if zipFile.FileInfo().IsDir() {
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

		extractedFiles = append(extractedFiles, types.File{
			Path:    zipFile.Name,
			Content: content,
			Size:    int64(zipFile.UncompressedSize64),
		})
	}

	return extractedFiles, nil
}
