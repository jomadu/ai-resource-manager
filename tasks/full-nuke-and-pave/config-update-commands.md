# Configuration Update Commands Design

## Overview

Add support for overwriting existing configurations and updating individual fields in registry and sink configurations.

## Requirements

### Registry Configuration Updates

1. **Overwrite existing registries** with `--force` flag
2. **Update individual fields** using dot notation
3. **Type-specific validation** for git registries

### Sink Configuration Updates

1. **Overwrite existing sinks** with `--force` flag
2. **Update individual fields** using dot notation
3. **Array field handling** for directories, include, exclude

## Command Interface

### Registry Commands

```bash
# Overwrite existing registry (requires --force)
arm config registry add ai-rules https://new-url --type git --force

# Update individual fields
arm config registry update ai-rules url https://new-url
arm config registry update ai-rules branches main,develop
arm config registry update ai-rules type git
```

### Sink Commands

```bash
# Overwrite existing sink (requires --force)
arm config sink add q --directories .amazonq/rules --force

# Update individual fields
arm config sink update q directories .amazonq/rules,.amazonq/shared
arm config sink update q include "ai-rules/amazonq-*","ai-rules/shared-*"
arm config sink update q layout flat
```

## Implementation Design

### Manager Interface Extensions

**Manifest Manager**
```go
type Manager interface {
    // Existing methods...
    UpdateGitRegistry(ctx context.Context, name, field, value string) error
}
```

**Config Manager**
```go
type Manager interface {
    // Existing methods...
    UpdateSink(ctx context.Context, name, field, value string) error
}
```

### Field Validation

**Registry Fields (Git Type)**
- `url` - string, must be valid Git URL
- `branches` - comma-separated string, converts to []string
- `type` - string, must be "git"

**Sink Fields**
- `directories` - comma-separated string, converts to []string
- `include` - comma-separated string, converts to []string
- `exclude` - comma-separated string, converts to []string
- `layout` - string, must be "hierarchical" or "flat"

### Command Structure

**Registry Update Command**
```go
var registryUpdateCmd = &cobra.Command{
    Use:   "update <name> <field> <value>",
    Short: "Update registry field",
    Args:  cobra.ExactArgs(3),
    RunE: func(cmd *cobra.Command, args []string) error {
        name, field, value := args[0], args[1], args[2]
        manifestManager := manifest.NewFileManager()
        return manifestManager.UpdateGitRegistry(ctx, name, field, value)
    },
}
```

**Sink Update Command**
```go
var sinkUpdateCmd = &cobra.Command{
    Use:   "update <name> <field> <value>",
    Short: "Update sink field",
    Args:  cobra.ExactArgs(3),
    RunE: func(cmd *cobra.Command, args []string) error {
        name, field, value := args[0], args[1], args[2]
        configManager := config.NewFileManager()
        return configManager.UpdateSink(ctx, name, field, value)
    },
}
```

### Force Flag Implementation

**Registry Add with Force**
```go
func (f *FileManager) AddGitRegistry(ctx context.Context, name string, config registry.GitRegistryConfig, force bool) error {
    manifest, err := f.loadManifest()
    if err != nil {
        manifest = &Manifest{
            Registries: make(map[string]map[string]interface{}),
            Rulesets:   make(map[string]map[string]Entry),
        }
    }

    if _, exists := manifest.Registries[name]; exists && !force {
        return fmt.Errorf("registry %s already exists (use --force to overwrite)", name)
    }

    // Continue with existing logic...
}
```

## Error Handling

### Validation Errors
- Invalid field names return specific error messages
- Type validation errors explain expected format
- Missing registry/sink errors suggest available options

### Safety Checks
- Prevent overwriting without `--force` flag
- Validate field values before persistence
- Atomic updates (load -> validate -> save)

## Testing Strategy

### Unit Tests
- Field validation for all supported types
- Force flag behavior verification
- Error message accuracy

### Integration Tests
- End-to-end command execution
- File persistence verification
- Configuration consistency checks

## Files Modified

### New Commands
- `cmd/arm/config.go` - Add update commands and force flag
- `internal/manifest/manager.go` - Add UpdateGitRegistry method
- `internal/config/manager.go` - Add UpdateSink method

### Enhanced Validation
- Field-specific validation functions
- Type-safe value conversion utilities
- Error message standardization

## Backward Compatibility

- Existing `add` commands work unchanged (no force required for new configs)
- Configuration file formats remain identical
- No breaking changes to existing APIs
