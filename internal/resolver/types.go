package resolver

import "github.com/jomadu/ai-rules-manager/internal/types"

// ConstraintType defines the type of version constraint.
type ConstraintType int

const (
	Exact ConstraintType = iota
	Major
	Minor
	BranchHead
	Latest
)

// Constraint represents a version constraint.
type Constraint struct {
	Type    ConstraintType
	Version string
	Major   int
	Minor   int
	Patch   int
}

// ResolvedVersion combines constraint and resolved version information.
type ResolvedVersion struct {
	Constraint Constraint    // Original constraint struct
	Version    types.Version // Resolved version with display
}
