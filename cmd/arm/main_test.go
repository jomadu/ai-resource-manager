package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/core"
)

func TestPrintVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printVersion()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "arm") {
		t.Errorf("Expected output to contain 'arm', got: %s", output)
	}
	if !strings.Contains(output, "build-id:") {
		t.Errorf("Expected output to contain 'build-id:', got: %s", output)
	}
	if !strings.Contains(output, "build-timestamp:") {
		t.Errorf("Expected output to contain 'build-timestamp:', got: %s", output)
	}
	if !strings.Contains(output, "build-platform:") {
		t.Errorf("Expected output to contain 'build-platform:', got: %s", output)
	}
}

func TestPrintHelp(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printHelp()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	if !strings.Contains(output, "AI Resource Manager") {
		t.Errorf("Expected output to contain 'AI Resource Manager', got: %s", output)
	}
	if !strings.Contains(output, "version") {
		t.Errorf("Expected output to contain 'version', got: %s", output)
	}
	if !strings.Contains(output, "help") {
		t.Errorf("Expected output to contain 'help', got: %s", output)
	}
}

func TestPrintCommandHelp(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected []string
	}{
		{
			name:    "version help",
			command: "version",
			expected: []string{
				"Display version information",
				"arm version",
				"Build ID",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w

			printCommandHelp(tt.command)

			w.Close()
			os.Stdout = old

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			for _, exp := range tt.expected {
				if !strings.Contains(output, exp) {
					t.Errorf("Expected output to contain '%s', got: %s", exp, output)
				}
			}
		})
	}
}

func TestBuildInfoIntegration(t *testing.T) {
	info := core.GetBuildInfo()
	
	if info.Version.Version == "" {
		t.Error("Expected version to be set")
	}
	if info.Arch == "" {
		t.Error("Expected arch to be set")
	}
}
