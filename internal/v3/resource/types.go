package resource

import "github.com/jomadu/ai-rules-manager/internal/types"

// Ruleset represents a Universal Rule Format file
type Ruleset struct {
	APIVersion string      `yaml:"apiVersion" validate:"required"`
	Kind       string      `yaml:"kind" validate:"required,eq=Ruleset"`
	Metadata   Metadata    `yaml:"metadata" validate:"required"`
	Spec       RulesetSpec `yaml:"spec" validate:"required"`
}

// RulesetSpec contains the ruleset specification
type RulesetSpec struct {
	Rules map[string]Rule `yaml:"rules" validate:"required,min=1"`
}

// Promptset represents a Universal Prompt Format file
type Promptset struct {
	APIVersion string        `yaml:"apiVersion" validate:"required"`
	Kind       string        `yaml:"kind" validate:"required,eq=Promptset"`
	Metadata   Metadata      `yaml:"metadata" validate:"required"`
	Spec       PromptsetSpec `yaml:"spec" validate:"required"`
}

// PromptsetSpec contains the promptset specification
type PromptsetSpec struct {
	Prompts map[string]Prompt `yaml:"prompts" validate:"required,min=1"`
}

// Metadata contains ruleset metadata
type Metadata struct {
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

// Parser interface for resource file parsing
type Parser interface {
	IsRuleset(file *types.File) bool
	IsPromptset(file *types.File) bool
	IsRulesetFile(path string) bool
	IsPromptsetFile(path string) bool
	ParseRuleset(file *types.File) (*Ruleset, error)
	ParsePromptset(file *types.File) (*Promptset, error)
	ParseRulesets(dirs []string, recursive bool, include, exclude []string) ([]*Ruleset, error)
	ParsePromptsets(dirs []string, recursive bool, include, exclude []string) ([]*Promptset, error)
}

// Compiler interface for compiling resource files to tool-specific formats
type Compiler interface {
	CompileRuleset(namespace string, ruleset *Ruleset) ([]*types.File, error)
	CompilePromptset(namespace string, promptset *Promptset) ([]*types.File, error)
}

// CompileTarget represents different AI tool formats
type CompileTarget string

const (
	TargetCursor   CompileTarget = "cursor"
	TargetMarkdown CompileTarget = "markdown"
	TargetAmazonQ  CompileTarget = "amazonq"
	TargetCopilot  CompileTarget = "copilot"
)

// RuleGenerator interface for generating tool-specific rule files
type RuleGenerator interface {
	GenerateRule(namespace string, ruleset *Ruleset, ruleID string, rule *Rule) string
}

// PromptGenerator interface for generating tool-specific prompt files
type PromptGenerator interface {
	GeneratePrompt(namespace string, promptset *Promptset, promptID string, prompt *Prompt) string
}

// RuleGeneratorFactory interface for creating rule generators
type RuleGeneratorFactory interface {
	NewRuleGenerator(target CompileTarget) (RuleGenerator, error)
}

// FilenameGenerator interface for generating filenames
type FilenameGenerator interface {
	GenerateFilename(rulesetID, ruleID string) string
}

// FilenameGeneratorFactory interface for creating filename generators
type FilenameGeneratorFactory interface {
	NewFilenameGenerator(target CompileTarget) (FilenameGenerator, error)
}

// RuleMetadataGenerator interface for generating metadata blocks
type RuleMetadataGenerator interface {
	GenerateRuleMetadata(namespace string, ruleset *Ruleset, ruleID string, rule *Rule) string
}
