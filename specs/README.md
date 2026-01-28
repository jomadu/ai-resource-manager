# ARM Specifications Index

## Jobs to be Done (JTBDs)

### Core Functionality
1. **Manage AI Resources Across Projects** - Install, update, upgrade, and uninstall AI rulesets and promptsets from remote registries with version tracking and integrity verification
2. **Compile Resources for Multiple Tools** - Transform ARM resources into tool-specific formats (Cursor, Amazon Q, Copilot, Markdown) with hierarchical or flat layouts
3. **Develop Resources Locally** - Compile local ARM resource files to tool formats without registry installation for testing and development
4. **Resolve Rule Conflicts** - Apply priority-based resolution when multiple rulesets define overlapping rules with arm_index.* files
5. **Cache and Optimize Performance** - Store packages locally to avoid redundant downloads, enable offline usage, and provide cleanup mechanisms (by age, access time, or nuke)
6. **Authenticate with Registries** - Securely access private registries using token-based authentication via .armrc files with section matching
7. **Filter Package Contents** - Selectively install files from packages using glob patterns with archive extraction (zip, tar.gz)
8. **Verify Package Integrity** - Ensure downloaded packages match expected content using SHA256 hashing stored in lock file
9. **Test in Isolation** - Support environment variables (ARM_HOME, ARM_CONFIG_PATH, ARM_MANIFEST_PATH) for test isolation and custom configurations
10. **Query Package Information** - List installed packages, check for outdated dependencies, view package details, and list available versions from registries

### Infrastructure & Distribution
11. **Build Cross-Platform Binaries** - Compile ARM for Linux (amd64, arm64), macOS (amd64, arm64), and Windows (amd64) with version metadata injection via LDFLAGS
12. **Automate Releases** - Use semantic-release with conventional commits to automatically version, tag, create GitHub releases, and upload binaries
13. **Install and Uninstall** - Provide shell scripts for easy installation and removal on Linux, macOS, and Windows with platform detection
14. **Ensure Code Quality** - Run linting (13 linters), formatting (gofmt, goimports), and pre-commit hooks with conventional commit validation
15. **Test Continuously** - Execute unit and E2E tests (75 test files, 120 total Go files) on every push and PR with race detection and coverage reporting
16. **Scan for Security Issues** - Run CodeQL analysis on Go code and GitHub Actions, review dependencies on PRs, weekly scheduled scans
17. **Manage Dependencies** - Automated dependency updates via Dependabot for Go modules

## Topics of Concern

### Resource Management
- **Package Installation** - Install rulesets and promptsets from registries to sinks with integrity verification and SHA256 hashing
- **Version Resolution** - Resolve semantic versions, constraints (^, ~), branches, tags, and "latest" with prerelease precedence
- **Dependency Tracking** - Track installed packages in arm.json (manifest) and arm-lock.json (lock file, colocated with manifest)
- **Update Workflows** - Update within constraints vs upgrade to latest (ignoring constraints), with manifest and lock file updates
- **Uninstall Cleanup** - Remove packages from sinks, clean up empty directories recursively, remove arm-index.json and arm_index.* files
- **Query Operations** - List installed packages, check for outdated dependencies, view package details, list available versions from registries

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

### Build & Distribution
- **Cross-Platform Builds** - Makefile targets for building Linux, macOS, Windows binaries with LDFLAGS for version injection
- **Version Metadata** - Inject version, commit hash, build timestamp via LDFLAGS at build time (buildVersion, buildCommit, buildTimestamp)
- **Installation Scripts** - Shell scripts (install.sh, uninstall.sh) with platform detection, GitHub API integration, and PATH configuration
- **Release Automation** - GitHub Actions workflows for semantic-release, binary building, and asset uploading with .tar.gz and SHA256 checksums
- **Artifact Management** - Package binaries as .tar.gz with SHA256 checksums for integrity verification

### Code Quality & CI/CD
- **Linting** - golangci-lint with 13 enabled linters (errcheck, gosimple, govet, ineffassign, staticcheck, typecheck, unused, gofmt, goimports, misspell, gocritic, unconvert, unparam)
- **Formatting** - gofmt and goimports for consistent code style
- **Pre-commit Hooks** - Automated checks for trailing whitespace, end-of-file-fixer, YAML/JSON validation, merge conflicts, large files, conventional commits, go fmt, go imports, go mod tidy, golangci-lint
- **Conventional Commits** - Enforce commit message format (feat, fix, docs, refactor, test, chore) via commitlint workflow on push and PRs
- **Security Scanning** - CodeQL analysis for Go and GitHub Actions (languages: actions, go), dependency review on PRs (dependency-review-action), weekly scheduled scans
- **Test Coverage** - Upload coverage reports to Codecov on every build (75 test files, 120 total Go files)
- **Build Matrix** - Test on multiple platforms (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- **Dependency Management** - Dependabot for Go modules with weekly updates, max 5 open PRs

### Documentation
- **User Documentation** - 12 docs files covering concepts, commands, registries, sinks, storage, resource schemas, publishing guide, migration guide
- **Root Files** - README.md (user-facing), CONTRIBUTING.md (development setup), SECURITY.md (vulnerability reporting), LICENSE.txt (GPL-3.0), go.mod/go.sum (dependencies), .gitignore (exclusions)
- **Configuration Files** - .golangci.yml (linting), .pre-commit-config.yaml (hooks), .releaserc.json (semantic-release), package.json (npm deps), .github/dependabot.yml (dependency updates)

## Specification Documents

### Core Workflows
- [package-installation.md](package-installation.md) - Install, update, upgrade, uninstall workflows
- [version-resolution.md](version-resolution.md) - Semantic versioning, constraints, branch resolution
- [query-operations.md](query-operations.md) - List packages, check outdated, view dependency info, list available versions

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

### Testing & Development
- [constructor-injection.md](constructor-injection.md) - Test isolation via environment variables
- [e2e-testing.md](e2e-testing.md) - End-to-end test specifications

### Infrastructure
- [build-system.md](build-system.md) - Makefile, cross-platform builds, version injection
- [ci-cd-workflows.md](ci-cd-workflows.md) - GitHub Actions for build, test, lint, security, release
- [installation-scripts.md](installation-scripts.md) - Install and uninstall shell scripts
- [code-quality.md](code-quality.md) - Linting, formatting, pre-commit hooks, conventional commits
- [root-files.md](root-files.md) - go.mod, .gitignore, CONTRIBUTING.md, SECURITY.md, LICENSE.txt, package.json, .releaserc.json, dependabot.yml
- [user-documentation.md](user-documentation.md) - README.md, docs/ structure, concepts, commands, registries, sinks, publishing guide, migration guide

## Specification Documents

### Core Workflows
- [package-installation.md](package-installation.md) - Install, update, upgrade, uninstall workflows
- [version-resolution.md](version-resolution.md) - Semantic versioning, constraints, branch resolution
- [query-operations.md](query-operations.md) - List packages, check outdated, view dependency info, list available versions

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

### Testing & Development
- [constructor-injection.md](constructor-injection.md) - Test isolation via environment variables
- [e2e-testing.md](e2e-testing.md) - End-to-end test specifications

### Infrastructure
- [build-system.md](build-system.md) - Makefile, cross-platform builds, version injection
- [ci-cd-workflows.md](ci-cd-workflows.md) - GitHub Actions for build, test, lint, security, release
- [installation-scripts.md](installation-scripts.md) - Install and uninstall shell scripts
- [code-quality.md](code-quality.md) - Linting, formatting, pre-commit hooks, conventional commits
- [root-files.md](root-files.md) - go.mod, .gitignore, CONTRIBUTING.md, SECURITY.md, LICENSE.txt, package.json, .releaserc.json, dependabot.yml
- [user-documentation.md](user-documentation.md) - README.md, docs/ structure, concepts, commands, registries, sinks, publishing guide, migration guide
