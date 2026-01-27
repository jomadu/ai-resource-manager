# ARM Specifications Index

## Jobs to be Done (JTBDs)

1. **Manage AI Resources Across Projects** - Install, update, and version AI rulesets and promptsets from remote registries
2. **Compile Resources for Multiple Tools** - Transform ARM resources into tool-specific formats (Cursor, Amazon Q, Copilot, Markdown)
3. **Resolve Rule Conflicts** - Apply priority-based resolution when multiple rulesets define overlapping rules
4. **Cache and Optimize Performance** - Store packages locally to avoid redundant downloads and enable offline usage
5. **Authenticate with Registries** - Securely access private registries using token-based authentication
6. **Filter Package Contents** - Selectively install files from packages using glob patterns

## Topics of Concern

### Resource Management
- **Package Installation** - Install rulesets and promptsets from registries to sinks
- **Version Resolution** - Resolve semantic versions, constraints, branches, and tags
- **Dependency Tracking** - Track installed packages in manifest and lock files
- **Update Workflows** - Update within constraints vs upgrade to latest

### Registry Integration
- **Registry Types** - Git, GitLab Package Registry, Cloudsmith
- **Authentication** - Token-based auth via .armrc files
- **Package Discovery** - List packages and versions from registries

### Compilation & Output
- **Sink Management** - Configure output destinations with tool-specific formats
- **Tool Compilation** - Generate Cursor (.mdc), Amazon Q (.md), Copilot (.instructions.md), Markdown (.md)
- **Priority Resolution** - Generate index files for conflict resolution
- **Layout Modes** - Hierarchical (preserves structure) vs Flat (single directory)

### Performance & Caching
- **Storage Structure** - Organize cached packages by registry and version
- **Cache Cleanup** - Remove old versions by age or access time
- **File Locking** - Prevent concurrent access corruption

### Filtering & Patterns
- **Glob Patterns** - Include/exclude files using ** wildcards
- **Archive Extraction** - Automatically extract .zip and .tar.gz files
- **Pattern Precedence** - Exclude overrides include

### Testing & Isolation
- **Environment Variables** - ARM_HOME, ARM_CONFIG_PATH, ARM_MANIFEST_PATH
- **Constructor Injection** - Test-friendly constructors accepting paths
- **E2E Testing** - Comprehensive end-to-end test coverage

## Specification Documents

### Core Workflows
- [package-installation.md](package-installation.md) - Install, update, upgrade, uninstall workflows
- [version-resolution.md](version-resolution.md) - Semantic versioning, constraints, branch resolution

### Registry & Authentication
- [registry-management.md](registry-management.md) - Git, GitLab, Cloudsmith registry types
- [authentication.md](authentication.md) - Token-based authentication via .armrc

### Compilation & Output
- [sink-compilation.md](sink-compilation.md) - Tool-specific compilation and sink management
- [priority-resolution.md](priority-resolution.md) - Priority-based rule conflict resolution

### Performance & Filtering
- [cache-management.md](cache-management.md) - Storage structure, cleanup, file locking
- [pattern-filtering.md](pattern-filtering.md) - Glob patterns, archive extraction

### Testing
- [constructor-injection.md](constructor-injection.md) - Test isolation via environment variables
- [e2e-testing.md](e2e-testing.md) - End-to-end test specifications
