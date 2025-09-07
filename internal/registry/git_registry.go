package registry

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// GitRegistry implements Git-based registry access with caching.
// It abstracts all git repository operations - users simply create the registry
// and call methods like GetTags() without worrying about cloning, fetching, etc.
type GitRegistry struct {
	cache    cache.RegistryRulesetCache
	repo     cache.GitRepoCache
	config   GitRegistryConfig
	resolver resolver.ConstraintResolver
}

// NewGitRegistry creates a new Git-based registry that handles all git operations internally.
// The registry will automatically clone the repository on first use and fetch updates as needed.
func NewGitRegistry(config GitRegistryConfig, rulesetCache cache.RegistryRulesetCache, repoCache cache.GitRepoCache) *GitRegistry {
	return &GitRegistry{
		cache:    rulesetCache,
		repo:     repoCache,
		config:   config,
		resolver: resolver.NewGitConstraintResolver(),
	}
}

// isBranchConstraint checks if a constraint refers to a branch name.
// Uses permissive detection: anything that's not a semantic version or "latest" is treated as a potential branch.
func (g *GitRegistry) isBranchConstraint(constraint string) bool {
	// Exclude "latest" special keyword
	if constraint == "latest" {
		return false
	}

	// Exclude semantic version patterns (contains dots and follows semver format)
	if strings.Contains(constraint, ".") {
		// Check if it's a valid semantic version
		if g.isSemverTag(constraint) || g.isSemverTag("v"+constraint) {
			return false
		}
		// Also exclude patterns with dots that aren't valid semver (these are invalid)
		return false
	}

	// Exclude version constraint prefixes
	if strings.HasPrefix(constraint, "^") || strings.HasPrefix(constraint, "~") {
		return false
	}

	// Everything else is treated as a potential branch name
	return true
}

// sortTagsBySemver sorts tags by semantic version in descending order.
func (g *GitRegistry) sortTagsBySemver(tags []string) []string {
	var semverTags []string
	var otherTags []string

	for _, tag := range tags {
		if g.isSemverTag(tag) {
			semverTags = append(semverTags, tag)
		} else {
			otherTags = append(otherTags, tag)
		}
	}

	// Sort semver tags by version (descending)
	sort.Slice(semverTags, func(i, j int) bool {
		return g.isHigherVersion(semverTags[i], semverTags[j])
	})

	// Combine semver tags first, then other tags
	result := make([]string, 0, len(tags))
	result = append(result, semverTags...)
	result = append(result, otherTags...)
	return result
}

// isSemverTag checks if a tag follows semantic versioning.
func (g *GitRegistry) isSemverTag(tag string) bool {
	normalized := strings.TrimPrefix(tag, "v")
	matched, _ := regexp.MatchString(`^\d+\.\d+\.\d+`, normalized)
	return matched
}

// isHigherVersion compares two semantic versions.
func (g *GitRegistry) isHigherVersion(v1, v2 string) bool {
	major1, minor1, patch1, err1 := g.parseVersion(v1)
	major2, minor2, patch2, err2 := g.parseVersion(v2)
	if err1 != nil || err2 != nil {
		return false
	}

	if major1 != major2 {
		return major1 > major2
	}
	if minor1 != minor2 {
		return minor1 > minor2
	}
	return patch1 > patch2
}

// parseVersion parses a semantic version string.
func (g *GitRegistry) parseVersion(version string) (major, minor, patch int, err error) {
	version = strings.TrimPrefix(version, "v")
	re := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)`)
	matches := re.FindStringSubmatch(version)
	if len(matches) < 4 {
		return 0, 0, 0, fmt.Errorf("invalid version format")
	}

	major, _ = strconv.Atoi(matches[1])
	minor, _ = strconv.Atoi(matches[2])
	patch, _ = strconv.Atoi(matches[3])
	return
}

func (g *GitRegistry) ListVersions(ctx context.Context) ([]types.Version, error) {
	tags, err := g.repo.GetTags(ctx)
	if err != nil {
		return nil, err
	}

	// Sort tags by semantic version (descending)
	sortedTags := g.sortTagsBySemver(tags)

	var versions []types.Version
	// Add tags first (priority ordering)
	for _, tag := range sortedTags {
		versions = append(versions, types.Version{Version: tag, Display: tag})
	}

	// Add branch commits in configuration order
	for _, branch := range g.config.Branches {
		hash, err := g.repo.GetBranchHeadCommitHash(ctx, branch)
		if err != nil {
			continue // Skip branches that don't exist
		}
		versions = append(versions, types.Version{Version: hash, Display: hash[:7]})
	}

	return versions, nil
}

func (g *GitRegistry) GetContent(ctx context.Context, version types.Version, selector types.ContentSelector) ([]types.File, error) {
	// Try cache first
	files, err := g.cache.GetRulesetVersion(ctx, selector, version.Version)
	if err == nil {
		return files, nil
	}

	// Get files from git repo
	files, err = g.repo.GetFilesFromCommit(ctx, version.Version, selector)
	if err != nil {
		return nil, err
	}

	// Cache the result
	_ = g.cache.SetRulesetVersion(ctx, selector, version.Version, files)

	return files, nil
}

// GetTags returns all available tags from the registry.
func (g *GitRegistry) GetTags(ctx context.Context) ([]string, error) {
	return g.repo.GetTags(ctx)
}

// GetBranches returns all available branches from the registry.
func (g *GitRegistry) GetBranches(ctx context.Context) ([]string, error) {
	return g.repo.GetBranches(ctx)
}

func (g *GitRegistry) ResolveVersion(ctx context.Context, constraint string) (*resolver.ResolvedVersion, error) {
	// Parse constraint first
	parsedConstraint, err := g.resolver.ParseConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint %s: %w", constraint, err)
	}

	// Handle branch constraints by resolving to commit hash
	if g.isBranchConstraint(constraint) {
		hash, err := g.repo.GetBranchHeadCommitHash(ctx, constraint)
		if err != nil {
			// Provide helpful error with available branches
			branches, branchErr := g.repo.GetBranches(ctx)
			if branchErr == nil && len(branches) > 0 {
				return nil, fmt.Errorf("branch %s not found. Available branches: %v", constraint, branches)
			}
			return nil, fmt.Errorf("failed to resolve branch %s: %w", constraint, err)
		}
		return &resolver.ResolvedVersion{
			Constraint: parsedConstraint,
			Version: types.Version{
				Version: hash,
				Display: hash[:7], // Use first 7 chars for display
			},
		}, nil
	}

	// Handle semantic version constraints
	versions, err := g.ListVersions(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	resolvedVersion, err := g.resolver.FindBestMatch(parsedConstraint, versions)
	if err != nil {
		return nil, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return &resolver.ResolvedVersion{
		Constraint: parsedConstraint,
		Version:    *resolvedVersion,
	}, nil
}
