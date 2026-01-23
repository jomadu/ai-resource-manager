package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestCleanSinks(t *testing.T) {
	// Build the binary once
	tmpBin := filepath.Join(t.TempDir(), "arm")
	cmd := exec.Command("go", "build", "-o", tmpBin, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		setupFunc   func(string) error
		wantErr     bool
		wantContain string
	}{
		{
			name:        "clean sinks without manifest",
			args:        []string{"clean", "sinks"},
			wantErr:     false,
			wantContain: "Sinks cleaned successfully",
		},
		{
			name: "clean sinks with empty manifest",
			args: []string{"clean", "sinks"},
			setupFunc: func(dir string) error {
				manifestPath := filepath.Join(dir, "arm-manifest.json")
				return os.WriteFile(manifestPath, []byte(`{"version":1,"registries":{},"sinks":{},"rulesets":{},"promptsets":{}}`), 0o644)
			},
			wantErr:     false,
			wantContain: "Sinks cleaned successfully",
		},
		{
			name: "clean sinks with configured sinks",
			args: []string{"clean", "sinks"},
			setupFunc: func(dir string) error {
				manifestPath := filepath.Join(dir, "arm-manifest.json")
				sinkDir := filepath.Join(dir, "test-sink")
				if err := os.MkdirAll(sinkDir, 0o755); err != nil {
					return err
				}
				manifest := `{
					"version": 1,
					"registries": {},
					"sinks": {
						"test-sink": {
							"tool": "cursor",
							"directory": "` + sinkDir + `"
						}
					},
					"rulesets": {},
					"promptsets": {}
				}`
				return os.WriteFile(manifestPath, []byte(manifest), 0o644)
			},
			wantErr:     false,
			wantContain: "Sinks cleaned successfully",
		},
		{
			name: "nuke sinks",
			args: []string{"clean", "sinks", "--nuke"},
			setupFunc: func(dir string) error {
				manifestPath := filepath.Join(dir, "arm-manifest.json")
				sinkDir := filepath.Join(dir, "test-sink")
				if err := os.MkdirAll(sinkDir, 0o755); err != nil {
					return err
				}
				manifest := `{
					"version": 1,
					"registries": {},
					"sinks": {
						"test-sink": {
							"tool": "cursor",
							"directory": "` + sinkDir + `"
						}
					},
					"rulesets": {},
					"promptsets": {}
				}`
				return os.WriteFile(manifestPath, []byte(manifest), 0o644)
			},
			wantErr:     false,
			wantContain: "Sinks nuked successfully",
		},
		{
			name:        "clean sinks with unknown flag",
			args:        []string{"clean", "sinks", "--unknown"},
			wantErr:     true,
			wantContain: "Unknown flag: --unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			if tt.setupFunc != nil {
				if err := tt.setupFunc(tmpDir); err != nil {
					t.Fatalf("Setup failed: %v", err)
				}
			}

			cmd := exec.Command(tmpBin, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+filepath.Join(tmpDir, "arm-manifest.json"))
			output, err := cmd.CombinedOutput()

			if tt.wantErr && err == nil {
				t.Errorf("Expected error but got none. Output: %s", output)
			}
			if !tt.wantErr && err != nil {
				t.Errorf("Unexpected error: %v. Output: %s", err, output)
			}
			if tt.wantContain != "" && !strings.Contains(string(output), tt.wantContain) {
				t.Errorf("Expected output to contain %q, got: %s", tt.wantContain, output)
			}
		})
	}
}
