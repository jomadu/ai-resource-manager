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

	"github.com/jomadu/ai-rules-manager/internal/archive"
	"github.com/jomadu/ai-rules-manager/internal/cache"
	"github.com/jomadu/ai-rules-manager/internal/rcfile"
	"github.com/jomadu/ai-rules-manager/internal/registry/common"
	"github.com/jomadu/ai-rules-manager/internal/resolver"
	"github.com/jomadu/ai-rules-manager/internal/types"
)

// CloudsmithRegistry implements the Registry interface for Cloudsmith package registries
type CloudsmithRegistry struct {
	cache        cache.RegistryPackageCache
	config       CloudsmithRegistryConfig
	resolver     resolver.ConstraintResolver
	client       *CloudsmithClient
	registryName string
	semver       *common.SemverHelper
	rcService    *rcfile.Service
	extractor    *archive.Extractor
}

// CloudsmithClient handles HTTP communication with Cloudsmith API
type CloudsmithClient struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// Cloudsmith API endpoint templates
const (
	CloudsmithPackageListTemplate     = "/v1/packages/%s/%s/"
	CloudsmithPackageDownloadTemplate = "/v1/packages/%s/%s/%s/download/"
)

// Cloudsmith package types
type CloudsmithPackage struct {
	Slug        string    `json:"slug"`
	Name        string    `json:"name"`
	Version     string    `json:"version"`
	Format      string    `json:"format"`
	Size        int64     `json:"size"`
	Filename    string    `json:"filename"`
	UploadedAt  time.Time `json:"uploaded_at"`
	DownloadURL string    `json:"cdn_url"`
	// Additional fields from actual API response
	DisplayName string              `json:"display_name"`
	Description string              `json:"description"`
	Summary     string              `json:"summary"`
	Namespace   string              `json:"namespace"`
	Repository  string              `json:"repository"`
	Tags        map[string][]string `json:"tags"`
}

// Note: The actual API returns an array directly, not wrapped in a "results" object
// We'll handle this in the parsing logic

// NewCloudsmithRegistry creates a new Cloudsmith-based registry
func NewCloudsmithRegistry(registryName string, config *CloudsmithRegistryConfig, packageCache cache.RegistryPackageCache) *CloudsmithRegistry {
	client := &CloudsmithClient{
		baseURL:    "https://api.cloudsmith.io", // Always use API base URL for client
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	return &CloudsmithRegistry{
		cache:        packageCache,
		config:       *config,
		resolver:     resolver.NewGitConstraintResolver(),
		client:       client,
		registryName: registryName,
		semver:       common.NewSemverHelper(),
		rcService:    rcfile.NewService(),
		extractor:    archive.NewExtractor(),
	}
}

// NewCloudsmithRegistryNoCache creates a new Cloudsmith-based registry without caching for testing
func NewCloudsmithRegistryNoCache(registryName string, config *CloudsmithRegistryConfig) *CloudsmithRegistry {
	client := &CloudsmithClient{
		baseURL:    "https://api.cloudsmith.io", // Always use API base URL for client
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	return &CloudsmithRegistry{
		cache:        cache.NewNoopRegistryPackageCache(),
		config:       *config,
		resolver:     resolver.NewGitConstraintResolver(),
		client:       client,
		registryName: registryName,
		semver:       common.NewSemverHelper(),
		rcService:    rcfile.NewService(),
		extractor:    archive.NewExtractor(),
	}
}

// loadToken loads authentication token from .armrc file
func (c *CloudsmithRegistry) loadToken() error {
	// Construct the full URL with owner/repository path for .armrc lookup
	fullURL := fmt.Sprintf("%s/%s/%s", c.config.URL, c.config.Owner, c.config.Repository)
	token, err := c.rcService.GetValue("registry "+fullURL, "token")
	if err != nil {
		return fmt.Errorf("failed to load token from .armrc: %w", err)
	}
	c.client.token = token
	return nil
}

func (c *CloudsmithRegistry) ListVersions(ctx context.Context, ruleset string) ([]types.Version, error) {
	if err := c.loadToken(); err != nil {
		return nil, err
	}

	packages, err := c.client.ListPackages(ctx, c.config.Owner, c.config.Repository, ruleset)
	if err != nil {
		return nil, err
	}

	// Extract unique versions from packages with matching name
	versionMap := make(map[string]bool)
	for i := range packages {
		pkg := &packages[i]
		if pkg.Format == "raw" && (pkg.Name == ruleset || strings.HasPrefix(pkg.Filename, ruleset)) {
			versionMap[pkg.Version] = true
		}
	}

	var versionStrings []string
	for version := range versionMap {
		versionStrings = append(versionStrings, version)
	}

	sortedVersions := c.semver.SortVersionsBySemver(versionStrings)

	var versions []types.Version
	for _, version := range sortedVersions {
		versions = append(versions, types.Version{Version: version, Display: version})
	}

	return versions, nil
}

func (c *CloudsmithRegistry) ResolveVersion(ctx context.Context, ruleset, constraint string) (*resolver.ResolvedVersion, error) {
	parsedConstraint, err := c.resolver.ParseConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint %s: %w", constraint, err)
	}

	versions, err := c.ListVersions(ctx, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	resolvedVersion, err := c.resolver.FindBestMatch(parsedConstraint, versions)
	if err != nil {
		return nil, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return &resolver.ResolvedVersion{
		Constraint: parsedConstraint,
		Version:    *resolvedVersion,
	}, nil
}

func (c *CloudsmithRegistry) GetContent(ctx context.Context, ruleset string, version types.Version, selector types.ContentSelector) ([]types.File, error) {
	// Try cache first
	files, err := c.cache.GetPackageVersion(ctx, selector, version.Version)
	if err == nil {
		return files, nil
	}

	if err := c.loadToken(); err != nil {
		return nil, err
	}

	// Download files from Cloudsmith
	rawFiles, err := c.client.DownloadPackages(ctx, c.config.Owner, c.config.Repository, ruleset, version.Version)
	if err != nil {
		return nil, err
	}

	// Extract and merge archives with loose files
	mergedFiles, err := c.extractor.ExtractAndMerge(rawFiles)
	if err != nil {
		return nil, fmt.Errorf("failed to extract and merge content: %w", err)
	}

	// Apply selector patterns to merged content
	var filteredFiles []types.File
	for _, file := range mergedFiles {
		if selector.Matches(file.Path) {
			filteredFiles = append(filteredFiles, file)
		}
	}

	// Cache the result
	_ = c.cache.SetPackageVersion(ctx, selector, version.Version, filteredFiles)

	return filteredFiles, nil
}

// Cloudsmith Client methods

func (c *CloudsmithClient) ListPackages(ctx context.Context, owner, repo, packageName string) ([]CloudsmithPackage, error) {
	url := c.buildPackageListURL(owner, repo, packageName)

	var allPackages []CloudsmithPackage
	nextURL := &url

	for nextURL != nil {
		packages, next, err := c.listPackagesPage(ctx, *nextURL)
		if err != nil {
			return nil, err
		}
		allPackages = append(allPackages, packages...)
		nextURL = next
	}

	return allPackages, nil
}

func (c *CloudsmithClient) listPackagesPage(ctx context.Context, url string) ([]CloudsmithPackage, *string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return nil, nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	// The API returns an array directly, not wrapped in a "results" object
	var packages []CloudsmithPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, nil, err
	}

	// Parse pagination from Link header
	nextURL := c.parseNextURLFromLinkHeader(resp.Header.Get("Link"))

	return packages, nextURL, nil
}

func (c *CloudsmithClient) DownloadPackages(ctx context.Context, owner, repo, packageName, version string) ([]types.File, error) {
	// First get package list to find matching packages
	packages, err := c.ListPackages(ctx, owner, repo, packageName)
	if err != nil {
		return nil, err
	}

	var files []types.File
	for i := range packages {
		pkg := &packages[i]
		if pkg.Version == version && (pkg.Name == packageName || strings.HasPrefix(pkg.Filename, packageName)) {
			content, err := c.downloadFile(ctx, pkg.DownloadURL)
			if err != nil {
				return nil, fmt.Errorf("failed to download %s: %w", pkg.Filename, err)
			}

			files = append(files, types.File{
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

// URL builders
func (c *CloudsmithClient) buildPackageListURL(owner, repo, packageName string) string {
	baseURL := fmt.Sprintf(c.baseURL+CloudsmithPackageListTemplate, url.QueryEscape(owner), url.QueryEscape(repo))

	// Add query parameter if packageName is specified
	if packageName != "" {
		params := url.Values{}
		params.Add("query", packageName)
		return baseURL + "?" + params.Encode()
	}

	return baseURL
}

// parseNextURLFromLinkHeader parses the Link header and extracts the "next" URL
func (c *CloudsmithClient) parseNextURLFromLinkHeader(linkHeader string) *string {
	if linkHeader == "" {
		return nil
	}

	// Parse Link header format: <url>; rel="next", <url>; rel="prev"
	links := strings.Split(linkHeader, ",")
	for _, link := range links {
		link = strings.TrimSpace(link)
		if strings.Contains(link, `rel="next"`) {
			// Extract URL from <url>; rel="next"
			start := strings.Index(link, "<")
			end := strings.Index(link, ">")
			if start != -1 && end != -1 && start < end {
				nextURL := link[start+1 : end]
				return &nextURL
			}
		}
	}

	return nil
}

// HTTP helpers
func (c *CloudsmithClient) downloadFile(ctx context.Context, downloadURL string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	return io.ReadAll(resp.Body)
}

func (c *CloudsmithClient) makeRequest(_ context.Context, req *http.Request) (*http.Response, error) {
	if c.token != "" {
		req.Header.Set("Authorization", "Token "+c.token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return nil, fmt.Errorf("Cloudsmith API error %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}
