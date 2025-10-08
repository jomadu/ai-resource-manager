package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newCleanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean utilities",
		Long:  "Clean utilities for ARM including cache and sink cleanup operations.",
	}

	cmd.AddCommand(cleanCacheCmd)
	cmd.AddCommand(cleanSinksCmd)

	return cmd
}

var cleanCacheCmd = &cobra.Command{
	Use:   "cache [--nuke | --max-age DURATION]",
	Short: "Clean the local cache directory",
	Long: `Clean the local cache directory.

This command removes cached registry data and downloaded packages from the local cache. The --nuke flag performs a more aggressive cleanup, removing all cached data including registry indexes and package archives. The --max-age flag allows you to specify how old cached data should be before it's removed.

Flags:
  --nuke     Aggressive cleanup (remove all cached data)
  --max-age  Remove cached data older than specified duration (e.g., "30m", "2h", "7d")

Examples:
  # Standard cache cleanup (removes data older than 7 days)
  arm clean cache

  # Remove data older than 2 hours
  arm clean cache --max-age 2h

  # Remove data older than 30 minutes
  arm clean cache --max-age 30m

  # Remove data older than 1 day and 6 hours
  arm clean cache --max-age 1d6h

  # Aggressive cleanup (remove all cached data)
  arm clean cache --nuke`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nuke, _ := cmd.Flags().GetBool("nuke")
		maxAgeStr, _ := cmd.Flags().GetString("max-age")

		// Validate mutual exclusivity
		if nuke && maxAgeStr != "" {
			return fmt.Errorf("--nuke and --max-age are mutually exclusive")
		}

		if nuke {
			return armService.NukeCache(context.Background())
		}

		// Parse duration string to time.Duration
		var maxAge time.Duration
		if maxAgeStr != "" {
			var err error
			maxAge, err = parseDuration(maxAgeStr)
			if err != nil {
				return fmt.Errorf("invalid max-age format: %w", err)
			}
		} else {
			// Default to 7 days if no max-age specified
			maxAge = 7 * 24 * time.Hour
		}

		return armService.CleanCacheWithAge(context.Background(), maxAge)
	},
}

var cleanSinksCmd = &cobra.Command{
	Use:   "sinks [--nuke]",
	Short: "Clean sink directories based on the ARM index",
	Long: `Clean sink directories based on the ARM index.

This command removes files from sink directories that shouldn't be there according to the arm-index.json file. The --nuke flag performs a more aggressive cleanup, clearing out the entire ARM directory entirely. Without the flag, it performs a selective cleanup based on the index.

Flags:
  --nuke  Complete cleanup (remove entire ARM directory)

Examples:
  # Selective cleanup based on ARM index
  arm clean sinks

  # Complete cleanup (remove entire ARM directory)
  arm clean sinks --nuke`,
	RunE: func(cmd *cobra.Command, args []string) error {
		nuke, _ := cmd.Flags().GetBool("nuke")
		if nuke {
			return armService.NukeSinks(context.Background())
		}
		return armService.CleanSinks(context.Background())
	},
}

func init() {
	// Add flags to clean commands
	cleanCacheCmd.Flags().Bool("nuke", false, "Aggressive cleanup (remove all cached data)")
	cleanCacheCmd.Flags().String("max-age", "", "Remove cached data older than specified duration (e.g., '30m', '2h', '7d')")
	cleanSinksCmd.Flags().Bool("nuke", false, "Complete cleanup (remove entire ARM directory)")
}
