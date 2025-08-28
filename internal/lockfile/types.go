package lockfile

// Entry represents a single ruleset entry in the lock file.
type Entry struct {
	URL        string   `json:"url"`
	Type       string   `json:"type"`
	Constraint string   `json:"constraint"`
	Resolved   string   `json:"resolved"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
}

// LockFile represents the arm.lock file structure.
type LockFile struct {
	Rulesets map[string]map[string]Entry `json:"rulesets"`
}
