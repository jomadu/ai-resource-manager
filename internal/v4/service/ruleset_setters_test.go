package service

import (
	"context"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
)

func TestSetRulesetVersion(t *testing.T) {
	t.Run("update version for existing ruleset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetVersion(context.Background(), "test-reg", "ruleset1", "2.0.0")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Dependencies["test-reg/ruleset1"]["version"] != "2.0.0" {
			t.Errorf("expected version 2.0.0, got %v", mgr.manifest.Dependencies["test-reg/ruleset1"]["version"])
		}
	})

	t.Run("fail when ruleset does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetVersion(context.Background(), "test-reg", "nonexistent", "2.0.0")

		if err == nil {
			t.Fatal("expected error when ruleset does not exist")
		}
	})
}

func TestSetRulesetPriority(t *testing.T) {
	t.Run("update priority for existing ruleset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetPriority(context.Background(), "test-reg", "ruleset1", 200)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		priority := mgr.manifest.Dependencies["test-reg/ruleset1"]["priority"]
		if priority != 200 && priority != float64(200) {
			t.Errorf("expected priority 200, got %v (type %T)", priority, priority)
		}
	})

	t.Run("fail when ruleset does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetPriority(context.Background(), "test-reg", "nonexistent", 200)

		if err == nil {
			t.Fatal("expected error when ruleset does not exist")
		}
	})
}

func TestSetRulesetInclude(t *testing.T) {
	t.Run("update include patterns for existing ruleset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetInclude(context.Background(), "test-reg", "ruleset1", []string{"*.yml", "*.yaml"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestSetRulesetExclude(t *testing.T) {
	t.Run("update exclude patterns for existing ruleset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetExclude(context.Background(), "test-reg", "ruleset1", []string{"test/**"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestSetRulesetSinks(t *testing.T) {
	t.Run("update sinks for existing ruleset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"sink1": {Directory: "/path1", Tool: compiler.Cursor},
					"sink2": {Directory: "/path2", Tool: compiler.AmazonQ},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"sink1"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetSinks(context.Background(), "test-reg", "ruleset1", []string{"sink1", "sink2"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("fail when sink does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"sink1": {Directory: "/path1", Tool: compiler.Cursor},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/ruleset1": {
						"type":     "ruleset",
						"version":  "1.0.0",
						"priority": float64(100),
						"sinks":    []interface{}{"sink1"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetRulesetSinks(context.Background(), "test-reg", "ruleset1", []string{"nonexistent"})

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})
}
