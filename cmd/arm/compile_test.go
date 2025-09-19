package main

import (
	"testing"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/jomadu/ai-rules-manager/internal/urf"
)

func TestParseTargets(t *testing.T) {
	tests := []struct {
		name        string
		targetStr   string
		expected    []urf.CompileTarget
		expectError bool
	}{
		{
			name:        "single target",
			targetStr:   "cursor",
			expected:    []urf.CompileTarget{urf.TargetCursor},
			expectError: false,
		},
		{
			name:        "multiple targets",
			targetStr:   "cursor,amazonq,markdown",
			expected:    []urf.CompileTarget{urf.TargetCursor, urf.TargetAmazonQ, urf.TargetMarkdown},
			expectError: false,
		},
		{
			name:        "targets with spaces",
			targetStr:   "cursor, amazonq , markdown",
			expected:    []urf.CompileTarget{urf.TargetCursor, urf.TargetAmazonQ, urf.TargetMarkdown},
			expectError: false,
		},
		{
			name:        "empty target string",
			targetStr:   "",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "duplicate targets",
			targetStr:   "cursor,cursor,amazonq",
			expected:    nil,
			expectError: true,
		},
		{
			name:        "empty target in list",
			targetStr:   "cursor,,amazonq",
			expected:    []urf.CompileTarget{urf.TargetCursor, urf.TargetAmazonQ},
			expectError: false,
		},
		{
			name:        "all supported targets",
			targetStr:   "cursor,amazonq,markdown,copilot",
			expected:    []urf.CompileTarget{urf.TargetCursor, urf.TargetAmazonQ, urf.TargetMarkdown, urf.TargetCopilot},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseTargets(tt.targetStr)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.targetStr)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input %q: %v", tt.targetStr, err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d targets, got %d", len(tt.expected), len(result))
				return
			}

			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("Expected target %s at position %d, got %s", expected, i, result[i])
				}
			}
		})
	}
}

func TestDisplayCompileResults(t *testing.T) {
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
			name: "empty result",
			result: &arm.CompileResult{
				CompiledFiles: make([]arm.CompiledFile, 0),
				Skipped:       make([]arm.SkippedFile, 0),
				Errors:        make([]arm.CompileError, 0),
				Stats: arm.CompileStats{
					FilesProcessed: 0,
					FilesCompiled:  0,
					FilesSkipped:   0,
					RulesGenerated: 0,
					Errors:         0,
					TargetStats:    make(map[string]int),
				},
			},
			verbose:     false,
			dryRun:      false,
			expectError: false,
		},
		{
			name: "result with errors",
			result: &arm.CompileResult{
				CompiledFiles: make([]arm.CompiledFile, 0),
				Skipped:       make([]arm.SkippedFile, 0),
				Errors: []arm.CompileError{
					{FilePath: "test.yaml", Error: "test error"},
				},
				Stats: arm.CompileStats{
					FilesProcessed: 1,
					FilesCompiled:  0,
					FilesSkipped:   0,
					RulesGenerated: 0,
					Errors:         1,
					TargetStats:    make(map[string]int),
				},
			},
			verbose:     false,
			dryRun:      false,
			expectError: true, // Should return error for exit code handling
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
