package service

import (
	"context"
	"errors"
	"testing"

	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
)

func TestAddGitRegistry(t *testing.T) {
	t.Run("add new git registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://github.com/test/repo", []string{"main"}, false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		reg, exists := mgr.manifest.Registries["test"]
		if !exists {
			t.Fatal("registry not added")
		}
		if reg["url"] != "https://github.com/test/repo" {
			t.Errorf("expected url https://github.com/test/repo, got %v", reg["url"])
		}
		if reg["type"] != "git" {
			t.Errorf("expected type git, got %v", reg["type"])
		}
	})

	t.Run("add git registry with no branches", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://github.com/test/repo", nil, false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		reg := mgr.manifest.Registries["test"]
		if reg["branches"] != nil {
			t.Error("expected no branches")
		}
	})

	t.Run("add when registry exists without force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://new.com", nil, false)

		if err == nil {
			t.Fatal("expected error when registry exists")
		}
		if mgr.manifest.Registries["test"]["url"] != "https://old.com" {
			t.Error("registry should not be modified")
		}
	})

	t.Run("add when registry exists with force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://new.com", []string{"dev"}, true)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["url"] != "https://new.com" {
			t.Errorf("expected url https://new.com, got %v", mgr.manifest.Registries["test"]["url"])
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://github.com/test/repo", nil, false)

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
			saveErr:  errors.New("save error"),
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitRegistry(context.Background(), "test", "https://github.com/test/repo", nil, false)

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestAddGitLabRegistry(t *testing.T) {
	t.Run("add new gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitLabRegistry(context.Background(), "test", "https://gitlab.com", "123", "456", "v4", false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		reg := mgr.manifest.Registries["test"]
		if reg["type"] != "gitlab" {
			t.Errorf("expected type gitlab, got %v", reg["type"])
		}
		if reg["projectId"] != "123" {
			t.Errorf("expected projectId 123, got %v", reg["projectId"])
		}
	})

	t.Run("add gitlab registry with empty optional fields", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitLabRegistry(context.Background(), "test", "https://gitlab.com", "", "", "", false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("add when registry exists without force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitLabRegistry(context.Background(), "test", "https://new.com", "123", "", "", false)

		if err == nil {
			t.Fatal("expected error when registry exists")
		}
	})

	t.Run("add when registry exists with force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddGitLabRegistry(context.Background(), "test", "https://new.com", "999", "", "", true)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["projectId"] != "999" {
			t.Error("registry should be updated")
		}
	})
}

func TestAddCloudsmithRegistry(t *testing.T) {
	t.Run("add new cloudsmith registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddCloudsmithRegistry(context.Background(), "test", "https://cloudsmith.io", "myorg", "myrepo", false)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		reg := mgr.manifest.Registries["test"]
		if reg["type"] != "cloudsmith" {
			t.Errorf("expected type cloudsmith, got %v", reg["type"])
		}
		if reg["owner"] != "myorg" {
			t.Errorf("expected owner myorg, got %v", reg["owner"])
		}
		if reg["repository"] != "myrepo" {
			t.Errorf("expected repository myrepo, got %v", reg["repository"])
		}
	})

	t.Run("add when registry exists without force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "cloudsmith"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddCloudsmithRegistry(context.Background(), "test", "https://new.com", "org", "repo", false)

		if err == nil {
			t.Fatal("expected error when registry exists")
		}
	})

	t.Run("add when registry exists with force", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "cloudsmith", "owner": "oldorg"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.AddCloudsmithRegistry(context.Background(), "test", "https://new.com", "neworg", "newrepo", true)

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["owner"] != "neworg" {
			t.Error("registry should be updated")
		}
	})
}

func TestRemoveRegistry(t *testing.T) {
	t.Run("remove existing registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.RemoveRegistry(context.Background(), "test")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, exists := mgr.manifest.Registries["test"]; exists {
			t.Error("registry should be removed")
		}
	})

	t.Run("remove non-existent registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.RemoveRegistry(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil)

		err := svc.RemoveRegistry(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})

	t.Run("save fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
			saveErr: errors.New("save error"),
		}
		svc := NewArmService(mgr, nil)

		err := svc.RemoveRegistry(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when save fails")
		}
	})
}

func TestGetRegistryConfig(t *testing.T) {
	t.Run("get existing registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		cfg, err := svc.GetRegistryConfig(context.Background(), "test")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if cfg["url"] != "https://test.com" {
			t.Errorf("expected url https://test.com, got %v", cfg["url"])
		}
	})

	t.Run("get non-existent registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		_, err := svc.GetRegistryConfig(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil)

		_, err := svc.GetRegistryConfig(context.Background(), "test")

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})
}

func TestGetAllRegistriesConfig(t *testing.T) {
	t.Run("get all when registries exist", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test1": {"url": "https://test1.com", "type": "git"},
					"test2": {"url": "https://test2.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		cfgs, err := svc.GetAllRegistriesConfig(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(cfgs) != 2 {
			t.Errorf("expected 2 registries, got %d", len(cfgs))
		}
	})

	t.Run("get all when no registries", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		cfgs, err := svc.GetAllRegistriesConfig(context.Background())

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(cfgs) != 0 {
			t.Errorf("expected 0 registries, got %d", len(cfgs))
		}
	})

	t.Run("load fails", func(t *testing.T) {
		mgr := &mockManifestManager{
			loadErr: errors.New("load error"),
		}
		svc := NewArmService(mgr, nil)

		_, err := svc.GetAllRegistriesConfig(context.Background())

		if err == nil {
			t.Fatal("expected error when load fails")
		}
	})
}

func TestSetRegistryName(t *testing.T) {
	t.Run("rename existing registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"old": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetRegistryName(context.Background(), "old", "new")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if _, exists := mgr.manifest.Registries["old"]; exists {
			t.Error("old registry should be removed")
		}
		if _, exists := mgr.manifest.Registries["new"]; !exists {
			t.Fatal("new registry should exist")
		}
		if mgr.manifest.Registries["new"]["url"] != "https://test.com" {
			t.Error("registry config should be preserved")
		}
	})

	t.Run("rename non-existent registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetRegistryName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})

	t.Run("rename to existing name", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"old": {"url": "https://old.com", "type": "git"},
					"new": {"url": "https://new.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetRegistryName(context.Background(), "old", "new")

		if err == nil {
			t.Fatal("expected error when new name already exists")
		}
	})
}

func TestSetRegistryURL(t *testing.T) {
	t.Run("set url for existing registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://old.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetRegistryURL(context.Background(), "test", "https://new.com")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["url"] != "https://new.com" {
			t.Errorf("expected url https://new.com, got %v", mgr.manifest.Registries["test"]["url"])
		}
	})

	t.Run("set url for non-existent registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetRegistryURL(context.Background(), "test", "https://new.com")

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})
}

func TestSetGitRegistryBranches(t *testing.T) {
	t.Run("set branches for git registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitRegistryBranches(context.Background(), "test", []string{"main", "dev"})

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		branches := mgr.manifest.Registries["test"]["branches"]
		if branches == nil {
			t.Fatal("branches should be set")
		}
	})

	t.Run("set branches for non-git registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitRegistryBranches(context.Background(), "test", []string{"main"})

		if err == nil {
			t.Fatal("expected error when registry is not git type")
		}
	})

	t.Run("set branches for non-existent registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{Registries: make(map[string]map[string]interface{})},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitRegistryBranches(context.Background(), "test", []string{"main"})

		if err == nil {
			t.Fatal("expected error when registry does not exist")
		}
	})
}

func TestSetGitLabRegistryProjectID(t *testing.T) {
	t.Run("set project id for gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://gitlab.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryProjectID(context.Background(), "test", "999")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["projectId"] != "999" {
			t.Error("projectId should be updated")
		}
	})

	t.Run("set project id for non-gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryProjectID(context.Background(), "test", "999")

		if err == nil {
			t.Fatal("expected error when registry is not gitlab type")
		}
	})
}

func TestSetGitLabRegistryGroupID(t *testing.T) {
	t.Run("set group id for gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://gitlab.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryGroupID(context.Background(), "test", "888")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["groupId"] != "888" {
			t.Error("groupId should be updated")
		}
	})

	t.Run("set group id for non-gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryGroupID(context.Background(), "test", "888")

		if err == nil {
			t.Fatal("expected error when registry is not gitlab type")
		}
	})
}

func TestSetGitLabRegistryAPIVersion(t *testing.T) {
	t.Run("set api version for gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://gitlab.com", "type": "gitlab"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryAPIVersion(context.Background(), "test", "v5")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["apiVersion"] != "v5" {
			t.Error("apiVersion should be updated")
		}
	})

	t.Run("set api version for non-gitlab registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetGitLabRegistryAPIVersion(context.Background(), "test", "v5")

		if err == nil {
			t.Fatal("expected error when registry is not gitlab type")
		}
	})
}

func TestSetCloudsmithRegistryOwner(t *testing.T) {
	t.Run("set owner for cloudsmith registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://cloudsmith.io", "type": "cloudsmith", "owner": "oldorg"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetCloudsmithRegistryOwner(context.Background(), "test", "neworg")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["owner"] != "neworg" {
			t.Error("owner should be updated")
		}
	})

	t.Run("set owner for non-cloudsmith registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetCloudsmithRegistryOwner(context.Background(), "test", "neworg")

		if err == nil {
			t.Fatal("expected error when registry is not cloudsmith type")
		}
	})
}

func TestSetCloudsmithRegistryRepository(t *testing.T) {
	t.Run("set repository for cloudsmith registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://cloudsmith.io", "type": "cloudsmith", "repository": "oldrepo"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetCloudsmithRegistryRepository(context.Background(), "test", "newrepo")

		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if mgr.manifest.Registries["test"]["repository"] != "newrepo" {
			t.Error("repository should be updated")
		}
	})

	t.Run("set repository for non-cloudsmith registry", func(t *testing.T) {
		mgr := &mockManifestManager{
			manifest: &manifest.Manifest{
				Registries: map[string]map[string]interface{}{
					"test": {"url": "https://test.com", "type": "git"},
				},
			},
		}
		svc := NewArmService(mgr, nil)

		err := svc.SetCloudsmithRegistryRepository(context.Background(), "test", "newrepo")

		if err == nil {
			t.Fatal("expected error when registry is not cloudsmith type")
		}
	})
}
