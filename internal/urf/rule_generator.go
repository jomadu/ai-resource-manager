package urf

import "fmt"

// DefaultRuleGeneratorFactory creates rule generators
type DefaultRuleGeneratorFactory struct{}

// NewRuleGenerator creates a rule generator for the specified target
func (f *DefaultRuleGeneratorFactory) NewRuleGenerator(target CompileTarget) (RuleGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorRuleGenerator{
			metadataGen: NewRuleMetadataGenerator(),
		}, nil
	case TargetAmazonQ:
		return &AmazonQRuleGenerator{
			metadataGen: NewRuleMetadataGenerator(),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// NewRuleGeneratorFactory creates a new rule generator factory
func NewRuleGeneratorFactory() RuleGeneratorFactory {
	return &DefaultRuleGeneratorFactory{}
}
