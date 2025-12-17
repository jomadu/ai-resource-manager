package cache

import (
	"context"
	"time"

	"github.com/jomadu/ai-rules-manager/internal/types"
)

// RegistryPackageCache provides registry-scoped storage for cached packages.
type RegistryPackageCache interface {
	ListVersions(ctx context.Context, keyObj interface{}) ([]string, error)
	GetPackageVersion(ctx context.Context, keyObj interface{}, version string) ([]types.File, error)
	SetPackageVersion(ctx context.Context, keyObj interface{}, version string, files []types.File) error
	Cleanup(maxAge time.Duration) error
}
