# Cloudsmith Registry v4 Port

## 1. Introduction

Port the existing v3 Cloudsmith registry implementation to the v4 architecture. Cloudsmith registries use Cloudsmith's package repository service to store and distribute ARM packages as versioned raw packages with semantic versioning.

**Problem**: The v3 Cloudsmith implementation exists but needs to be ported to v4's new architecture (manifest, service layer, registry factory pattern).

## 2. Goals

- Port v3 Cloudsmith registry to v4 architecture
- Maintain feature parity with Git registry (install, update, list, resolve versions)
- Use API key authentication via .armrc
- Support Cloudsmith's native semantic versioning
- Use raw package format for ARM resources
- Pass existing integration tests with minimal changes

## 3. User Stories

### US-001: Add Cloudsmith Registry
**Description:** As a user, I want to add a Cloudsmith registry so that I can install packages from Cloudsmith repositories.

**Acceptance Criteria:**
- [ ] `arm add registry cloudsmith --url <url> --owner <owner> --repo <repo> <name>` adds registry to manifest
- [ ] Registry config stored with type "cloudsmith"
- [ ] Command validates required flags (url, owner, repo, name)
- [ ] Typecheck passes
- [ ] Tests pass

### US-002: Authenticate with Cloudsmith
**Description:** As a user, I want to authenticate using .armrc so that I can access private Cloudsmith repositories.

**Acceptance Criteria:**
- [ ] Registry reads token from .armrc using format `[registry https://api.cloudsmith.io/owner/repo]`
- [ ] Token sent as `Authorization: Token <token>` header
- [ ] Clear error message when token missing or invalid
- [ ] Typecheck passes
- [ ] Tests pass

### US-003: List Package Versions
**Description:** As a user, I want to list available versions so that I can see what versions exist.

**Acceptance Criteria:**
- [ ] `arm list versions <registry>/<package>` shows all semantic versions
- [ ] Versions sorted by semver (newest first)
- [ ] Only shows raw format packages matching package name
- [ ] Handles paginated API responses
- [ ] Typecheck passes
- [ ] Tests pass

### US-004: Resolve Version Constraints
**Description:** As a user, I want to use version constraints so that I can install compatible versions.

**Acceptance Criteria:**
- [ ] Supports exact versions: `package@1.0.0`
- [ ] Supports major constraints: `package@1` (>= 1.0.0, < 2.0.0)
- [ ] Supports minor constraints: `package@1.1` (>= 1.1.0, < 1.2.0)
- [ ] Latest version used when no constraint specified
- [ ] Typecheck passes
- [ ] Tests pass

### US-005: Download Package Content
**Description:** As a user, I want to download package files so that I can install them to sinks.

**Acceptance Criteria:**
- [ ] Downloads raw package files from Cloudsmith CDN
- [ ] Extracts archives (tar.gz, zip) automatically
- [ ] Applies include/exclude patterns to filter files
- [ ] Caches downloaded content
- [ ] Typecheck passes
- [ ] Tests pass

### US-006: Install Packages
**Description:** As a user, I want to install packages so that I can use Cloudsmith-hosted resources.

**Acceptance Criteria:**
- [ ] `arm install ruleset <registry>/<package> <sink>` installs from Cloudsmith
- [ ] `arm install promptset <registry>/<package> <sink>` installs from Cloudsmith
- [ ] Supports version constraints in package spec
- [ ] Updates manifest and lockfile
- [ ] Typecheck passes
- [ ] Tests pass

### US-007: Update Packages
**Description:** As a user, I want to update packages so that I can get newer versions.

**Acceptance Criteria:**
- [ ] `arm update` checks for newer versions matching constraints
- [ ] Shows available updates
- [ ] Updates lockfile with new resolved versions
- [ ] Typecheck passes
- [ ] Tests pass

### US-008: Registry Factory Integration
**Description:** As a developer, I want the registry factory to create Cloudsmith registries so that the service layer works correctly.

**Acceptance Criteria:**
- [ ] Factory creates CloudsmithRegistry when type is "cloudsmith"
- [ ] Registry initialized with config (url, owner, repo)
- [ ] Registry uses shared cache instance
- [ ] Typecheck passes
- [ ] Tests pass

## 4. Functional Requirements

**FR-1:** CLI must accept `--url`, `--owner`, `--repo` flags and NAME positional arg for `arm add registry cloudsmith`

**FR-2:** Manifest must store Cloudsmith registry with fields: name, type="cloudsmith", url, owner, repository

**FR-3:** Registry must load API token from .armrc using key format `[registry <url>/<owner>/<repo>]`

**FR-4:** Registry must call Cloudsmith API v1 endpoints:
- List packages: `/v1/packages/<owner>/<repo>/?query=<name>`
- Download: Use `cdn_url` from package metadata

**FR-5:** Registry must handle paginated responses using Link header

**FR-6:** Registry must filter packages by format="raw" and matching name/filename

**FR-7:** Registry must extract versions from package metadata and sort by semver

**FR-8:** Registry must use semver constraint resolver (not git-style branches/commits)

**FR-9:** Registry must download files using CDN URLs from package metadata

**FR-10:** Registry must extract tar.gz and zip archives automatically

**FR-11:** Registry must apply ContentSelector patterns to filter files

**FR-12:** Registry must cache downloaded content using RegistryPackageCache

**FR-13:** Service layer must call AddCloudsmithRegistry with (name, url, owner, repo, force)

**FR-14:** Factory must create CloudsmithRegistry when manifest type is "cloudsmith"

## 5. Non-Goals

- Publishing packages to Cloudsmith (use Cloudsmith CLI)
- Supporting non-raw package formats (npm, python, etc.)
- Branch or commit-based versioning (Cloudsmith uses semver only)
- Entitlement tokens or OIDC authentication (API key only)
- Offline mode or enhanced caching beyond standard cache
- Custom Cloudsmith API endpoints beyond standard v1

## 6. Technical Considerations

**Architecture:**
- Port v3 implementation to v4 structure: `internal/v4/registry/cloudsmith_registry.go`
- Reuse v3 HTTP client, API models, and pagination logic
- Integrate with v4 service layer and registry factory
- Use v4 manifest and lockfile managers

**Dependencies:**
- Existing v3 code in `internal/v3/registry/cloudsmith_registry.go`
- v4 registry interface and factory pattern
- Shared cache, rcfile, archive extractor packages
- Semver constraint resolver

**Authentication:**
- .armrc format: `[registry https://api.cloudsmith.io/owner/repo]` with `token = <key>`
- HTTP header: `Authorization: Token <token>`
- Error when token missing or API returns 401/403

**API Details:**
- Base URL: `https://api.cloudsmith.io` (configurable via --url)
- Pagination: Link header with `rel="next"`
- Response: JSON array of packages (not wrapped in "results")
- Package fields: slug, name, version, format, filename, cdn_url, size

**File Handling:**
- Single file per package (typical Cloudsmith raw pattern)
- Auto-extract .tar.gz and .zip archives
- Merge extracted files with loose files
- Apply include/exclude patterns after extraction

## 7. Success Metrics

- All v3 Cloudsmith integration tests pass with v4 implementation
- `arm add registry cloudsmith` successfully adds registry to manifest
- `arm install ruleset cloudsmith-reg/package` downloads and installs from Cloudsmith
- `arm update` detects and applies updates for Cloudsmith packages
- Authentication errors provide clear guidance about .armrc configuration
- Performance matches or exceeds v3 implementation
