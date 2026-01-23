package service

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
)

type mockManifestManager struct {
	manifest *manifest.Manifest
	loadErr  error
	saveErr  error
}

func (m *mockManifestManager) GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.manifest.Registries, nil
}

func (m *mockManifestManager) GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	cfg, exists := m.manifest.Registries[name]
	if !exists {
		return nil, errors.New("registry does not exist")
	}
	return cfg, nil
}

func (m *mockManifestManager) GetGitRegistryConfig(ctx context.Context, name string) (manifest.GitRegistryConfig, error) {
	if m.loadErr != nil {
		return manifest.GitRegistryConfig{}, m.loadErr
	}
	cfg, exists := m.manifest.Registries[name]
	if !exists {
		return manifest.GitRegistryConfig{}, errors.New("registry does not exist")
	}
	regType, ok := cfg["type"].(string)
	if !ok || regType != "git" {
		return manifest.GitRegistryConfig{}, errors.New("registry is not a git registry")
	}
	configMap, _ := json.Marshal(cfg)
	var result manifest.GitRegistryConfig
	_ = json.Unmarshal(configMap, &result)
	return result, nil
}

func (m *mockManifestManager) GetGitLabRegistryConfig(ctx context.Context, name string) (manifest.GitLabRegistryConfig, error) {
	if m.loadErr != nil {
		return manifest.GitLabRegistryConfig{}, m.loadErr
	}
	cfg, exists := m.manifest.Registries[name]
	if !exists {
		return manifest.GitLabRegistryConfig{}, errors.New("registry does not exist")
	}
	regType, ok := cfg["type"].(string)
	if !ok || regType != "gitlab" {
		return manifest.GitLabRegistryConfig{}, errors.New("registry is not a gitlab registry")
	}
	configMap, _ := json.Marshal(cfg)
	var result manifest.GitLabRegistryConfig
	_ = json.Unmarshal(configMap, &result)
	return result, nil
}

func (m *mockManifestManager) GetCloudsmithRegistryConfig(ctx context.Context, name string) (manifest.CloudsmithRegistryConfig, error) {
	if m.loadErr != nil {
		return manifest.CloudsmithRegistryConfig{}, m.loadErr
	}
	cfg, exists := m.manifest.Registries[name]
	if !exists {
		return manifest.CloudsmithRegistryConfig{}, errors.New("registry does not exist")
	}
	regType, ok := cfg["type"].(string)
	if !ok || regType != "cloudsmith" {
		return manifest.CloudsmithRegistryConfig{}, errors.New("registry is not a cloudsmith registry")
	}
	configMap, _ := json.Marshal(cfg)
	var result manifest.CloudsmithRegistryConfig
	_ = json.Unmarshal(configMap, &result)
	return result, nil
}

func (m *mockManifestManager) UpsertRegistryConfig(ctx context.Context, name string, config map[string]interface{}) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.manifest.Registries[name] = config
	return nil
}

func (m *mockManifestManager) UpsertGitRegistryConfig(ctx context.Context, name string, config manifest.GitRegistryConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	config.Type = "git"
	configMap, _ := json.Marshal(config)
	var result map[string]interface{}
	_ = json.Unmarshal(configMap, &result)
	m.manifest.Registries[name] = result
	return nil
}

func (m *mockManifestManager) UpsertGitLabRegistryConfig(ctx context.Context, name string, config *manifest.GitLabRegistryConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	config.Type = "gitlab"
	configMap, _ := json.Marshal(config)
	var result map[string]interface{}
	_ = json.Unmarshal(configMap, &result)
	m.manifest.Registries[name] = result
	return nil
}

func (m *mockManifestManager) UpsertCloudsmithRegistryConfig(ctx context.Context, name string, config manifest.CloudsmithRegistryConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	config.Type = "cloudsmith"
	configMap, _ := json.Marshal(config)
	var result map[string]interface{}
	_ = json.Unmarshal(configMap, &result)
	m.manifest.Registries[name] = result
	return nil
}

func (m *mockManifestManager) UpdateRegistryConfigName(ctx context.Context, name, newName string) error {
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	cfg, exists := m.manifest.Registries[name]
	if !exists {
		return errors.New("registry does not exist")
	}
	if _, exists := m.manifest.Registries[newName]; exists {
		return errors.New("registry with new name already exists")
	}
	m.manifest.Registries[newName] = cfg
	delete(m.manifest.Registries, name)
	return nil
}

func (m *mockManifestManager) RemoveRegistryConfig(ctx context.Context, name string) error {
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	if _, exists := m.manifest.Registries[name]; !exists {
		return errors.New("registry does not exist")
	}
	delete(m.manifest.Registries, name)
	return nil
}

func (m *mockManifestManager) GetAllSinksConfig(ctx context.Context) (map[string]manifest.SinkConfig, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.manifest.Sinks, nil
}

func (m *mockManifestManager) GetSinkConfig(ctx context.Context, name string) (manifest.SinkConfig, error) {
	if m.loadErr != nil {
		return manifest.SinkConfig{}, m.loadErr
	}
	cfg, exists := m.manifest.Sinks[name]
	if !exists {
		return manifest.SinkConfig{}, errors.New("sink does not exist")
	}
	return cfg, nil
}

func (m *mockManifestManager) UpsertSinkConfig(ctx context.Context, name string, config manifest.SinkConfig) error {
	if m.saveErr != nil {
		return m.saveErr
	}
	m.manifest.Sinks[name] = config
	return nil
}

func (m *mockManifestManager) UpdateSinkConfigName(ctx context.Context, name, newName string) error {
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	cfg, exists := m.manifest.Sinks[name]
	if !exists {
		return errors.New("sink does not exist")
	}
	if _, exists := m.manifest.Sinks[newName]; exists {
		return errors.New("sink with new name already exists")
	}
	m.manifest.Sinks[newName] = cfg
	delete(m.manifest.Sinks, name)
	return nil
}

func (m *mockManifestManager) RemoveSinkConfig(ctx context.Context, name string) error {
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	if _, exists := m.manifest.Sinks[name]; !exists {
		return errors.New("sink does not exist")
	}
	delete(m.manifest.Sinks, name)
	return nil
}

func (m *mockManifestManager) GetAllDependenciesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	return m.manifest.Dependencies, nil
}

func (m *mockManifestManager) GetDependencyConfig(ctx context.Context, registry, packageName string) (map[string]interface{}, error) {
	key := registry + "/" + packageName
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	cfg, exists := m.manifest.Dependencies[key]
	if !exists {
		return nil, errors.New("dependency does not exist")
	}
	return cfg, nil
}

func (m *mockManifestManager) UpsertDependencyConfig(ctx context.Context, registry, packageName string, config map[string]interface{}) error {
	key := registry + "/" + packageName
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.manifest.Dependencies == nil {
		m.manifest.Dependencies = make(map[string]map[string]interface{})
	}
	m.manifest.Dependencies[key] = config
	return nil
}

func (m *mockManifestManager) UpsertRulesetDependencyConfig(ctx context.Context, registry, packageName string, config *manifest.RulesetDependencyConfig) error {
	key := registry + "/" + packageName
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.manifest.Dependencies == nil {
		m.manifest.Dependencies = make(map[string]map[string]interface{})
	}
	config.Type = manifest.ResourceTypeRuleset
	configMap, _ := json.Marshal(config)
	var result map[string]interface{}
	_ = json.Unmarshal(configMap, &result)
	m.manifest.Dependencies[key] = result
	return nil
}

func (m *mockManifestManager) UpsertPromptsetDependencyConfig(ctx context.Context, registry, packageName string, config *manifest.PromptsetDependencyConfig) error {
	key := registry + "/" + packageName
	if m.saveErr != nil {
		return m.saveErr
	}
	if m.manifest.Dependencies == nil {
		m.manifest.Dependencies = make(map[string]map[string]interface{})
	}
	config.Type = manifest.ResourceTypePromptset
	configMap, _ := json.Marshal(config)
	var result map[string]interface{}
	_ = json.Unmarshal(configMap, &result)
	m.manifest.Dependencies[key] = result
	return nil
}

func (m *mockManifestManager) UpdateDependencyConfigName(ctx context.Context, registry, packageName, newRegistry, newPackageName string) error {
	key := registry + "/" + packageName
	newKey := newRegistry + "/" + newPackageName
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	cfg, exists := m.manifest.Dependencies[key]
	if !exists {
		return errors.New("dependency does not exist")
	}
	m.manifest.Dependencies[newKey] = cfg
	delete(m.manifest.Dependencies, key)
	return nil
}

func (m *mockManifestManager) RemoveDependencyConfig(ctx context.Context, registry, packageName string) error {
	key := registry + "/" + packageName
	if m.loadErr != nil {
		return m.loadErr
	}
	if m.saveErr != nil {
		return m.saveErr
	}
	if _, exists := m.manifest.Dependencies[key]; !exists {
		return errors.New("dependency does not exist")
	}
	delete(m.manifest.Dependencies, key)
	return nil
}

func (m *mockManifestManager) GetRulesetDependencyConfig(ctx context.Context, registry, packageName string) (*manifest.RulesetDependencyConfig, error) {
	key := registry + "/" + packageName
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	cfg, exists := m.manifest.Dependencies[key]
	if !exists {
		return nil, errors.New("dependency does not exist")
	}
	configMap, _ := json.Marshal(cfg)
	var result manifest.RulesetDependencyConfig
	_ = json.Unmarshal(configMap, &result)
	return &result, nil
}

func (m *mockManifestManager) GetPromptsetDependencyConfig(ctx context.Context, registry, packageName string) (*manifest.PromptsetDependencyConfig, error) {
	key := registry + "/" + packageName
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	cfg, exists := m.manifest.Dependencies[key]
	if !exists {
		return nil, errors.New("dependency does not exist")
	}
	configMap, _ := json.Marshal(cfg)
	var result manifest.PromptsetDependencyConfig
	_ = json.Unmarshal(configMap, &result)
	return &result, nil
}

func (m *mockManifestManager) GetAllRulesetDependenciesConfig(ctx context.Context) (map[string]*manifest.RulesetDependencyConfig, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	rulesets := make(map[string]*manifest.RulesetDependencyConfig)
	for key, rawConfig := range m.manifest.Dependencies {
		depType, ok := rawConfig["type"].(string)
		if !ok || depType != "ruleset" {
			continue
		}
		configMap, _ := json.Marshal(rawConfig)
		var result manifest.RulesetDependencyConfig
		_ = json.Unmarshal(configMap, &result)
		rulesets[key] = &result
	}
	return rulesets, nil
}

func (m *mockManifestManager) GetAllPromptsetDependenciesConfig(ctx context.Context) (map[string]*manifest.PromptsetDependencyConfig, error) {
	if m.loadErr != nil {
		return nil, m.loadErr
	}
	promptsets := make(map[string]*manifest.PromptsetDependencyConfig)
	for key, rawConfig := range m.manifest.Dependencies {
		depType, ok := rawConfig["type"].(string)
		if !ok || depType != "promptset" {
			continue
		}
		configMap, _ := json.Marshal(rawConfig)
		var result manifest.PromptsetDependencyConfig
		_ = json.Unmarshal(configMap, &result)
		promptsets[key] = &result
	}
	return promptsets, nil
}

func TestAddSink(t *testing.T) {
	t.Run("add new sink", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Sinks: make(map[string]manifest.SinkConfig)},
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		err := svc.RemoveSink(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		_, err := svc.GetSinkConfig(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		err := svc.SetSinkName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when new name already exists")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		err := svc.SetSinkDirectory(context.Background(), "test", "/new")

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err == nil {
			t.Fatal("expected error when sink does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil, nil)

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
		svc := NewArmService(mgr, nil, nil)

		err := svc.SetSinkTool(context.Background(), "test", compiler.AmazonQ)

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}
