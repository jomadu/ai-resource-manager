# CI/CD Workflows

## Job to be Done
Automate building, testing, linting, security scanning, and releasing of ARM binaries through GitHub Actions workflows.

## Activities
1. Build and test on every push and PR
2. Lint code with golangci-lint
3. Scan for security vulnerabilities with CodeQL
4. Review dependencies on PRs
5. Enforce conventional commit messages
6. Automatically version and release using semantic-release
7. Build and upload release binaries for all platforms

## Acceptance Criteria
- [x] Run tests on push to main/rc and PRs
- [x] Run tests with race detection and coverage reporting
- [x] Upload coverage to Codecov
- [x] Lint code with golangci-lint on every build
- [x] Build binaries for all platforms (Linux amd64/arm64, macOS amd64/arm64, Windows amd64)
- [x] Upload build artifacts for each platform
- [x] Run CodeQL security analysis on push, PR, and weekly schedule
- [x] Scan for GitHub Actions vulnerabilities
- [x] Review dependencies on PRs for security issues
- [x] Validate commit messages follow conventional commit format
- [x] Automatically determine version using semantic-release
- [x] Create GitHub release with changelog
- [x] Build and upload release binaries with SHA256 checksums
- [x] Package binaries as .tar.gz archives

## Data Structures

### Build Matrix
```yaml
strategy:
  matrix:
    goos: [linux, darwin, windows]
    goarch: [amd64, arm64]
    exclude:
      - goos: windows
        goarch: arm64
```

### Release Output
```yaml
outputs:
  released: ${{ steps.release.outputs.released }}
  version: ${{ steps.release.outputs.version }}
```

### Semantic Release Config (.releaserc.json)
```json
{
  "branches": ["main", {"name": "rc", "prerelease": true}],
  "plugins": [
    "@semantic-release/commit-analyzer",
    "@semantic-release/release-notes-generator",
    "@semantic-release/github"
  ]
}
```

## Algorithm

### Build Workflow (.github/workflows/build.yml)
1. **Test Job:**
   - Checkout code
   - Setup Go 1.22
   - Cache Go modules
   - Download dependencies
   - Configure Git for E2E tests
   - Build ARM binary for E2E tests
   - Run tests with race detection and coverage
   - Upload coverage to Codecov

2. **Lint Job:**
   - Checkout code
   - Setup Go 1.22
   - Run golangci-lint with latest version

3. **Build Job (depends on test + lint):**
   - For each platform in matrix:
     - Checkout code
     - Setup Go 1.22
     - Build binary with GOOS/GOARCH
     - Add .exe extension for Windows
     - Upload artifact

### Release Workflow (.github/workflows/release.yml)
1. **Release Job:**
   - Checkout code with full history
   - Setup Node.js 20
   - Install npm dependencies (semantic-release)
   - Run semantic-release
   - Parse output to determine if release published
   - Extract version from latest git tag
   - Set outputs: released (true/false), version

2. **Build Binaries Job (if released):**
   - For each platform in matrix:
     - Checkout code
     - Setup Go 1.22
     - Build binary with version metadata injected via LDFLAGS
     - Create .tar.gz package
     - Generate SHA256 checksum
     - Upload package and checksum to GitHub release

### Security Workflow (.github/workflows/security.yml)
1. **Dependency Review (PRs only):**
   - Checkout code
   - Run dependency-review-action

2. **CodeQL:**
   - Checkout code
   - Initialize CodeQL for Go
   - Perform CodeQL analysis

### CodeQL Advanced Workflow (.github/workflows/codeql.yml)
1. **Analyze Job:**
   - Matrix: [actions, go]
   - Checkout code
   - Initialize CodeQL with language-specific build mode
   - Perform CodeQL analysis with category tagging

### Commitlint Workflow (.github/workflows/commitlint.yml)
1. **Push Event:**
   - Checkout code with full history
   - Setup Python
   - Install pre-commit
   - Validate last commit message

2. **PR Event:**
   - Checkout code with full history
   - Setup Python
   - Install pre-commit
   - Validate all commit messages in PR

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No release needed | Release job outputs released=false, build-binaries skipped |
| Prerelease (rc branch) | Version tagged with -rc.X suffix |
| Test failure | Build job blocked, no artifacts uploaded |
| Lint failure | Build job blocked, no artifacts uploaded |
| Security vulnerability | CodeQL creates alert, workflow continues |
| Invalid commit message | Commitlint fails, PR blocked |
| Missing GITHUB_TOKEN | Release fails with authentication error |

## Dependencies

- GitHub Actions runners (ubuntu-latest)
- Go 1.22
- Node.js 20
- semantic-release npm packages
- golangci-lint
- CodeQL
- Codecov

## Implementation Mapping

**Source files:**
- `.github/workflows/build.yml` - Build, test, lint workflow
- `.github/workflows/release.yml` - Semantic release and binary publishing
- `.github/workflows/security.yml` - Dependency review and CodeQL
- `.github/workflows/codeql.yml` - Advanced CodeQL analysis
- `.github/workflows/commitlint.yml` - Conventional commit validation
- `.releaserc.json` - Semantic release configuration
- `package.json` - npm dependencies for semantic-release

**Related specs:**
- `build-system.md` - Makefile targets used by workflows
- `code-quality.md` - Linting and formatting standards
- `installation-scripts.md` - Scripts that download workflow artifacts

## Examples

### Example 1: Successful Build on PR

**Input:**
```bash
git push origin feature/new-feature
# Create PR to main
```

**Expected Output:**
```
✓ Test job passes (all tests green, coverage uploaded)
✓ Lint job passes (no linting errors)
✓ Build job passes (5 platform binaries uploaded as artifacts)
✓ Commitlint passes (all commits follow conventional format)
✓ Security passes (no vulnerabilities found)
```

**Verification:**
- All checks green on PR
- Artifacts available for download
- Coverage report on Codecov

### Example 2: Automatic Release on Main

**Input:**
```bash
git push origin main
# Commit message: "feat: add new registry type"
```

**Expected Output:**
```
✓ Build workflow passes
✓ Release job runs semantic-release
  - Analyzes commits since last release
  - Determines new version (minor bump for feat)
  - Creates git tag v3.1.0
  - Generates changelog
  - Creates GitHub release
✓ Build-binaries job runs
  - Builds 5 platform binaries with v3.1.0 metadata
  - Creates .tar.gz packages
  - Generates SHA256 checksums
  - Uploads to GitHub release
```

**Verification:**
- New tag v3.1.0 exists
- GitHub release created with changelog
- 10 assets uploaded (5 .tar.gz + 5 .sha256)

### Example 3: Prerelease on RC Branch

**Input:**
```bash
git push origin rc
# Commit message: "feat: experimental feature"
```

**Expected Output:**
```
✓ Release creates prerelease tag v3.1.0-rc.1
✓ GitHub release marked as prerelease
✓ Binaries built with v3.1.0-rc.1 metadata
```

**Verification:**
- Tag v3.1.0-rc.1 exists
- Release marked as prerelease on GitHub
- Binaries contain prerelease version

## Notes

- Build workflow runs on both push and PR to catch issues early
- Release workflow only runs on push to main/rc (not on PRs)
- CodeQL runs on push, PR, and weekly schedule (Sundays at 2:41 AM)
- Dependency review only runs on PRs (not on push)
- Commitlint validates all commits in PR, not just the merge commit
- Semantic-release uses conventional commits to determine version bump (feat=minor, fix=patch, BREAKING CHANGE=major)
- Build artifacts are temporary (90 days), release assets are permanent
- GITHUB_TOKEN is automatically provided by GitHub Actions

## Known Issues

None - all workflows functioning as expected.

## Areas for Improvement

- Add workflow to automatically update documentation on release
- Add workflow to publish binaries to package managers (Homebrew, Chocolatey)
- Add workflow to run E2E tests against real registries (requires secrets)
- Add workflow to benchmark performance and track regressions
- Add workflow to automatically update CHANGELOG.md
- Consider adding workflow dispatch for manual releases
