package compiler

import (
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

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

type ResourceCompilerFactory interface {
	NewResourceCompiler(target CompileTarget) (ResourceCompiler, error)
}

type ResourceCompiler interface {
	CompileRuleset(namespace string, ruleset *resource.RulesetResource) ([]*core.File, error)
	CompilePromptset(namespace string, promptset *resource.PromptsetResource) ([]*core.File, error)
}

// RuleGenerator interface for generating tool-specific rule files
type RuleGenerator interface {
	GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) string
}

// PromptGenerator interface for generating tool-specific prompt files
type PromptGenerator interface {
	GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) string
}

// RuleGeneratorFactory interface for creating rule generators
type RuleGeneratorFactory interface {
	NewRuleGenerator(target CompileTarget) (RuleGenerator, error)
}

// PromptGeneratorFactory interface for creating prompt generators
type PromptGeneratorFactory interface {
	NewPromptGenerator(target CompileTarget) (PromptGenerator, error)
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
	GenerateRuleMetadata(namespace string, ruleset *resource.RulesetResource, ruleID string, rule *resource.Rule) string
}
