package registry

import (
	"testing"
)

func TestDefaultFactory_CreateGitLabRegistry(t *testing.T) {
	factory := &DefaultFactory{}

	config := map[string]interface{}{
		"type":      "gitlab",
		"url":       "https://gitlab.example.com",
		"projectId": "123",
	}

	registry, err := factory.CreateRegistry("test-gitlab", config)
	if err != nil {
		t.Fatalf("failed to create gitlab registry: %v", err)
	}

	if registry == nil {
		t.Error("expected registry, got nil")
	}
}

func TestDefaultFactory_CreateGitRegistry(t *testing.T) {
	factory := &DefaultFactory{}

	config := map[string]interface{}{
		"type": "git",
		"url":  "https://github.com/test/repo",
	}

	registry, err := factory.CreateRegistry("test-git", config)
	if err != nil {
		t.Fatalf("failed to create git registry: %v", err)
	}

	if registry == nil {
		t.Error("expected registry, got nil")
	}
}

func TestDefaultFactory_CreateCloudsmithRegistry(t *testing.T) {
	factory := &DefaultFactory{}

	config := map[string]interface{}{
		"type":       "cloudsmith",
		"url":        "https://api.cloudsmith.io",
		"owner":      "test-owner",
		"repository": "test-repo",
	}

	registry, err := factory.CreateRegistry("test-cloudsmith", config)
	if err != nil {
		t.Fatalf("failed to create cloudsmith registry: %v", err)
	}

	if registry == nil {
		t.Error("expected registry, got nil")
	}
}

func TestDefaultFactory_UnsupportedType(t *testing.T) {
	factory := &DefaultFactory{}

	config := map[string]interface{}{
		"type": "unsupported",
		"url":  "https://example.com",
	}

	_, err := factory.CreateRegistry("test", config)
	if err == nil {
		t.Error("expected error for unsupported type, got nil")
	}
}
