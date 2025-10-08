package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

func newCacheCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cache",
		Short: "Manage cache",
		Long:  "Manage ARM cache including cleanup operations.",
	}

	cmd.AddCommand(cacheCleanCmd)
	cmd.AddCommand(cacheNukeCmd)

	return cmd
}

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache",
	Long:  "Remove old cached versions based on age criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		nuke, _ := cmd.Flags().GetBool("nuke")
		maxAgeStr, _ := cmd.Flags().GetString("max-age")

		// Validate mutual exclusivity
		if nuke && maxAgeStr != "" {
			return fmt.Errorf("--nuke and --max-age are mutually exclusive")
		}

		if nuke {
			return armService.NukeCache(ctx)
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

		return armService.CleanCacheWithAge(ctx, maxAge)
	},
}

var cacheNukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Remove entire cache directory",
	Long:  "Remove the entire ~/.arm/cache directory and all cached data.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		return armService.NukeCache(ctx)
	},
}

func init() {
	cacheCleanCmd.Flags().Bool("nuke", false, "Aggressive cleanup (remove all cached data)")
	cacheCleanCmd.Flags().String("max-age", "", "Remove cached data older than specified duration (e.g., '30m', '2h', '7d')")
}
