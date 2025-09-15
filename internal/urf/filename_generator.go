package urf

import "fmt"

// DefaultFilenameGeneratorFactory creates filename generators
type DefaultFilenameGeneratorFactory struct{}

// NewFilenameGenerator creates a filename generator for the specified target
func (f *DefaultFilenameGeneratorFactory) NewFilenameGenerator(target CompileTarget) (FilenameGenerator, error) {
	switch target {
	case TargetCursor:
		return &CursorFilenameGenerator{}, nil
	case TargetAmazonQ:
		return &AmazonQFilenameGenerator{}, nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// NewFilenameGeneratorFactory creates a new filename generator factory
func NewFilenameGeneratorFactory() FilenameGeneratorFactory {
	return &DefaultFilenameGeneratorFactory{}
}
