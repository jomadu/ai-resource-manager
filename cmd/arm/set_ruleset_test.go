package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetRuleset(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		setupFunc   func(string)
		wantErr     bool
		wantContain string
	}{
		{
			name:        "missing arguments",
			args:        []string{"set", "ruleset"},
			wantErr:     true,
			wantContain: "Usage: arm set ruleset",
		},
		{
			name:        "invalid package spec",
			args:        []string{"set", "ruleset", "invalid", "version", "1.0.0"},
			wantErr:     true,
			wantContain: "Invalid package spec",
		},
		{
			name: "set version",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "version", "2.0.0"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated ruleset 'test-registry/test-ruleset' version",
		},
		{
			name: "set priority",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "priority", "200"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated ruleset 'test-registry/test-ruleset' priority",
		},
		{
			name: "set priority invalid",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "priority", "invalid"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     true,
			wantContain: "Invalid priority value",
		},
		{
			name: "set sinks",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "sinks", "sink1,sink2"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"sinks":{"sink1":{"tool":"cursor","directory":".cursor"},"sink2":{"tool":"amazonq","directory":".amazonq"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated ruleset 'test-registry/test-ruleset' sinks",
		},
		{
			name: "set include",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "include", "*.yml,*.yaml"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated ruleset 'test-registry/test-ruleset' include",
		},
		{
			name: "set exclude",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "exclude", "test/**"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated ruleset 'test-registry/test-ruleset' exclude",
		},
		{
			name: "unknown key",
			args: []string{"set", "ruleset", "test-registry/test-ruleset", "unknown", "value"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-ruleset":{"type":"ruleset","version":"1.0.0","priority":100,"sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     true,
			wantContain: "Unknown key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if tt.setupFunc != nil {
				tt.setupFunc(manifestPath)
			}

			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
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

func TestSetPromptset(t *testing.T) {
	// Build the binary once
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	tests := []struct {
		name        string
		args        []string
		setupFunc   func(string)
		wantErr     bool
		wantContain string
	}{
		{
			name:        "missing arguments",
			args:        []string{"set", "promptset"},
			wantErr:     true,
			wantContain: "Usage: arm set promptset",
		},
		{
			name:        "invalid package spec",
			args:        []string{"set", "promptset", "invalid", "version", "1.0.0"},
			wantErr:     true,
			wantContain: "Invalid package spec",
		},
		{
			name: "set version",
			args: []string{"set", "promptset", "test-registry/test-promptset", "version", "2.0.0"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-promptset":{"type":"promptset","version":"1.0.0","sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated promptset 'test-registry/test-promptset' version",
		},
		{
			name: "set sinks",
			args: []string{"set", "promptset", "test-registry/test-promptset", "sinks", "sink1,sink2"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"sinks":{"sink1":{"tool":"cursor","directory":".cursor"},"sink2":{"tool":"amazonq","directory":".amazonq"}},"dependencies":{"test-registry/test-promptset":{"type":"promptset","version":"1.0.0","sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated promptset 'test-registry/test-promptset' sinks",
		},
		{
			name: "set include",
			args: []string{"set", "promptset", "test-registry/test-promptset", "include", "*.yml,*.yaml"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-promptset":{"type":"promptset","version":"1.0.0","sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated promptset 'test-registry/test-promptset' include",
		},
		{
			name: "set exclude",
			args: []string{"set", "promptset", "test-registry/test-promptset", "exclude", "test/**"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-promptset":{"type":"promptset","version":"1.0.0","sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     false,
			wantContain: "Updated promptset 'test-registry/test-promptset' exclude",
		},
		{
			name: "unknown key",
			args: []string{"set", "promptset", "test-registry/test-promptset", "unknown", "value"},
			setupFunc: func(manifestPath string) {
				content := `{"version":1,"registries":{"test-registry":{"type":"git","url":"https://example.com"}},"dependencies":{"test-registry/test-promptset":{"type":"promptset","version":"1.0.0","sinks":["sink1"]}}}`
				_ = os.WriteFile(manifestPath, []byte(content), 0o644)
			},
			wantErr:     true,
			wantContain: "Unknown key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm.json")

			if tt.setupFunc != nil {
				tt.setupFunc(manifestPath)
			}

			cmd := exec.Command(binaryPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
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
