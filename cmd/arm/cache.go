package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/spf13/cobra"
)

func init() {
	cacheCmd.AddCommand(cacheCleanCmd)
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
		maxAge, err := time.ParseDuration(maxAgeStr)
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

func init() {
	cacheCleanCmd.Flags().String("max-age", "30d", "Remove versions not accessed within this duration (e.g., 30d, 7d, 24h)")
}
