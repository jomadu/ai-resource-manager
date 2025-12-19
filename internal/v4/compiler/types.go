package compiler

import (
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



// RuleGenerator interface for generating tool-specific rule files
type RuleGenerator interface {
	GenerateRule(namespace string, ruleset *resource.RulesetResource, ruleID string) (string, error)
}

// PromptGenerator interface for generating tool-specific prompt files
type PromptGenerator interface {
	GeneratePrompt(namespace string, promptset *resource.PromptsetResource, promptID string) (string, error)
}

// RuleFilenameGenerator interface for generating rule filenames
type RuleFilenameGenerator interface {
	GenerateRuleFilename(rulesetID, ruleID string) (string, error)
}

// PromptFilenameGenerator interface for generating prompt filenames
type PromptFilenameGenerator interface {
	GeneratePromptFilename(promptsetID, promptID string) (string, error)
}

// RuleGeneratorFactory interface for creating rule generators
type RuleGeneratorFactory interface {
	NewRuleGenerator(target CompileTarget) (RuleGenerator, error)
}

// PromptGeneratorFactory interface for creating prompt generators
type PromptGeneratorFactory interface {
	NewPromptGenerator(target CompileTarget) (PromptGenerator, error)
}

// RuleFilenameGeneratorFactory interface for creating rule filename generators
type RuleFilenameGeneratorFactory interface {
	NewRuleFilenameGenerator(target CompileTarget) (RuleFilenameGenerator, error)
}

// PromptFilenameGeneratorFactory interface for creating prompt filename generators
type PromptFilenameGeneratorFactory interface {
	NewPromptFilenameGenerator(target CompileTarget) (PromptFilenameGenerator, error)
}

