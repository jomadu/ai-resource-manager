package registry

import (
	"net/http"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/storage"
)

type CloudsmithRegistry struct {
	name         string
	config       CloudsmithRegistryConfig
	client       *cloudsmithClient
	packageCache *storage.PackageCache
}

type cloudsmithClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

func NewCloudsmithRegistry(name string, cfg CloudsmithRegistryConfig) (*CloudsmithRegistry, error) {
	registry, err := storage.NewRegistry(cfg)
	if err != nil {
		return nil, err
	}

	baseURL := cfg.URL
	if baseURL == "" {
		baseURL = "https://api.cloudsmith.io"
	}

	return &CloudsmithRegistry{
		name:   name,
		config: cfg,
		client: &cloudsmithClient{
			baseURL:    baseURL,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		},
		packageCache: storage.NewPackageCache(registry.GetPackagesDir()),
	}, nil
}
