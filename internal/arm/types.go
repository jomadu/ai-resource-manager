package arm

import "github.com/jomadu/ai-rules-manager/internal/ui"

// Type aliases to avoid duplication
type OutdatedRuleset = ui.OutdatedRuleset
type ManifestInfo = ui.ManifestInfo
type InstallationInfo = ui.InstallationInfo
type RulesetInfo = ui.RulesetInfo

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
