package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/config"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/storage"
)

type GitLabRegistry struct {
	name         string
	config       GitLabRegistryConfig
	configMgr    config.Manager
	client       *gitLabClient
	packageCache *storage.PackageCache
}

type gitLabClient struct {
	baseURL    string
	apiVersion string
	token      string
	httpClient *http.Client
}

type gitLabPackage struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	PackageType string    `json:"package_type"`
	CreatedAt   time.Time `json:"created_at"`
}

type gitLabPackageFile struct {
	ID       int    `json:"id"`
	FileName string `json:"file_name"`
	Size     int64  `json:"size"`
}

func NewGitLabRegistry(name string, cfg *GitLabRegistryConfig, configMgr config.Manager) (*GitLabRegistry, error) {
	registry, err := storage.NewRegistry(cfg)
	if err != nil {
		return nil, err
	}

	apiVersion := cfg.APIVersion
	if apiVersion == "" {
		apiVersion = "v4"
	}

	return &GitLabRegistry{
		name:      name,
		config:    *cfg,
		configMgr: configMgr,
		client: &gitLabClient{
			baseURL:    ensureProtocol(cfg.URL),
			apiVersion: apiVersion,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		},
		packageCache: storage.NewPackageCache(registry.GetPackagesDir()),
	}, nil
}

func NewGitLabRegistryWithPath(baseDir, name string, cfg *GitLabRegistryConfig, configMgr config.Manager) (*GitLabRegistry, error) {
	registry, err := storage.NewRegistryWithPath(baseDir, cfg)
	if err != nil {
		return nil, err
	}

	apiVersion := cfg.APIVersion
	if apiVersion == "" {
		apiVersion = "v4"
	}

	return &GitLabRegistry{
		name:      name,
		config:    *cfg,
		configMgr: configMgr,
		client: &gitLabClient{
			baseURL:    ensureProtocol(cfg.URL),
			apiVersion: apiVersion,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		},
		packageCache: storage.NewPackageCache(registry.GetPackagesDir()),
	}, nil
}

func (g *GitLabRegistry) loadToken(ctx context.Context) error {
	if g.client.token != "" {
		return nil
	}

	if g.configMgr == nil {
		return nil
	}

	authKey := g.getAuthKey()
	token, err := g.configMgr.GetValue(ctx, "registry "+authKey, "token")
	if err != nil {
		return fmt.Errorf("failed to load token from .armrc: %w", err)
	}

	g.client.token = token
	return nil
}

func (g *GitLabRegistry) getAuthKey() string {
	baseURL := g.config.URL
	if g.config.ProjectID != "" {
		return fmt.Sprintf("%s/project/%s", baseURL, g.config.ProjectID)
	}
	return fmt.Sprintf("%s/group/%s", baseURL, g.config.GroupID)
}

func (g *GitLabRegistry) ListPackages(ctx context.Context) ([]*core.PackageMetadata, error) {
	if err := g.loadToken(ctx); err != nil {
		return nil, err
	}

	packages, err := g.client.listPackages(ctx, g.config.ProjectID, g.config.GroupID)
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	var result []*core.PackageMetadata
	for _, pkg := range packages {
		if pkg.PackageType == "generic" && !seen[pkg.Name] {
			seen[pkg.Name] = true
			result = append(result, &core.PackageMetadata{
				RegistryName: g.name,
				Name:         pkg.Name,
			})
		}
	}

	return result, nil
}

func (g *GitLabRegistry) ListPackageVersions(ctx context.Context, packageName string) ([]core.Version, error) {
	if err := g.loadToken(ctx); err != nil {
		return nil, err
	}

	packages, err := g.client.listPackages(ctx, g.config.ProjectID, g.config.GroupID)
	if err != nil {
		return nil, err
	}

	var versions []core.Version
	for _, pkg := range packages {
		if pkg.PackageType == "generic" && pkg.Name == packageName {
			version, _ := core.ParseVersion(pkg.Version)
			versions = append(versions, version)
		}
	}

	return versions, nil
}

func (g *GitLabRegistry) GetPackage(ctx context.Context, packageName string, version *core.Version, include, exclude []string) (*core.Package, error) {
	cacheKey := struct {
		Version core.Version `json:"version"`
		Include []string     `json:"include"`
		Exclude []string     `json:"exclude"`
	}{*version, normalizePatterns(include), normalizePatterns(exclude)}

	if files, err := g.packageCache.GetPackageVersion(ctx, cacheKey, version); err == nil {
		return &core.Package{
			Metadata: core.PackageMetadata{
				RegistryName: g.name,
				Name:         packageName,
				Version:      *version,
			},
			Files:     files,
			Integrity: calculateIntegrity(files),
		}, nil
	}

	if err := g.loadToken(ctx); err != nil {
		return nil, err
	}

	packages, err := g.client.listPackages(ctx, g.config.ProjectID, g.config.GroupID)
	if err != nil {
		return nil, err
	}

	var targetPackage *gitLabPackage
	for _, pkg := range packages {
		if pkg.PackageType == "generic" && pkg.Name == packageName && pkg.Version == version.Version {
			targetPackage = &pkg
			break
		}
	}

	if targetPackage == nil {
		return nil, fmt.Errorf("package %s version %s not found", packageName, version.Version)
	}

	files, err := g.client.downloadPackage(ctx, g.config.ProjectID, g.config.GroupID, targetPackage.ID, packageName, version.Version)
	if err != nil {
		return nil, err
	}

	extractor := core.NewExtractor()
	files, err = extractor.Extract(files)
	if err != nil {
		return nil, err
	}

	var filteredFiles []*core.File
	for _, file := range files {
		if matchesPatterns(file.Path, include, exclude) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	_ = g.packageCache.SetPackageVersion(ctx, cacheKey, version, filteredFiles)

	return &core.Package{
		Metadata: core.PackageMetadata{
			RegistryName: g.name,
			Name:         packageName,
			Version:      *version,
		},
		Files:     filteredFiles,
		Integrity: calculateIntegrity(filteredFiles),
	}, nil
}

func (c *gitLabClient) listPackages(ctx context.Context, projectID, groupID string) ([]gitLabPackage, error) {
	var apiURL string
	switch {
	case projectID != "":
		apiURL = fmt.Sprintf("%s/api/%s/projects/%s/packages", c.baseURL, c.apiVersion, url.QueryEscape(projectID))
	case groupID != "":
		apiURL = fmt.Sprintf("%s/api/%s/groups/%s/packages", c.baseURL, c.apiVersion, url.QueryEscape(groupID))
	default:
		return nil, fmt.Errorf("either project_id or group_id must be specified")
	}

	var allPackages []gitLabPackage
	page := 1
	perPage := 100

	for {
		paginatedURL := fmt.Sprintf("%s?page=%d&per_page=%d", apiURL, page, perPage)
		req, err := http.NewRequestWithContext(ctx, "GET", paginatedURL, http.NoBody)
		if err != nil {
			return nil, err
		}

		resp, err := c.makeRequest(req)
		if err != nil {
			return nil, err
		}

		var packages []gitLabPackage
		if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
			_ = resp.Body.Close()
			return nil, err
		}
		_ = resp.Body.Close()

		allPackages = append(allPackages, packages...)

		if len(packages) < perPage {
			break
		}
		page++
	}

	return allPackages, nil
}

func (c *gitLabClient) downloadPackage(ctx context.Context, projectID, groupID string, packageID int, packageName, version string) ([]*core.File, error) {
	filesURL := c.buildPackageFilesURL(projectID, groupID, packageID)

	req, err := http.NewRequestWithContext(ctx, "GET", filesURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var packageFiles []gitLabPackageFile
	if err := json.NewDecoder(resp.Body).Decode(&packageFiles); err != nil {
		return nil, err
	}

	var files []*core.File
	for _, pkgFile := range packageFiles {
		downloadURL := c.buildDownloadURL(projectID, groupID, packageName, version, pkgFile.FileName)

		req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, http.NoBody)
		if err != nil {
			return nil, fmt.Errorf("failed to create request for %s: %w", pkgFile.FileName, err)
		}

		resp, err := c.makeRequest(req)
		if err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", pkgFile.FileName, err)
		}

		content, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %w", pkgFile.FileName, err)
		}

		files = append(files, &core.File{
			Path:    pkgFile.FileName,
			Content: content,
			Size:    pkgFile.Size,
		})
	}

	return files, nil
}

func (c *gitLabClient) buildPackageFilesURL(projectID, groupID string, packageID int) string {
	if projectID != "" {
		return fmt.Sprintf("%s/api/%s/projects/%s/packages/%d/package_files", c.baseURL, c.apiVersion, url.QueryEscape(projectID), packageID)
	}
	return fmt.Sprintf("%s/api/%s/groups/%s/-/packages/%d/package_files", c.baseURL, c.apiVersion, url.QueryEscape(groupID), packageID)
}

func (c *gitLabClient) buildDownloadURL(projectID, groupID, packageName, version, fileName string) string {
	if projectID != "" {
		return fmt.Sprintf("%s/api/%s/projects/%s/packages/generic/%s/%s/%s", c.baseURL, c.apiVersion, url.QueryEscape(projectID), url.QueryEscape(packageName), url.QueryEscape(version), url.QueryEscape(fileName))
	}
	return fmt.Sprintf("%s/api/%s/groups/%s/-/packages/generic/%s/%s/%s", c.baseURL, c.apiVersion, url.QueryEscape(groupID), url.QueryEscape(packageName), url.QueryEscape(version), url.QueryEscape(fileName))
}

func (c *gitLabClient) makeRequest(req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("GitLab API error %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

func ensureProtocol(baseURL string) string {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return "https://" + baseURL
	}
	return baseURL
}

func matchesPatterns(filePath string, include, exclude []string) bool {
	// Default to YAML files if no patterns specified
	if len(include) == 0 && len(exclude) == 0 {
		include = []string{"**/*.yml", "**/*.yaml"}
	}

	for _, pattern := range exclude {
		if core.MatchPattern(pattern, filePath) {
			return false
		}
	}

	if len(include) == 0 {
		return true
	}

	for _, pattern := range include {
		if core.MatchPattern(pattern, filePath) {
			return true
		}
	}

	return false
}
