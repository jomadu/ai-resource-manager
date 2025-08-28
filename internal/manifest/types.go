package manifest

// Entry represents a single ruleset entry in the manifest.
type Entry struct {
	Version string   `json:"version"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// Manifest represents the arm.json manifest file structure.
type Manifest struct {
	Rulesets map[string]map[string]Entry `json:"rulesets"`
}
