package main

import (
	"os"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/ui"
	"github.com/spf13/cobra"
)

var (
	armService arm.Service
	uiInstance *ui.UI
)

func main() {
	uiInstance = ui.New(false) // TODO: Add debug flag support
	armService = arm.NewArmService(uiInstance)

	if err := rootCmd.Execute(); err != nil {
		WriteError(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "arm",
	Short: "AI Rules Manager - Manage AI rule rulesets",
	Long:  "ARM helps you install, manage, and organize AI rule rulesets from various registries.",
}

func init() {
	rootCmd.AddCommand(newInstallCmd())
	rootCmd.AddCommand(newUninstallCmd())
	rootCmd.AddCommand(newUpdateCmd())
	rootCmd.AddCommand(newOutdatedCmd())
	rootCmd.AddCommand(newListCmd())
	rootCmd.AddCommand(newInfoCmd())
	rootCmd.AddCommand(newConfigCmd())
	rootCmd.AddCommand(newCacheCmd())
	rootCmd.AddCommand(newVersionCmd())
}
