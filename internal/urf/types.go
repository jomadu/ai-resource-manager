package urf

import "github.com/jomadu/ai-rules-manager/internal/types"

// Ruleset represents a Universal Rule Format file
type Ruleset struct {
	Version  string   `yaml:"version"`
	Metadata Metadata `yaml:"metadata"`
	Rules    []Rule   `yaml:"rules"`
}

// Metadata contains ruleset metadata
type Metadata struct {
	ID          string `yaml:"id"`
	Name        string `yaml:"name"`
	Version     string `yaml:"version"`
	Description string `yaml:"description"`
}

// Rule represents a single rule within a URF file
type Rule struct {
	ID          string  `yaml:"id"`
	Name        string  `yaml:"name"`
	Description string  `yaml:"description"`
	Priority    int     `yaml:"priority"`
	Enforcement string  `yaml:"enforcement"`
	Scope       []Scope `yaml:"scope"`
	Body        string  `yaml:"body"`
}

// Scope defines where a rule applies
type Scope struct {
	Files []string `yaml:"files"`
}

// Parser interface for URF file parsing
type Parser interface {
	IsURF(file *types.File) bool
	Parse(file *types.File) (*Ruleset, error)
}

// Compiler interface for compiling URF to tool-specific formats
type Compiler interface {
	Compile(namespace string, ruleset *Ruleset) ([]*types.File, error)
}

// CompilerFactory interface for creating compilers
type CompilerFactory interface {
	GetCompiler(target CompileTarget) (Compiler, error)
	SupportedTargets() []CompileTarget
}

// CompileTarget represents different AI tool formats
type CompileTarget string

const (
	TargetCursor  CompileTarget = "cursor"
	TargetAmazonQ CompileTarget = "amazonq"
)

// DefaultCompilerFactory creates compilers
type DefaultCompilerFactory struct{}

// GetCompiler creates a compiler for the specified target
func (f *DefaultCompilerFactory) GetCompiler(target CompileTarget) (Compiler, error) {
	return NewCompiler(target)
}

// SupportedTargets returns supported compilation targets
func (f *DefaultCompilerFactory) SupportedTargets() []CompileTarget {
	return []CompileTarget{TargetCursor, TargetAmazonQ}
}

// NewDefaultCompilerFactory creates a new compiler factory
func NewDefaultCompilerFactory() CompilerFactory {
	return &DefaultCompilerFactory{}
}

// RuleGenerator interface for generating tool-specific rule files
type RuleGenerator interface {
	GenerateRule(namespace string, ruleset *Ruleset, rule *Rule) string
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
	GenerateRuleMetadata(namespace string, ruleset *Ruleset, rule *Rule) string
}
