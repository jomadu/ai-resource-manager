# Root Files & Configuration

## Job to be Done
Provide essential project metadata, configuration, and documentation files that define the project structure, dependencies, and development guidelines.

## Activities
1. Define Go module and dependencies
2. Specify files to exclude from version control
3. Document contribution guidelines
4. Define security policy
5. Specify license terms
6. Configure npm dependencies for semantic-release
7. Configure semantic-release behavior
8. Configure Dependabot for automated dependency updates

## Acceptance Criteria
- [x] Define Go module with go.mod (module path, Go version 1.24.5)
- [x] Lock dependencies with go.sum (transitive dependencies)
- [x] Exclude build artifacts, IDE files, OS files, temp files from git (.gitignore)
- [x] Document development setup, code quality, commit format, workflow (CONTRIBUTING.md)
- [x] Define supported versions and vulnerability reporting process (SECURITY.md)
- [x] Specify GPL-3.0 license terms (LICENSE.txt)
- [x] Define npm dependencies for semantic-release (package.json)
- [x] Configure semantic-release branches and plugins (.releaserc.json)
- [x] Configure Dependabot for weekly Go module updates (.github/dependabot.yml)

## Data Structures

### go.mod
```go
module github.com/jomadu/ai-resource-manager

go 1.24.5

require (
    github.com/go-playground/validator/v10 v10.29.0
    github.com/stretchr/testify v1.8.4
    gopkg.in/ini.v1 v1.67.0
    gopkg.in/yaml.v3 v3.0.1
)
```

### .gitignore Categories
```
# Binaries: /bin, /arm, *.exe, /dist
# Test coverage: coverage.out, *.cover
# IDE: .vscode/, .idea/, *.swp, *.swo, *~
# OS: .DS_Store, Thumbs.db
# Temporary: *.tmp, *.temp, *.log
# Python: .venv/, __pycache__/, *.pyc
# Node: node_modules
# Testing: sandbox/, test-sandbox/
# Secrets: .env, .armrc
# Ralph: /prd.json, /progress.txt, /.last-branch, /archive/
```

### package.json
```json
{
  "name": "arm",
  "private": true,
  "devDependencies": {
    "semantic-release": "^22.0.0",
    "@semantic-release/github": "^9.0.0",
    "@semantic-release/exec": "^6.0.0",
    "conventional-changelog-conventionalcommits": "^7.0.0"
  }
}
```

### .releaserc.json
```json
{
  "branches": [
    "main",
    {"name": "rc", "prerelease": true}
  ],
  "plugins": [
    ["@semantic-release/commit-analyzer", {"preset": "conventionalcommits"}],
    "@semantic-release/release-notes-generator",
    ["@semantic-release/exec", {"prepareCmd": "sed -i 's/var version = \".*\"/var version = \"${nextRelease.version}\"/g' cmd/arm/main.go"}],
    ["@semantic-release/github", {"assets": []}]
  ]
}
```

### .github/dependabot.yml
```yaml
version: 2
updates:
  - package-ecosystem: "gomod"
    directory: "/"
    schedule:
      interval: "weekly"
    open-pull-requests-limit: 5
```

## Algorithm

### Dependency Management
1. **Add dependency:**
   - Run `go get <package>@<version>`
   - go.mod updated with direct dependency
   - go.sum updated with checksums

2. **Update dependencies:**
   - Run `go get -u ./...` (update all)
   - Or `go get -u <package>` (update specific)
   - go.mod and go.sum updated

3. **Tidy dependencies:**
   - Run `go mod tidy`
   - Remove unused dependencies
   - Add missing dependencies
   - Update go.sum

### Semantic Release
1. **Analyze commits:**
   - Parse commits since last release
   - Determine version bump (major, minor, patch)
   - Generate changelog

2. **Create release:**
   - Update version in code (via exec plugin)
   - Create git tag
   - Create GitHub release with changelog
   - Upload assets (handled by release workflow)

### Dependabot Updates
1. **Weekly check:**
   - Scan go.mod for outdated dependencies
   - Create PR for each update (max 5 open)
   - Run CI tests on PR
   - Auto-merge if tests pass (requires configuration)

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Incompatible Go version | go.mod specifies minimum version, build fails if too old |
| Missing dependency | go.sum verification fails, build fails |
| Corrupted go.sum | go mod verify fails, must run go mod tidy |
| Breaking change in dependency | Dependabot PR fails tests, manual intervention required |
| No commits since last release | Semantic-release skips release |
| Multiple commit types | Semantic-release uses highest version bump (major > minor > patch) |

## Dependencies

- Go 1.24.5 or later
- npm (for semantic-release)
- Git (for version control)

## Implementation Mapping

**Source files:**
- `go.mod` - Go module definition
- `go.sum` - Dependency checksums
- `.gitignore` - Git exclusions
- `CONTRIBUTING.md` - Contribution guidelines
- `SECURITY.md` - Security policy
- `LICENSE.txt` - GPL-3.0 license
- `package.json` - npm dependencies
- `.releaserc.json` - Semantic-release config
- `.github/dependabot.yml` - Dependabot config

**Related specs:**
- `build-system.md` - Uses go.mod for building
- `ci-cd-workflows.md` - Uses semantic-release for versioning
- `code-quality.md` - Uses CONTRIBUTING.md guidelines

## Examples

### Example 1: Add New Dependency

**Input:**
```bash
go get gopkg.in/yaml.v3
```

**Expected Output:**
```
go: downloading gopkg.in/yaml.v3 v3.0.1
go: added gopkg.in/yaml.v3 v3.0.1
```

**Verification:**
- go.mod contains `gopkg.in/yaml.v3 v3.0.1` in require block
- go.sum contains checksums for yaml.v3

### Example 2: Tidy Dependencies

**Input:**
```bash
go mod tidy
```

**Expected Output:**
```
# No output if already tidy
# Or removes unused dependencies
```

**Verification:**
- go.mod only contains used dependencies
- go.sum matches go.mod

### Example 3: Semantic Release (feat commit)

**Input:**
```bash
git commit -m "feat: add new feature"
git push origin main
```

**Expected Output:**
```
semantic-release output:
- Analyzing commits
- Determined version bump: minor (0.1.0 → 0.2.0)
- Creating tag v0.2.0
- Creating GitHub release
- Uploading assets
```

**Verification:**
- Tag v0.2.0 exists
- GitHub release created with changelog
- Version in code updated

### Example 4: Dependabot Update

**Input:**
```
# Automatic weekly check
```

**Expected Output:**
```
Dependabot creates PR:
- Title: "Bump github.com/stretchr/testify from 1.8.4 to 1.9.0"
- Changes: go.mod and go.sum updated
- CI runs tests
```

**Verification:**
- PR created with dependency update
- Tests pass
- Ready for review/merge

## Notes

- go.mod uses Go 1.24.5 (latest stable as of implementation)
- go.sum is automatically maintained by Go toolchain
- .gitignore excludes common development artifacts
- CONTRIBUTING.md references removed gosec (simplified setup)
- SECURITY.md only supports 1.x.x versions (current major version)
- LICENSE.txt is GPL-3.0 (copyleft license)
- package.json is private (not published to npm)
- .releaserc.json supports main (stable) and rc (prerelease) branches
- Dependabot limits to 5 open PRs to avoid noise

## Known Issues

None - all root files and configurations functioning as expected.

## Areas for Improvement

- Add CODE_OF_CONDUCT.md for community guidelines
- Add CHANGELOG.md for manual changelog tracking
- Add .editorconfig for consistent editor settings
- Add .nvmrc for Node.js version pinning
- Add renovate.json as alternative to Dependabot
- Add .dockerignore for Docker builds
- Add .npmrc for npm configuration
- Consider adding go.work for multi-module workspace
