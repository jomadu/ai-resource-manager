# Unified Non-Git Registry Support

## Implementation Specification

### Overview
Support mixed-content packages in non-Git registries (GitLab, HTTP) containing URF files, pre-built rules, and compressed archives. Use extract-and-merge strategy to flatten package content before applying include patterns.

### Architecture Decisions

**Package Structure**: Flexible parsing - handle mixed content types within single versioned package
**Ruleset Identification**: Package-level naming (`registry/package-name`) with include patterns for selective installation
**Archive Handling**: Extract and merge - unpack archives into package namespace for uniform pattern matching
**Pattern Resolution**: Post-extraction - extract all archives first, then apply include patterns to flattened content
**Version Conflicts**: Package version wins over individual URF metadata versions

### Extraction Engine Specifications

**Archive Detection**: Extension-based (`.tar.gz`, `.zip`)
**Extraction Method**: Streaming extraction to prevent memory issues
**Cleanup Strategy**: Immediate cleanup after pattern matching
**Path Handling**: Preserve directory structure from archives
**Error Handling**: Fail fast on extraction errors
**Size Limits**: No limits imposed
**Concurrency**: Sequential processing

### Processing Flow

1. **Download Package**: Fetch versioned package from registry
2. **Identify Archives**: Scan for `.tar.gz` and `.zip` files by extension
3. **Extract Archives**: Stream extract to temporary directory, preserve paths
4. **Merge Content**: Combine loose files + extracted files (archives win on conflicts)
5. **Apply Patterns**: Run include patterns against merged content
6. **Install Files**: Copy matched files to sinks
7. **Cleanup**: Remove temporary extraction directories

### Test Cases

**Archive Contents**:
```
ruleset-zip-1.tar.gz:
├── ruleset-zip-1.yml
└── build/
    ├── amazonq/ruleset-zip-1_rule-1.md
    └── cursor/ruleset-zip-1_rule-1.mdc

ruleset-zip-2.tar.gz:
├── ruleset-zip-2.yml
└── build/
    ├── amazonq/ruleset-zip-2_rule-1.md
    └── cursor/ruleset-zip-2_rule-1.mdc
```

**Package Scenarios**:
1. **Single Archive**: `ruleset-zip-1.tar.gz`
2. **Multiple Archives**: `ruleset-zip-1.tar.gz`, `ruleset-zip-2.tar.gz`
3. **Single URF**: `ruleset-a.yml`
4. **Multiple URF**: `ruleset-a.yml`, `ruleset-b.yml`
5. **URF + Builds**: `ruleset-a.yml` + `build/amazonq/`, `build/cursor/`
6. **Mixed Content**: Archives + URF + Builds

### Example Transformations

**Input Package**:
```
ai-rules@1.0.0/
├── ruleset-zip-1.tar.gz
├── ruleset-a.yml
└── build/cursor/existing.mdc
```

**After Extract-and-Merge**:
```
merged-content/
├── ruleset-zip-1.yml          # from archive
├── build/amazonq/ruleset-zip-1_rule-1.md  # from archive
├── build/cursor/ruleset-zip-1_rule-1.mdc  # from archive
├── ruleset-a.yml              # loose file
└── build/cursor/existing.mdc  # loose file
```

**Pattern Examples**:
- Default URF: `*.yml` → `ruleset-zip-1.yml`, `ruleset-a.yml`
- Pre-built Cursor: `build/cursor/**` → `ruleset-zip-1_rule-1.mdc`, `existing.mdc`
- Specific file: `ruleset-a.yml` → `ruleset-a.yml`

### Implementation Notes

- Archives are processed sequentially to avoid I/O contention
- Extraction failures abort entire installation to ensure integrity
- Directory structure from archives is preserved in merged content
- No metadata tracking of file origins - treat all as package content post-merge
- Compatible with existing GitLab registry authentication and versioning
