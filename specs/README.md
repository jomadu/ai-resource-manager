# ARM Specifications Index

## Jobs to be Done (JTBDs)

1. **Manage AI Resources Across Projects** - Install, update, upgrade, and uninstall AI rulesets and promptsets from remote registries with version tracking
2. **Compile Resources for Multiple Tools** - Transform ARM resources into tool-specific formats (Cursor, Amazon Q, Copilot, Markdown) with proper directory layouts
3. **Develop Resources Locally** - Compile local ARM resource files to tool formats without registry installation for testing and development
4. **Resolve Rule Conflicts** - Apply priority-based resolution when multiple rulesets define overlapping rules
5. **Cache and Optimize Performance** - Store packages locally to avoid redundant downloads, enable offline usage, and provide cleanup mechanisms
6. **Authenticate with Registries** - Securely access private registries using token-based authentication via .armrc files
7. **Filter Package Contents** - Selectively install files from packages using glob patterns with archive extraction
8. **Verify Package Integrity** - Ensure downloaded packages match expected content using SHA256 hashing

## Topics of Concern

### Resource Management
- **Package Installation** - Install rulesets and promptsets from registries to sinks with integrity verification
- **Version Resolution** - Resolve semantic versions, constraints, branches, tags, and "latest"
- **Dependency Tracking** - Track installed packages in arm.json (manifest) and arm-lock.json (lock file)
- **Update Workflows** - Update within constraints vs upgrade to latest (ignoring constraints)
- **Uninstall Cleanup** - Remove packages from sinks and clean up empty directories

### Registry Integration
- **Registry Types** - Git, GitLab Package Registry, Cloudsmith
- **Authentication** - Token-based auth via .armrc files
- **Package Discovery** - List packages and versions from registries

### Compilation & Output
- **Sink Management** - Configure output destinations with tool-specific formats and directory paths
- **Standalone Compilation** - Compile local ARM resource files without registry installation for development and testing
- **Tool Compilation** - Generate Cursor (.mdc with frontmatter), Amazon Q (.md), Copilot (.instructions.md), Markdown (.md)
- **Priority Resolution** - Generate arm_index.* files for conflict resolution with priority ordering
- **Layout Modes** - Hierarchical (preserves structure in arm/ subdirectory) vs Flat (single directory with hash prefixes for Copilot)
- **Index Management** - Track installed packages in arm-index.json for each sink

### Performance & Caching
- **Storage Structure** - Organize cached packages in ~/.arm/storage/registries/{key}/packages/{package}/{version}/
- **Cache Keys** - Generate consistent SHA256-based keys from registry configuration
- **Cache Cleanup** - Remove old versions by age or last access time, or nuke entire cache
- **File Locking** - Prevent concurrent access corruption using OS-level file locks
- **Metadata Tracking** - Track creation time and last access time for each cached version

### Filtering & Patterns
- **Glob Patterns** - Include/exclude files using ** (recursive) and * (single component) wildcards
- **Archive Extraction** - Automatically extract .zip and .tar.gz files with security checks
- **Pattern Precedence** - Exclude overrides include; archives take precedence over loose files
- **Default Patterns** - Use **/*.yml and **/*.yaml if no patterns specified

### Testing & Isolation
- **Environment Variables** - ARM_HOME (custom cache location), ARM_CONFIG_PATH (custom .armrc), ARM_MANIFEST_PATH (custom arm.json)
- **Constructor Injection** - Test-friendly *WithPath and *WithHomeDir constructors
- **Lock File Colocation** - arm-lock.json always colocated with arm.json
- **E2E Testing** - 14 comprehensive test suites covering all workflows with 100% pass rate

## Specification Documents

### Core Workflows
- [package-installation.md](package-installation.md) - Install, update, upgrade, uninstall workflows
- [version-resolution.md](version-resolution.md) - Semantic versioning, constraints, branch resolution

### Registry & Authentication
- [registry-management.md](registry-management.md) - Git, GitLab, Cloudsmith registry types
- [authentication.md](authentication.md) - Token-based authentication via .armrc

### Compilation & Output
- [sink-compilation.md](sink-compilation.md) - Tool-specific compilation and sink management
- [standalone-compilation.md](standalone-compilation.md) - Local file compilation without registry installation
- [priority-resolution.md](priority-resolution.md) - Priority-based rule conflict resolution

### Performance & Filtering
- [cache-management.md](cache-management.md) - Storage structure, cleanup, file locking
- [pattern-filtering.md](pattern-filtering.md) - Glob patterns, archive extraction

### Testing
- [constructor-injection.md](constructor-injection.md) - Test isolation via environment variables
- [e2e-testing.md](e2e-testing.md) - End-to-end test specifications
