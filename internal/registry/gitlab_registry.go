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

// GitLabRegistry implements the Registry interface for GitLab package registries
type GitLabRegistry struct {
	cache        cache.RegistryRulesetCache
	config       GitLabRegistryConfig
	resolver     resolver.ConstraintResolver
	client       *GitLabClient
	registryName string
	semver       *common.SemverHelper
	rcService    *rcfile.Service
	extractor    *archive.Extractor
}

// GitLabClient handles HTTP communication with GitLab API
type GitLabClient struct {
	baseURL    string
	apiVersion string
	httpClient *http.Client
	token      string
}

// GitLab API endpoint templates
const (
	ProjectPackageListTemplate     = "/api/%s/projects/%s/packages"
	ProjectPackageFilesTemplate    = "/api/%s/projects/%s/packages/%d/package_files"
	ProjectPackageDownloadTemplate = "/api/%s/projects/%s/packages/generic/%s/%s/%s"
	GroupPackageListTemplate       = "/api/%s/groups/%s/packages"
	GroupPackageFilesTemplate      = "/api/%s/groups/%s/-/packages/%d/package_files"
	GroupPackageDownloadTemplate   = "/api/%s/groups/%s/-/packages/generic/%s/%s/%s"
)

// GitLab package types
type GitLabPackage struct {
	ID          int                 `json:"id"`
	Name        string              `json:"name"`
	Version     string              `json:"version"`
	PackageType string              `json:"packageType"`
	CreatedAt   time.Time           `json:"createdAt"`
	Files       []GitLabPackageFile `json:"packageFiles"`
}

type GitLabPackageFile struct {
	ID       int    `json:"id"`
	FileName string `json:"fileName"`
	Size     int64  `json:"size"`
}

// NewGitLabRegistry creates a new GitLab-based registry
func NewGitLabRegistry(registryName string, config *GitLabRegistryConfig, rulesetCache cache.RegistryRulesetCache) *GitLabRegistry {
	client := &GitLabClient{
		baseURL:    config.URL,
		apiVersion: config.GetAPIVersion(),
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}

	return &GitLabRegistry{
		cache:        rulesetCache,
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
func (g *GitLabRegistry) loadToken() error {
	authKey := g.getAuthKey()
	token, err := g.rcService.GetValue("registry "+authKey, "token")
	if err != nil {
		return fmt.Errorf("failed to load token from .armrc: %w", err)
	}
	g.client.token = token
	return nil
}

// getAuthKey constructs the composite authentication key
func (g *GitLabRegistry) getAuthKey() string {
	host := strings.TrimPrefix(g.config.URL, "https://")
	host = strings.TrimPrefix(host, "http://")

	if g.config.ProjectID != "" {
		return fmt.Sprintf("%s/project/%s", host, g.config.ProjectID)
	}
	return fmt.Sprintf("%s/group/%s", host, g.config.GroupID)
}

func (g *GitLabRegistry) ListVersions(ctx context.Context, ruleset string) ([]types.Version, error) {
	// Use ruleset as package name for GitLab Package Registry
	if err := g.loadToken(); err != nil {
		return nil, err
	}

	var packages []GitLabPackage
	var err error

	switch {
	case g.config.ProjectID != "":
		packages, err = g.client.ListProjectPackages(ctx, g.config.ProjectID)
	case g.config.GroupID != "":
		packages, err = g.client.ListGroupPackages(ctx, g.config.GroupID)
	default:
		return nil, fmt.Errorf("either project_id or group_id must be specified")
	}

	if err != nil {
		return nil, err
	}

	// Extract unique versions from the specific ruleset package
	versionMap := make(map[string]bool)
	for _, pkg := range packages {
		if pkg.PackageType == "generic" && pkg.Name == ruleset {
			versionMap[pkg.Version] = true
		}
	}

	var versionStrings []string
	for version := range versionMap {
		versionStrings = append(versionStrings, version)
	}

	sortedVersions := g.semver.SortVersionsBySemver(versionStrings)

	var versions []types.Version
	for _, version := range sortedVersions {
		versions = append(versions, types.Version{Version: version, Display: version})
	}

	return versions, nil
}

func (g *GitLabRegistry) ResolveVersion(ctx context.Context, ruleset, constraint string) (*resolver.ResolvedVersion, error) {
	parsedConstraint, err := g.resolver.ParseConstraint(constraint)
	if err != nil {
		return nil, fmt.Errorf("invalid version constraint %s: %w", constraint, err)
	}

	versions, err := g.ListVersions(ctx, ruleset)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}

	resolvedVersion, err := g.resolver.FindBestMatch(parsedConstraint, versions)
	if err != nil {
		return nil, fmt.Errorf("no matching version found for %s: %w", constraint, err)
	}

	return &resolver.ResolvedVersion{
		Constraint: parsedConstraint,
		Version:    *resolvedVersion,
	}, nil
}

func (g *GitLabRegistry) GetContent(ctx context.Context, ruleset string, version types.Version, selector types.ContentSelector) ([]types.File, error) {
	// Try cache first
	files, err := g.cache.GetRulesetVersion(ctx, selector, version.Version)
	if err == nil {
		return files, nil
	}

	if err := g.loadToken(); err != nil {
		return nil, err
	}

	// Download all files from package
	var rawFiles []types.File
	switch {
	case g.config.ProjectID != "":
		rawFiles, err = g.client.DownloadProjectPackage(ctx, g.config.ProjectID, ruleset, version.Version)
	case g.config.GroupID != "":
		rawFiles, err = g.client.DownloadGroupPackage(ctx, g.config.GroupID, ruleset, version.Version)
	default:
		return nil, fmt.Errorf("either project_id or group_id must be specified")
	}

	if err != nil {
		return nil, err
	}

	// Extract and merge archives with loose files
	mergedFiles, err := g.extractor.ExtractAndMerge(rawFiles)
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
	_ = g.cache.SetRulesetVersion(ctx, selector, version.Version, filteredFiles)

	return filteredFiles, nil
}

// GitLab Client methods

func (c *GitLabClient) ListProjectPackages(ctx context.Context, projectID string) ([]GitLabPackage, error) {
	url := c.buildProjectPackageListURL(projectID)
	return c.listPackages(ctx, url)
}

func (c *GitLabClient) ListGroupPackages(ctx context.Context, groupID string) ([]GitLabPackage, error) {
	url := c.buildGroupPackageListURL(groupID)
	return c.listPackages(ctx, url)
}

func (c *GitLabClient) DownloadProjectPackage(ctx context.Context, projectID, packageName, version string) ([]types.File, error) {
	// First get package info to get package ID
	packages, err := c.ListProjectPackages(ctx, projectID)
	if err != nil {
		return nil, err
	}

	var targetPackage *GitLabPackage
	for _, pkg := range packages {
		if pkg.Name == packageName && pkg.Version == version && pkg.PackageType == "generic" {
			targetPackage = &pkg
			break
		}
	}

	if targetPackage == nil {
		return nil, fmt.Errorf("package %s version %s not found", packageName, version)
	}

	// Get package files
	filesURL := c.buildProjectPackageFilesURL(projectID, targetPackage.ID)
	packageFiles, err := c.getPackageFiles(ctx, filesURL)
	if err != nil {
		return nil, err
	}

	// Download all files
	var files []types.File
	for _, pkgFile := range packageFiles {
		downloadURL := c.buildProjectPackageDownloadURL(projectID, packageName, version, pkgFile.FileName)
		content, err := c.downloadFile(ctx, downloadURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", pkgFile.FileName, err)
		}

		files = append(files, types.File{
			Path:    pkgFile.FileName,
			Content: content,
			Size:    pkgFile.Size,
		})
	}

	return files, nil
}

func (c *GitLabClient) DownloadGroupPackage(ctx context.Context, groupID, packageName, version string) ([]types.File, error) {
	// First get package info to get package ID
	packages, err := c.ListGroupPackages(ctx, groupID)
	if err != nil {
		return nil, err
	}

	var targetPackage *GitLabPackage
	for _, pkg := range packages {
		if pkg.Name == packageName && pkg.Version == version && pkg.PackageType == "generic" {
			targetPackage = &pkg
			break
		}
	}

	if targetPackage == nil {
		return nil, fmt.Errorf("package %s version %s not found", packageName, version)
	}

	// Get package files
	filesURL := c.buildGroupPackageFilesURL(groupID, targetPackage.ID)
	packageFiles, err := c.getPackageFiles(ctx, filesURL)
	if err != nil {
		return nil, err
	}

	// Download all files
	var files []types.File
	for _, pkgFile := range packageFiles {
		downloadURL := c.buildGroupPackageDownloadURL(groupID, packageName, version, pkgFile.FileName)
		content, err := c.downloadFile(ctx, downloadURL)
		if err != nil {
			return nil, fmt.Errorf("failed to download %s: %w", pkgFile.FileName, err)
		}

		files = append(files, types.File{
			Path:    pkgFile.FileName,
			Content: content,
			Size:    pkgFile.Size,
		})
	}

	return files, nil
}

// URL builders
func (c *GitLabClient) buildProjectPackageListURL(projectID string) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+ProjectPackageListTemplate, c.apiVersion, url.QueryEscape(projectID))
}

func (c *GitLabClient) buildProjectPackageFilesURL(projectID string, packageID int) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+ProjectPackageFilesTemplate, c.apiVersion, url.QueryEscape(projectID), packageID)
}

func (c *GitLabClient) buildProjectPackageDownloadURL(projectID, packageName, version, fileName string) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+ProjectPackageDownloadTemplate, c.apiVersion, url.QueryEscape(projectID), url.QueryEscape(packageName), url.QueryEscape(version), url.QueryEscape(fileName))
}

func (c *GitLabClient) buildGroupPackageListURL(groupID string) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+GroupPackageListTemplate, c.apiVersion, url.QueryEscape(groupID))
}

func (c *GitLabClient) buildGroupPackageFilesURL(groupID string, packageID int) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+GroupPackageFilesTemplate, c.apiVersion, url.QueryEscape(groupID), packageID)
}

func (c *GitLabClient) buildGroupPackageDownloadURL(groupID, packageName, version, fileName string) string {
	baseURL := c.ensureProtocol(c.baseURL)
	return fmt.Sprintf(baseURL+GroupPackageDownloadTemplate, c.apiVersion, url.QueryEscape(groupID), url.QueryEscape(packageName), url.QueryEscape(version), url.QueryEscape(fileName))
}

// ensureProtocol adds https:// if no protocol is present
func (c *GitLabClient) ensureProtocol(baseURL string) string {
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		return "https://" + baseURL
	}
	return baseURL
}

// HTTP helpers
func (c *GitLabClient) listPackages(ctx context.Context, url string) ([]GitLabPackage, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var packages []GitLabPackage
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func (c *GitLabClient) getPackageFiles(ctx context.Context, url string) ([]GitLabPackageFile, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
	if err != nil {
		return nil, err
	}

	resp, err := c.makeRequest(ctx, req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	var files []GitLabPackageFile
	if err := json.NewDecoder(resp.Body).Decode(&files); err != nil {
		return nil, err
	}

	return files, nil
}

func (c *GitLabClient) downloadFile(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, http.NoBody)
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

func (c *GitLabClient) makeRequest(_ context.Context, req *http.Request) (*http.Response, error) {
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
