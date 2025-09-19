package types

// File represents a file with its content and metadata.
type File struct {
	Path    string
	Content []byte
	Size    int64
}
