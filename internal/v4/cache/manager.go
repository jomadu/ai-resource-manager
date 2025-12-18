package cache

import (
	"context"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type PackageRegistryCacheManager interface {
	ListPackageVersions(ctx context.Context, key string) ([]core.Version, error)
	GetPackageVersion(ctx context.Context, key string, version core.Version) ([]*core.File, error)
	SetPackageVersion(ctx context.Context, key string, version core.Version, files []*core.File) error
	RemoveOldPackagesVersions(ctx context.Context, maxAge time.Duration) error
	RemoveUnusedPackagesVersions(ctx context.Context, maxTimeSinceLastAccess time.Duration) error
}
