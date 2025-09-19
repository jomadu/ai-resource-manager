package main

import (
	"fmt"
	"strings"
)

// RulesetRef represents a parsed ruleset reference
type RulesetRef struct {
	Registry string
	Name     string
	Version  string
}

// ParseRulesetArg parses registry/ruleset[@version] format
func ParseRulesetArg(arg string) (RulesetRef, error) {
	if arg == "" {
		return RulesetRef{}, fmt.Errorf("ruleset argument cannot be empty")
	}

	parts := strings.SplitN(arg, "/", 2)
	if len(parts) != 2 {
		return RulesetRef{}, fmt.Errorf("invalid ruleset format: %s (expected registry/ruleset[@version])", arg)
	}

	registry := parts[0]
	rulesetAndVersion := parts[1]

	versionParts := strings.SplitN(rulesetAndVersion, "@", 2)
	ruleset := versionParts[0]
	version := ""
	if len(versionParts) > 1 {
		version = versionParts[1]
	}

	if registry == "" {
		return RulesetRef{}, fmt.Errorf("registry name cannot be empty")
	}
	if ruleset == "" {
		return RulesetRef{}, fmt.Errorf("ruleset name cannot be empty")
	}

	return RulesetRef{
		Registry: registry,
		Name:     ruleset,
		Version:  version,
	}, nil
}

// ParseRulesetArgs parses multiple ruleset arguments
func ParseRulesetArgs(args []string) ([]RulesetRef, error) {
	refs := make([]RulesetRef, len(args))
	for i, arg := range args {
		ref, err := ParseRulesetArg(arg)
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
		return []string{"**/*.yml", "**/*.yaml"}
	}
	return include
}
