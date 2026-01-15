package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
)

type mockManifestManager struct {
	manifest *manifest.Manifest
	loadErr  error
	saveErr  error
}

func (m *mockManifestManager) Load() (*manifest.Manifest, error) {
	return m.manifest, m.loadErr
}

func (m *mockManifestManager) Save(mf *manifest.Manifest) error {
	m.manifest = mf
	return m.saveErr
}

func TestAddSink(t *testing.T) {
	t.Run("add new sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		err := svc.AddSink(context.Background(), "test", "/path", compiler.Cursor, false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		sink, exists := mgr.manifest.Sinks["test"]
		if !exists {
			t.Fatal("sink not added")
		}
		if sink.Directory != "/path" {
			t.Errorf("expected directory /path, got %s", sink.Directory)
		}
		if sink.Tool != compiler.Cursor {
			t.Errorf("expected tool cursor, got %v", sink.Tool)
		}
	})

	t.Run("add when sink exists without force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/old", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.AddSink(context.Background(), "test", "/new", compiler.AmazonQ, false)

		if err == nil {
			t.Fatal("expected error when sink exists")
		}
		if mgr.manifest.Sinks["test"].Directory != "/old" {
			t.Error("sink should not be modified")
		}
	})

	t.Run("add when sink exists with force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/old", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.AddSink(context.Background(), "test", "/new", compiler.AmazonQ, true)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Sinks["test"].Directory != "/new" {
			t.Errorf("expected directory /new, got %s", mgr.manifest.Sinks["test"].Directory)
		}
		if mgr.manifest.Sinks["test"].Tool != compiler.AmazonQ {
			t.Errorf("expected tool amazonq, got %v", mgr.manifest.Sinks["test"].Tool)
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		err := svc.AddSink(context.Background(), "test", "/path", compiler.Cursor, false)

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
			saveErr:  errors.New("save error"),
		}
		svc := NewArmService(mgr)

		err := svc.AddSink(context.Background(), "test", "/path", compiler.Cursor, false)

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestRemoveSink(t *testing.T) {
	t.Run("remove existing sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.RemoveSink(context.Background(), "test")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, exists := mgr.manifest.Sinks["test"]; exists {
			t.Error("sink should be removed")
		}
	})

	t.Run("remove non-existent sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		err := svc.RemoveSink(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		err := svc.RemoveSink(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
			saveErr: errors.New("save error"),
		}
		svc := NewArmService(mgr)

		err := svc.RemoveSink(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestGetSinkConfig(t *testing.T) {
	t.Run("get existing sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		cfg, err := svc.GetSinkConfig(context.Background(), "test")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg.Directory != "/path" {
			t.Errorf("expected directory /path, got %s", cfg.Directory)
		}
	})

	t.Run("get non-existent sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		_, err := svc.GetSinkConfig(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		_, err := svc.GetSinkConfig(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})
}

func TestGetAllSinkConfigs(t *testing.T) {
	t.Run("get all when sinks exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test1": {Directory: "/path1", Tool: compiler.Cursor},
					"test2": {Directory: "/path2", Tool: compiler.AmazonQ},
				},
			},
		}
		svc := NewArmService(mgr)

		cfgs, err := svc.GetAllSinkConfigs(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(cfgs) != 2 {
			t.Errorf("expected 2 sinks, got %d", len(cfgs))
		}
	})

	t.Run("get all when no sinks", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		cfgs, err := svc.GetAllSinkConfigs(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(cfgs) != 0 {
			t.Errorf("expected 0 sinks, got %d", len(cfgs))
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		_, err := svc.GetAllSinkConfigs(context.Background())

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})
}

func TestSetSinkName(t *testing.T) {
	t.Run("rename existing sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"old": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, exists := mgr.manifest.Sinks["old"]; exists {
			t.Error("old sink should be removed")
		}
		if _, exists := mgr.manifest.Sinks["new"]; !exists {
			t.Fatal("new sink should exist")
		}
		if mgr.manifest.Sinks["new"].Directory != "/path" {
			t.Error("sink config should be preserved")
		}
	})

	t.Run("rename non-existent sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("rename to existing name", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"old": {Directory: "/path1", Tool: compiler.Cursor},
					"new": {Directory: "/path2", Tool: compiler.AmazonQ},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when new name already exists")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"old": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
			saveErr: errors.New("save error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestSetSinkDirectory(t *testing.T) {
	t.Run("update existing sink directory", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/old", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkDirectory(context.Background(), "test", "/new")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Sinks["test"].Directory != "/new" {
			t.Errorf("expected directory /new, got %s", mgr.manifest.Sinks["test"].Directory)
		}
	})

	t.Run("update non-existent sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkDirectory(context.Background(), "test", "/new")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkDirectory(context.Background(), "test", "/new")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/old", Tool: compiler.Cursor},
				},
			},
			saveErr: errors.New("save error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkDirectory(context.Background(), "test", "/new")

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestSetSinkTool(t *testing.T) {
	t.Run("update existing sink tool", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Sinks["test"].Tool != compiler.AmazonQ {
			t.Errorf("expected tool amazonq, got %v", mgr.manifest.Sinks["test"].Tool)
		}
	})

	t.Run("update non-existent sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Sinks: map[string]manifest.SinkConfig{
					"test": {Directory: "/path", Tool: compiler.Cursor},
				},
			},
			saveErr: errors.New("save error"),
		}
		svc := NewArmService(mgr)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}
