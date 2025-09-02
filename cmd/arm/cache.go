package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(cacheCleanCmd)
	cacheCmd.AddCommand(cacheNukeCmd)
}

var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cache",
	Long:  "Manage ARM cache including cleanup operations.",
}

var cacheCleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Clean cache",
	Long:  "Remove old cached versions based on age criteria.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()

		maxAgeStr, _ := cmd.Flags().GetString("max-age")
		maxAge, err := parseDuration(maxAgeStr)
		if err != nil {
			return fmt.Errorf("invalid max-age format: %w", err)
		}

		cacheManager := cache.NewManager()
		if err := cacheManager.CleanupOldVersions(ctx, maxAge); err != nil {
			return fmt.Errorf("failed to clean cache: %w", err)
		}

		fmt.Printf("Cache cleaned: removed versions older than %s\n", maxAgeStr)
		return nil
	},
}

var cacheNukeCmd = &cobra.Command{
	Use:   "nuke",
	Short: "Remove entire cache directory",
	Long:  "Remove the entire ~/.arm/cache directory and all cached data.",
	RunE: func(cmd *cobra.Command, args []string) error {
		cacheDir := cache.GetCacheDir()
		if err := os.RemoveAll(cacheDir); err != nil {
			return fmt.Errorf("failed to remove cache directory: %w", err)
		}

		fmt.Println("Cache directory removed successfully")
		return nil
	},
}

// parseDuration extends time.ParseDuration to support days (d)
func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		daysStr := strings.TrimSuffix(s, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("invalid days value: %s", daysStr)
		}
		return time.Duration(days) * 24 * time.Hour, nil
	}
	if strings.HasSuffix(s, "h") {
		hoursStr := strings.TrimSuffix(s, "h")
		hours, err := strconv.Atoi(hoursStr)
		if err != nil {
			return 0, fmt.Errorf("invalid hours value: %s", hoursStr)
		}
		return time.Duration(hours) * time.Hour, nil
	}
	if strings.HasSuffix(s, "m") {
		minutesStr := strings.TrimSuffix(s, "m")
		minutes, err := strconv.Atoi(minutesStr)
		if err != nil {
			return 0, fmt.Errorf("invalid minutes value: %s", minutesStr)
		}
		return time.Duration(minutes) * time.Minute, nil
	}
	return 0, fmt.Errorf("unsupported duration format: %s (use format like 30d, 24h, or 60m)", s)
}

func init() {
	cacheCleanCmd.Flags().String("max-age", "30d", "Remove versions not accessed within this duration (e.g., 30d, 7d, 24h)")
}
