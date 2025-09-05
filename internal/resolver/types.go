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
	Type    ConstraintType `json:"type"`
	Version string         `json:"version"`
	Major   int            `json:"major"`
	Minor   int            `json:"minor"`
	Patch   int            `json:"patch"`
}

// ResolvedVersion combines constraint and resolved version information.
type ResolvedVersion struct {
	Constraint Constraint    `json:"constraint"` // Original constraint struct
	Version    types.Version `json:"version"`    // Resolved version with display
}
