# Resource Manager Design Tradeoffs

## Key Decision Points

### 1. Command Interface: Explicit vs Implicit Resource Types

#### Option A: Explicit Types (Chosen)
```bash
arm install ruleset my-reg/clean-code
arm install promptset my-reg/code-review
```

**Pros:**
- Clear intent and type safety
- Extensible to new resource types
- No ambiguity in commands
- Better error messages

**Cons:**
- More verbose commands
- Breaking change from current interface
- Higher cognitive load

#### Option B: Implicit Detection
```bash
arm install my-reg/clean-code  # Auto-detects as ruleset
arm install my-reg/code-review # Auto-detects as promptset
```

**Pros:**
- Backward compatible
- Shorter commands
- Familiar interface

**Cons:**
- Requires registry lookup for type detection
- Ambiguous error cases
- Harder to extend
- Performance overhead

**Decision:** Option A for clarity and extensibility

### 2. Schema Design: Separate vs Unified

#### Option A: Separate Schemas (Chosen)
```yaml
# ruleset.yml
version: "1.0"
metadata: { ... }
rules: { ... }

# promptset.yml
version: "1.0"
metadata: { ... }
prompts: { ... }
```

**Pros:**
- Type-specific validation
- Clear separation of concerns
- Easier to extend each type independently
- No conditional logic in parsing

**Cons:**
- Schema duplication in metadata
- More complex tooling

#### Option B: Unified Schema with Discriminator
```yaml
version: "1.0"
type: "ruleset"  # or "promptset"
metadata: { ... }
rules: { ... }     # Only if type=ruleset
prompts: { ... }   # Only if type=promptset
```

**Pros:**
- Single schema to maintain
- Explicit type declaration
- Unified parsing logic

**Cons:**
- Conditional validation complexity
- Risk of mixed content
- Less type safety

**Decision:** Option A for better type safety and maintainability

### 3. Registry Organization: Flat vs Hierarchical

#### Option A: Hierarchical (Chosen)
```
registry/
├── rulesets/
│   ├── clean-code.yml
│   └── security.yml
└── promptsets/
    ├── code-review.yml
    └── documentation.yml
```

**Pros:**
- Clear organization by type
- Easy to browse and understand
- Natural discovery patterns
- Scales well with many resources

**Cons:**
- More directory structure
- Requires type-aware discovery

#### Option B: Flat with Naming Convention
```
registry/
├── ruleset-clean-code.yml
├── ruleset-security.yml
├── promptset-code-review.yml
└── promptset-documentation.yml
```

**Pros:**
- Simple flat structure
- Type visible in filename
- Easy to list all resources

**Cons:**
- Naming convention enforcement
- Harder to browse by type
- Filename pollution

**Decision:** Option A for better organization and discoverability

### 4. Sink Specialization: Shared vs Dedicated

#### Option A: Dedicated Sinks (Chosen)
```json
{
  "sinks": {
    "cursor-rules": {
      "directory": ".cursor/rules",
      "type": "cursor",
      "resourceTypes": ["ruleset"]
    },
    "cursor-prompts": {
      "directory": ".cursor/prompts",
      "type": "cursor",
      "resourceTypes": ["promptset"]
    }
  }
}
```

**Pros:**
- Clear separation of resource types
- Tool-specific optimizations
- No mixing of rules and prompts
- Better organization

**Cons:**
- More sink configuration
- Potential confusion about sink types

#### Option B: Shared Sinks with Subdirectories
```json
{
  "sinks": {
    "cursor": {
      "directory": ".cursor",
      "type": "cursor",
      "resourceTypes": ["ruleset", "promptset"]
    }
  }
}
```

**Pros:**
- Fewer sink definitions
- Automatic subdirectory organization
- Simpler configuration

**Cons:**
- Less control over organization
- Potential conflicts between types
- Tool limitations on mixed content

**Decision:** Option A for better control and tool compatibility

### 5. Compilation Strategy: Per-Resource vs Batch

#### Option A: Per-Resource Compilation (Chosen)
- Each resource compiles independently
- Clear ownership of compiled files
- Easier debugging and troubleshooting

#### Option B: Batch Compilation
- All resources of a type compile together
- Potential for cross-resource optimization
- More complex dependency management

**Decision:** Option A for simplicity and maintainability

## Implementation Complexity Analysis

### Low Complexity
- Schema validation (separate schemas are simpler)
- Resource type detection (explicit in schema)
- Command parsing (clear resource type prefix)

### Medium Complexity
- Configuration migration (additive changes)
- Sink management (new sink types)
- Registry discovery (type-aware patterns)

### High Complexity
- Unified resource interface (common operations across types)
- Cache system refactoring (type-aware caching)
- Compilation pipeline (type-specific compilers)

## Migration Strategy

### Backward Compatibility
1. **Phase 1**: Support both old and new command formats
2. **Phase 2**: Deprecation warnings for old format
3. **Phase 3**: Remove old format support

### Configuration Evolution
```json
// v2.x (current)
{
  "version": "2.0",
  "rulesets": { ... }
}

// v3.0 (new)
{
  "version": "3.0",
  "rulesets": { ... },    // Migrated
  "promptsets": { ... }   // New
}
```

### User Communication
- Clear migration guide with examples
- Automated migration tooling
- Version-specific documentation
