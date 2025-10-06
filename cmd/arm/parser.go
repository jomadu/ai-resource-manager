package main

import (
	"fmt"
	"strings"
)

// PackageRef represents a parsed package reference (ruleset or promptset)
type PackageRef struct {
	Registry string
	Name     string
	Version  string
}

// ParsePackageArg parses registry/package[@version] format
func ParsePackageArg(arg string) (PackageRef, error) {
	if arg == "" {
		return PackageRef{}, fmt.Errorf("package argument cannot be empty")
	}

	parts := strings.SplitN(arg, "/", 2)
	if len(parts) != 2 {
		return PackageRef{}, fmt.Errorf("invalid package format: %s (expected registry/package[@version])", arg)
	}

	registry := parts[0]
	packageAndVersion := parts[1]

	versionParts := strings.SplitN(packageAndVersion, "@", 2)
	packageName := versionParts[0]
	version := ""
	if len(versionParts) > 1 {
		version = versionParts[1]
	}

	if registry == "" {
		return PackageRef{}, fmt.Errorf("registry name cannot be empty")
	}
	if packageName == "" {
		return PackageRef{}, fmt.Errorf("package name cannot be empty")
	}

	return PackageRef{
		Registry: registry,
		Name:     packageName,
		Version:  version,
	}, nil
}

// ParsePackageArgs parses multiple package arguments
func ParsePackageArgs(args []string) ([]PackageRef, error) {
	refs := make([]PackageRef, len(args))
	for i, arg := range args {
		ref, err := ParsePackageArg(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to parse argument %d (%s): %w", i+1, arg, err)
		}
		refs[i] = ref
	}
	return refs, nil
}

// GetDefaultIncludePatterns returns default include patterns if none provided
func GetDefaultIncludePatterns(include []string) []string {
	if len(include) == 0 {
		return []string{"*.yml", "*.yaml"}
	}
	return include
}
