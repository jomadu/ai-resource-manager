package service

import (
	"context"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
)

func TestSetPromptsetVersion(t *testing.T) {
	t.Run("update version for existing promptset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "1.0.0",
						"sinks":   []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetVersion(context.Background(), "test-reg", "promptset1", "2.0.0")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Dependencies["test-reg/promptset1"]["version"] != "2.0.0" {
			t.Errorf("expected version 2.0.0, got %v", mgr.manifest.Dependencies["test-reg/promptset1"]["version"])
		}
	})

	t.Run("fail when promptset does not exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: make(map[string]map[string]interface{}),
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetVersion(context.Background(), "test-reg", "nonexistent", "2.0.0")

		if err == nil {
			t.Fatal("expected error when promptset does not exist")
		}
	})
}

func TestSetPromptsetInclude(t *testing.T) {
	t.Run("update include patterns for existing promptset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "1.0.0",
						"sinks":   []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetInclude(context.Background(), "test-reg", "promptset1", []string{"*.yml", "*.yaml"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestSetPromptsetExclude(t *testing.T) {
	t.Run("update exclude patterns for existing promptset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Dependencies: map[string]map[string]interface{}{
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "1.0.0",
						"sinks":   []interface{}{"test-sink"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetExclude(context.Background(), "test-reg", "promptset1", []string{"test/**"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})
}

func TestSetPromptsetSinks(t *testing.T) {
	t.Run("update sinks for existing promptset", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"sink1": {Directory: "/path1", Tool: compiler.Cursor},
					"sink2": {Directory: "/path2", Tool: compiler.AmazonQ},
				},
				Dependencies: map[string]map[string]interface{}{
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "1.0.0",
						"sinks":   []interface{}{"sink1"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetSinks(context.Background(), "test-reg", "promptset1", []string{"sink1", "sink2"})

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
					"test-reg/promptset1": {
						"type":    "promptset",
						"version": "1.0.0",
						"sinks":   []interface{}{"sink1"},
					},
				},
			},
		}

		svc := NewArmService(mgr, nil, nil)

		err := svc.SetPromptsetSinks(context.Background(), "test-reg", "promptset1", []string{"nonexistent"})

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})
}
