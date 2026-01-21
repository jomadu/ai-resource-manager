package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestAddSink(t *testing.T) {
	// Build the binary once
	tmpBin := filepath.Join(t.TempDir(), "arm")
	cmd := exec.Command("go", "build", "-o", tmpBin, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
		validate    func(t *testing.T, manifestPath string)
	}{
		{
			name:        "missing tool flag",
			args:        []string{"add", "sink", "my-sink", ".cursor/rules"},
			wantErr:     true,
			errContains: "--tool is required",
		},
		{
			name:        "missing name",
			args:        []string{"add", "sink", "--tool", "cursor"},
			wantErr:     true,
			errContains: "NAME is required",
		},
		{
			name:        "missing path",
			args:        []string{"add", "sink", "--tool", "cursor", "my-sink"},
			wantErr:     true,
			errContains: "PATH is required",
		},
		{
			name:        "invalid tool",
			args:        []string{"add", "sink", "--tool", "invalid", "my-sink", ".cursor/rules"},
			wantErr:     true,
			errContains: "Invalid tool",
		},
		{
			name:    "add cursor sink",
			args:    []string{"add", "sink", "--tool", "cursor", "my-sink", ".cursor/rules"},
			wantErr: false,
			validate: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatalf("Failed to parse manifest: %v", err)
				}
				sinks := manifest["sinks"].(map[string]interface{})
				sink := sinks["my-sink"].(map[string]interface{})
				if sink["directory"] != ".cursor/rules" {
					t.Errorf("Expected directory .cursor/rules, got %v", sink["directory"])
				}
				if sink["tool"] != "cursor" {
					t.Errorf("Expected tool cursor, got %v", sink["tool"])
				}
			},
		},
		{
			name:    "add amazonq sink",
			args:    []string{"add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules"},
			wantErr: false,
			validate: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatalf("Failed to parse manifest: %v", err)
				}
				sinks := manifest["sinks"].(map[string]interface{})
				sink := sinks["q-rules"].(map[string]interface{})
				if sink["tool"] != "amazonq" {
					t.Errorf("Expected tool amazonq, got %v", sink["tool"])
				}
			},
		},
		{
			name:    "add copilot sink",
			args:    []string{"add", "sink", "--tool", "copilot", "copilot-rules", ".github/copilot"},
			wantErr: false,
			validate: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatalf("Failed to parse manifest: %v", err)
				}
				sinks := manifest["sinks"].(map[string]interface{})
				sink := sinks["copilot-rules"].(map[string]interface{})
				if sink["tool"] != "copilot" {
					t.Errorf("Expected tool copilot, got %v", sink["tool"])
				}
			},
		},
		{
			name:    "add markdown sink",
			args:    []string{"add", "sink", "--tool", "markdown", "md-rules", "./docs"},
			wantErr: false,
			validate: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatalf("Failed to parse manifest: %v", err)
				}
				sinks := manifest["sinks"].(map[string]interface{})
				sink := sinks["md-rules"].(map[string]interface{})
				if sink["tool"] != "markdown" {
					t.Errorf("Expected tool markdown, got %v", sink["tool"])
				}
			},
		},
		{
			name:    "add sink with force flag",
			args:    []string{"add", "sink", "--tool", "cursor", "--force", "my-sink", ".cursor/rules"},
			wantErr: false,
			validate: func(t *testing.T, manifestPath string) {
				data, err := os.ReadFile(manifestPath)
				if err != nil {
					t.Fatalf("Failed to read manifest: %v", err)
				}
				var manifest map[string]interface{}
				if err := json.Unmarshal(data, &manifest); err != nil {
					t.Fatalf("Failed to parse manifest: %v", err)
				}
				sinks := manifest["sinks"].(map[string]interface{})
				if _, ok := sinks["my-sink"]; !ok {
					t.Error("Expected sink my-sink to exist")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			manifestPath := filepath.Join(tmpDir, "arm.json")

			// Create empty manifest
			if err := os.WriteFile(manifestPath, []byte(`{"sinks":{},"registries":{},"dependencies":{}}`), 0644); err != nil {
				t.Fatalf("Failed to create manifest: %v", err)
			}

			cmd := exec.Command(tmpBin, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
				if tt.validate != nil {
					tt.validate(t, manifestPath)
				}
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && containsAt(s, substr))
}

func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
