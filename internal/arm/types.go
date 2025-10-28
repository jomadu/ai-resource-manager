package arm

import "github.com/jomadu/ai-rules-manager/internal/ui"

// Type aliases to avoid duplication
type OutdatedRuleset = ui.OutdatedRuleset
type OutdatedPromptset = ui.OutdatedPromptset
type ManifestInfo = ui.ManifestInfo
type InstallationInfo = ui.InstallationInfo
type RulesetInfo = ui.RulesetInfo
type PromptsetInfo = ui.PromptsetInfo
type CompileStats = ui.CompileStats

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
