package main

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/v3/arm"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List resources",
	Long:  "List registries, sinks, rulesets, and promptsets",
	Run: func(cmd *cobra.Command, args []string) {
		listAll()
	},
}

var listRegistryCmd = &cobra.Command{
	Use:   "registry",
	Short: "List all registries",
	Long:  "List all configured registries.",
	Run: func(cmd *cobra.Command, args []string) {
		listRegistries()
	},
}

var listSinkCmd = &cobra.Command{
	Use:   "sink",
	Short: "List all sinks",
	Long:  "List all configured sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		listSinks()
	},
}

var listRulesetCmd = &cobra.Command{
	Use:   "ruleset",
	Short: "List all rulesets",
	Long:  "List all installed rulesets.",
	Run: func(cmd *cobra.Command, args []string) {
		listRulesets()
	},
}

var listPromptsetCmd = &cobra.Command{
	Use:   "promptset",
	Short: "List all promptsets",
	Long:  "List all installed promptsets.",
	Run: func(cmd *cobra.Command, args []string) {
		listPromptsets()
	},
}

func init() {
	// Add subcommands
	listCmd.AddCommand(listRegistryCmd)
	listCmd.AddCommand(listSinkCmd)
	listCmd.AddCommand(listRulesetCmd)
	listCmd.AddCommand(listPromptsetCmd)
}

func listRegistries() {
	registries, err := armService.GetRegistries(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printRegistriesList(registries)
}

func listSinks() {
	sinks, err := armService.GetSinks(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printSinksList(sinks)
}

func listRulesets() {
	rulesets, err := armService.GetRulesets(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printRulesetsList(rulesets)
}

func listPromptsets() {
	promptsets, err := armService.GetPromptsets(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printPromptsetsList(promptsets)
}

func listAll() {
	registries, sinks, rulesets, promptsets, err := armService.GetAllResources(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printAllList(registries, sinks, rulesets, promptsets)
}

func printRegistriesList(registries map[string]map[string]interface{}) {
	if len(registries) == 0 {
		fmt.Println("No registries configured")
		return
	}
	for name := range registries {
		fmt.Printf("- %s\n", name)
	}
}

func printSinksList(sinks map[string]manifest.SinkConfig) {
	if len(sinks) == 0 {
		fmt.Println("No sinks configured")
		return
	}
	for name := range sinks {
		fmt.Printf("- %s\n", name)
	}
}

func printRulesetsList(rulesets []*arm.RulesetInfo) {
	if len(rulesets) == 0 {
		fmt.Println("No rulesets installed")
		return
	}
	for _, ruleset := range rulesets {
		fmt.Printf("- %s/%s@%s\n", ruleset.Registry, ruleset.Name, ruleset.Installation.Version)
	}
}

func printPromptsetsList(promptsets []*arm.PromptsetInfo) {
	if len(promptsets) == 0 {
		fmt.Println("No promptsets installed")
		return
	}
	for _, promptset := range promptsets {
		fmt.Printf("- %s/%s@%s\n", promptset.Registry, promptset.Name, promptset.Installation.Version)
	}
}

func printAllList(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig, rulesets []*arm.RulesetInfo, promptsets []*arm.PromptsetInfo) {
	if len(registries) > 0 {
		fmt.Println("registries:")
		for name := range registries {
			fmt.Printf("    - %s\n", name)
		}
	}
	if len(sinks) > 0 {
		fmt.Println("sinks:")
		for name := range sinks {
			fmt.Printf("    - %s\n", name)
		}
	}
	if len(rulesets) > 0 {
		fmt.Println("rulesets:")
		for _, ruleset := range rulesets {
			fmt.Printf("    - %s/%s@%s\n", ruleset.Registry, ruleset.Name, ruleset.Installation.Version)
		}
	}
	if len(promptsets) > 0 {
		fmt.Println("promptsets:")
		for _, promptset := range promptsets {
			fmt.Printf("    - %s/%s@%s\n", promptset.Registry, promptset.Name, promptset.Installation.Version)
		}
	}
	if len(registries) == 0 && len(sinks) == 0 && len(rulesets) == 0 && len(promptsets) == 0 {
		fmt.Println("No resources configured")
	}
}
