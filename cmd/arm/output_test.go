package main

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

func TestNewCompileOutputFormatter(t *testing.T) {
	tests := []struct {
		name    string
		verbose bool
		dryRun  bool
	}{
		{"default", false, false},
		{"verbose", true, false},
		{"dry-run", false, true},
		{"verbose dry-run", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewCompileOutputFormatter(tt.verbose, tt.dryRun)
			if formatter == nil {
				t.Error("Expected formatter, got nil")
				return
			}
			if formatter.verbose != tt.verbose {
				t.Errorf("Expected verbose=%v, got %v", tt.verbose, formatter.verbose)
			}
			if formatter.dryRun != tt.dryRun {
				t.Errorf("Expected dryRun=%v, got %v", tt.dryRun, formatter.dryRun)
			}
		})
	}
}

func TestCompileOutputFormatter_DisplayResults(t *testing.T) {
	tests := []struct {
		name        string
		result      *arm.CompileResult
		verbose     bool
		dryRun      bool
		expectError bool
	}{
		{
			name:        "nil result",
			result:      nil,
			verbose:     false,
			dryRun:      false,
			expectError: true,
		},
		{
			name: "success result",
			result: &arm.CompileResult{
				CompiledFiles: []arm.CompiledFile{
					{
						SourcePath: "test.yaml",
						TargetPath: "test.mdc",
						Target:     urf.TargetCursor,
						RuleCount:  2,
					},
				},
				Skipped: []arm.SkippedFile{},
				Errors:  []arm.CompileError{},
				Stats: arm.CompileStats{
					FilesProcessed: 1,
					FilesCompiled:  1,
					FilesSkipped:   0,
					RulesGenerated: 2,
					Errors:         0,
					TargetStats:    map[string]int{"cursor": 1},
				},
			},
			verbose:     false,
			dryRun:      false,
			expectError: false,
		},
		{
			name: "result with errors",
			result: &arm.CompileResult{
				CompiledFiles: []arm.CompiledFile{},
				Skipped:       []arm.SkippedFile{},
				Errors: []arm.CompileError{
					{
						FilePath: "invalid.yaml",
						Target:   "cursor",
						Error:    "parse error",
					},
				},
				Stats: arm.CompileStats{
					FilesProcessed: 1,
					FilesCompiled:  0,
					FilesSkipped:   0,
					RulesGenerated: 0,
					Errors:         1,
					TargetStats:    map[string]int{},
				},
			},
			verbose:     false,
			dryRun:      false,
			expectError: true,
		},
		{
			name: "multi-target result",
			result: &arm.CompileResult{
				CompiledFiles: []arm.CompiledFile{
					{SourcePath: "test.yaml", TargetPath: "cursor/test.mdc", Target: urf.TargetCursor},
					{SourcePath: "test.yaml", TargetPath: "amazonq/test.md", Target: urf.TargetAmazonQ},
				},
				Skipped: []arm.SkippedFile{},
				Errors:  []arm.CompileError{},
				Stats: arm.CompileStats{
					FilesProcessed: 1,
					FilesCompiled:  2,
					FilesSkipped:   0,
					RulesGenerated: 2,
					Errors:         0,
					TargetStats:    map[string]int{"cursor": 1, "amazonq": 1},
				},
			},
			verbose:     true,
			dryRun:      false,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewCompileOutputFormatter(tt.verbose, tt.dryRun)
			err := formatter.DisplayResults(tt.result)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestCompileOutputFormatter_DisplayValidationResults(t *testing.T) {
	formatter := NewCompileOutputFormatter(true, false)

	tests := []struct {
		name   string
		result *arm.CompileResult
	}{
		{
			name: "valid files",
			result: &arm.CompileResult{
				CompiledFiles: []arm.CompiledFile{},
				Skipped:       []arm.SkippedFile{},
				Errors:        []arm.CompileError{},
				Stats: arm.CompileStats{
					FilesProcessed: 2,
					FilesCompiled:  0,
					FilesSkipped:   0,
					RulesGenerated: 0,
					Errors:         0,
					TargetStats:    map[string]int{},
				},
			},
		},
		{
			name: "invalid files",
			result: &arm.CompileResult{
				CompiledFiles: []arm.CompiledFile{},
				Skipped:       []arm.SkippedFile{},
				Errors: []arm.CompileError{
					{FilePath: "invalid.yaml", Error: "validation error"},
				},
				Stats: arm.CompileStats{
					FilesProcessed: 2,
					FilesCompiled:  0,
					FilesSkipped:   0,
					RulesGenerated: 0,
					Errors:         1,
					TargetStats:    map[string]int{},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: This function prints to stdout, so we're just testing it doesn't panic
			formatter.DisplayValidationResults(tt.result)
		})
	}
}

func TestCompileOutputFormatter_DisplayDryRunPlan(_ *testing.T) {
	formatter := NewCompileOutputFormatter(true, true)

	result := &arm.CompileResult{
		CompiledFiles: []arm.CompiledFile{
			{SourcePath: "test.yaml", TargetPath: "test.mdc", Target: urf.TargetCursor},
		},
		Skipped: []arm.SkippedFile{
			{Path: "non-urf.txt", Reason: "not a URF file"},
		},
		Errors: []arm.CompileError{},
		Stats: arm.CompileStats{
			FilesProcessed: 2,
			FilesCompiled:  1,
			FilesSkipped:   1,
			RulesGenerated: 1,
			Errors:         0,
			TargetStats:    map[string]int{"cursor": 1},
		},
	}

	// Test that it doesn't panic
	formatter.DisplayDryRunPlan(result)
}

func TestCompileOutputFormatter_DisplayCompileProgress(t *testing.T) {
	tests := []struct {
		name       string
		verbose    bool
		sourceFile string
		targetFile string
		target     string
	}{
		{
			name:       "verbose mode",
			verbose:    true,
			sourceFile: "rules.yaml",
			targetFile: "rules.mdc",
			target:     "cursor",
		},
		{
			name:       "non-verbose mode",
			verbose:    false,
			sourceFile: "rules.yaml",
			targetFile: "rules.mdc",
			target:     "cursor",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewCompileOutputFormatter(tt.verbose, false)
			// Test that it doesn't panic
			formatter.DisplayCompileProgress(tt.sourceFile, tt.targetFile, tt.target)
		})
	}
}

func TestCompileOutputFormatter_getExitError(t *testing.T) {
	formatter := NewCompileOutputFormatter(false, false)

	tests := []struct {
		name           string
		result         *arm.CompileResult
		expectError    bool
		expectExitCall bool // For cases where os.Exit would be called
	}{
		{
			name: "no errors",
			result: &arm.CompileResult{
				Stats: arm.CompileStats{
					FilesCompiled: 2,
					Errors:        0,
				},
			},
			expectError:    false,
			expectExitCall: false,
		},
		{
			name: "partial failure",
			result: &arm.CompileResult{
				Stats: arm.CompileStats{
					FilesCompiled: 1,
					Errors:        1,
				},
			},
			expectError:    true,
			expectExitCall: false,
		},
		// Note: We can't easily test the total failure case (os.Exit(2))
		// without modifying the function to be more testable
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatter.getExitError(tt.result)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}
