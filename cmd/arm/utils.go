package main

import "strings"

// parseRegistryPackage parses a package name like "registry/package"
func parseRegistryPackage(packageName string) (registry, pkgName string) {
	parts := strings.Split(packageName, "/")
	if len(parts) != 2 {
		// TODO: Handle error
		return "", ""
	}

	return parts[0], parts[1]
}
