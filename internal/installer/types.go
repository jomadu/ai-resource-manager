package installer

// Installation represents an installed ruleset in a directory.
type Installation struct {
	Ruleset string `json:"ruleset"`
	Version string `json:"version"`
	Path    string `json:"path"`
}
