package main

import (
	"time"

	"github.com/spf13/cobra"
)

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache and sinks",
	Long:  "Clean cache directories and sink directories",
}

var cleanCacheCmd = &cobra.Command{
	Use:   "cache [--nuke | --max-age DURATION]",
	Short: "Clean cache directory",
	Long: `Clean the local cache directory. This command removes cached registry data and downloaded packages from the local cache.

The --nuke flag performs a more aggressive cleanup, removing all cached data including registry indexes and package archives.
The --max-age flag allows you to specify how old cached data should be before it's removed.`,
	Run: func(cmd *cobra.Command, args []string) {
		cleanCache(cmd)
	},
}

var cleanSinksCmd = &cobra.Command{
	Use:   "sinks [--nuke]",
	Short: "Clean sink directories",
	Long: `Clean sink directories based on the ARM index. This command removes files from sink directories that shouldn't be there according to the arm-index.json file.

The --nuke flag performs a more aggressive cleanup, clearing out the entire ARM directory entirely.`,
	Run: func(cmd *cobra.Command, args []string) {
		cleanSinks(cmd)
	},
}

func init() {
	// Add subcommands
	cleanCmd.AddCommand(cleanCacheCmd)
	cleanCmd.AddCommand(cleanSinksCmd)

	// Add clean flags
	cleanCacheCmd.Flags().Bool("nuke", false, "Aggressive cleanup (remove all cached data)")
	cleanCacheCmd.Flags().String("max-age", "7d", "Remove cached data older than specified duration")
	cleanSinksCmd.Flags().Bool("nuke", false, "Complete cleanup (remove entire ARM directory)")

	// Mark flags as mutually exclusive
	cleanCacheCmd.MarkFlagsMutuallyExclusive("nuke", "max-age")
}

func cleanCache(cmd *cobra.Command) {
	nuke, _ := cmd.Flags().GetBool("nuke")
	maxAgeStr, _ := cmd.Flags().GetString("max-age")

	if nuke {
		if err := armService.NukeCache(ctx); err != nil {
			handleCommandError(err)
		}
	} else {
		// Parse duration
		maxAge, err := time.ParseDuration(maxAgeStr)
		if err != nil {
			handleCommandError(err)
		}

		if err := armService.CleanCacheWithAge(ctx, maxAge); err != nil {
			handleCommandError(err)
		}
	}
}

func cleanSinks(cmd *cobra.Command) {
	nuke, _ := cmd.Flags().GetBool("nuke")

	if nuke {
		if err := armService.NukeSinks(ctx); err != nil {
			handleCommandError(err)
		}
	} else {
		if err := armService.CleanSinks(ctx); err != nil {
			handleCommandError(err)
		}
	}
}
