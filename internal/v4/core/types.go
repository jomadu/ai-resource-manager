package core

type BuildInfo struct {
	Arch      string
	Version   Version
	Commit    string
	BuildTime string
}

type File struct {
	Path    string
	Content []byte
	Size    int64
}

type Version struct {
	Major      int
	Minor      int
	Patch      int
	Prerelease string
	Build      string
	Version    string
}

type PackageMetadata struct {
	RegistryName string
	Name         string
	Version      Version
}

type Package struct {
	Metadata PackageMetadata
	Files    []*File
}

type ConstraintType string

const (
	Exact      ConstraintType = "exact"
	Major      ConstraintType = "major"
	Minor      ConstraintType = "minor"
	BranchHead ConstraintType = "branch-head"
	Latest     ConstraintType = "latest"
)

type Constraint struct {
	Type    ConstraintType
	Version *Version
}
