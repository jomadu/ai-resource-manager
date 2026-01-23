package config

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestFileManager_GetAllSections(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFiles  func(t *testing.T) (string, string) // workingDir, userHomeDir
		want        map[string]map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - project file only",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123

[registry https://gitlab.example.com/group/456]
token = project-token-456
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/123": {
					"token": "project-token-123",
				},
				"registry https://gitlab.example.com/group/456": {
					"token": "project-token-456",
				},
			},
			wantErr: false,
		},
		{
			name: "success - user home file only",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				userRc := filepath.Join(userHomeDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/789]
token = user-token-789
`
				if err := os.WriteFile(userRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/789": {
					"token": "user-token-789",
				},
			},
			wantErr: false,
		},
		{
			name: "success - project overrides user home (same section)",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				projectContent := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(projectContent), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				userRc := filepath.Join(userHomeDir, ".armrc")
				userContent := `[registry https://gitlab.example.com/project/123]
token = user-token-123

[registry https://gitlab.example.com/group/456]
token = user-token-456
`
				if err := os.WriteFile(userRc, []byte(userContent), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/123": {
					"token": "project-token-123",
				},
				"registry https://gitlab.example.com/group/456": {
					"token": "user-token-456",
				},
			},
			wantErr: false,
		},
		{
			name: "success - both files with different sections",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				projectContent := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(projectContent), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				userRc := filepath.Join(userHomeDir, ".armrc")
				userContent := `[registry https://gitlab.example.com/group/456]
token = user-token-456
`
				if err := os.WriteFile(userRc, []byte(userContent), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/123": {
					"token": "project-token-123",
				},
				"registry https://gitlab.example.com/group/456": {
					"token": "user-token-456",
				},
			},
			wantErr: false,
		},
		{
			name: "success - no files exist, returns empty map",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()
				return workingDir, userHomeDir
			},
			want:    make(map[string]map[string]string),
			wantErr: false,
		},
		{
			name: "success - environment variable expansion",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				_ = os.Setenv("TEST_TOKEN", "env-token-123")
				t.Cleanup(func() { _ = os.Unsetenv("TEST_TOKEN") })

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = ${TEST_TOKEN}
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/123": {
					"token": "env-token-123",
				},
			},
			wantErr: false,
		},
		{
			name: "success - missing env var expands to empty",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = ${MISSING_VAR}
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			want: map[string]map[string]string{
				"registry https://gitlab.example.com/project/123": {
					"token": "",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workingDir, userHomeDir := tt.setupFiles(t)
			fm := NewFileManagerWithPaths(workingDir, userHomeDir)

			got, err := fm.GetAllSections(ctx)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSections() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetAllSections() expected error but got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetAllSections() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_GetSection(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFiles  func(t *testing.T) (string, string) // workingDir, userHomeDir
		section     string
		want        map[string]string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - section in project file",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			want: map[string]string{
				"token": "project-token-123",
			},
			wantErr: false,
		},
		{
			name: "success - section in user home file",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				userRc := filepath.Join(userHomeDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/789]
token = user-token-789
`
				if err := os.WriteFile(userRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/789",
			want: map[string]string{
				"token": "user-token-789",
			},
			wantErr: false,
		},
		{
			name: "success - project overrides user home (same section)",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				projectContent := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(projectContent), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				userRc := filepath.Join(userHomeDir, ".armrc")
				userContent := `[registry https://gitlab.example.com/project/123]
token = user-token-123
`
				if err := os.WriteFile(userRc, []byte(userContent), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			want: map[string]string{
				"token": "project-token-123",
			},
			wantErr: false,
		},
		{
			name: "success - multiple keys in section",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
api_version = v4
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			want: map[string]string{
				"token":       "project-token-123",
				"api_version": "v4",
			},
			wantErr: false,
		},
		{
			name: "success - environment variable expansion",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				_ = os.Setenv("TEST_TOKEN", "env-token-123")
				t.Cleanup(func() { _ = os.Unsetenv("TEST_TOKEN") })

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = ${TEST_TOKEN}
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			want: map[string]string{
				"token": "env-token-123",
			},
			wantErr: false,
		},
		{
			name: "error - section not found in either file",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section:     "registry https://gitlab.example.com/project/999",
			want:        nil,
			wantErr:     true,
			errContains: "section not found",
		},
		{
			name: "error - no files exist",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()
				return workingDir, userHomeDir
			},
			section:     "registry https://gitlab.example.com/project/123",
			want:        nil,
			wantErr:     true,
			errContains: "section not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workingDir, userHomeDir := tt.setupFiles(t)
			fm := NewFileManagerWithPaths(workingDir, userHomeDir)

			got, err := fm.GetSection(ctx, tt.section)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetSection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetSection() expected error but got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetSection() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSection() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_GetValue(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		setupFiles  func(t *testing.T) (string, string) // workingDir, userHomeDir
		section     string
		key         string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name: "success - key exists in project file",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			key:     "token",
			want:    "project-token-123",
			wantErr: false,
		},
		{
			name: "success - key exists in user home file",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				userRc := filepath.Join(userHomeDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/789]
token = user-token-789
`
				if err := os.WriteFile(userRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/789",
			key:     "token",
			want:    "user-token-789",
			wantErr: false,
		},
		{
			name: "success - project overrides user home",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				projectContent := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(projectContent), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				userRc := filepath.Join(userHomeDir, ".armrc")
				userContent := `[registry https://gitlab.example.com/project/123]
token = user-token-123
`
				if err := os.WriteFile(userRc, []byte(userContent), 0o644); err != nil {
					t.Fatalf("failed to write user .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			key:     "token",
			want:    "project-token-123",
			wantErr: false,
		},
		{
			name: "success - environment variable expansion",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				_ = os.Setenv("TEST_TOKEN", "env-token-123")
				t.Cleanup(func() { _ = os.Unsetenv("TEST_TOKEN") })

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = ${TEST_TOKEN}
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section: "registry https://gitlab.example.com/project/123",
			key:     "token",
			want:    "env-token-123",
			wantErr: false,
		},
		{
			name: "error - section not found",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section:     "registry https://gitlab.example.com/project/999",
			key:         "token",
			want:        "",
			wantErr:     true,
			errContains: "section not found",
		},
		{
			name: "error - key not found in section",
			setupFiles: func(t *testing.T) (string, string) {
				workingDir := t.TempDir()
				userHomeDir := t.TempDir()

				projectRc := filepath.Join(workingDir, ".armrc")
				content := `[registry https://gitlab.example.com/project/123]
token = project-token-123
`
				if err := os.WriteFile(projectRc, []byte(content), 0o644); err != nil {
					t.Fatalf("failed to write project .armrc: %v", err)
				}

				return workingDir, userHomeDir
			},
			section:     "registry https://gitlab.example.com/project/123",
			key:         "missing_key",
			want:        "",
			wantErr:     true,
			errContains: "key not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workingDir, userHomeDir := tt.setupFiles(t)
			fm := NewFileManagerWithPaths(workingDir, userHomeDir)

			got, err := fm.GetValue(ctx, tt.section, tt.key)

			if (err != nil) != tt.wantErr {
				t.Errorf("GetValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				if err == nil {
					t.Errorf("GetValue() expected error but got nil")
				} else if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("GetValue() error = %v, should contain %v", err, tt.errContains)
				}
				return
			}

			if got != tt.want {
				t.Errorf("GetValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileManager_NewFileManager(t *testing.T) {
	fm := NewFileManager()
	if fm == nil {
		t.Fatal("NewFileManager() returned nil")
	}
	if fm.workingDir == "" && fm.userHomeDir == "" {
		t.Error("NewFileManager() should have at least one directory set")
	}
}

func TestFileManager_NewFileManagerWithPaths(t *testing.T) {
	workingDir := "/test/working"
	userHomeDir := "/test/home"

	fm := NewFileManagerWithPaths(workingDir, userHomeDir)
	if fm == nil {
		t.Fatal("NewFileManagerWithPaths() returned nil")
	}
	if fm.workingDir != workingDir {
		t.Errorf("NewFileManagerWithPaths() workingDir = %v, want %v", fm.workingDir, workingDir)
	}
	if fm.userHomeDir != userHomeDir {
		t.Errorf("NewFileManagerWithPaths() userHomeDir = %v, want %v", fm.userHomeDir, userHomeDir)
	}
}
