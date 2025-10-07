package lockfile

// Entry represents a single resource entry in the lock file.
type Entry struct {
	Version  string `json:"version"`
	Display  string `json:"display"`
	Checksum string `json:"checksum"`
}

// LockFile represents the arm.lock file structure.
type LockFile struct {
	Rulesets   map[string]map[string]Entry `json:"rulesets"`
	Promptsets map[string]map[string]Entry `json:"promptsets"`
}
