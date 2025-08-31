package resolver

// ConstraintType defines the type of version constraint.
type ConstraintType int

const (
	Pin ConstraintType = iota
	Caret
	Tilde
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
