package core

import "errors"

func GetBestMatching(versions []Version, constraint Constraint) (*Version, error) {
	if len(versions) == 0 {
		return nil, errors.New("no versions available")

	}
	var candidates []*Version
	for i := range versions {
		if constraint.IsSatisfiedBy(versions[i]) {
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
