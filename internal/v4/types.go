package v4






type BuildInfo struct {
	Arch      string
	Version   Version
	Commit    string
	BuildTime string
}



type OutdatedPackage struct {
	Current PackageMetadata
	Constraint string
	Wanted PackageMetadata
	Latest PackageMetadata
}



type PackageInfo struct {
	Installation PackageInstallation
	LockInfo PackageLockInfo
	Config map[string]interface{}
}

type ConstraintType string

const (
	Exact ConstraintType = "exact"
	Major ConstraintType = "major"
	Minor ConstraintType = "minor"
	BranchHead ConstraintType = "branch-head"
	Latest ConstraintType = "latest"
)

type Constraint struct {
	Type ConstraintType
	Version *Version
}

// Ruleset represents a Universal Rule Format file
type RulesetResource struct {
	APIVersion string      `yaml:"apiVersion" validate:"required"`
	Kind       string      `yaml:"kind" validate:"required,eq=Ruleset"`
	Metadata   ResourceMetadata    `yaml:"metadata" validate:"required"`
	Spec       RulesetSpec `yaml:"spec" validate:"required"`
}

// RulesetSpec contains the ruleset specification
type RulesetSpec struct {
	Rules map[string]Rule `yaml:"rules" validate:"required,min=1"`
}

// Promptset represents a Universal Prompt Format file
type PromptsetResource struct {
	APIVersion string        `yaml:"apiVersion" validate:"required"`
	Kind       string        `yaml:"kind" validate:"required,eq=Promptset"`
	Metadata   ResourceMetadata      `yaml:"metadata" validate:"required"`
	Spec       PromptsetSpec `yaml:"spec" validate:"required"`
}

// PromptsetSpec contains the promptset specification
type PromptsetSpec struct {
	Prompts map[string]Prompt `yaml:"prompts" validate:"required,min=1"`
}

// Metadata contains ruleset metadata
type ResourceMetadata struct {
	ID          string `yaml:"id" validate:"required"`
	Name        string `yaml:"name" validate:"required"`
	Description string `yaml:"description,omitempty"`
}

// Rule represents a single rule within a resource file
type Rule struct {
	Name        string  `yaml:"name" validate:"required"`
	Description string  `yaml:"description,omitempty"`
	Priority    int     `yaml:"priority,omitempty"`
	Enforcement string  `yaml:"enforcement,omitempty" validate:"omitempty,oneof=may should must"`
	Scope       []Scope `yaml:"scope,omitempty"`
	Body        string  `yaml:"body" validate:"required"`
}

// Prompt represents a single prompt within a promptset
type Prompt struct {
	Name        string `yaml:"name" validate:"required"`
	Description string `yaml:"description,omitempty"`
	Body        string `yaml:"body" validate:"required"`
}

// Scope defines where a rule applies
type Scope struct {
	Files []string `yaml:"files"`
}

// CompileTarget represents different AI tool formats
type CompileTarget string

const (
	TargetCursor   CompileTarget = "cursor"
	TargetMarkdown CompileTarget = "markdown"
	TargetAmazonQ  CompileTarget = "amazonq"
	TargetCopilot  CompileTarget = "copilot"
)

// CompileRequest groups compile parameters following ARM patterns
type CompileRequest struct {
	Paths        []string
	Targets      []string
	OutputDir    string
	Namespace    string
	Force        bool
	Recursive    bool
	Verbose      bool
	ValidateOnly bool
	Include      []string
	Exclude      []string
	FailFast     bool
}

// PackageLockInfo represents lock information for a package
type PackageLockInfo struct {
	Version  string `json:"version"`
	Checksum string `json:"checksum"`
}


