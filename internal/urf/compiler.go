package urf

import "fmt"

// DefaultCompilerFactory implements CompilerFactory
type DefaultCompilerFactory struct{}

// NewCompilerFactory creates a new compiler factory
func NewCompilerFactory() CompilerFactory {
	return &DefaultCompilerFactory{}
}

// GetCompiler returns a compiler for the specified target
func (f *DefaultCompilerFactory) GetCompiler(target CompileTarget) (Compiler, error) {
	switch target {
	case TargetCursor:
		return NewCursorCompiler(), nil
	case TargetAmazonQ:
		return NewAmazonQCompiler(), nil
	default:
		return nil, fmt.Errorf("unsupported compile target: %s", target)
	}
}

// SupportedTargets returns list of supported compilation targets
func (f *DefaultCompilerFactory) SupportedTargets() []CompileTarget {
	return []CompileTarget{TargetCursor, TargetAmazonQ}
}
