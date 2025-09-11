package arm

// OutdatedRuleset represents a ruleset that has newer versions available.
type OutdatedRuleset struct {
	Registry   string `json:"registry"`
	Name       string `json:"name"`
	Constraint string `json:"constraint"`
	Current    string `json:"current"`
	Wanted     string `json:"wanted"`
	Latest     string `json:"latest"`
}

// InstalledRuleset represents a currently installed ruleset.
type InstalledRuleset struct {
	Registry   string   `json:"registry"`
	Name       string   `json:"name"`
	Version    string   `json:"version"`
	Constraint string   `json:"constraint"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// RulesetInfo provides detailed information about a ruleset.
type RulesetInfo struct {
	Registry       string   `json:"registry"`
	Name           string   `json:"name"`
	Include        []string `json:"include"`
	Exclude        []string `json:"exclude"`
	InstalledPaths []string `json:"installed_paths"`
	Sinks          []string `json:"sinks"`
	Constraint     string   `json:"constraint"`
	Resolved       string   `json:"resolved"`
}

// Installation represents an installed ruleset in a directory.
type Installation struct {
	Ruleset string `json:"ruleset"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

// InstallRequest groups install parameters to avoid repetitive parameter passing.
type InstallRequest struct {
	Registry string
	Ruleset  string
	Version  string
	Include  []string
	Exclude  []string
	Sinks    []string
}
