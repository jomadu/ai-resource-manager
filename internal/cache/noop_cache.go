package cache

import (
	"context"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// NoopRegistryPackageCache is a no-op implementation of RegistryPackageCache for testing
type NoopRegistryPackageCache struct{}

// NewNoopRegistryPackageCache creates a new no-op cache
func NewNoopRegistryPackageCache() *NoopRegistryPackageCache {
	return &NoopRegistryPackageCache{}
}

func (n *NoopRegistryPackageCache) ListVersions(ctx context.Context, keyObj interface{}) ([]string, error) {
	return []string{}, nil
}

func (n *NoopRegistryPackageCache) GetPackageVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error) {
	// Always return cache miss to force fresh fetches
	return nil, context.DeadlineExceeded // Standard cache miss error
}

func (n *NoopRegistryPackageCache) SetPackageVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error {
	// Do nothing - don't cache
	return nil
}

func (n *NoopRegistryPackageCache) Cleanup(maxAge time.Duration) error {
	// Do nothing - no cache to cleanup
	return nil
}
