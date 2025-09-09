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
		fmt.Printf("Ruleset: %s/%s@%s (%s)\n", info.Registry, info.Name, info.Resolved, info.Constraint)
		fmt.Println("include:")
		for _, pattern := range info.Include {
			fmt.Printf("  - %s\n", pattern)
		}
		if len(info.Exclude) > 0 {
			fmt.Println("exclude:")
			for _, pattern := range info.Exclude {
				fmt.Printf("  - %s\n", pattern)
			}
		}
		fmt.Println("Installed:")
		for _, path := range info.InstalledPaths {
			fmt.Printf("  - %s\n", path)
		}
		fmt.Println("Sinks:")
		for _, sink := range info.Sinks {
			fmt.Printf("  - %s\n", sink)
		}
		fmt.Printf("Constraint: %s\n", info.Constraint)
		fmt.Printf("Resolved: %s\n", info.Resolved)
	} else {
		fmt.Printf("%s/%s@%s (%s)\n", info.Registry, info.Name, info.Resolved, info.Constraint)
		fmt.Println("  include:")
		for _, pattern := range info.Include {
			fmt.Printf("    - %s\n", pattern)
		}
		if len(info.Exclude) > 0 {
			fmt.Println("  exclude:")
			for _, pattern := range info.Exclude {
				fmt.Printf("    - %s\n", pattern)
			}
		}
		fmt.Println("  Installed:")
		for _, path := range info.InstalledPaths {
			fmt.Printf("    - %s\n", path)
		}
		fmt.Println("  Sinks:")
		for _, sink := range info.Sinks {
			fmt.Printf("    - %s\n", sink)
		}
		fmt.Printf("  Constraint: %s | Resolved: %s\n", info.Constraint, info.Resolved)
	}
}

// FormatInstalledRulesets formats the list of installed rulesets
func FormatInstalledRulesets(rulesets []arm.InstalledRuleset) {
	for i := range rulesets {
		ruleset := &rulesets[i]
		if ruleset.Constraint != "" {
			fmt.Printf("%s/%s@%s (%s)\n", ruleset.Registry, ruleset.Name, ruleset.Version, ruleset.Constraint)
		} else {
			fmt.Printf("%s/%s@%s\n", ruleset.Registry, ruleset.Name, ruleset.Version)
		}
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
		if err := table.Append(r.Registry, r.Name, r.Constraint, r.Current, r.Wanted, r.Latest); err != nil {
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
