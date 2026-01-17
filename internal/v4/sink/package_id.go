package sink

import (
	"fmt"
	"strings"
)

// pkgKey creates a package key in format "registry/name@version"
func pkgKey(registry, name, version string) string {
	return fmt.Sprintf("%s/%s@%s", registry, name, version)
}

// parsePkgKey parses a package key and returns registry, name, version
func parsePkgKey(key string) (registry, name, version string, err error) {
	parts := strings.Split(key, "@")
	if len(parts) != 2 {
		return "", "", "", fmt.Errorf("invalid package key format: %s", key)
	}

	version = parts[1]
	registryAndName := parts[0]

	nameParts := strings.Split(registryAndName, "/")
	if len(nameParts) != 2 {
		return "", "", "", fmt.Errorf("invalid package key format: %s", key)
	}

	return nameParts[0], nameParts[1], version, nil
}
