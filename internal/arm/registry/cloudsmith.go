package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/config"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
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

func (c *CloudsmithRegistry) ListPackages(ctx context.Context) ([]*core.PackageMetadata, error) {
	if err := c.loadToken(ctx); err != nil {
		return nil, err
	}

	packages, err := c.client.listAllPackages(ctx, c.config.Owner, c.config.Repository)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var result []*core.PackageMetadata
	for _, pkg := range packages {
		if pkg.Format == "raw" && !seen[pkg.Name] {
			seen[pkg.Name] = true
			result = append(result, &core.PackageMetadata{
				RegistryName: c.name,
				Name:         pkg.Name,
			})
		}
	}

	return result, nil
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
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, http.NoBody)
	if err != nil {
		return nil, err
	}

	if c.token != "" {
		req.Header.Set("Authorization", "Token "+c.token)
	}

	return c.httpClient.Do(req)
}

type cloudsmithPackage struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Format      string `json:"format"`
	Filename    string `json:"filename"`
	DownloadURL string `json:"cdn_url"`
	Size        int64  `json:"size"`
}

func (c *CloudsmithRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
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
		return versions[i].Compare(&versions[j]) > 0
	})

	return versions, nil
}

func (c *cloudsmithClient) listAllPackages(ctx context.Context, owner, repo string) ([]cloudsmithPackage, error) {
	path := fmt.Sprintf("/v1/packages/%s/%s/", owner, repo)

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
	defer func() { _ = resp.Body.Close() }()

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
	versions, err := c.ListPackageVersions(ctx, packageName)
	if err != nil {
		return core.Version{}, fmt.Errorf("failed to list versions: %w", err)
	}

	resolved, err := core.ResolveVersion(constraint, versions)
	if err != nil {
		return core.Version{}, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return resolved, nil
}

func (c *CloudsmithRegistry) GetPackage(ctx context.Context, packageName string, version *core.Version, include, exclude []string) (*core.Package, error) {
	cacheKey := map[string]interface{}{
		"registry": c.name,
		"package":  packageName,
		"include":  include,
		"exclude":  exclude,
	}

	// Check cache first
	files, err := c.packageCache.GetPackageVersion(ctx, cacheKey, version)
	if err == nil {
		integrity := calculateIntegrity(files)
		return &core.Package{
			Metadata: core.PackageMetadata{
				RegistryName: c.name,
				Name:         packageName,
				Version:      *version,
			},
			Files:     files,
			Integrity: integrity,
		}, nil
	}

	if err := c.loadToken(ctx); err != nil {
		return nil, err
	}

	// Download files from Cloudsmith
	rawFiles, err := c.client.downloadPackages(ctx, c.config.Owner, c.config.Repository, packageName, version.Version)
	if err != nil {
		return nil, err
	}

	// Extract archives and merge with loose files
	extractor := core.NewExtractor()
	files, err = extractor.ExtractAndMerge(rawFiles)
	if err != nil {
		return nil, err
	}

	// Apply include/exclude filtering
	var filteredFiles []*core.File
	for _, file := range files {
		if c.matchesPatterns(file.Path, include, exclude) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Cache the filtered result
	_ = c.packageCache.SetPackageVersion(ctx, cacheKey, version, filteredFiles)

	integrity := calculateIntegrity(filteredFiles)

	return &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: c.name,
			Name:         packageName,
			Version:      *version,
		},
		Files:     filteredFiles,
		Integrity: integrity,
	}, nil
}

func (c *cloudsmithClient) downloadPackages(ctx context.Context, owner, repo, packageName, version string) ([]*core.File, error) {
	packages, err := c.listPackages(ctx, owner, repo, packageName)
	if err != nil {
		return nil, err
	}

	var files []*core.File
	for _, pkg := range packages {
		if pkg.Version == version && (pkg.Name == packageName || strings.HasPrefix(pkg.Filename, packageName)) {
			content, err := c.downloadFile(ctx, pkg.DownloadURL)
			if err != nil {
				return nil, fmt.Errorf("failed to download %s: %w", pkg.Filename, err)
			}

			files = append(files, &core.File{
				Path:    pkg.Filename,
				Content: content,
				Size:    pkg.Size,
			})
		}
	}

	if len(files) == 0 {
		return nil, fmt.Errorf("no packages found for %s version %s", packageName, version)
	}

	return files, nil
}

func (c *cloudsmithClient) downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("download failed: %d", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return content, nil
}

func (c *CloudsmithRegistry) matchesPatterns(path string, include, exclude []string) bool {
	if len(include) == 0 && len(exclude) == 0 {
		return true
	}

	for _, pattern := range exclude {
		if matchPattern(pattern, path) {
			return false
		}
	}

	if len(include) == 0 {
		return true
	}

	for _, pattern := range include {
		if matchPattern(pattern, path) {
			return true
		}
	}

	return false
}

func matchPattern(pattern, path string) bool {
	pattern = strings.ReplaceAll(pattern, "\\", "/")
	path = strings.ReplaceAll(path, "\\", "/")

	if strings.HasPrefix(pattern, "**/") {
		suffix := pattern[3:]
		return strings.HasSuffix(path, suffix) || strings.Contains(path, "/"+suffix)
	}

	if strings.HasSuffix(pattern, "/**") {
		prefix := pattern[:len(pattern)-3]
		return strings.HasPrefix(path, prefix+"/") || path == prefix
	}

	if !strings.Contains(pattern, "*") {
		return pattern == path
	}

	parts := strings.Split(pattern, "*")
	pos := 0
	for i, part := range parts {
		if part == "" {
			continue
		}
		idx := strings.Index(path[pos:], part)
		if idx == -1 || (i == 0 && idx != 0) {
			return false
		}
		pos += idx + len(part)
	}
	if len(parts) > 0 && parts[len(parts)-1] != "" && !strings.HasSuffix(path, parts[len(parts)-1]) {
		return false
	}
	return true
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
