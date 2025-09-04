package lockfile

// Entry represents a single ruleset entry in the lock file.
type Entry struct {
	Resolved string `json:"resolved"`
	Checksum string `json:"checksum"`
}

// LockFile represents the arm.lock file structure.
type LockFile struct {
	Rulesets map[string]map[string]Entry `json:"rulesets"`
}
