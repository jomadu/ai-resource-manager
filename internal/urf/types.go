package urf

import "github.com/jomadu/ai-rules-manager/internal/types"

// URFFile represents a Universal Rule Format file
type URFFile struct {
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
	Parse(file *types.File) (*URFFile, error)
}

// Compiler interface for compiling URF to tool-specific formats
type Compiler interface {
	Compile(urf *URFFile) ([]*types.File, error)
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
