package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAddGitRegistry(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid git registry",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg"},
			wantErr: false,
		},
		{
			name:    "with branches",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "--branches", "main,dev", "test-reg2"},
			wantErr: false,
		},
		{
			name:        "missing url",
			args:        []string{"add", "registry", "git", "test-reg"},
			wantErr:     true,
			errContains: "--url is required",
		},
		{
			name:        "missing name",
			args:        []string{"add", "registry", "git", "--url", "https://github.com/test/repo"},
			wantErr:     true,
			errContains: "NAME is required",
		},
		{
			name:        "duplicate without force",
			args:        []string{"add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg"},
			wantErr:     true,
			errContains: "registry already exists",
		},
		{
			name:    "duplicate with force",
			args:    []string{"add", "registry", "git", "--url", "https://github.com/test/repo2", "--force", "test-reg"},
			wantErr: false,
		},
	}

	// Build the binary once
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(".") // Current directory is cmd/arm
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each test gets its own manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm-manifest.json")

			// For duplicate tests, create the registry first
			if strings.Contains(tt.name, "duplicate") {
				setupCmd := exec.Command(binPath, "add", "registry", "git", "--url", "https://github.com/test/repo", "test-reg")
				setupCmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
				if err := setupCmd.Run(); err != nil {
					t.Fatalf("Failed to setup duplicate test: %v", err)
				}
			}

			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}
		})
	}
}

func TestAddGitLabRegistry(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid gitlab registry with url only",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "test-gitlab"},
			wantErr: false,
		},
		{
			name:    "with project-id",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "--project-id", "123", "test-gitlab2"},
			wantErr: false,
		},
		{
			name:    "with group-id",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "--group-id", "456", "test-gitlab3"},
			wantErr: false,
		},
		{
			name:    "with api-version",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "--api-version", "v4", "test-gitlab4"},
			wantErr: false,
		},
		{
			name:    "with all options",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "--project-id", "123", "--group-id", "456", "--api-version", "v4", "test-gitlab5"},
			wantErr: false,
		},
		{
			name:        "missing url",
			args:        []string{"add", "registry", "gitlab", "test-gitlab"},
			wantErr:     true,
			errContains: "--url is required",
		},
		{
			name:        "missing name",
			args:        []string{"add", "registry", "gitlab", "--url", "https://gitlab.com"},
			wantErr:     true,
			errContains: "NAME is required",
		},
		{
			name:        "duplicate without force",
			args:        []string{"add", "registry", "gitlab", "--url", "https://gitlab.com", "test-gitlab"},
			wantErr:     true,
			errContains: "registry already exists",
		},
		{
			name:    "duplicate with force",
			args:    []string{"add", "registry", "gitlab", "--url", "https://gitlab.example.com", "--force", "test-gitlab"},
			wantErr: false,
		},
	}

	// Build the binary once
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(".") // Current directory is cmd/arm
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each test gets its own manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm-manifest.json")

			// For duplicate tests, create the registry first
			if strings.Contains(tt.name, "duplicate") {
				setupCmd := exec.Command(binPath, "add", "registry", "gitlab", "--url", "https://gitlab.com", "test-gitlab")
				setupCmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
				if err := setupCmd.Run(); err != nil {
					t.Fatalf("Failed to setup duplicate test: %v", err)
				}
			}

			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}
		})
	}
}


func TestAddCloudsmithRegistry(t *testing.T) {
	tests := []struct {
		name        string
		args        []string
		wantErr     bool
		errContains string
	}{
		{
			name:    "valid cloudsmith registry",
			args:    []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "myorg", "--repo", "myrepo", "test-cs"},
			wantErr: false,
		},
		{
			name:        "missing url",
			args:        []string{"add", "registry", "cloudsmith", "--owner", "myorg", "--repo", "myrepo", "test-cs"},
			wantErr:     true,
			errContains: "--url is required",
		},
		{
			name:        "missing owner",
			args:        []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--repo", "myrepo", "test-cs"},
			wantErr:     true,
			errContains: "--owner is required",
		},
		{
			name:        "missing repo",
			args:        []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "myorg", "test-cs"},
			wantErr:     true,
			errContains: "--repo is required",
		},
		{
			name:        "missing name",
			args:        []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "myorg", "--repo", "myrepo"},
			wantErr:     true,
			errContains: "NAME is required",
		},
		{
			name:        "duplicate without force",
			args:        []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "myorg", "--repo", "myrepo", "test-cs"},
			wantErr:     true,
			errContains: "registry already exists",
		},
		{
			name:    "duplicate with force",
			args:    []string{"add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "neworg", "--repo", "newrepo", "--force", "test-cs"},
			wantErr: false,
		},
	}

	// Build the binary once
	tmpDir := t.TempDir()
	binPath := filepath.Join(tmpDir, "arm")
	cmd := exec.Command("go", "build", "-o", binPath, ".")
	cmd.Dir = filepath.Join(".") // Current directory is cmd/arm
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Each test gets its own manifest
			testDir := t.TempDir()
			manifestPath := filepath.Join(testDir, "arm-manifest.json")

			// For duplicate tests, create the registry first
			if strings.Contains(tt.name, "duplicate") {
				setupCmd := exec.Command(binPath, "add", "registry", "cloudsmith", "--url", "https://cloudsmith.io", "--owner", "myorg", "--repo", "myrepo", "test-cs")
				setupCmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
				if err := setupCmd.Run(); err != nil {
					t.Fatalf("Failed to setup duplicate test: %v", err)
				}
			}

			cmd := exec.Command(binPath, tt.args...)
			cmd.Env = append(os.Environ(), "ARM_MANIFEST_PATH="+manifestPath)
			output, err := cmd.CombinedOutput()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none. Output: %s", output)
				}
				if tt.errContains != "" && !strings.Contains(string(output), tt.errContains) {
					t.Errorf("Expected error containing %q, got: %s", tt.errContains, output)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v. Output: %s", err, output)
				}
			}
		})
	}
}
