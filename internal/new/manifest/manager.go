package manifest

import "context"

type Manifest struct {
	Version string `json:"version"`
	Registries map[string]map[string]interface{} `json:"registries,omitempty"`
	Sinks      map[string]SinkConfig             `json:"sinks,omitempty"`
	Packages   map[string]map[string]interface{} `json:"packages"`
}

type SinkConfig struct {
	Directory string `json:"directory"`
	Layout string `json:"layout,omitempty"`
	CompileTarget string `json:"compileTarget"`
}

type PackageConfig struct {
	Version string `json:"version"`
	Include []string `json:"include,omitempty"`
	Exclude []string `json:"exclude,omitempty"`
	Sinks []string `json:"sinks"`
	ResourceType ResourceType `json:"resourceType,omitempty"`
}

type RulesetConfig struct {
	PackageConfig
	Priority int `json:"priority,omitempty"`
}

type PromptsetConfig struct {
	PackageConfig
}

type RegistryConfig struct {
	URL string `json:"url"`
	Type string `json:"type"`
}

type GitRegistryConfig struct {
	RegistryConfig
	Branches []string `json:"branches,omitempty"`
}

type GitLabRegistryConfig struct {
	RegistryConfig
	ProjectID string `json:"projectId,omitempty"`
	GroupID string `json:"groupId,omitempty"`
	APIVersion string `json:"apiVersion,omitempty"`
}

type CloudsmithRegistryConfig struct {
	RegistryConfig
	Owner string `json:"owner"`
	Repository string `json:"repository"`
}
type Manager interface {
	// Registry operations (generic)
	GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error)
	GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error)
	AddRegistryConfig(ctx context.Context, name string, config map[string]interface{}) error
	UpdateRegistryConfigName(ctx context.Context, name string, newName string) error
	UpdateRegistryConfig(ctx context.Context, name string, config map[string]interface{}) error
	RemoveRegistryConfig(ctx context.Context, name string) error

	// Registry operations (type-safe helpers)
	GetGitRegistryConfig(ctx context.Context, name string) (*GitRegistryConfig, error)
	GetGitLabRegistryConfig(ctx context.Context, name string) (*GitLabRegistryConfig, error)
	GetCloudsmithRegistryConfig(ctx context.Context, name string) (*CloudsmithRegistryConfig, error)
	AddGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error
	AddGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error
	AddCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error
	UpdateGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error
	UpdateGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error
	UpdateCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error

	// Sink operations
	GetAllSinksConfig(ctx context.Context) (map[string]*SinkConfig, error)
	GetSinkConfig(ctx context.Context, name string) (*SinkConfig, error)
	AddSinkConfig(ctx context.Context, name string, config *SinkConfig) error
	UpdateSinkConfigName(ctx context.Context, name string, newName string) error
	UpdateSinkConfig(ctx context.Context, name string, config *SinkConfig) error
	RemoveSinkConfig(ctx context.Context, name string) error

	// Package operations (generic)
	GetAllPackagesConfig(ctx context.Context) (map[string]map[string]interface{}, error)
	GetPackageConfig(ctx context.Context, registryName, packageName string) (map[string]interface{}, error)
	AddPackageConfig(ctx context.Context, registryName, packageName string, config map[string]interface{}) error
	UpdatePackageConfigName(ctx context.Context, registryName, packageName string, newPackageName string) error
	UpdatePackageConfig(ctx context.Context, registryName, packageName string, config map[string]interface{}) error
	RemovePackageConfig(ctx context.Context, registryName, packageName string) error

	// Package operations (type-safe helpers)
	GetRulesetConfig(ctx context.Context, registryName, packageName string) (*RulesetConfig, error)
	GetPromptsetConfig(ctx context.Context, registryName, packageName string) (*PromptsetConfig, error)
	AddRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error
	AddPromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error
	UpdateRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error
	UpdatePromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error
}

// FileManager implements file-based manifest management.
// It reads from and writes to arm.json in the current directory.
type FileManager struct{}

// NewFileManager creates a new file-based manifest manager.
func NewFileManager() *FileManager {
	return &FileManager{}
}

// Registry operations (generic)

func (f *FileManager) GetAllRegistriesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetRegistryConfig(ctx context.Context, name string) (map[string]interface{}, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) AddRegistryConfig(ctx context.Context, name string, config map[string]interface{}) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateRegistryConfigName(ctx context.Context, name string, newName string) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateRegistryConfig(ctx context.Context, name string, config map[string]interface{}) error {
	// TODO: implement
	return nil
}

func (f *FileManager) RemoveRegistryConfig(ctx context.Context, name string) error {
	// TODO: implement
	return nil
}

// Registry operations (type-safe helpers)

func (f *FileManager) GetGitRegistryConfig(ctx context.Context, name string) (*GitRegistryConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetGitLabRegistryConfig(ctx context.Context, name string) (*GitLabRegistryConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetCloudsmithRegistryConfig(ctx context.Context, name string) (*CloudsmithRegistryConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) AddGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) AddGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) AddCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateGitRegistryConfig(ctx context.Context, name string, config *GitRegistryConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateGitLabRegistryConfig(ctx context.Context, name string, config *GitLabRegistryConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateCloudsmithRegistryConfig(ctx context.Context, name string, config *CloudsmithRegistryConfig) error {
	// TODO: implement
	return nil
}

// Sink operations

func (f *FileManager) GetAllSinksConfig(ctx context.Context) (map[string]*SinkConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetSinkConfig(ctx context.Context, name string) (*SinkConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) AddSinkConfig(ctx context.Context, name string, config *SinkConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateSinkConfigName(ctx context.Context, name string, newName string) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateSinkConfig(ctx context.Context, name string, config *SinkConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) RemoveSinkConfig(ctx context.Context, name string) error {
	// TODO: implement
	return nil
}

// Package operations (generic)

func (f *FileManager) GetAllPackagesConfig(ctx context.Context) (map[string]map[string]interface{}, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetPackageConfig(ctx context.Context, registryName, packageName string) (map[string]interface{}, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) AddPackageConfig(ctx context.Context, registryName, packageName string, config map[string]interface{}) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdatePackageConfigName(ctx context.Context, registryName, packageName string, newPackageName string) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdatePackageConfig(ctx context.Context, registryName, packageName string, config map[string]interface{}) error {
	// TODO: implement
	return nil
}

func (f *FileManager) RemovePackageConfig(ctx context.Context, registryName, packageName string) error {
	// TODO: implement
	return nil
}

// Package operations (type-safe helpers)

func (f *FileManager) GetRulesetConfig(ctx context.Context, registryName, packageName string) (*RulesetConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) GetPromptsetConfig(ctx context.Context, registryName, packageName string) (*PromptsetConfig, error) {
	// TODO: implement
	return nil, nil
}

func (f *FileManager) AddRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) AddPromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdateRulesetConfig(ctx context.Context, registryName, packageName string, config *RulesetConfig) error {
	// TODO: implement
	return nil
}

func (f *FileManager) UpdatePromptsetConfig(ctx context.Context, registryName, packageName string, config *PromptsetConfig) error {
	// TODO: implement
	return nil
}