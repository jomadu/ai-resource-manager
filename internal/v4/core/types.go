type File struct {
	Path    string
	Content []byte
	Size    int64
}

type Version struct {
	Major int
	Minor int
	Patch int
	Prerelease string
	Build string
	Version string
}

type RegistryId struct {
	ID string
	Name string
}
type PackageId struct {
	ID string
	Name string
}

type PackageMetadata struct {
	PackageId PackageId
	RegistryId RegistryId
	Version Version
}

type Package struct {
	Metadata PackageMetadata
	Files []File
}