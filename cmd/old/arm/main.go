package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jomadu/ai-rules-manager/internal/v3/arm"
	"github.com/spf13/cobra"
)

var (
	armService arm.Service
	ctx        context.Context
)

var rootCmd = &cobra.Command{
	Use:   "arm",
	Short: "AI Resource Manager - Manage rulesets and promptsets for AI coding assistants",
	Long: `AI Resource Manager (ARM) is a tool for managing rulesets and promptsets
for AI coding assistants like Cursor, GitHub Copilot, and Amazon Q.

ARM allows you to:
- Manage registries (Git, GitLab, Cloudsmith)
- Configure sinks (output destinations)
- Install and manage rulesets and promptsets
- Compile source files to platform-specific formats
- Clean cache and manage installations

For more information, visit: https://github.com/jomadu/ai-resource-manager`,
}

func init() {
	// Initialize context
	ctx = context.Background()

	// Initialize ARM service
	armService = arm.NewArmService()

	// Add all subcommands
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(removeCmd)
	rootCmd.AddCommand(setCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(infoCmd)
	rootCmd.AddCommand(installCmd)
	rootCmd.AddCommand(uninstallCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(upgradeCmd)
	rootCmd.AddCommand(outdatedCmd)
	rootCmd.AddCommand(cleanCmd)
	rootCmd.AddCommand(compileCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
