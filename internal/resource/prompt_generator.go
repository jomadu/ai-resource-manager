package resource

import "fmt"

// DefaultPromptGeneratorFactory creates prompt generators
type DefaultPromptGeneratorFactory struct{}

// NewPromptGenerator creates a prompt generator for the specified target
func (f *DefaultPromptGeneratorFactory) NewPromptGenerator(target CompileTarget) (PromptGenerator, error) {
	switch target {
	case TargetCursor, TargetMarkdown, TargetAmazonQ, TargetCopilot:
		return &DefaultPromptGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// NewPromptGeneratorFactory creates a new prompt generator factory
func NewPromptGeneratorFactory() PromptGeneratorFactory {
	return &DefaultPromptGeneratorFactory{}
}

// DefaultPromptGenerator generates simple prompt content for all targets
type DefaultPromptGenerator struct{}

// GeneratePrompt generates prompt content (same for all targets - just the body)
func (g *DefaultPromptGenerator) GeneratePrompt(namespace string, promptset *Promptset, promptID string, prompt *Prompt) string {
	return prompt.Body
}

// PromptGeneratorFactory interface for creating prompt generators
type PromptGeneratorFactory interface {
	NewPromptGenerator(target CompileTarget) (PromptGenerator, error)
}
