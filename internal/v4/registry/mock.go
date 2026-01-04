package registry

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
)

type MockRegistry struct {
	packages map[string]map[string]*core.Package
}

func NewMockRegistry() *MockRegistry {
	return &MockRegistry{
		packages: make(map[string]map[string]*core.Package),
	}
}

func (m *MockRegistry) AddPackage(pkg *core.Package) {
	if m.packages[pkg.Metadata.Name] == nil {
		m.packages[pkg.Metadata.Name] = make(map[string]*core.Package)
	}
	m.packages[pkg.Metadata.Name][pkg.Metadata.Version.Version] = pkg
}

func (m *MockRegistry) ListPackages(ctx context.Context) ([]*core.PackageMetadata, error) {
	var result []*core.PackageMetadata
	for _, versions := range m.packages {
		for _, pkg := range versions {
			result = append(result, &pkg.Metadata)
		}
	}
	return result, nil
}

func (m *MockRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	versions, exists := m.packages[packageName]
	if !exists {
		return nil, fmt.Errorf("package %s not found", packageName)
	}
	
	var result []core.Version
	for _, pkg := range versions {
		result = append(result, pkg.Metadata.Version)
	}
	return result, nil
}

func (m *MockRegistry) GetPackage(ctx context.Context, packageName string, version core.Version, include []string, exclude []string) (*core.Package, error) {
	versions, exists := m.packages[packageName]
	if !exists {
		return nil, fmt.Errorf("package %s not found", packageName)
	}
	
	pkg, exists := versions[version.Version]
	if !exists {
		return nil, fmt.Errorf("package %s version %s not found", packageName, version.Version)
	}
	
	return pkg, nil
}