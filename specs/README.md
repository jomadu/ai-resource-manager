# ARM Specifications Index

## Jobs to be Done (JTBDs)

1. **Manage AI Resources Across Projects** - Install, update, upgrade, and uninstall AI rulesets and promptsets from remote registries with version tracking and integrity verification
2. **Compile Resources for Multiple Tools** - Transform ARM resources into tool-specific formats (Cursor, Amazon Q, Copilot, Markdown) with hierarchical or flat layouts
3. **Develop Resources Locally** - Compile local ARM resource files to tool formats without registry installation for testing and development
4. **Resolve Rule Conflicts** - Apply priority-based resolution when multiple rulesets define overlapping rules with arm_index.* files
5. **Cache and Optimize Performance** - Store packages locally to avoid redundant downloads, enable offline usage, and provide cleanup mechanisms (by age, access time, or nuke)
6. **Authenticate with Registries** - Securely access private registries using token-based authentication via .armrc files with section matching
7. **Filter Package Contents** - Selectively install files from packages using glob patterns with archive extraction (zip, tar.gz)
8. **Verify Package Integrity** - Ensure downloaded packages match expected content using SHA256 hashing stored in lock file
9. **Test in Isolation** - Support environment variables (ARM_HOME, ARM_CONFIG_PATH, ARM_MANIFEST_PATH) for test isolation and custom configurations

## Topics of Concern

### Resource Management
- **Package Installation** - Install rulesets and promptsets from registries to sinks with integrity verification and SHA256 hashing
- **Version Resolution** - Resolve semantic versions, constraints (^, ~), branches, tags, and "latest" with prerelease precedence
- **Dependency Tracking** - Track installed packages in arm.json (manifest) and arm-lock.json (lock file, colocated with manifest)
- **Update Workflows** - Update within constraints vs upgrade to latest (ignoring constraints), with manifest and lock file updates
- **Uninstall Cleanup** - Remove packages from sinks, clean up empty directories recursively, remove arm-index.json and arm_index.* files

### Registry Integration
- **Registry Types** - Git (GitHub, GitLab, Git remotes), GitLab Package Registry, Cloudsmith
- **Authentication** - Token-based auth via .armrc files with section matching ([registry.name] or [registry.*])
- **Package Discovery** - List packages and versions from registries with pagination support

### Compilation & Output
- **Sink Management** - Configure output destinations with tool-specific formats and directory paths
- **Standalone Compilation** - Compile local ARM resource files without registry installation for development and testing
- **Tool Compilation** - Generate Cursor (.mdc with frontmatter), Amazon Q (.md), Copilot (.instructions.md), Markdown (.md)
- **Priority Resolution** - Generate arm_index.* files for conflict resolution with priority ordering (higher priority wins)
- **Layout Modes** - Hierarchical (preserves structure in arm/ subdirectory) vs Flat (single directory with hash prefixes for Copilot)
- **Index Management** - Track installed packages in arm-index.json for each sink, remove when all packages uninstalled

### Performance & Caching
- **Storage Structure** - Organize cached packages in ~/.arm/storage/registries/{key}/packages/{package}/{version}/
- **Cache Keys** - Generate consistent SHA256-based keys from registry configuration (normalized, order-independent)
- **Cache Cleanup** - Remove old versions by age (--max-age), last access time, or nuke entire cache
- **File Locking** - Prevent concurrent access corruption using OS-level file locks with context cancellation
- **Metadata Tracking** - Track creation time and last access time for each cached version

### Filtering & Patterns
- **Glob Patterns** - Include/exclude files using ** (recursive) and * (single component) wildcards
- **Archive Extraction** - Automatically extract .zip and .tar.gz files with security checks (path traversal prevention)
- **Pattern Precedence** - Exclude overrides include; archives take precedence over loose files
- **Default Patterns** - Use **/*.yml and **/*.yaml if no patterns specified

### Testing & Isolation
- **Environment Variables** - ARM_HOME (custom cache location), ARM_CONFIG_PATH (custom .armrc), ARM_MANIFEST_PATH (custom arm.json)
- **Constructor Injection** - Test-friendly *WithPath and *WithHomeDir constructors for all components
- **Lock File Colocation** - arm-lock.json always colocated with arm.json (derived from manifest path)
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
