package main

import (
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information",
	Long: `Display the current version, build information, and build datetime of the AI Rules Manager tool.

This information is useful for:
- Verifying which version is installed
- Debugging compatibility issues with specific builds
- Checking if updates are available
- Reporting issues with precise version context
- Understanding when the binary was built (useful for troubleshooting time-sensitive issues)`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := armService.ShowVersion(); err != nil {
			cmd.PrintErrln("Error:", err)
			return
		}
	},
}
