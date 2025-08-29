package arm

// File represents a file with its content and metadata.
type File struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}

// VersionRefType defines the type of version reference.
type VersionRefType int

const (
	Tag VersionRefType = iota
	Branch
	Commit
)

// VersionRef represents a version reference in a repository.
type VersionRef struct {
	ID   string         `json:"id"`
	Type VersionRefType `json:"type"`
}

// ContentSelector defines include/exclude patterns for content filtering.
type ContentSelector struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// OutdatedRuleset represents a ruleset that has newer versions available.
type OutdatedRuleset struct {
	Registry string `json:"registry"`
	Name     string `json:"name"`
	Current  string `json:"current"`
	Wanted   string `json:"wanted"`
	Latest   string `json:"latest"`
}

// InstalledRuleset represents a currently installed ruleset.
type InstalledRuleset struct {
	Registry string   `json:"registry"`
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Include  []string `json:"include"`
	Exclude  []string `json:"exclude"`
	Sinks    []string `json:"sinks"`
}

// RulesetInfo provides detailed information about a ruleset.
type RulesetInfo struct {
	Registry       string   `json:"registry"`
	Name           string   `json:"name"`
	RegistryURL    string   `json:"registry_url"`
	RegistryType   string   `json:"registry_type"`
	Include        []string `json:"include"`
	Exclude        []string `json:"exclude"`
	InstalledPaths []string `json:"installed_paths"`
	Sinks          []string `json:"sinks"`
	Constraint     string   `json:"constraint"`
	Resolved       string   `json:"resolved"`
	Wanted         string   `json:"wanted"`
	Latest         string   `json:"latest"`
}

// Installation represents an installed ruleset in a directory.
type Installation struct {
	Ruleset string `json:"ruleset"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

// VersionInfo provides version information about the ARM tool itself.
type VersionInfo struct {
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Timestamp string `json:"timestamp"`
}
