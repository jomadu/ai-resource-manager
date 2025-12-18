package registry

import (
	"context"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type Registry interface {
	ListPackages(ctx context.Context) ([]*core.PackageMetadata, error)
	ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error)
	GetPackage(ctx context.Context, packageName string, version core.Version, include []string, exclude []string) (*core.Package, error)
}
