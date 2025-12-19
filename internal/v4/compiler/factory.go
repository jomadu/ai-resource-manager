package compiler

import "fmt"

// DefaultRuleGeneratorFactory creates rule generators
type DefaultRuleGeneratorFactory struct{}

func NewRuleGeneratorFactory() RuleGeneratorFactory {
	return &DefaultRuleGeneratorFactory{}
}

func (f *DefaultRuleGeneratorFactory) NewRuleGenerator(target CompileTarget) (RuleGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorRuleGenerator{}, nil
	case TargetMarkdown:
		return &MarkdownRuleGenerator{}, nil
	case TargetAmazonQ:
		return &AmazonQRuleGenerator{}, nil
	case TargetCopilot:
		return &CopilotRuleGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// DefaultPromptGeneratorFactory creates prompt generators
type DefaultPromptGeneratorFactory struct{}

func NewPromptGeneratorFactory() PromptGeneratorFactory {
	return &DefaultPromptGeneratorFactory{}
}

func (f *DefaultPromptGeneratorFactory) NewPromptGenerator(target CompileTarget) (PromptGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorPromptGenerator{}, nil
	case TargetMarkdown:
		return &MarkdownPromptGenerator{}, nil
	case TargetAmazonQ:
		return &AmazonQPromptGenerator{}, nil
	case TargetCopilot:
		return &CopilotPromptGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// DefaultRuleFilenameGeneratorFactory creates rule filename generators
type DefaultRuleFilenameGeneratorFactory struct{}

func NewRuleFilenameGeneratorFactory() RuleFilenameGeneratorFactory {
	return &DefaultRuleFilenameGeneratorFactory{}
}

func (f *DefaultRuleFilenameGeneratorFactory) NewRuleFilenameGenerator(target CompileTarget) (RuleFilenameGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorRuleFilenameGenerator{}, nil
	case TargetMarkdown:
		return &MarkdownRuleFilenameGenerator{}, nil
	case TargetAmazonQ:
		return &AmazonQRuleFilenameGenerator{}, nil
	case TargetCopilot:
		return &CopilotRuleFilenameGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported rule filename target: %s", target)
	}
}

// DefaultPromptFilenameGeneratorFactory creates prompt filename generators
type DefaultPromptFilenameGeneratorFactory struct{}

func NewPromptFilenameGeneratorFactory() PromptFilenameGeneratorFactory {
	return &DefaultPromptFilenameGeneratorFactory{}
}

func (f *DefaultPromptFilenameGeneratorFactory) NewPromptFilenameGenerator(target CompileTarget) (PromptFilenameGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorPromptFilenameGenerator{}, nil
	case TargetMarkdown:
		return &MarkdownPromptFilenameGenerator{}, nil
	case TargetAmazonQ:
		return &AmazonQPromptFilenameGenerator{}, nil
	case TargetCopilot:
		return &CopilotPromptFilenameGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported prompt filename target: %s", target)
	}
}