package installer

import "testing"

func TestFlatInstallerHashFile(t *testing.T) {
	installer := NewFlatInstaller()

	tests := []struct {
		name     string
		registry string
		ruleset  string
		version  string
		filePath string
	}{
		{
			name:     "basic hash",
			registry: "test-registry",
			ruleset:  "test-ruleset",
			version:  "1.0.0",
			filePath: "rules/test.md",
		},
		{
			name:     "different registry",
			registry: "other-registry",
			ruleset:  "test-ruleset",
			version:  "1.0.0",
			filePath: "rules/test.md",
		},
		{
			name:     "different version",
			registry: "test-registry",
			ruleset:  "test-ruleset",
			version:  "2.0.0",
			filePath: "rules/test.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := installer.hashFile(tt.registry, tt.ruleset, tt.version, tt.filePath)
			if len(got) != 8 {
				t.Errorf("hashFile() = %v, want 8 character hash", got)
			}
		})
	}
}

func TestFlatInstallerHashFileDeterministic(t *testing.T) {
	installer := NewFlatInstaller()

	hash1 := installer.hashFile("registry", "ruleset", "1.0.0", "file.md")
	hash2 := installer.hashFile("registry", "ruleset", "1.0.0", "file.md")

	if hash1 != hash2 {
		t.Errorf("hashFile() not deterministic: %v != %v", hash1, hash2)
	}
}

func TestFlatInstallerHashFileUnique(t *testing.T) {
	installer := NewFlatInstaller()

	hash1 := installer.hashFile("registry1", "ruleset", "1.0.0", "file.md")
	hash2 := installer.hashFile("registry2", "ruleset", "1.0.0", "file.md")

	if hash1 == hash2 {
		t.Errorf("hashFile() not unique for different registries: %v == %v", hash1, hash2)
	}
}
