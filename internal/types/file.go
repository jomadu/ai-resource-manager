package types

// File represents a file with its content and metadata.
type File struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}
