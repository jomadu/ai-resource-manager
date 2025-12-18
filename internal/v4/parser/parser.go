package parser

import (
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/resource"
)

type Parser interface {
	IsRuleset(file *core.File) bool
	IsPromptset(file *core.File) bool
	IsRulesetFile(path string) bool
	IsPromptsetFile(path string) bool
	ParseRuleset(file *core.File) (*resource.RulesetResource, error)
	ParsePromptset(file *core.File) (*resource.PromptsetResource, error)
	ParseRulesets(dirs []string, recursive bool, include, exclude []string) ([]*resource.RulesetResource, error)
	ParsePromptsets(dirs []string, recursive bool, include, exclude []string) ([]*resource.PromptsetResource, error)
}

// TODO: Implement
