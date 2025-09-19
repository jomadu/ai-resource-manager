package arm

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/cache"
)

func (a *ArmService) CleanCacheWithAge(ctx context.Context, maxAge time.Duration) error {
	cacheManager := cache.NewManager()
	if err := cacheManager.CleanupOldVersions(ctx, maxAge); err != nil {
		return fmt.Errorf("failed to clean cache: %w", err)
	}
	a.ui.Success(fmt.Sprintf("Cache cleaned: removed versions older than %v", maxAge))
	return nil
}

func (a *ArmService) NukeCache(ctx context.Context) error {
	cacheDir := cache.GetCacheDir()
	err := os.RemoveAll(cacheDir)
	if err != nil {
		return fmt.Errorf("failed to remove cache directory: %w", err)
	}
	a.ui.Success("Cache directory removed successfully")
	return nil
}
