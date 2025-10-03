# Resource Manager Implementation Plan

## Core Abstractions

### Resource Interface
```go
type ResourceType string

const (
    ResourceTypeRuleset   ResourceType = "ruleset"
    ResourceTypePromptset ResourceType = "promptset"
)

type Resource interface {
    GetType() ResourceType
    GetID() string
    GetMetadata() Metadata
    Validate() error
    Compile(tool string, options CompileOptions) ([]CompiledFile, error)
}

type Metadata struct {
    ID          string            `yaml:"id"`
    Name        string            `yaml:"name"`
    Description string            `yaml:"description,omitempty"`
    Version     string            `yaml:"version"`
    Tags        []string          `yaml:"tags,omitempty"`
    Author      string            `yaml:"author,omitempty"`
    License     string            `yaml:"license,omitempty"`
    Extra       map[string]any    `yaml:",inline"`
}
```

### Resource Implementations
```go
type Ruleset struct {
    Version  string            `yaml:"version"`
    Metadata Metadata          `yaml:"metadata"`
    Rules    map[string]Rule   `yaml:"rules"`
}

type Rule struct {
    Name        string   `yaml:"name"`
    Priority    int      `yaml:"priority,omitempty"`
    Enforcement string   `yaml:"enforcement,omitempty"`
    Body        string   `yaml:"body"`
    Tags        []string `yaml:"tags,omitempty"`
}

type Promptset struct {
    Version  string              `yaml:"version"`
    Metadata Metadata            `yaml:"metadata"`
    Prompts  map[string]Prompt   `yaml:"prompts"`
}

type Prompt struct {
    Name        string         `yaml:"name"`
    Description string         `yaml:"description,omitempty"`
    Body        string         `yaml:"body"`
    Parameters  map[string]any `yaml:"parameters,omitempty"`
    Tags        []string       `yaml:"tags,omitempty"`
}
```

## Command Structure Refactoring

### Current Command Handler
```go
// cmd/install.go (current)
func installCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "install <ruleset>",
        RunE:  installRuleset,
    }
}
```

### New Command Structure
```go
// cmd/install.go (new)
func installCommand() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "install",
        Short: "Install resources",
    }

    cmd.AddCommand(installRulesetCommand())
    cmd.AddCommand(installPromptsetCommand())

    return cmd
}

func installRulesetCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "ruleset <name>",
        RunE:  func(cmd *cobra.Command, args []string) error {
            return installResource(ResourceTypeRuleset, args[0], cmd)
        },
    }
}

func installPromptsetCommand() *cobra.Command {
    return &cobra.Command{
        Use:   "promptset <name>",
        RunE:  func(cmd *cobra.Command, args []string) error {
            return installResource(ResourceTypePromptset, args[0], cmd)
        },
    }
}
```

## Configuration Evolution

### Current Config (v2)
```go
type Config struct {
    Version    string                    `json:"version"`
    Registries map[string]Registry       `json:"registries"`
    Sinks      map[string]Sink           `json:"sinks"`
    Rulesets   map[string]RulesetConfig  `json:"rulesets"`
}
```

### New Config (v3)
```go
type Config struct {
    Version     string                     `json:"version"`
    Registries  map[string]Registry        `json:"registries"`
    Sinks       map[string]Sink            `json:"sinks"`
    Rulesets    map[string]RulesetConfig   `json:"rulesets"`
    Promptsets  map[string]PromptsetConfig `json:"promptsets"`
}

type ResourceConfig struct {
    Version  string   `json:"version"`
    Sinks    []string `json:"sinks"`
    Include  []string `json:"include,omitempty"`
}

type RulesetConfig struct {
    ResourceConfig
    Priority int `json:"priority,omitempty"`
}

type PromptsetConfig struct {
    ResourceConfig
    Parameters map[string]any `json:"parameters,omitempty"`
}
```

## Registry System Updates

### Resource Discovery
```go
type RegistryResourceCache interface {
    GetResource(resourceType ResourceType, name string, version string) (Resource, error)
    ListResources(resourceType ResourceType, name string) ([]string, error)
    InvalidateResource(resourceType ResourceType, name string, version string)
}

type GitRegistryCache struct {
    // ... existing fields
}

func (g *GitRegistryCache) GetResource(resourceType ResourceType, name string, version string) (Resource, error) {
    // Discover resources by type-specific patterns
    patterns := getDiscoveryPatterns(resourceType)
    files := g.findMatchingFiles(patterns, name)

    for _, file := range files {
        content, err := g.getFileContent(file, version)
        if err != nil {
            continue
        }

        resource, err := parseResource(content)
        if err != nil {
            continue
        }

        if resource.GetType() == resourceType && resource.GetID() == name {
            return resource, nil
        }
    }

    return nil, fmt.Errorf("resource not found: %s/%s", resourceType, name)
}

func getDiscoveryPatterns(resourceType ResourceType) []string {
    switch resourceType {
    case ResourceTypeRuleset:
        return []string{"rulesets/**/*.yml", "rulesets/**/*.yaml", "**/*ruleset*.yml"}
    case ResourceTypePromptset:
        return []string{"promptsets/**/*.yml", "promptsets/**/*.yaml", "**/*promptset*.yml"}
    default:
        return []string{"**/*.yml", "**/*.yaml"}
    }
}
```

### Resource Parsing
```go
func parseResource(content []byte) (Resource, error) {
    // First, detect the resource type
    resourceType, err := detectResourceType(content)
    if err != nil {
        return nil, err
    }

    // Parse based on detected type
    switch resourceType {
    case ResourceTypeRuleset:
        var ruleset Ruleset
        if err := yaml.Unmarshal(content, &ruleset); err != nil {
            return nil, err
        }
        return &ruleset, nil

    case ResourceTypePromptset:
        var promptset Promptset
        if err := yaml.Unmarshal(content, &promptset); err != nil {
            return nil, err
        }
        return &promptset, nil

    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
}

func detectResourceType(content []byte) (ResourceType, error) {
    var detector struct {
        Rules   map[string]any `yaml:"rules"`
        Prompts map[string]any `yaml:"prompts"`
    }

    if err := yaml.Unmarshal(content, &detector); err != nil {
        return "", err
    }

    hasRules := len(detector.Rules) > 0
    hasPrompts := len(detector.Prompts) > 0

    if hasRules && hasPrompts {
        return "", errors.New("resource cannot contain both rules and prompts")
    }
    if hasRules {
        return ResourceTypeRuleset, nil
    }
    if hasPrompts {
        return ResourceTypePromptset, nil
    }

    return "", errors.New("resource must contain either rules or prompts")
}
```

## Compilation Pipeline

### Unified Compiler Interface
```go
type Compiler interface {
    Compile(resource Resource, tool string, options CompileOptions) ([]CompiledFile, error)
    SupportedTools() []string
    SupportedResourceTypes() []ResourceType
}

type CompileOptions struct {
    Namespace string
    Priority  int
    Metadata  map[string]any
}

type CompiledFile struct {
    Path    string
    Content []byte
    Mode    os.FileMode
}
```

### Resource-Specific Compilers
```go
type RulesetCompiler struct{}

func (c *RulesetCompiler) Compile(resource Resource, tool string, options CompileOptions) ([]CompiledFile, error) {
    ruleset, ok := resource.(*Ruleset)
    if !ok {
        return nil, errors.New("expected ruleset resource")
    }

    switch tool {
    case "cursor":
        return c.compileCursor(ruleset, options)
    case "amazonq":
        return c.compileAmazonQ(ruleset, options)
    default:
        return nil, fmt.Errorf("unsupported tool: %s", tool)
    }
}

type PromptsetCompiler struct{}

func (c *PromptsetCompiler) Compile(resource Resource, tool string, options CompileOptions) ([]CompiledFile, error) {
    promptset, ok := resource.(*Promptset)
    if !ok {
        return nil, errors.New("expected promptset resource")
    }

    switch tool {
    case "cursor":
        return c.compileCursor(promptset, options)
    case "amazonq":
        return c.compileAmazonQ(promptset, options)
    default:
        return nil, fmt.Errorf("unsupported tool: %s", tool)
    }
}
```

## Migration Strategy

### Configuration Migration
```go
func MigrateConfig(configPath string) error {
    content, err := os.ReadFile(configPath)
    if err != nil {
        return err
    }

    var versionCheck struct {
        Version string `json:"version"`
    }

    if err := json.Unmarshal(content, &versionCheck); err != nil {
        return err
    }

    switch versionCheck.Version {
    case "2.0":
        return migrateV2ToV3(configPath, content)
    case "3.0":
        return nil // Already migrated
    default:
        return fmt.Errorf("unsupported config version: %s", versionCheck.Version)
    }
}

func migrateV2ToV3(configPath string, content []byte) error {
    var v2Config ConfigV2
    if err := json.Unmarshal(content, &v2Config); err != nil {
        return err
    }

    v3Config := ConfigV3{
        Version:     "3.0",
        Registries:  v2Config.Registries,
        Sinks:       v2Config.Sinks,
        Rulesets:    v2Config.Rulesets,
        Promptsets:  make(map[string]PromptsetConfig),
    }

    return writeConfig(configPath, v3Config)
}
```

### Command Compatibility Layer
```go
// Temporary compatibility for old commands
func legacyInstallCommand() *cobra.Command {
    return &cobra.Command{
        Use:        "install <name>",
        Deprecated: "use 'arm install ruleset <name>' instead",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Fprintf(os.Stderr, "Warning: 'arm install <name>' is deprecated. Use 'arm install ruleset <name>' instead.\n")
            return installResource(ResourceTypeRuleset, args[0], cmd)
        },
    }
}
```

## Testing Strategy

### Unit Tests
- Resource parsing and validation
- Compilation for each tool/resource type combination
- Configuration migration
- Command parsing

### Integration Tests
- End-to-end resource installation
- Multi-resource project scenarios
- Registry discovery across resource types
- Sink management with mixed resource types

### Migration Tests
- V2 to V3 configuration migration
- Backward compatibility scenarios
- Command deprecation warnings
