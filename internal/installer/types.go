package installer

// File represents a file with its content and metadata.
type File struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}

// Installation represents an installed ruleset in a directory.
type Installation struct {
	Ruleset string `json:"ruleset"`
	Version string `json:"version"`
	Path    string `json:"path"`
}
