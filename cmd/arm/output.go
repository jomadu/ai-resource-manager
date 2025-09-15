package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/version"
	"github.com/olekukonko/tablewriter"
)

// FormatRulesetInfo formats ruleset information for display
func FormatRulesetInfo(info *arm.RulesetInfo, detailed bool) {
	if detailed {
		fmt.Printf("Ruleset: %s/%s@%s (%s)\n", info.Registry, info.Name, info.Installation.Version, info.Manifest.Constraint)
		fmt.Println("include:")
		for _, pattern := range info.Manifest.Include {
			fmt.Printf("  - %s\n", pattern)
		}
		if len(info.Manifest.Exclude) > 0 {
			fmt.Println("exclude:")
			for _, pattern := range info.Manifest.Exclude {
				fmt.Printf("  - %s\n", pattern)
			}
		}
		fmt.Println("Installed:")
		for _, path := range info.Installation.InstalledPaths {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Println("Sinks:")
		for _, sink := range info.Manifest.Sinks {
			fmt.Printf("  - %s\n", sink)
		}
		fmt.Printf("Priority: %d\n", info.Manifest.Priority)
		fmt.Printf("Constraint: %s\n", info.Manifest.Constraint)
		fmt.Printf("Resolved: %s\n", info.Installation.Version)
	} else {
		fmt.Printf("%s/%s@%s (%s)\n", info.Registry, info.Name, info.Installation.Version, info.Manifest.Constraint)
		fmt.Println("  include:")
		for _, pattern := range info.Manifest.Include {
			fmt.Printf("    - %s\n", pattern)
		}
		if len(info.Manifest.Exclude) > 0 {
			fmt.Println("  exclude:")
			for _, pattern := range info.Manifest.Exclude {
				fmt.Printf("    - %s\n", pattern)
			}
		}
		fmt.Println("  Installed:")
		for _, path := range info.Installation.InstalledPaths {
			fmt.Printf("    - %s\n", path)
		}
		fmt.Println("  Sinks:")
		for _, sink := range info.Manifest.Sinks {
			fmt.Printf("    - %s\n", sink)
		}
		fmt.Printf("  Priority: %d | Constraint: %s | Resolved: %s\n", info.Manifest.Priority, info.Manifest.Constraint, info.Installation.Version)
	}
}

// FormatInstalledRulesets formats the list of installed rulesets
func FormatInstalledRulesets(rulesets []*arm.RulesetInfo, showPriority, sortPriority bool) {
	// Sort by priority if requested
	if sortPriority {
		// Sort by priority (highest first), then by name for stable ordering
		for i := 0; i < len(rulesets)-1; i++ {
			for j := i + 1; j < len(rulesets); j++ {
				if rulesets[i].Manifest.Priority < rulesets[j].Manifest.Priority {
					rulesets[i], rulesets[j] = rulesets[j], rulesets[i]
				} else if rulesets[i].Manifest.Priority == rulesets[j].Manifest.Priority {
					// Same priority, sort alphabetically
					iName := rulesets[i].Registry + "/" + rulesets[i].Name
					jName := rulesets[j].Registry + "/" + rulesets[j].Name
					if iName > jName {
						rulesets[i], rulesets[j] = rulesets[j], rulesets[i]
					}
				}
			}
		}
	}

	for _, ruleset := range rulesets {
		if ruleset.Manifest.Constraint != "" {
			fmt.Printf("%s/%s@%s (%s)", ruleset.Registry, ruleset.Name, ruleset.Installation.Version, ruleset.Manifest.Constraint)
		} else {
			fmt.Printf("%s/%s@%s", ruleset.Registry, ruleset.Name, ruleset.Installation.Version)
		}
		if showPriority {
			fmt.Printf(" (priority: %d)", ruleset.Manifest.Priority)
		}
		if len(ruleset.Manifest.Sinks) > 0 {
			fmt.Printf(" (sinks: %v)", ruleset.Manifest.Sinks)
		}
		fmt.Println()
	}
}

// FormatOutdatedRulesets formats outdated rulesets as table or JSON
func FormatOutdatedRulesets(outdated []arm.OutdatedRuleset, format string) error {
	if len(outdated) == 0 {
		if format == "json" {
			fmt.Println("[]")
		} else {
			fmt.Println("All rulesets are up to date!")
		}
		return nil
	}

	if format == "json" {
		return FormatJSON(outdated)
	}

	return FormatOutdatedTable(outdated)
}

// FormatOutdatedTable formats outdated rulesets as a table
func FormatOutdatedTable(outdated []arm.OutdatedRuleset) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header("Registry", "Ruleset", "Constraint", "Current", "Wanted", "Latest")

	for _, r := range outdated {
		if err := table.Append(r.RulesetInfo.Registry, r.RulesetInfo.Name, r.RulesetInfo.Manifest.Constraint, r.RulesetInfo.Installation.Version, r.Wanted, r.Latest); err != nil {
			return fmt.Errorf("failed to add table row: %w", err)
		}
	}

	return table.Render()
}

// FormatJSON formats any data structure as JSON
func FormatJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// FormatVersionInfo formats version information
func FormatVersionInfo(versionInfo version.VersionInfo) {
	fmt.Printf("arm %s\n", versionInfo.Version)
	if versionInfo.Commit != "" {
		fmt.Printf("commit: %s\n", versionInfo.Commit)
	}
	if versionInfo.Arch != "" {
		fmt.Printf("arch: %s\n", versionInfo.Arch)
	}
}

// WriteError writes an error message to stderr
func WriteError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// WriteErrorf writes a formatted error message to stderr
func WriteErrorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+format+"\n", args...)
}
