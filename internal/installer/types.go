package installer

// Ruleset represents an installed ruleset in a directory.
type Ruleset struct {
	Registry  string   `json:"registry"`
	Ruleset   string   `json:"ruleset"`
	Version   string   `json:"version"`
	Priority  int      `json:"priority"`
	Path      string   `json:"path"`
	FilePaths []string `json:"filePaths"`
}
