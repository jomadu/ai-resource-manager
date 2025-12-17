package arm

type OutdatedPackage struct {
	Package    string `json:"package"`
	Type       string `json:"type"` // "ruleset" or "promptset"
	Constraint string `json:"constraint"`
	Current    string `json:"current"`
	Wanted     string `json:"wanted"`
	Latest     string `json:"latest"`
}

type ManifestInfo struct {
	Constraint string   `json:"constraint"`
	Priority   int      `json:"priority"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

type InstallationInfo struct {
	Version        string   `json:"version"`
	InstalledPaths []string `json:"installedPaths"`
}

type RulesetInfo struct {
	Registry     string           `json:"registry"`
	Name         string           `json:"name"`
	Manifest     ManifestInfo     `json:"manifest"`
	Installation InstallationInfo `json:"installation"`
}

type PromptsetInfo struct {
	Registry     string                `json:"registry"`
	Name         string                `json:"name"`
	Manifest     PromptsetManifestInfo `json:"manifest"`
	Installation InstallationInfo      `json:"installation"`
}

type PromptsetManifestInfo struct {
	Constraint string   `json:"constraint"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// CompileStats tracks compilation statistics
type CompileStats struct {
	FilesProcessed int `json:"filesProcessed"`
	FilesCompiled  int `json:"filesCompiled"`
	FilesSkipped   int `json:"filesSkipped"`
	RulesGenerated int `json:"rulesGenerated"`
	Errors         int `json:"errors"`
}
// InstallRulesetRequest groups ruleset install parameters to avoid repetitive parameter passing.
type InstallRulesetRequest struct {
	Registry string
	Ruleset  string
	Version  string
	Priority int // Defaults to 100 via constructor
	Include  []string
	Exclude  []string
	Sinks    []string
}

// NewInstallRulesetRequest creates a new request with default priority of 100
func NewInstallRulesetRequest(registry, ruleset, version string, sinks []string) *InstallRulesetRequest {
	return &InstallRulesetRequest{
		Registry: registry,
		Ruleset:  ruleset,
		Version:  version,
		Priority: 100, // Default priority
		Sinks:    sinks,
	}
}

// WithPriority sets a custom priority (for fluent API)
func (r *InstallRulesetRequest) WithPriority(priority int) *InstallRulesetRequest {
	r.Priority = priority
	return r
}

// WithInclude sets include patterns (for fluent API)
func (r *InstallRulesetRequest) WithInclude(include []string) *InstallRulesetRequest {
	r.Include = include
	return r
}

// WithExclude sets exclude patterns (for fluent API)
func (r *InstallRulesetRequest) WithExclude(exclude []string) *InstallRulesetRequest {
	r.Exclude = exclude
	return r
}

// InstallPromptsetRequest groups promptset install parameters.
type InstallPromptsetRequest struct {
	Registry  string
	Promptset string
	Version   string
	Include   []string
	Exclude   []string
	Sinks     []string
}

// CompileRequest groups compile parameters following ARM patterns
type CompileRequest struct {
	Paths        []string
	Targets      []string
	OutputDir    string
	Namespace    string
	Force        bool
	Recursive    bool
	Verbose      bool
	ValidateOnly bool
	Include      []string
	Exclude      []string
	FailFast     bool
}
