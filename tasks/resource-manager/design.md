# AI Resource Manager (ARM) - Resource Extension Design

## Overview

Extend AI Rules Manager to AI Resource Manager, supporting both rulesets and promptsets as first-class resources while maintaining the same package management approach.

## Command Changes

### Current vs New Commands

| Current | New |
|---------|-----|
| `arm install <ruleset>` | `arm install ruleset <ruleset>` |
| N/A | `arm install promptset <promptset>` |
| `arm info [ruleset]` | `arm info [ruleset\|promptset <name>]` |
| N/A | `arm info promptset <promptset>` |
| `arm list` | `arm list` (shows both, separate sections) |
| `arm outdated` | `arm outdated` (shows both, separate sections) |
| `arm uninstall <ruleset>` | `arm uninstall ruleset <ruleset>` |
| N/A | `arm uninstall promptset <promptset>` |
| `arm update [ruleset]` | `arm update [ruleset\|promptset <name>]` |

### Command Specifications

#### Install (Breaking Change)
```bash
arm install ruleset my-reg/my-ruleset --sinks cursor-rules --include "**/*.yml"
arm install promptset my-reg/my-promptset --sinks cursor-prompts --include "**/*.yml"
```

#### Info (Enhanced)
```bash
arm info                            # shows all installed resources
arm info ruleset my-reg/my-ruleset  # shows specific ruleset
arm info promptset my-reg/my-promptset # shows specific promptset
```

#### List (Enhanced)
```bash
arm list
```
Displays all installed resources organized into separate sections for rulesets and promptsets. Each resource shows its registry/name and installed version. May include additional metadata like installation date, sink assignments, or status indicators.

#### Outdated (Enhanced)
```bash
arm outdated
```
Shows resources with available updates, grouped by type. Displays current and available versions with clear upgrade paths. May include release notes summaries or breaking change warnings.

When called without arguments, displays summary information for all installed resources grouped by type. When called with a specific resource, provides detailed information including metadata, version history, dependencies, sink assignments, file contents, and registry information. May show compilation status for URF resources.

#### Uninstall (Breaking Change)
```bash
arm uninstall ruleset my-ruleset
arm uninstall promptset my-promptset
```

#### Update (Enhanced)
```bash
arm update                          # updates all outdated resources
arm update ruleset my-ruleset       # updates specific ruleset
arm update promptset my-promptset   # updates specific promptset
```

## File Structure Changes

### Configuration Files

#### arm.json (Enhanced)
```json
{
  "version": "3.0",
  "registries": { ... },
  "sinks": { ... },
  "rulesets": {
    "my-reg/clean-code": {
      "version": "^1.0.0",
      "sinks": ["cursor-rules"],
      "priority": 100
    }
  },
  "promptsets": {
    "my-reg/code-review": {
      "version": "^1.2.0",
      "sinks": ["cursor-prompts"]
    }
  }
}
```

#### arm-lock.json (Enhanced)
```json
{
  "rulesets": {
    "my-reg": {
      "clean-code": {
        "version": "v1.0.0",
        "display": "v1.0.0",
        "checksum": "sha256:..."
      }
    }
  },
  "promptsets": {
    "my-reg": {
      "code-review": {
        "version": "v1.2.0",
        "display": "v1.2.0",
        "checksum": "sha256:..."
      }
    }
  }
}
```

### Directory Structure

```
.cursor/
├── rules/           # Rulesets only
│   └── arm/
│       ├── my-reg/
│       ├── arm_index.mdc
│       └── arm-index.json
└── commands/        # Promptsets only
    └── arm/
        ├── my-reg/
        └── arm-index.json  # No priority ordering
```

## URF Extension for Prompts

### Prompt URF Format
```yaml
version: "1.0"
metadata:
  id: "code-review-prompts"
  name: "Code Review Prompt Set"
  description: "Prompts for code review tasks"
prompts:
  security-review:
    name: "Security Review"
    description: "Review code for security vulnerabilities"
    body: |
      Review this code for security vulnerabilities.
      Focus on input validation, authentication, and data exposure.
  performance-review:
    name: "Performance Review"
    description: "Review code for performance issues"
    body: |
      Analyze this code for performance bottlenecks.
      Look for inefficient algorithms and resource usage.
```

### Compilation Targets

#### Cursor (.md)
```markdown
---
namespace: my-reg
metadata:
  id: code-review-prompts
  name: Code Review Prompt Set
  version: 1.0.0
  prompts:
    - security-review
    - performance-review
prompt:
  id: security-review
  name: Security Review
---

# Security Review

Review this code for security vulnerabilities.
Focus on input validation, authentication, and data exposure.
```

#### Amazon Q (.md)
```markdown
---
namespace: my-reg
metadata:
  id: code-review-prompts
  name: Code Review Prompt Set
  version: 1.0.0
  prompts:
    - security-review
    - performance-review
prompt:
  id: security-review
  name: Security Review
---

# Security Review

Review this code for security vulnerabilities.
Focus on input validation, authentication, and data exposure.
```

## Implementation Plan

### Phase 1: Ruleset Migration
1. Update command parsing to require `ruleset` keyword
2. Update configuration file structure
3. Maintain backward compatibility warnings
4. Update documentation

### Phase 2: Promptset Support
1. Add promptset command handling
2. Extend URF compiler for prompts
3. Add promptset installation logic
4. Update list/outdated commands for dual display
5. Update metadata block compilation to use `metadata` section header instead of `ruleset`
6. Refactor cache system for unified resource support:
   - Rename `rulesets/` directory to `resources/` with type subdirectories
   - Generalize `RegistryIndex` to handle resource types (rulesets, promptsets)
   - Create generic `RegistryResourceCache` interface replacing `RegistryRulesetCache`
   - Update cache key generation to include resource type

### Phase 3: Polish
1. Enhanced error messages
2. Migration tooling
3. Updated examples and documentation

## Resource Type Design Deep Dive

### Core Principles

1. **Type Safety**: Each resource type has a distinct schema and purpose
2. **Extensibility**: New resource types can be added without breaking existing ones
3. **Tool Specialization**: Each AI tool can handle resource types differently
4. **Clear Semantics**: Commands explicitly state what type of resource they're operating on

### Resource Type Characteristics

| Aspect | Rulesets | Promptsets |
|--------|----------|------------|
| **Purpose** | Modify AI behavior | Provide task templates |
| **Priority** | Yes (conflict resolution) | No (independent prompts) |
| **Compilation** | Tool-specific formats | Tool-specific formats |
| **Sink Types** | `cursor`, `amazonq`, `copilot` | `cursor-prompts`, `amazonq-prompts` |
| **Index File** | `arm_index.*` with priorities | `arm-index.json` metadata only |
| **Versioning** | Semantic versioning | Semantic versioning |

### Schema Evolution Strategy

#### Base Resource Schema
```yaml
version: "1.0"  # URF version
metadata:
  id: string      # Unique identifier
  name: string    # Human-readable name
  description?: string
  version: string # Resource version
  tags?: string[]
  author?: string
  license?: string
```

#### Ruleset Schema Extension
```yaml
# Extends base schema
rules:
  rule-id:        # Map key as rule ID
    name: string
    priority?: number
    enforcement?: "must" | "should" | "may"
    body: string
    tags?: string[]
```

#### Promptset Schema Extension
```yaml
# Extends base schema
prompts:
  prompt-id:      # Map key as prompt ID
    name: string
    description?: string
    body: string
    parameters?: object  # Future: parameterized prompts
    tags?: string[]
```

### Registry Organization

#### Mixed Resource Repositories
```
my-ai-resources/
├── rulesets/
│   ├── clean-code.yml
│   └── security.yml
├── promptsets/
│   ├── code-review.yml
│   └── documentation.yml
└── build/           # Pre-compiled resources
    ├── cursor/
    │   ├── rules/
    │   └── prompts/
    └── amazonq/
        ├── rules/
        └── prompts/
```

#### Resource Discovery
- Default patterns by type:
  - Rulesets: `rulesets/**/*.yml`, `rulesets/**/*.yaml`
  - Promptsets: `promptsets/**/*.yml`, `promptsets/**/*.yaml`
- Explicit patterns: `--include "custom-path/**/*.yml"`
- Type inference from schema content

### Command Interface Refinements

#### Resource-Specific Options
```bash
# Ruleset-specific options
arm install ruleset my-reg/clean-code --priority 100 --enforcement must

# Promptset-specific options (future)
arm install promptset my-reg/code-review --parameters '{"language": "go"}'
```

#### Bulk Operations
```bash
# Install all resources of a type
arm install ruleset my-reg/* --sinks cursor-rules
arm install promptset my-reg/* --sinks cursor-prompts

# Mixed operations
arm update  # Updates all resource types
arm outdated --type ruleset  # Filter by type
```

### Sink Type Specialization

#### Ruleset Sinks
- `cursor` → `.cursor/rules/`
- `amazonq` → `.amazonq/rules/`
- `copilot` → `.github/copilot/`

#### Promptset Sinks
- `cursor-prompts` → `.cursor/prompts/` or `.cursor/commands/`
- `amazonq-prompts` → `.amazonq/prompts/`
- `copilot-prompts` → `.github/copilot/prompts/`

### Implementation Challenges

#### 1. Resource Type Detection
```go
type ResourceType string

const (
    ResourceTypeRuleset   ResourceType = "ruleset"
    ResourceTypePromptset ResourceType = "promptset"
)

func DetectResourceType(content []byte) (ResourceType, error) {
    var base struct {
        Rules   map[string]interface{} `yaml:"rules"`
        Prompts map[string]interface{} `yaml:"prompts"`
    }

    if err := yaml.Unmarshal(content, &base); err != nil {
        return "", err
    }

    hasRules := len(base.Rules) > 0
    hasPrompts := len(base.Prompts) > 0

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

#### 2. Unified Resource Interface
```go
type Resource interface {
    GetType() ResourceType
    GetMetadata() Metadata
    Compile(tool string) ([]CompiledFile, error)
    Validate() error
}

type Ruleset struct {
    Metadata Metadata            `yaml:"metadata"`
    Rules    map[string]Rule     `yaml:"rules"`
}

type Promptset struct {
    Metadata Metadata            `yaml:"metadata"`
    Prompts  map[string]Prompt   `yaml:"prompts"`
}
```

#### 3. Configuration Migration
```go
// Migrate from v2 to v3 configuration
func MigrateConfig(v2Config *ConfigV2) *ConfigV3 {
    v3Config := &ConfigV3{
        Version:    "3.0",
        Registries: v2Config.Registries,
        Sinks:      v2Config.Sinks,
        Rulesets:   v2Config.Rulesets,  // Direct copy
        Promptsets: make(map[string]PromptsetConfig),  // Empty initially
    }
    return v3Config
}
```

### Future Extensions

#### Additional Resource Types
- **Templates**: Code scaffolding and boilerplate
- **Workflows**: Multi-step AI processes
- **Contexts**: Domain-specific knowledge bases
- **Agents**: Complete AI assistant configurations

#### Resource Dependencies
```yaml
metadata:
  id: "advanced-security"
  dependencies:
    - type: ruleset
      name: "my-reg/basic-security"
      version: "^1.0.0"
```

## Breaking Changes

### Command Interface
- All install/info/uninstall commands now require resource type
- `arm install <name>` becomes `arm install ruleset <name>`

### Configuration Format
- `arm.json` version bumped to "3.0"
- Rulesets moved under `rulesets` key
- New `promptsets` key added

### Migration Path
1. **Phase 1**: Add deprecation warnings for old commands
2. **Phase 2**: Support both old and new formats
3. **Phase 3**: Remove backward compatibility
