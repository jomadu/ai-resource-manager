# Registry Type-Specific Configuration Design

## Problem
Registry config is currently flat with only URL and type fields. Need type-specific configuration (e.g., git registries need branch specifications) while maintaining type safety and extensibility.

## Solution

### 1. Config Structure - Embedded Structs
```go
// Base registry config
type RegistryConfig struct {
    URL  string `json:"url"`
    Type string `json:"type"`
}

// Git-specific config
type GitRegistryConfig struct {
    RegistryConfig
    Branches []string `json:"branches,omitempty"`
}
```

### 2. JSON Serialization - Flat Structure
```json
{
  "url": "https://github.com/user/repo",
  "type": "git",
  "branches": ["main", "develop"]
}
```

### 3. Factory Pattern - Type Switch
```go
func CreateRegistry(rawConfig map[string]interface{}) (Registry, error) {
    // Parse base config first
    // Switch on type field
    // Parse type-specific config
    // Return appropriate registry
}
```

### 4. Parsing Location - Factory
- Manifest manager returns raw `map[string]interface{}`
- Factory handles all type-specific parsing and validation
- Keeps parsing logic with registry creation

## Implementation Steps

1. Create shared base `RegistryConfig` type
2. Create `GitRegistryConfig` with embedded base
3. Update factory to parse type-specific configs
4. Update manifest manager to return raw config maps
5. Remove duplicate `RegistryConfig` definitions

## Files to Modify

- `internal/config/types.go` - Remove duplicate, use shared type
- `internal/manifest/types.go` - Use shared base type
- `internal/registry/factory.go` - Add type-specific parsing
- `internal/registry/git_registry.go` - Accept `GitRegistryConfig`
- `cmd/arm/config.go` - Handle additional git flags

## Backward Compatibility
No backward compatibility required - breaking change acceptable.
