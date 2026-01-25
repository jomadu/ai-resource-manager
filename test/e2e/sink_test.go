package e2e

import (
	"path/filepath"
	"testing"

	"github.com/jomadu/ai-resource-manager/test/e2e/helpers"
)

func TestSinkManagement(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Test: Add sink for each tool type
	t.Run("AddCursorSink", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/rules")

		// Verify arm.json was created
		armJSON := filepath.Join(workDir, "arm.json")
		helpers.AssertFileExists(t, armJSON)

		// Verify sink is in the manifest
		manifest := helpers.ReadJSON(t, armJSON)
		sinks, ok := manifest["sinks"].(map[string]interface{})
		if !ok {
			t.Fatal("sinks field not found or invalid type")
		}

		sink, ok := sinks["cursor-rules"].(map[string]interface{})
		if !ok {
			t.Fatal("cursor-rules sink not found")
		}

		if sink["tool"] != "cursor" {
			t.Errorf("expected tool 'cursor', got %v", sink["tool"])
		}
		if sink["directory"] != ".cursor/rules" {
			t.Errorf("expected directory '.cursor/rules', got %v", sink["directory"])
		}
	})

	t.Run("AddAmazonQSink", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "amazonq", "q-rules", ".amazonq/rules")

		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})

		sink, ok := sinks["q-rules"].(map[string]interface{})
		if !ok {
			t.Fatal("q-rules sink not found")
		}

		if sink["tool"] != "amazonq" {
			t.Errorf("expected tool 'amazonq', got %v", sink["tool"])
		}
	})

	t.Run("AddCopilotSink", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "copilot", "copilot-rules", ".github/copilot")

		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})

		sink, ok := sinks["copilot-rules"].(map[string]interface{})
		if !ok {
			t.Fatal("copilot-rules sink not found")
		}

		if sink["tool"] != "copilot" {
			t.Errorf("expected tool 'copilot', got %v", sink["tool"])
		}
	})

	t.Run("AddMarkdownSink", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "markdown", "md-rules", "docs/rules")

		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})

		sink, ok := sinks["md-rules"].(map[string]interface{})
		if !ok {
			t.Fatal("md-rules sink not found")
		}

		if sink["tool"] != "markdown" {
			t.Errorf("expected tool 'markdown', got %v", sink["tool"])
		}
	})

	// Test: List sinks
	t.Run("ListSinks", func(t *testing.T) {
		output := arm.MustRun("list", "sink")

		// Should contain all sink names
		expectedSinks := []string{"copilot-rules", "cursor-rules", "md-rules", "q-rules"}
		for _, sink := range expectedSinks {
			if !contains(output, sink) {
				t.Errorf("list sink output should contain '%s', got: %s", sink, output)
			}
		}
	})

	// Test: Info for specific sink
	t.Run("InfoSink", func(t *testing.T) {
		output := arm.MustRun("info", "sink", "cursor-rules")

		// Should contain sink details
		if !contains(output, "cursor-rules") {
			t.Errorf("info sink output should contain 'cursor-rules', got: %s", output)
		}
		if !contains(output, "cursor") {
			t.Errorf("info sink output should contain tool 'cursor', got: %s", output)
		}
	})

	// Test: Set sink configuration
	t.Run("SetSink", func(t *testing.T) {
		arm.MustRun("set", "sink", "cursor-rules", "directory", ".cursor/new-rules")

		// Verify the directory was updated
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})
		sink := sinks["cursor-rules"].(map[string]interface{})

		if sink["directory"] != ".cursor/new-rules" {
			t.Errorf("sink directory not updated: expected '.cursor/new-rules', got %v", sink["directory"])
		}
	})

	// Test: Remove sink
	t.Run("RemoveSink", func(t *testing.T) {
		arm.MustRun("remove", "sink", "md-rules")

		// Verify sink is removed from manifest
		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})

		if _, ok := sinks["md-rules"]; ok {
			t.Error("md-rules should be removed from sinks")
		}
	})

	// Test: Add duplicate sink should fail without --force
	t.Run("AddDuplicateSinkFails", func(t *testing.T) {
		stderr := arm.MustFail("add", "sink", "--tool", "cursor", "cursor-rules", ".cursor/dup")

		if !contains(stderr, "already exists") && !contains(stderr, "duplicate") {
			t.Errorf("expected error about duplicate sink, got: %s", stderr)
		}
	})
}

func TestSinkLayoutModes(t *testing.T) {
	workDir := t.TempDir()
	arm := helpers.NewARMRunner(t, workDir)

	// Test hierarchical layout (default)
	t.Run("HierarchicalLayout", func(t *testing.T) {
		arm.MustRun("add", "sink", "--tool", "cursor", "hierarchical-sink", ".cursor/hierarchical")

		armJSON := filepath.Join(workDir, "arm.json")
		manifest := helpers.ReadJSON(t, armJSON)
		sinks := manifest["sinks"].(map[string]interface{})
		sink := sinks["hierarchical-sink"].(map[string]interface{})

		// Default layout should be hierarchical or not specified
		layout, ok := sink["layout"]
		if ok && layout != "hierarchical" && layout != "" {
			t.Errorf("expected hierarchical layout, got %v", layout)
		}
	})

	// Note: Flat layout is not exposed via CLI, it's an internal implementation detail
	// The layout is determined by the sink manager based on the tool and configuration
}
