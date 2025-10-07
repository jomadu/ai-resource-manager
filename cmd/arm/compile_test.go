package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCompileCmd(t *testing.T) {
	cmd := newCompileCmd()

	assert.Equal(t, "compile [file...]", cmd.Use)
	assert.Equal(t, "Compile resource files to target format", cmd.Short)
	assert.True(t, cmd.HasFlags())

	// Check that target flag exists
	targetFlag := cmd.Flag("target")
	assert.NotNil(t, targetFlag)

	// Check optional flags exist
	flags := []string{"output", "namespace", "force", "recursive", "verbose", "validate-only", "include", "exclude", "fail-fast"}
	for _, flagName := range flags {
		flag := cmd.Flag(flagName)
		assert.NotNil(t, flag, "Flag %s should exist", flagName)
	}
}

func TestCompileCmdArgs(t *testing.T) {
	cmd := newCompileCmd()

	// Should require at least one argument
	err := cmd.Args(cmd, []string{})
	assert.Error(t, err)

	// Should accept one or more arguments
	err = cmd.Args(cmd, []string{"file1.yaml"})
	assert.NoError(t, err)

	err = cmd.Args(cmd, []string{"file1.yaml", "file2.yml"})
	assert.NoError(t, err)
}

func TestCompileCmdFlagDefaults(t *testing.T) {
	cmd := newCompileCmd()

	// Set required flag to avoid validation error
	cmd.SetArgs([]string{"--target", "cursor", "test.yaml"})

	// Parse flags
	err := cmd.ParseFlags([]string{"--target", "cursor"})
	assert.NoError(t, err)

	// Check defaults
	output, _ := cmd.Flags().GetString("output")
	assert.Equal(t, ".", output)

	force, _ := cmd.Flags().GetBool("force")
	assert.False(t, force)

	recursive, _ := cmd.Flags().GetBool("recursive")
	assert.False(t, recursive)

	verbose, _ := cmd.Flags().GetBool("verbose")
	assert.False(t, verbose)

	validateOnly, _ := cmd.Flags().GetBool("validate-only")
	assert.False(t, validateOnly)

	failFast, _ := cmd.Flags().GetBool("fail-fast")
	assert.False(t, failFast)
}
