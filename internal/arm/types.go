package arm

import "github.com/jomadu/ai-rules-manager/internal/ui"

// Type aliases to avoid duplication
type OutdatedRuleset = ui.OutdatedRuleset
type ManifestInfo = ui.ManifestInfo
type InstallationInfo = ui.InstallationInfo
type RulesetInfo = ui.RulesetInfo
type CompileStats = ui.CompileStats

// InstallRulesetRequest groups ruleset install parameters to avoid repetitive parameter passing.
type InstallRulesetRequest struct {
	Registry string
	Ruleset  string
	Version  string
	Priority int
	Include  []string
	Exclude  []string
	Sinks    []string
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
	Files        []string
	Target       string
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
