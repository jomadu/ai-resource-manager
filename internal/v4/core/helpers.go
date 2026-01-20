package core

import (
	"errors"
	"strings"
)

// PackageKey creates a package key from registry and package names
func PackageKey(registryName, packageName string) string {
	return registryName + "/" + packageName
}

// ParsePackageKey splits a package key into registry and package names
func ParsePackageKey(key string) (registryName, packageName string) {
	parts := strings.SplitN(key, "/", 2)
	if len(parts) != 2 {
		return "", ""
	}
	return parts[0], parts[1]
}

func GetBestMatching(versions []Version, constraint Constraint) (*Version, error) {
	if len(versions) == 0 {
		return nil, errors.New("no versions available")

	}
	var candidates []*Version
	for i := range versions {
		satisfied, err := constraint.IsSatisfiedBy(versions[i])
		if err != nil {
			return nil, err
		}
		if satisfied {
			candidates = append(candidates, &versions[i])
		}
	}

	if len(candidates) == 0 {
		return nil, errors.New("no matching version found")
	}

	// Find highest - versions know how to compare
	best := candidates[0]
	for _, candidate := range candidates[1:] {
		if candidate.IsNewerThan(*best) {
			best = candidate
		}
	}

	return best, nil
}
