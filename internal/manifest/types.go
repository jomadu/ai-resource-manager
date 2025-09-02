package manifest

// RegistryConfig defines a registry configuration.
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

// Entry represents a single ruleset entry in the manifest.
type Entry struct {
	Version string   `json:"version"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// Manifest represents the arm.json manifest file structure.
type Manifest struct {
	Registries map[string]RegistryConfig   `json:"registries"`
	Rulesets   map[string]map[string]Entry `json:"rulesets"`
}
