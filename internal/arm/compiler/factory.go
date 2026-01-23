package compiler

import "fmt"

// DefaultRuleGeneratorFactory creates rule generators
type DefaultRuleGeneratorFactory struct{}

func NewRuleGeneratorFactory() RuleGeneratorFactory {
	return &DefaultRuleGeneratorFactory{}
}

func (f *DefaultRuleGeneratorFactory) NewRuleGenerator(tool Tool) (RuleGenerator, error) {
	switch tool {
	case Cursor:
		return &CursorRuleGenerator{}, nil
	case Markdown:
		return &MarkdownRuleGenerator{}, nil
	case AmazonQ:
		return &AmazonQRuleGenerator{}, nil
	case Copilot:
		return &CopilotRuleGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool)
	}
}

// DefaultPromptGeneratorFactory creates prompt generators
type DefaultPromptGeneratorFactory struct{}

func NewPromptGeneratorFactory() PromptGeneratorFactory {
	return &DefaultPromptGeneratorFactory{}
}

func (f *DefaultPromptGeneratorFactory) NewPromptGenerator(tool Tool) (PromptGenerator, error) {
	switch tool {
	case Cursor:
		return &CursorPromptGenerator{}, nil
	case Markdown:
		return &MarkdownPromptGenerator{}, nil
	case AmazonQ:
		return &AmazonQPromptGenerator{}, nil
	case Copilot:
		return &CopilotPromptGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported tool: %s", tool)
	}
}

// DefaultRuleFilenameGeneratorFactory creates rule filename generators
type DefaultRuleFilenameGeneratorFactory struct{}

func NewRuleFilenameGeneratorFactory() RuleFilenameGeneratorFactory {
	return &DefaultRuleFilenameGeneratorFactory{}
}

func (f *DefaultRuleFilenameGeneratorFactory) NewRuleFilenameGenerator(tool Tool) (RuleFilenameGenerator, error) {
	switch tool {
	case Cursor:
		return &CursorRuleFilenameGenerator{}, nil
	case Markdown:
		return &MarkdownRuleFilenameGenerator{}, nil
	case AmazonQ:
		return &AmazonQRuleFilenameGenerator{}, nil
	case Copilot:
		return &CopilotRuleFilenameGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported rule filename tool: %s", tool)
	}
}

// DefaultPromptFilenameGeneratorFactory creates prompt filename generators
type DefaultPromptFilenameGeneratorFactory struct{}

func NewPromptFilenameGeneratorFactory() PromptFilenameGeneratorFactory {
	return &DefaultPromptFilenameGeneratorFactory{}
}

func (f *DefaultPromptFilenameGeneratorFactory) NewPromptFilenameGenerator(tool Tool) (PromptFilenameGenerator, error) {
	switch tool {
	case Cursor:
		return &CursorPromptFilenameGenerator{}, nil
	case Markdown:
		return &MarkdownPromptFilenameGenerator{}, nil
	case AmazonQ:
		return &AmazonQPromptFilenameGenerator{}, nil
	case Copilot:
		return &CopilotPromptFilenameGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported prompt filename tool: %s", tool)
	}
}
