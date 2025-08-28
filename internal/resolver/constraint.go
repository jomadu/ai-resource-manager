package resolver

import (
	"errors"

	"github.com/jomadu/ai-rules-manager/internal/arm"
)

// ConstraintResolver handles semantic versioning and constraint resolution.
type ConstraintResolver interface {
	ParseConstraint(constraint string) (Constraint, error)
	SatisfiesConstraint(version string, constraint Constraint) bool
	FindBestMatch(constraint Constraint, versions []arm.VersionRef) (*arm.VersionRef, error)
}

// GitConstraintResolver implements semantic versioning constraint resolution.
type GitConstraintResolver struct{}

// NewGitConstraintResolver creates a new Git-based constraint resolver.
func NewGitConstraintResolver() *GitConstraintResolver {
	return &GitConstraintResolver{}
}

func (g *GitConstraintResolver) ParseConstraint(constraint string) (Constraint, error) {
	return Constraint{}, errors.New("not implemented")
}

func (g *GitConstraintResolver) SatisfiesConstraint(version string, constraint Constraint) bool {
	return false // TODO: implement
}

func (g *GitConstraintResolver) FindBestMatch(constraint Constraint, versions []arm.VersionRef) (*arm.VersionRef, error) {
	return nil, errors.New("not implemented")
}
