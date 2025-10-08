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

// Promptset represents an installed promptset in a directory.
type Promptset struct {
	Registry  string
	Promptset string
	Version   string
	Path      string
	FilePaths []string
}
