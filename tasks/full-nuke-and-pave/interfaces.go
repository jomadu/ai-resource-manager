package main

import "context"

// Configuration structs
type RegistryConfig struct {
	URL  string `json:"url"`
	Type string `json:"type"`
}

type SinkConfig struct {
	Directories []string `json:"directories"`
	Rulesets    []string `json:"rulesets"`
}

type RCConfig struct {
	Registries map[string]RegistryConfig `json:"registries"`
	Sinks      map[string]SinkConfig     `json:"sinks"`
}

type ManifestEntry struct {
	Version string   `json:"version"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

type Manifest struct {
	Rulesets map[string]map[string]ManifestEntry `json:"rulesets"`
}

type LockFileEntry struct {
	URL        string   `json:"url"`
	Type       string   `json:"type"`
	Constraint string   `json:"constraint"`
	Resolved   string   `json:"resolved"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
}

type LockFile struct {
	Rulesets map[string]map[string]LockFileEntry `json:"rulesets"`
}

// File and version structs
type File struct {
	Path    string `json:"path"`
	Content []byte `json:"content"`
	Size    int64  `json:"size"`
}

type VersionRefType int

const (
	Tag VersionRefType = iota
	Branch
	Commit
)

type VersionRef struct {
	ID   string         `json:"id"`
	Type VersionRefType `json:"type"`
}

type ContentSelector struct {
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

// Constraint handling
type ConstraintType int

const (
	Pin ConstraintType = iota
	Caret
	Tilde
	BranchHead
)

type Constraint struct {
	Type    ConstraintType `json:"type"`
	Version string         `json:"version"`
	Major   int            `json:"major"`
	Minor   int            `json:"minor"`
	Patch   int            `json:"patch"`
}

// Output data structures
type OutdatedRuleset struct {
	Registry string `json:"registry"`
	Ruleset  string `json:"ruleset"`
	Current  string `json:"current"`
	Wanted   string `json:"wanted"`
	Latest   string `json:"latest"`
}

type InstalledRuleset struct {
	Registry string   `json:"registry"`
	Ruleset  string   `json:"ruleset"`
	Version  string   `json:"version"`
	Include  []string `json:"include"`
	Exclude  []string `json:"exclude"`
	Sinks    []string `json:"sinks"`
}

type RulesetInfo struct {
	Registry   string   `json:"registry"`
	URL        string   `json:"url"`
	Type       string   `json:"type"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Installed  []string `json:"installed"`
	Sinks      []string `json:"sinks"`
	Constraint string   `json:"constraint"`
	Resolved   string   `json:"resolved"`
	Wanted     string   `json:"wanted"`
	Latest     string   `json:"latest"`
}

type Installation struct {
	Ruleset string `json:"ruleset"`
	Version string `json:"version"`
	Path    string `json:"path"`
}

type VersionInfo struct {
	Arch      string `json:"arch"`
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	Timestamp string `json:"timestamp"`
}

// Service provides the main ARM functionality for managing AI rule rulesets.
// It handles installation, updates, and querying of rulesets from configured registries.
type Service interface {
	Install(ctx context.Context, ruleset, version string, include, exclude []string) error
	Uninstall(ctx context.Context, ruleset string) error
	Update(ctx context.Context, ruleset string) error
	Outdated(ctx context.Context) ([]OutdatedRuleset, error)
	List(ctx context.Context) ([]InstalledRuleset, error)
	Info(ctx context.Context, ruleset string) (*RulesetInfo, error)
	Version() VersionInfo
}

// ConfigManager handles .armrc.json configuration file operations.
// Manages registry definitions and sink configurations for rule deployment.
type ConfigManager interface {
	GetRegistries(ctx context.Context) (map[string]RegistryConfig, error)
	GetSinks(ctx context.Context) (map[string]SinkConfig, error)
	AddRegistry(ctx context.Context, name, url, registryType string) error
	AddSink(ctx context.Context, name string, dirs []string, include, exclude []string) error
	RemoveRegistry(ctx context.Context, name string) error
	RemoveSink(ctx context.Context, name string) error
}

// ManifestManager handles arm.json manifest file operations.
// Tracks user-specified version constraints and include/exclude patterns.
type ManifestManager interface {
	GetEntry(ctx context.Context, registry, ruleset string) (*ManifestEntry, error)
	GetEntries(ctx context.Context) (map[string]map[string]ManifestEntry, error)
	CreateEntry(ctx context.Context, registry, ruleset string, entry ManifestEntry) error
	UpdateEntry(ctx context.Context, registry, ruleset string, entry ManifestEntry) error
	RemoveEntry(ctx context.Context, registry, ruleset string) error
}

// LockFileManager handles arm.lock file operations.
// Maintains resolved versions and registry metadata for reproducible installs.
type LockFileManager interface {
	GetEntry(ctx context.Context, registry, ruleset string) (*LockFileEntry, error)
	GetEntries(ctx context.Context) (map[string]map[string]LockFileEntry, error)
	CreateEntry(ctx context.Context, registry, ruleset string, entry LockFileEntry) error
	UpdateEntry(ctx context.Context, registry, ruleset string, entry LockFileEntry) error
	RemoveEntry(ctx context.Context, registry, ruleset string) error
}

// Registry provides version-controlled access to ruleset repositories.
// Abstracts Git operations for fetching tags, branches, and file content.
type Registry interface {
	ListVersions(ctx context.Context) ([]VersionRef, error)
	GetContent(ctx context.Context, version VersionRef, selector ContentSelector) ([]File, error)
}

// Cache provides local storage for registry data and ruleset files.
// Implements ~/.arm/cache structure with index.json metadata tracking.
type Cache interface {
	ListVersions(ctx context.Context, registryKey, rulesetKey string) ([]string, error)
	Get(ctx context.Context, registryKey, rulesetKey, version string) ([]File, error)
	Set(ctx context.Context, registryKey, rulesetKey, version string, files []File) error
	InvalidateRegistry(ctx context.Context, registryKey string) error
	InvalidateRuleset(ctx context.Context, registryKey, rulesetKey string) error
	InvalidateVersion(ctx context.Context, registryKey, rulesetKey, version string) error
}

// Repository handles low-level Git operations for registry access.
// Maintains Git repositories within ~/.arm/cache/registries/ and manages cloning, fetching, and file extraction.
type Repository interface {
	Clone(ctx context.Context, url string) error
	Fetch(ctx context.Context) error
	Pull(ctx context.Context) error
	GetTags(ctx context.Context) ([]string, error)
	GetBranches(ctx context.Context) ([]string, error)
	Checkout(ctx context.Context, ref string) error
	GetFiles(ctx context.Context, selector ContentSelector) ([]File, error)
}

// KeyGenerator creates consistent hash keys for cache storage.
// Generates SHA256-based keys for registries and rulesets.
type KeyGenerator interface {
	RegistryKey(url, registryType string) string
	RulesetKey(selector ContentSelector) string
}

// Installer manages physical file deployment to sink directories.
// Handles the arm/ subdirectory structure and version-specific installations.
type Installer interface {
	Install(ctx context.Context, dir, ruleset, version string, files []File) error
	Uninstall(ctx context.Context, dir, ruleset string) error
	ListInstalled(ctx context.Context, dir string) ([]Installation, error)
}

// ConstraintResolver handles semantic versioning and constraint resolution.
// Supports semver patterns (^, ~, =) and branch tracking for version selection.
type ConstraintResolver interface {
	ParseConstraint(constraint string) (Constraint, error)
	SatisfiesConstraint(version string, constraint Constraint) bool
	FindBestMatch(constraint Constraint, versions []VersionRef) (*VersionRef, error)
}

// ContentResolver applies include/exclude patterns to filter ruleset files.
// Implements glob pattern matching for selective rule installation.
type ContentResolver interface {
	ResolveContent(selector ContentSelector, files []File) ([]File, error)
}