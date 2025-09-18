package arm

// OutdatedRuleset represents a ruleset that has newer versions available.
type OutdatedRuleset struct {
	RulesetInfo *RulesetInfo `json:"rulesetInfo"`
	Wanted      string       `json:"wanted"`
	Latest      string       `json:"latest"`
}

// ManifestInfo contains information from the manifest file.
type ManifestInfo struct {
	Constraint string   `json:"constraint"`
	Priority   int      `json:"priority"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// InstallationInfo contains information about the actual installation.
type InstallationInfo struct {
	Version        string   `json:"version"`
	InstalledPaths []string `json:"installedPaths"`
}

// RulesetInfo provides detailed information about a ruleset.
type RulesetInfo struct {
	Registry     string           `json:"registry"`
	Name         string           `json:"name"`
	Manifest     ManifestInfo     `json:"manifest"`
	Installation InstallationInfo `json:"installation"`
}

// InstallRequest groups install parameters to avoid repetitive parameter passing.
type InstallRequest struct {
	Registry string
	Ruleset  string
	Version  string
	Priority int
	Include  []string
	Exclude  []string
	Sinks    []string
}
