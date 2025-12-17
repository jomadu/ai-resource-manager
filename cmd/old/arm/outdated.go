package main

import (
	"encoding/json"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/v3/arm"
	"github.com/spf13/cobra"
)

var outdatedCmd = &cobra.Command{
	Use:   "outdated",
	Short: "Check for outdated resources",
	Long:  "Check for outdated rulesets and promptsets across configured registries",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedAll(cmd)
	},
}

var outdatedRulesetCmd = &cobra.Command{
	Use:   "ruleset [--output <table|json|list>]",
	Short: "Check for outdated rulesets",
	Long:  "Check for outdated rulesets across all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedRulesets(cmd)
	},
}

var outdatedPromptsetCmd = &cobra.Command{
	Use:   "promptset [--output <table|json|list>]",
	Short: "Check for outdated promptsets",
	Long:  "Check for outdated promptsets across all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		checkOutdatedPromptsets(cmd)
	},
}

func init() {
	// Add subcommands
	outdatedCmd.AddCommand(outdatedRulesetCmd)
	outdatedCmd.AddCommand(outdatedPromptsetCmd)

	// Add output format flags
	outdatedRulesetCmd.Flags().String("output", "table", "Output format (table, json, list)")
	outdatedPromptsetCmd.Flags().String("output", "table", "Output format (table, json, list)")

	// Add output format flag to main command
	outdatedCmd.Flags().String("output", "table", "Output format (table, json, list)")
}

func checkOutdatedAll(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	outdated, err := armService.GetOutdatedPackages(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printOutdatedPackages(outdated, output)
}

func checkOutdatedRulesets(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	outdated, err := armService.GetOutdatedPackages(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}

	// Filter to only rulesets
	rulesetsOnly := make([]*arm.OutdatedPackage, 0)
	for _, pkg := range outdated {
		if pkg.Type == "ruleset" {
			rulesetsOnly = append(rulesetsOnly, pkg)
		}
	}

	printOutdatedPackages(rulesetsOnly, output)
}

func checkOutdatedPromptsets(cmd *cobra.Command) {
	output, _ := cmd.Flags().GetString("output")

	outdated, err := armService.GetOutdatedPackages(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}

	// Filter to only promptsets
	promptsetsOnly := make([]*arm.OutdatedPackage, 0)
	for _, pkg := range outdated {
		if pkg.Type == "promptset" {
			promptsetsOnly = append(promptsetsOnly, pkg)
		}
	}

	printOutdatedPackages(promptsetsOnly, output)
}


func printOutdatedPackages(outdated []*arm.OutdatedPackage, outputFormat string) {
	if len(outdated) == 0 {
		fmt.Println("All packages are up to date!")
		return
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.Marshal(outdated)
		if err != nil {
			fmt.Printf("Error: Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	case "list":
		for _, pkg := range outdated {
			fmt.Println(pkg.Package)
		}
	default: // table format
		fmt.Printf("%-40s %-10s %-15s %-15s %-15s %-15s\n", "Package", "Type", "Constraint", "Current", "Wanted", "Latest")
		for _, pkg := range outdated {
			fmt.Printf("%-40s %-10s %-15s %-15s %-15s %-15s\n",
				pkg.Package, pkg.Type, pkg.Constraint, pkg.Current, pkg.Wanted, pkg.Latest)
		}
	}
}