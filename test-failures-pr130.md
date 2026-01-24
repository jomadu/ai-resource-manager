# PR 130 Test Failures Investigation

**PR:** #130 (branch: `fix-tests`)  
**Workflow Run:** 21300784628  
**Job ID:** 61317705294  
**Date:** 2026-01-23T20:51:30Z

## Failed Tests

### 1. TestCleanCache/cache_with_nuke
- **Package:** `cmd/arm`
- **Duration:** 0.01s
- **Error:** No specific error message in logs
- **Context:** Subtest of `TestCleanCache` - all other subtests passed

### 2. TestGitRegistry_ListPackageVersions
- **Package:** `internal/arm/registry`
- **Duration:** 0.04s
- **Error:** `failed to list versions: exit status 128`
- **Context:** Git exit status 128 indicates command failure (auth, repo access, or config issue)

## Investigation Tasks

### Task 1: Fix TestCleanCache/cache_with_nuke
- [ ] Review test code in `cmd/arm/*_test.go` for `TestCleanCache`
- [ ] Identify the `cache_with_nuke` subtest implementation
- [ ] Run test locally: `go test -v -run TestCleanCache/cache_with_nuke ./cmd/arm`
- [ ] Check for race conditions or timing issues
- [ ] Verify cache directory cleanup logic

### Task 2: Fix TestGitRegistry_ListPackageVersions
- [ ] Review test code in `internal/arm/registry/git_test.go` line 47
- [ ] Identify which Git command is failing (likely `git ls-remote` or `git tag`)
- [ ] Run test locally: `go test -v -run TestGitRegistry_ListPackageVersions ./internal/arm/registry`
- [ ] Check test repository setup/initialization
- [ ] Verify Git operations have proper error handling
- [ ] Consider CI environment differences (permissions, Git config)

### Task 3: Verify Fix
- [ ] Run full test suite locally: `go test ./...`
- [ ] Push fix and verify CI passes
- [ ] Check coverage hasn't decreased

## Additional Context

**Passing Packages:**
- `internal/arm/service` - 54.3% coverage ✓
- `internal/arm/sink` - 84.6% coverage ✓
- `internal/arm/storage` - 84.4% coverage ✓

**Logs Location:**
```bash
gh api repos/jomadu/ai-resource-manager/actions/jobs/61317705294/logs
```
