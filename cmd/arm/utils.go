package main

import (
	"fmt"
	"os"
	"strings"
)

// handleCommandError outputs a formatted error message and exits with code 1
func handleCommandError(err error) {
	fmt.Fprintf(os.Stderr, "\nError: %v\n", err)
	os.Exit(1)
}

// parseRegistry extracts registry from "registry/package[@version]" format
func parseRegistry(input string) (string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: %s (expected registry/package)", input)
	}
	if parts[0] == "" {
		return "", fmt.Errorf("registry name cannot be empty")
	}
	return parts[0], nil
}

// parsePackage extracts package name from "registry/package[@version]" format
func parsePackage(input string) (string, error) {
	// First, get everything after registry/
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: %s (expected registry/package)", input)
	}
	
	pkgWithVersion := parts[1]
	if pkgWithVersion == "" {
		return "", fmt.Errorf("package name cannot be empty")
	}
	
	// Remove version if present
	if strings.Contains(pkgWithVersion, "@") {
		pkg := strings.SplitN(pkgWithVersion, "@", 2)[0]
		if pkg == "" {
			return "", fmt.Errorf("package name cannot be empty")
		}
		return pkg, nil
	}
	
	return pkgWithVersion, nil
}

// parseVersion extracts version from "registry/package@version" format
// Returns empty string if no version specified
func parseVersion(input string) (string, error) {
	// Find the @ symbol
	if !strings.Contains(input, "@") {
		return "", nil
	}
	
	parts := strings.SplitN(input, "@", 2)
	if len(parts) != 2 {
		return "", nil
	}
	
	version := parts[1]
	if version == "" {
		return "", fmt.Errorf("version cannot be empty after @")
	}
	
	return version, nil
}
