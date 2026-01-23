# Investigation Plan: PR #130 Test Failures

## Summary
5 tests failing in PR #130, all related to Git branch operations in the registry and storage packages.

## Failed Tests
1. `TestGitRegistry_BranchSupport` - exit status 1
2. `TestGitRegistry_BranchNotFound` - expected 1 version, got 0
3. `TestGitRegistry_VersionPriority` - expected 3 versions (2 tags + 1 branch), got 2
4. `TestGetBranches_RepoNotCloned` - exit status 1
5. `TestGetBranchHeadCommitHash_ValidBranch` - exit status 128 (Git error)

## Root Cause Hypothesis
Branch-related Git operations are failing in test setup or execution. All failures involve branch creation, fetching, or access.

## Investigation Steps

### 1. Examine Test Files
```bash
# Look at the failing test implementations
cat internal/arm/registry/git_test.go | grep -A 30 "TestGitRegistry_BranchSupport"
cat internal/arm/storage/repo_test.go | grep -A 30 "TestGetBranches_RepoNotCloned"
```

### 2. Check Test Helper
```bash
# The error traces point to repo_test_helper.go:125
cat internal/arm/storage/repo_test_helper.go
```

### 3. Review Recent Changes
```bash
# Check what changed in PR #130
gh pr diff 130 --name-only
gh pr diff 130 | grep -E "(branch|Branch)" -C 3
```

### 4. Check Git Version Requirements
```bash
# Test logs show git version 2.52.0
# Verify if branch operations require specific Git features
grep -r "git.*branch" internal/arm/storage/ internal/arm/registry/
```

### 5. Run Single Test Locally
```bash
# Isolate and run one failing test with verbose output
go test -v -run TestGitRegistry_BranchSupport ./internal/arm/registry/
go test -v -run TestGetBranches_RepoNotCloned ./internal/arm/storage/
```

### 6. Check for Race Conditions
```bash
# Branch operations might have timing issues
go test -race -run "Branch" ./internal/arm/...
```

### 7. Examine Git Command Execution
```bash
# Find where Git commands are executed
grep -r "exec.Command.*git" internal/arm/storage/
grep -r "git branch" internal/arm/storage/
```

## Likely Issues

1. **Test setup not creating branches properly** - Check `repo_test_helper.go:80` and `:125`
2. **Git command syntax error** - Exit status 128 indicates Git command failure
3. **Missing remote branch tracking** - Branches might not be pushed/tracked in test repos
4. **Git config issue** - Test environment might need specific Git configuration

## CI-Specific Issue Strategy

Since tests pass locally but fail in CI, this requires iterative debugging via commits:

### Debugging Approach
1. Add debug logging to understand CI environment
2. **Verify all tests pass locally** before committing
3. Push commits to PR #130 to trigger CI runs
4. **Wait for CI to complete** (typically 2-5 minutes)
5. Analyze logs from each run
6. Adjust and repeat until root cause found

**Critical Rule:** Never push commits with failing local tests. Debug output should not break existing functionality.

### Commit Strategy (Conventional Commits)

**Phase 1: Add Diagnostics**
```bash
# Add logging to see what's happening in CI
test: add debug output for branch operations

- Log Git version and config in test setup
- Print branch list before assertions
- Output Git command stderr on failure
```

**Phase 2: Environment Checks**
```bash
# Check for CI-specific Git configuration issues
test: verify git config in CI environment

- Check if git user.name/email are set
- Verify remote tracking branch setup
- Log working directory state
```

**Phase 3: Fix Attempt**
```bash
# Based on diagnostic output, apply fix
fix(test): configure git for branch operations in CI

- Set required Git config for test environment
- Ensure branches are properly tracked
- Add explicit remote setup if needed
```

**Phase 4: Cleanup**
```bash
# Remove debug code once fixed
test: remove debug logging after fix
```

## Resolution

### Root Cause Found
The issue was in `GetBranchHeadCommitHash` method in `internal/arm/storage/repo.go`. The method tried to resolve branch refs using:
1. `origin/branch` - fails because Git doesn't recognize this shorthand
2. `branch` - fails for remote tracking branches in cloned repos

### The Fix
Changed the branch resolution order to use full Git ref paths:
1. `refs/remotes/origin/branch` - for remote tracking branches (most common)
2. `refs/heads/branch` - for local branches
3. `branch` - final fallback for any other ref format

### Test Results
All 5 previously failing tests now pass:
- ✅ `TestGitRegistry_BranchSupport`
- ✅ `TestGitRegistry_BranchNotFound`
- ✅ `TestGitRegistry_VersionPriority`
- ✅ `TestGetBranches_RepoNotCloned`
- ✅ `TestGetBranchHeadCommitHash_ValidBranch`

Full test suite passes locally (all packages in `internal/arm/...`).

### Why This Fixes CI
The issue affected both local and CI environments, but was more apparent in CI because:
- CI runs tests in isolation with fresh clones
- Local development might have had working directories with different Git states
- The fix uses proper Git ref paths that work consistently in all environments

### Next Steps
Ready to commit and push to PR #130 for CI verification.

## Files to Focus On
- `internal/arm/storage/repo_test_helper.go` (lines 80, 125)
- `internal/arm/registry/git_test.go` (BranchSupport, BranchNotFound, VersionPriority tests)
- `internal/arm/storage/repo_test.go` (GetBranches tests)
- `internal/arm/storage/repo.go` (actual Git command implementations)

## Common CI Git Issues
- Missing `git config user.name/user.email`
- Default branch name differences (master vs main)
- Remote tracking not set up
- Shallow clones missing branch refs
- File permissions in temp directories
