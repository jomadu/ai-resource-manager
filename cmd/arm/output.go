package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/version"
)

// FormatRulesetInfo formats ruleset information for display
func FormatRulesetInfo(info *arm.RulesetInfo, detailed bool) {
	uiInstance.RulesetInfo(info, detailed)
}

// FormatInstalledRulesets formats the list of installed rulesets
func FormatInstalledRulesets(rulesets []*arm.RulesetInfo, sortPriority bool) {
	// Sort by priority if requested
	if sortPriority {
		for i := 0; i < len(rulesets)-1; i++ {
			for j := i + 1; j < len(rulesets); j++ {
				if rulesets[i].Manifest.Priority < rulesets[j].Manifest.Priority {
					rulesets[i], rulesets[j] = rulesets[j], rulesets[i]
				} else if rulesets[i].Manifest.Priority == rulesets[j].Manifest.Priority {
					iName := rulesets[i].Registry + "/" + rulesets[i].Name
					jName := rulesets[j].Registry + "/" + rulesets[j].Name
					if iName > jName {
						rulesets[i], rulesets[j] = rulesets[j], rulesets[i]
					}
				}
			}
		}
	}
	uiInstance.RulesetList(rulesets)
}

// FormatOutdatedRulesets formats outdated rulesets as table or JSON
func FormatOutdatedRulesets(outdated []arm.OutdatedRuleset, format string) error {
	if format == "json" {
		return FormatJSON(outdated)
	}
	uiInstance.OutdatedTable(outdated, format)
	return nil
}

// FormatOutdatedTable formats outdated rulesets as a table
func FormatOutdatedTable(outdated []arm.OutdatedRuleset) {
	uiInstance.OutdatedTable(outdated, "table")
}

// FormatJSON formats any data structure as JSON
func FormatJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FormatVersionInfo formats version information
func FormatVersionInfo(versionInfo version.VersionInfo) {
	uiInstance.VersionInfo(versionInfo)
}

// WriteError writes an error message to stderr
func WriteError(err error) {
	uiInstance.Error(err)
}

// WriteErrorf writes a formatted error message to stderr
func WriteErrorf(format string, args ...interface{}) {
	uiInstance.Error(fmt.Errorf(format, args...))
}
