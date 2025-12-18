package main

import (
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/v3/version"
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
		info := version.GetVersionInfo()
		printVersionInfo(info)
	},
}

func printVersionInfo(info version.VersionInfo) {
	fmt.Printf("arm %s\n", info.Version)
	if info.Commit != "" {
		fmt.Printf("commit: %s\n", info.Commit)
	}
	if info.Arch != "" {
		fmt.Printf("arch: %s\n", info.Arch)
	}
	if info.Timestamp != "" {
		fmt.Printf("timestamp: %s\n", info.Timestamp)
	}
}
