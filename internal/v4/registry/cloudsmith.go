package registry

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/config"
	"github.com/jomadu/ai-resource-manager/internal/v4/storage"
)

type CloudsmithRegistry struct {
	name         string
	config       CloudsmithRegistryConfig
	configMgr    config.Manager
	client       *cloudsmithClient
	packageCache *storage.PackageCache
}

type cloudsmithClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewCloudsmithRegistry(name string, cfg CloudsmithRegistryConfig, configMgr config.Manager) (*CloudsmithRegistry, error) {
	registry, err := storage.NewRegistry(cfg)
	if err != nil {
		return nil, err
	}

	baseURL := cfg.URL
	if baseURL == "" {
		baseURL = "https://api.cloudsmith.io"
	}

	return &CloudsmithRegistry{
		name:      name,
		config:    cfg,
		configMgr: configMgr,
		client: &cloudsmithClient{
			baseURL:    baseURL,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		},
		packageCache: storage.NewPackageCache(registry.GetPackagesDir()),
	}, nil
}

func (c *CloudsmithRegistry) loadToken(ctx context.Context) error {
	if c.client.token != "" {
		return nil
	}

	if c.configMgr == nil {
		return fmt.Errorf("no token configured for Cloudsmith registry")
	}

	authKey := fmt.Sprintf("%s/%s/%s", c.config.URL, c.config.Owner, c.config.Repository)
	token, err := c.configMgr.GetValue(ctx, "registry "+authKey, "token")
	if err != nil {
		return fmt.Errorf("failed to load token from .armrc: %w", err)
	}

	c.client.token = token
	return nil
}

func (c *cloudsmithClient) makeRequest(ctx context.Context, method, path string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, nil)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Token "+c.token)
	}

	return c.httpClient.Do(req)
}
