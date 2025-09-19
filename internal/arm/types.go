package arm

import "github.com/jomadu/ai-rules-manager/internal/urf"

// OutdatedRuleset represents a ruleset that has newer versions available.
type OutdatedRuleset struct {
	RulesetInfo *RulesetInfo `json:"ruleset_info"`
	Wanted      string       `json:"wanted"`
	Latest      string       `json:"latest"`
}

// ManifestInfo contains information from the manifest file.
type ManifestInfo struct {
	Constraint string   `json:"constraint"`
	Priority   int      `json:"priority"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// InstallationInfo contains information about the actual installation.
type InstallationInfo struct {
	Version        string   `json:"version"`
	InstalledPaths []string `json:"installed_paths"`
}

// RulesetInfo provides detailed information about a ruleset.
type RulesetInfo struct {
	Registry     string           `json:"registry"`
	Name         string           `json:"name"`
	Manifest     ManifestInfo     `json:"manifest"`
	Installation InstallationInfo `json:"installation"`
}

// InstallRequest groups install parameters to avoid repetitive parameter passing.
type InstallRequest struct {
	Registry string
	Ruleset  string
	Version  string
	Priority int
	Include  []string
	Exclude  []string
	Sinks    []string
}

// CompileRequest groups compile parameters to avoid repetitive parameter passing.
type CompileRequest struct {
	Files        []string            // Input file paths/patterns
	Targets      []urf.CompileTarget // Compilation targets (cursor, amazonq, etc.)
	OutputDir    string              // Output directory
	Namespace    string              // Optional namespace override
	Force        bool                // Overwrite existing files
	Recursive    bool                // Recursive directory processing
	DryRun       bool                // Show what would be compiled
	Verbose      bool                // Detailed output
	ValidateOnly bool                // Syntax validation only
	FailFast     bool                // Stop on first error
	Include      []string            // Include patterns for file filtering
	Exclude      []string            // Exclude patterns for file filtering
}

// CompileResult contains compilation results and statistics.
type CompileResult struct {
	CompiledFiles []CompiledFile // Successfully compiled files
	Skipped       []SkippedFile  // Files that were skipped
	Errors        []CompileError // Compilation errors
	Stats         CompileStats   // Summary statistics
}

// CompiledFile represents a successfully compiled file.
type CompiledFile struct {
	SourcePath string            `json:"source_path"` // Original URF file path
	TargetPath string            `json:"target_path"` // Generated file path
	Target     urf.CompileTarget `json:"target"`      // Compilation target used
	RuleCount  int               `json:"rule_count"`  // Number of rules in the file
}

// SkippedFile represents a file that was skipped during compilation.
type SkippedFile struct {
	Path   string `json:"path"`   // File path
	Reason string `json:"reason"` // Why it was skipped
}

// CompileError represents a compilation error.
type CompileError struct {
	FilePath string `json:"file_path"` // File that caused the error
	Target   string `json:"target"`    // Target format (if applicable)
	Error    string `json:"error"`     // Error message
}

// CompileStats provides compilation statistics.
type CompileStats struct {
	FilesProcessed int            `json:"files_processed"` // Total files examined
	FilesCompiled  int            `json:"files_compiled"`  // Successfully compiled files
	FilesSkipped   int            `json:"files_skipped"`   // Skipped files
	RulesGenerated int            `json:"rules_generated"` // Total rules generated
	Errors         int            `json:"errors"`          // Total errors
	TargetStats    map[string]int `json:"target_stats"`    // Per-target compilation counts
}
