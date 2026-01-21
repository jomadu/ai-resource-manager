package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/v4/config"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
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

type cloudsmithPackage struct {
	Name     string `json:"name"`
	Version  string `json:"version"`
	Format   string `json:"format"`
	Filename string `json:"filename"`
}

func (c *CloudsmithRegistry) ListVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	if err := c.loadToken(ctx); err != nil {
		return nil, err
	}

	packages, err := c.client.listPackages(ctx, c.config.Owner, c.config.Repository, packageName)
	if err != nil {
		return nil, err
	}

	versionMap := make(map[string]bool)
	for _, pkg := range packages {
		if pkg.Format == "raw" && (pkg.Name == packageName || strings.HasPrefix(pkg.Filename, packageName)) {
			versionMap[pkg.Version] = true
		}
	}

	var versions []core.Version
	for versionStr := range versionMap {
		version, _ := core.NewVersion(versionStr)
		versions = append(versions, version)
	}

	sort.Slice(versions, func(i, j int) bool {
		if !versions[i].IsSemver || !versions[j].IsSemver {
			return versions[i].Version < versions[j].Version
		}
		return versions[i].Compare(versions[j]) > 0
	})

	return versions, nil
}

func (c *cloudsmithClient) listPackages(ctx context.Context, owner, repo, packageName string) ([]cloudsmithPackage, error) {
	path := fmt.Sprintf("/v1/packages/%s/%s/?query=%s", owner, repo, packageName)
	
	var allPackages []cloudsmithPackage
	for path != "" {
		packages, nextPath, err := c.listPackagesPage(ctx, path)
		if err != nil {
			return nil, err
		}
		allPackages = append(allPackages, packages...)
		path = nextPath
	}

	return allPackages, nil
}

func (c *cloudsmithClient) listPackagesPage(ctx context.Context, path string) ([]cloudsmithPackage, string, error) {
	resp, err := c.makeRequest(ctx, "GET", path)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, "", fmt.Errorf("cloudsmith API error: %d", resp.StatusCode)
	}

	var packages []cloudsmithPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, "", err
	}

	nextPath := parseNextURLFromLinkHeader(resp.Header.Get("Link"))
	return packages, nextPath, nil
}

func (c *CloudsmithRegistry) ResolveVersion(ctx context.Context, packageName, constraint string) (core.Version, error) {
	versions, err := c.ListVersions(ctx, packageName)
	if err != nil {
		return core.Version{}, fmt.Errorf("failed to list versions: %w", err)
	}

	resolved, err := core.ResolveVersion(constraint, versions)
	if err != nil {
		return core.Version{}, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return resolved, nil
}

func parseNextURLFromLinkHeader(linkHeader string) string {
	if linkHeader == "" {
		return ""
	}

	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		link = strings.TrimSpace(link)
		if strings.Contains(link, `rel="next"`) {
			start := strings.Index(link, "<")
			end := strings.Index(link, ">")
			if start != -1 && end != -1 && start < end {
				fullURL := link[start+1 : end]
				if idx := strings.Index(fullURL, "/v1/"); idx != -1 {
					return fullURL[idx:]
				}
			}
		}
	}

	return ""
}
