package urf

import (
	"testing"
)

func TestDefaultCompilerFactory_GetCompiler(t *testing.T) {
	factory := NewCompilerFactory()

	tests := []struct {
		name        string
		target      CompileTarget
		expectError bool
	}{
		{
			name:        "cursor compiler",
			target:      TargetCursor,
			expectError: false,
		},
		{
			name:        "amazonq compiler",
			target:      TargetAmazonQ,
			expectError: false,
		},
		{
			name:        "unsupported target",
			target:      CompileTarget("unsupported"),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			compiler, err := factory.GetCompiler(tt.target)

			if tt.expectError {
				if err == nil {
					t.Error("Expected error but got none")
				}
				if compiler != nil {
					t.Error("Expected nil compiler on error")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if compiler == nil {
				t.Error("Expected compiler but got nil")
			}
		})
	}
}

func TestDefaultCompilerFactory_SupportedTargets(t *testing.T) {
	factory := NewCompilerFactory()
	targets := factory.SupportedTargets()

	expectedTargets := []CompileTarget{TargetCursor, TargetAmazonQ}

	if len(targets) != len(expectedTargets) {
		t.Errorf("Expected %d targets, got %d", len(expectedTargets), len(targets))
	}

	targetMap := make(map[CompileTarget]bool)
	for _, target := range targets {
		targetMap[target] = true
	}

	for _, expected := range expectedTargets {
		if !targetMap[expected] {
			t.Errorf("Missing expected target: %s", expected)
		}
	}
}
