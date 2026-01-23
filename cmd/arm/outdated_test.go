package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestOutdated(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name           string
		setupManifest  string
		setupLockfile  string
		args           []string
		wantErr        bool
		wantContains   []string
		wantNotContain []string
	}{
		{
			name: "no outdated packages",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated"},
			wantContains: []string{"All packages are up to date"},
		},
		{
			name: "table output format (default)",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated"},
			wantContains: []string{"All packages are up to date"},
		},
		{
			name: "json output format",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated", "--output", "json"},
			wantContains: []string{"All packages are up to date"},
		},
		{
			name: "list output format",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated", "--output", "list"},
			wantContains: []string{"All packages are up to date"},
		},
		{
			name: "invalid output format",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated", "--output", "invalid"},
			wantErr:      true,
			wantContains: []string{"Invalid output format"},
		},
		{
			name: "missing output value",
			setupManifest: `{
				"version": 1,
				"registries": {},
				"sinks": {},
				"rulesets": {},
				"promptsets": {}
			}`,
			setupLockfile: `{
				"version": 1,
				"rulesets": {},
				"promptsets": {}
			}`,
			args:         []string{"outdated", "--output"},
			wantErr:      true,
			wantContains: []string{"--output requires a value"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test directory
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm-manifest.json")
			lockfilePath := filepath.Join(testDir, "arm-manifest-lock.json")

			// Write manifest
			if err := os.WriteFile(manifestPath, []byte(tt.setupManifest), 0o644); err != nil {
				t.Fatalf("Failed to write manifest: %v", err)
			}

			// Write lockfile
			if err := os.WriteFile(lockfilePath, []byte(tt.setupLockfile), 0o644); err != nil {
				t.Fatalf("Failed to write lockfile: %v", err)
			}

			// Run command
			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}

			outputStr := string(output)
			for _, want := range tt.wantContains {
				if !strings.Contains(outputStr, want) {
					t.Errorf("Output missing expected string %q. Got: %s", want, outputStr)
				}
			}

			for _, notWant := range tt.wantNotContain {
				if strings.Contains(outputStr, notWant) {
					t.Errorf("Output contains unexpected string %q. Got: %s", notWant, outputStr)
				}
			}
		})
	}
}
