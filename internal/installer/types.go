package installer

// Ruleset represents an installed ruleset in a directory.
type Ruleset struct {
	Registry  string
	Ruleset   string
	Version   string
	Priority  int
	Path      string
	FilePaths []string
}
