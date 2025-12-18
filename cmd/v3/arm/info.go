package main

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/v3/arm"
	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Show detailed information",
	Long:  "Show detailed information about registries, sinks, rulesets, and promptsets",
	Run: func(cmd *cobra.Command, args []string) {
		infoAll()
	},
}

var infoRegistryCmd = &cobra.Command{
	Use:   "registry [NAME]...",
	Short: "Show registry information",
	Long:  "Display detailed information about one or more registries.",
	Run: func(cmd *cobra.Command, args []string) {
		infoRegistries(args)
	},
}

var infoSinkCmd = &cobra.Command{
	Use:   "sink [NAME]...",
	Short: "Show sink information",
	Long:  "Display detailed information about one or more sinks.",
	Run: func(cmd *cobra.Command, args []string) {
		infoSinks(args)
	},
}

var infoRulesetCmd = &cobra.Command{
	Use:   "ruleset [REGISTRY_NAME/RULESET_NAME...]",
	Short: "Show ruleset information",
	Long:  "Display detailed information about one or more rulesets.",
	Run: func(cmd *cobra.Command, args []string) {
		infoRulesets(args)
	},
}

var infoPromptsetCmd = &cobra.Command{
	Use:   "promptset [REGISTRY_NAME/PROMPTSET_NAME...]",
	Short: "Show promptset information",
	Long:  "Display detailed information about one or more promptsets.",
	Run: func(cmd *cobra.Command, args []string) {
		infoPromptsets(args)
	},
}

func init() {
	// Add subcommands
	infoCmd.AddCommand(infoRegistryCmd)
	infoCmd.AddCommand(infoSinkCmd)
	infoCmd.AddCommand(infoRulesetCmd)
	infoCmd.AddCommand(infoPromptsetCmd)
}

func infoRegistries(names []string) {
	allRegistries, err := armService.GetRegistries(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}

	if len(names) == 0 {
		printRegistriesInfo(allRegistries)
		return
	}

	filtered := make(map[string]map[string]interface{})
	for _, name := range names {
		if config, exists := allRegistries[name]; exists {
			filtered[name] = config
		}
	}
	printRegistriesInfo(filtered)
}

func infoSinks(names []string) {
	allSinks, err := armService.GetSinks(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}

	if len(names) == 0 {
		printSinksInfo(allSinks)
		return
	}

	filtered := make(map[string]manifest.SinkConfig)
	for _, name := range names {
		if config, exists := allSinks[name]; exists {
			filtered[name] = config
		}
	}
	printSinksInfo(filtered)
}

func infoRulesets(names []string) {
	var rulesets []*arm.RulesetInfo
	var err error

	if len(names) == 0 {
		rulesets, err = armService.GetRulesets(ctx)
	} else {
		rulesets, err = armService.GetRulesetsByNames(ctx, names)
	}

	if err != nil {
		handleCommandError(err)
		return
	}
	printRulesetsInfo(rulesets)
}

func infoPromptsets(names []string) {
	var promptsets []*arm.PromptsetInfo
	var err error

	if len(names) == 0 {
		promptsets, err = armService.GetPromptsets(ctx)
	} else {
		promptsets, err = armService.GetPromptsetsByNames(ctx, names)
	}

	if err != nil {
		handleCommandError(err)
		return
	}
	printPromptsetsInfo(promptsets)
}

func infoAll() {
	registries, sinks, rulesets, promptsets, err := armService.GetAllResources(ctx)
	if err != nil {
		handleCommandError(err)
		return
	}
	printAllInfo(registries, sinks, rulesets, promptsets)
}

func printRegistriesInfo(registries map[string]map[string]interface{}) {
	if len(registries) == 0 {
		fmt.Println("No registries configured")
		return
	}
	for name, config := range registries {
		fmt.Printf("%s:\n", name)
		if regType, ok := config["type"].(string); ok {
			fmt.Printf("    type: %s\n", regType)
		}
		if url, ok := config["url"].(string); ok {
			fmt.Printf("    url: %s\n", url)
		}
		if groupID, ok := config["group_id"].(string); ok && groupID != "" {
			fmt.Printf("    group_id: %s\n", groupID)
		}
		if projectID, ok := config["project_id"].(string); ok && projectID != "" {
			fmt.Printf("    project_id: %s\n", projectID)
		}
		if owner, ok := config["owner"].(string); ok && owner != "" {
			fmt.Printf("    owner: %s\n", owner)
		}
		if repo, ok := config["repository"].(string); ok && repo != "" {
			fmt.Printf("    repository: %s\n", repo)
		}
		if branches, ok := config["branches"].([]interface{}); ok && len(branches) > 0 {
			fmt.Printf("    branches:\n")
			for _, branch := range branches {
				fmt.Printf("        - %v\n", branch)
			}
		}
		if apiVersion, ok := config["api_version"].(string); ok && apiVersion != "" {
			fmt.Printf("    api_version: %s\n", apiVersion)
		}
	}
}

func printSinksInfo(sinks map[string]manifest.SinkConfig) {
	if len(sinks) == 0 {
		fmt.Println("No sinks configured")
		return
	}
	for name, config := range sinks {
		fmt.Printf("%s:\n", name)
		fmt.Printf("    directory: %s\n", config.Directory)
		layout := config.GetLayout()
		if layout == "" {
			layout = "hierarchical"
		}
		fmt.Printf("    layout: %s\n", layout)
		fmt.Printf("    compileTarget: %s\n", string(config.CompileTarget))
	}
}

func printRulesetsInfo(rulesets []*arm.RulesetInfo) {
	if len(rulesets) == 0 {
		fmt.Println("No rulesets installed")
		return
	}
	registryGroups := make(map[string][]*arm.RulesetInfo)
	for _, ruleset := range rulesets {
		registryGroups[ruleset.Registry] = append(registryGroups[ruleset.Registry], ruleset)
	}
	for registry, groupRulesets := range registryGroups {
		fmt.Printf("%s:\n", registry)
		for _, ruleset := range groupRulesets {
			fmt.Printf("    %s:\n", ruleset.Name)
			fmt.Printf("        version: %s\n", ruleset.Installation.Version)
			fmt.Printf("        constraint: %s\n", ruleset.Manifest.Constraint)
			fmt.Printf("        priority: %d\n", ruleset.Manifest.Priority)
			if len(ruleset.Manifest.Sinks) > 0 {
				fmt.Printf("        sinks:\n")
				for _, sink := range ruleset.Manifest.Sinks {
					fmt.Printf("            - %s\n", sink)
				}
			}
			if len(ruleset.Manifest.Include) > 0 {
				fmt.Printf("        include:\n")
				for _, pattern := range ruleset.Manifest.Include {
					fmt.Printf("            - %q\n", pattern)
				}
			}
			if len(ruleset.Manifest.Exclude) > 0 {
				fmt.Printf("        exclude:\n")
				for _, pattern := range ruleset.Manifest.Exclude {
					fmt.Printf("            - %q\n", pattern)
				}
			}
		}
	}
}

func printPromptsetsInfo(promptsets []*arm.PromptsetInfo) {
	if len(promptsets) == 0 {
		fmt.Println("No promptsets installed")
		return
	}
	registryGroups := make(map[string][]*arm.PromptsetInfo)
	for _, promptset := range promptsets {
		registryGroups[promptset.Registry] = append(registryGroups[promptset.Registry], promptset)
	}
	for registry, groupPromptsets := range registryGroups {
		fmt.Printf("%s:\n", registry)
		for _, promptset := range groupPromptsets {
			fmt.Printf("    %s:\n", promptset.Name)
			fmt.Printf("        version: %s\n", promptset.Installation.Version)
			fmt.Printf("        constraint: %s\n", promptset.Manifest.Constraint)
			if len(promptset.Manifest.Sinks) > 0 {
				fmt.Printf("        sinks:\n")
				for _, sink := range promptset.Manifest.Sinks {
					fmt.Printf("            - %s\n", sink)
				}
			}
			if len(promptset.Manifest.Include) > 0 {
				fmt.Printf("        include:\n")
				for _, pattern := range promptset.Manifest.Include {
					fmt.Printf("            - %q\n", pattern)
				}
			}
			if len(promptset.Manifest.Exclude) > 0 {
				fmt.Printf("        exclude:\n")
				for _, pattern := range promptset.Manifest.Exclude {
					fmt.Printf("            - %q\n", pattern)
				}
			}
		}
	}
}

func printAllInfo(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig, rulesets []*arm.RulesetInfo, promptsets []*arm.PromptsetInfo) {
	if len(registries) > 0 {
		fmt.Println("registries:")
		for name, config := range registries {
			fmt.Printf("    %s:\n", name)
			if regType, ok := config["type"].(string); ok {
				fmt.Printf("        type: %s\n", regType)
			}
			if url, ok := config["url"].(string); ok {
				fmt.Printf("        url: %s\n", url)
			}
			if groupID, ok := config["group_id"].(string); ok && groupID != "" {
				fmt.Printf("        group_id: %s\n", groupID)
			}
			if projectID, ok := config["project_id"].(string); ok && projectID != "" {
				fmt.Printf("        project_id: %s\n", projectID)
			}
			if owner, ok := config["owner"].(string); ok && owner != "" {
				fmt.Printf("        owner: %s\n", owner)
			}
			if repo, ok := config["repository"].(string); ok && repo != "" {
				fmt.Printf("        repository: %s\n", repo)
			}
			if branches, ok := config["branches"].([]interface{}); ok && len(branches) > 0 {
				fmt.Printf("        branches:\n")
				for _, branch := range branches {
					fmt.Printf("            - %v\n", branch)
				}
			}
		}
	}
	if len(sinks) > 0 {
		fmt.Println("sinks:")
		for name, config := range sinks {
			fmt.Printf("    %s:\n", name)
			fmt.Printf("        directory: %s\n", config.Directory)
			fmt.Printf("        layout: %s\n", config.Layout)
			fmt.Printf("        compileTarget: %s\n", string(config.CompileTarget))
		}
	}
	if len(rulesets) > 0 || len(promptsets) > 0 {
		fmt.Println("packages:")
		if len(rulesets) > 0 {
			fmt.Println("    rulesets:")
			rulesetGroups := make(map[string][]*arm.RulesetInfo)
			for _, ruleset := range rulesets {
				rulesetGroups[ruleset.Registry] = append(rulesetGroups[ruleset.Registry], ruleset)
			}
			for registry, groupRulesets := range rulesetGroups {
				fmt.Printf("        %s:\n", registry)
				for _, ruleset := range groupRulesets {
					fmt.Printf("            %s:\n", ruleset.Name)
					fmt.Printf("                version: %s\n", ruleset.Installation.Version)
					fmt.Printf("                constraint: %s\n", ruleset.Manifest.Constraint)
					fmt.Printf("                priority: %d\n", ruleset.Manifest.Priority)
					if len(ruleset.Manifest.Sinks) > 0 {
						fmt.Printf("                sinks:\n")
						for _, sink := range ruleset.Manifest.Sinks {
							fmt.Printf("                    - %s\n", sink)
						}
					}
					if len(ruleset.Manifest.Include) > 0 {
						fmt.Printf("                include:\n")
						for _, pattern := range ruleset.Manifest.Include {
							fmt.Printf("                    - %q\n", pattern)
						}
					}
					if len(ruleset.Manifest.Exclude) > 0 {
						fmt.Printf("                exclude:\n")
						for _, pattern := range ruleset.Manifest.Exclude {
							fmt.Printf("                    - %q\n", pattern)
						}
					}
				}
			}
		}
		if len(promptsets) > 0 {
			fmt.Println("    promptsets:")
			promptsetGroups := make(map[string][]*arm.PromptsetInfo)
			for _, promptset := range promptsets {
				promptsetGroups[promptset.Registry] = append(promptsetGroups[promptset.Registry], promptset)
			}
			for registry, groupPromptsets := range promptsetGroups {
				fmt.Printf("        %s:\n", registry)
				for _, promptset := range groupPromptsets {
					fmt.Printf("            %s:\n", promptset.Name)
					fmt.Printf("                version: %s\n", promptset.Installation.Version)
					fmt.Printf("                constraint: %s\n", promptset.Manifest.Constraint)
					if len(promptset.Manifest.Sinks) > 0 {
						fmt.Printf("                sinks:\n")
						for _, sink := range promptset.Manifest.Sinks {
							fmt.Printf("                    - %s\n", sink)
						}
					}
					if len(promptset.Manifest.Include) > 0 {
						fmt.Printf("                include:\n")
						for _, pattern := range promptset.Manifest.Include {
							fmt.Printf("                    - %q\n", pattern)
						}
					}
					if len(promptset.Manifest.Exclude) > 0 {
						fmt.Printf("                exclude:\n")
						for _, pattern := range promptset.Manifest.Exclude {
							fmt.Printf("                    - %q\n", pattern)
						}
					}
				}
			}
		}
	}
	if len(registries) == 0 && len(sinks) == 0 && len(rulesets) == 0 && len(promptsets) == 0 {
		fmt.Println("No resources configured")
	}
}
