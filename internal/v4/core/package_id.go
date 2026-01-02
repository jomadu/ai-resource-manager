package core

import (
	"fmt"
	"strings"
)

// PackageID creates a package ID in format "registry/name@version"
func PackageID(registry, name, version string) string {
	return fmt.Sprintf("%s/%s@%s", registry, name, version)
}

// ParsePackageID parses a package ID and returns registry, name, version
func ParsePackageID(id string) (registry, name, version string, err error) {
	// Split on @ to get version
	parts := strings.Split(id, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid package ID format: %s", id)
	}

	version = parts[1]
	registryAndName := parts[0]

	// Split on / to get registry and name
	nameParts := strings.Split(registryAndName, "/")
	if len(nameParts) != 2 {
		return "", "", "", fmt.Errorf("invalid package ID format: %s", id)
	}

	return nameParts[0], nameParts[1], version, nil
}