package e2e

import (
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestErrorHandling(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Setup: Create test Git repository
	repoDir := t.TempDir()
	repo := helpers.NewGitRepo(t, repoDir)
	repo.WriteFile("test-ruleset.yml", helpers.MinimalRuleset)
	repo.Commit("Initial commit")
	repo.Tag("v1.0.0")

	// Setup: Add registry and sink
	repoURL := "file://" + repoDir
	arm.MustRun("add", "registry", "git", "--url", repoURL, "test-registry")
	arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

	t.Run("InstallNonExistentVersion", func(t *testing.T) {
		stdout, stderr, err := arm.Run("install", "ruleset", "test-registry/test-ruleset@99.99.99", "cursor-rules")
		if err == nil {
			t.Error("expected error when installing non-existent version")
		}
		output := stdout + stderr
		if !strings.Contains(output, "not found") && !strings.Contains(output, "no matching version") && !strings.Contains(output, "no version satisfies") {
			t.Errorf("expected version error in message, got: %s", output)
		}
	})

	t.Run("InstallToNonExistentSink", func(t *testing.T) {
		_, stderr, err := arm.Run("install", "ruleset", "test-registry/test-ruleset@1.0.0", "non-existent-sink")
		if err == nil {
			t.Error("expected error when installing to non-existent sink")
		}
		if !strings.Contains(stderr, "sink") && !strings.Contains(stderr, "not found") {
			t.Errorf("expected sink error in message, got: %s", stderr)
		}
	})

	t.Run("InstallFromNonExistentRegistry", func(t *testing.T) {
		_, stderr, err := arm.Run("install", "ruleset", "non-existent-registry/test-ruleset@1.0.0", "cursor-rules")
		if err == nil {
			t.Error("expected error when installing from non-existent registry")
		}
		if !strings.Contains(stderr, "registry") && !strings.Contains(stderr, "not found") {
			t.Errorf("expected registry error in message, got: %s", stderr)
		}
	})

	t.Run("AddDuplicateRegistryWithoutForce", func(t *testing.T) {
		// Already added in setup
		_, stderr, err := arm.Run("add", "registry", "git", "--url", repoURL, "test-registry")
		if err == nil {
			t.Error("expected error when adding duplicate registry without --force")
		}
		if !strings.Contains(stderr, "already exists") && !strings.Contains(stderr, "duplicate") {
			t.Errorf("expected duplicate error in message, got: %s", stderr)
		}
	})

	t.Run("AddDuplicateSinkWithoutForce", func(t *testing.T) {
		// Already added in setup
		_, stderr, err := arm.Run("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")
		if err == nil {
			t.Error("expected error when adding duplicate sink without --force")
		}
		if !strings.Contains(stderr, "already exists") && !strings.Contains(stderr, "duplicate") {
			t.Errorf("expected duplicate error in message, got: %s", stderr)
		}
	})

	t.Run("RemoveNonExistentRegistry", func(t *testing.T) {
		_, stderr, err := arm.Run("remove", "registry", "non-existent-registry")
		if err == nil {
			t.Error("expected error when removing non-existent registry")
		}
		if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "does not exist") {
			t.Errorf("expected not found error in message, got: %s", stderr)
		}
	})

	t.Run("RemoveNonExistentSink", func(t *testing.T) {
		_, stderr, err := arm.Run("remove", "sink", "non-existent-sink")
		if err == nil {
			t.Error("expected error when removing non-existent sink")
		}
		if !strings.Contains(stderr, "not found") && !strings.Contains(stderr, "does not exist") {
			t.Errorf("expected not found error in message, got: %s", stderr)
		}
	})

	t.Run("InvalidVersionConstraint", func(t *testing.T) {
		stdout, stderr, err := arm.Run("install", "ruleset", "test-registry/test-ruleset@invalid-version", "cursor-rules")
		if err == nil {
			t.Error("expected error with invalid version constraint")
		}
		output := stdout + stderr
		if !strings.Contains(output, "version") && !strings.Contains(output, "invalid") {
			t.Errorf("expected version error in message, got: %s", output)
		}
	})
}
