package helpers

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// ARMRunner runs ARM commands in a test environment
type ARMRunner struct {
	BinaryPath string
	WorkDir    string
	t          *testing.T
}

// NewARMRunner creates a new ARM command runner
func NewARMRunner(t *testing.T, workDir string) *ARMRunner {
	t.Helper()

	// Find the ARM binary (should be in project root)
	binaryPath := findARMBinary(t)

	return &ARMRunner{
		BinaryPath: binaryPath,
		WorkDir:    workDir,
		t:          t,
	}
}

// Run executes an ARM command and returns stdout, stderr, and error
func (r *ARMRunner) Run(args ...string) (stdout, stderr string, err error) {
	r.t.Helper()

	cmd := exec.Command(r.BinaryPath, args...)
	cmd.Dir = r.WorkDir

	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	err = cmd.Run()
	return stdoutBuf.String(), stderrBuf.String(), err
}

// MustRun executes an ARM command and fails the test if it returns an error
func (r *ARMRunner) MustRun(args ...string) string {
	r.t.Helper()
	stdout, stderr, err := r.Run(args...)
	if err != nil {
		r.t.Fatalf("ARM command failed: arm %v\nStdout: %s\nStderr: %s\nError: %v",
			args, stdout, stderr, err)
	}
	return stdout
}

// MustFail executes an ARM command and fails the test if it succeeds
func (r *ARMRunner) MustFail(args ...string) string {
	r.t.Helper()

	stdout, stderr, err := r.Run(args...)
	if err == nil {
		r.t.Fatalf("ARM command should have failed but succeeded: arm %v\nStdout: %s\nStderr: %s",
			args, stdout, stderr)
	}
	return stderr
}

// findARMBinary locates the ARM binary for testing
func findARMBinary(t *testing.T) string {
	t.Helper()

	// Try to find the binary in the project root
	// Start from the test directory and walk up
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	// Walk up to find the project root (contains go.mod)
	for {
		binaryPath := filepath.Join(dir, "arm")
		if _, err := os.Stat(binaryPath); err == nil {
			return binaryPath
		}

		// Check if we've reached the root
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatal("ARM binary not found. Run 'go build -o arm cmd/arm/main.go' first")
	return ""
}
