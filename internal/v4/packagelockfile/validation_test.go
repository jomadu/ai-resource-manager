package packagelockfile

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestValidation_SampleLockFileFormat(t *testing.T) {
	ctx := context.Background()

	// Create sample lock file with exact format from docs/examples/demo/project/sample.arm-lock.json
	sampleContent := `{
    "version": "1.0.0",
    "packages": {
        "sample-registry/clean-code-ruleset": {
            "version": "1.1.0",
            "checksum": "sha256:a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456"
        },
        "sample-registry/code-review-promptset": {
            "version": "1.1.0",
            "checksum": "sha256:fedcba0987654321098765432109876543210fedcba098765432109876543210"
        }
    }
}`

	tmpDir := t.TempDir()
	lockPath := filepath.Join(tmpDir, "sample.arm-lock.json")

	err := os.WriteFile(lockPath, []byte(sampleContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write sample lockfile: %v", err)
	}

	fm := NewFileManagerWithPath(lockPath)

	// Test reading the sample lockfile
	lockfile, err := fm.GetPackageLockfile(ctx)
	if err != nil {
		t.Fatalf("GetPackageLockfile() error = %v", err)
	}

	// Validate structure
	if lockfile.Version != "1.0.0" {
		t.Errorf("Version = %v, want 1.0.0", lockfile.Version)
	}

	if len(lockfile.Packages) != 2 {
		t.Errorf("Packages length = %v, want 2", len(lockfile.Packages))
	}

	// Test getting specific packages
	cleanCodeInfo, err := fm.GetPackageLockInfo(ctx, "sample-registry", "clean-code-ruleset")
	if err != nil {
		t.Fatalf("GetPackageLockInfo() error = %v", err)
	}

	if cleanCodeInfo.Version != "1.1.0" {
		t.Errorf("clean-code-ruleset Version = %v, want 1.1.0", cleanCodeInfo.Version)
	}

	if cleanCodeInfo.Checksum != "sha256:a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456" {
		t.Errorf("clean-code-ruleset Checksum mismatch")
	}

	codeReviewInfo, err := fm.GetPackageLockInfo(ctx, "sample-registry", "code-review-promptset")
	if err != nil {
		t.Fatalf("GetPackageLockInfo() error = %v", err)
	}

	if codeReviewInfo.Version != "1.1.0" {
		t.Errorf("code-review-promptset Version = %v, want 1.1.0", codeReviewInfo.Version)
	}

	if codeReviewInfo.Checksum != "sha256:fedcba0987654321098765432109876543210fedcba098765432109876543210" {
		t.Errorf("code-review-promptset Checksum mismatch")
	}

	// Test modifying the lockfile
	newInfo := &PackageLockInfo{
		Version:  "1.2.0",
		Checksum: "sha256:newchecksum123",
	}

	err = fm.UpsertPackageLockInfo(ctx, "sample-registry", "clean-code-ruleset", newInfo)
	if err != nil {
		t.Fatalf("UpsertPackageLockInfo() error = %v", err)
	}

	// Verify the update worked
	updatedInfo, err := fm.GetPackageLockInfo(ctx, "sample-registry", "clean-code-ruleset")
	if err != nil {
		t.Fatalf("GetPackageLockInfo() after update error = %v", err)
	}

	if updatedInfo.Version != "1.2.0" {
		t.Errorf("Updated Version = %v, want 1.2.0", updatedInfo.Version)
	}
}