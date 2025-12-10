package new




type Manifest struct {
	Version string `json:"version"`
	Registries map[string]map[string]interface{} `json:"registries,omitempty"`
	Sinks      map[string]SinkConfig             `json:"sinks,omitempty"`
	Packages   map[string]map[string]interface{} `json:"packages"`
}

type SinkConfig struct {
	Directory string `json:"directory"`
	Layout string `json:"layout,omitempty"`
	CompileTarget string `json:"compileTarget"`
}

type PackageConfig struct {
	Version string `json:"version"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Sinks []string `json:"sinks"`
	ResourceType ResourceType `json:"resourceType,omitempty"`
}

type RulesetConfig struct {
	PackageConfig
	Priority int `json:"priority,omitempty"`
}

type PromptsetConfig struct {
	PackageConfig
}

type RegistryConfig struct {
	URL string `json:"url"`
	Type string `json:"type"`
}

type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches,omitempty"`
}

type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID string `json:"projectId,omitempty"`
	GroupID string `json:"groupId,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type CloudsmithRegistryConfig struct {
	RegistryConfig
	Owner string `json:"owner"`
	Repository string `json:"repository"`
}

type BuildInfo struct {
	Arch      string
	Version   Version
	Commit    string
	BuildTime string
}

type File struct {
	Path    string
	Content []byte
	Size    int64
}

type Version struct {
	Major int
	Minor int
	Patch int
	Prerelease string
	Build string
	Version string
}

type RegistryId struct {
	ID string
	Name string
}
type PackageId struct {
	ID string
	Name string
}

type PackageMetadata struct {
	PackageId PackageId
	RegistryId RegistryId
	Version Version
}

type Package struct {
	Metadata PackageMetadata
	Files []File
}

type OutdatedPackage struct {
	Current PackageMetadata
	Constraint string
	Wanted PackageMetadata
	Latest PackageMetadata
}

type ResourceType string

const (
	ResourceTypeRuleset ResourceType = "ruleset"
	ResourceTypePromptset ResourceType = "promptset"
)

type PackageInstallation struct {
	Metadata PackageMetadata
	ResourceType ResourceType
	FilePaths []string
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


