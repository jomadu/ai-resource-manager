# Code Quality

## Job to be Done
Maintain consistent code style, catch bugs early, and enforce best practices through automated linting, formatting, and commit message validation.

## Activities
1. Lint Go code with golangci-lint
2. Format Go code with gofmt and goimports
3. Validate commit messages follow conventional commit format
4. Run pre-commit hooks automatically
5. Tidy Go module dependencies
6. Check for common issues (trailing whitespace, YAML/JSON syntax, merge conflicts)

## Acceptance Criteria
- [x] Lint with 13 enabled linters (errcheck, gosimple, govet, ineffassign, staticcheck, typecheck, unused, gofmt, goimports, misspell, gocritic, unconvert, unparam)
- [x] Format code with gofmt (standard Go formatting)
- [x] Format imports with goimports (group and sort imports)
- [x] Validate commit messages follow conventional commit format (type: description)
- [x] Run pre-commit hooks on git commit
- [x] Run commit-msg hooks on git commit
- [x] Tidy go.mod and go.sum automatically
- [x] Check for trailing whitespace
- [x] Check for missing end-of-file newlines
- [x] Validate YAML and JSON syntax
- [x] Detect merge conflicts
- [x] Prevent large files from being committed
- [x] Configure golangci-lint timeout (5 minutes)
- [x] Enable gocritic tags (diagnostic, experimental, opinionated, performance, style)
- [x] Show all issues (no limits per linter or per issue type)

## Data Structures

### golangci-lint Config (.golangci.yml)
```yaml
run:
  timeout: 5m
  go: "1.23"

linters:
  enable:
    - errcheck      # Check for unchecked errors
    - gosimple      # Simplify code
    - govet         # Vet examines Go source code
    - ineffassign   # Detect ineffectual assignments
    - staticcheck   # Static analysis
    - typecheck     # Type checking
    - unused        # Check for unused code
    - gofmt         # Check formatting
    - goimports     # Check import formatting
    - misspell      # Check for misspellings
    - gocritic      # Opinionated linter
    - unconvert     # Remove unnecessary conversions
    - unparam       # Detect unused parameters

linters-settings:
  gocritic:
    enabled-tags:
      - diagnostic
      - experimental
      - opinionated
      - performance
      - style

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

### Pre-commit Config (.pre-commit-config.yaml)
```yaml
repos:
  - repo: https://github.com/pre-commit/pre-commit-hooks
    hooks:
      - trailing-whitespace
      - end-of-file-fixer
      - check-yaml
      - check-json
      - check-merge-conflict
      - check-added-large-files

  - repo: https://github.com/compilerla/conventional-pre-commit
    hooks:
      - conventional-pre-commit (commit-msg stage)

  - repo: local
    hooks:
      - go-fmt
      - go-imports
      - go-mod-tidy
      - golangci-lint
```

### Conventional Commit Format
```
<type>: <description>

Types: feat, fix, docs, refactor, test, chore
Breaking changes: feat!, fix!
```

## Algorithm

### Linting (make lint)
1. Run golangci-lint with config from .golangci.yml
2. Check all enabled linters against all Go files
3. Report all issues (no limits)
4. Exit with error if any issues found

### Formatting (make fmt)
1. Run gofmt -w . (format all Go files in place)
2. Run goimports -w . (format imports in all Go files)
3. Files are modified in place

### Pre-commit Hooks
1. **On git add:**
   - Check trailing whitespace
   - Check end-of-file newlines
   - Validate YAML/JSON syntax
   - Check for merge conflicts
   - Check for large files

2. **On git commit:**
   - Run go fmt on staged .go files
   - Run go imports on staged .go files
   - Run go mod tidy if go.mod/go.sum changed
   - Run golangci-lint on staged .go files

3. **On commit message:**
   - Validate message follows conventional commit format
   - Check type is valid (feat, fix, docs, refactor, test, chore)
   - Check description is present

### CI/CD Integration
1. **Build workflow:**
   - Run golangci-lint on all code
   - Fail build if linting errors found

2. **Commitlint workflow:**
   - Validate all commit messages in PR
   - Validate last commit on push
   - Fail if any message invalid

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Linting errors | Build fails, errors displayed |
| Formatting issues | Pre-commit hook fixes automatically |
| Invalid commit message | Commit rejected with error message |
| Large file added | Pre-commit hook rejects commit |
| Merge conflict markers | Pre-commit hook rejects commit |
| Invalid YAML/JSON | Pre-commit hook rejects commit |
| go.mod out of sync | Pre-commit hook runs go mod tidy |
| Timeout (>5min) | golangci-lint fails with timeout error |

## Dependencies

- golangci-lint (installed via make install-tools)
- goimports (installed via make install-tools)
- pre-commit (installed via make setup-hooks)
- Python 3.x (for pre-commit)

## Implementation Mapping

**Source files:**
- `.golangci.yml` - golangci-lint configuration
- `.pre-commit-config.yaml` - Pre-commit hooks configuration
- `Makefile` - Lint, format, and setup targets
- `.github/workflows/build.yml` - CI linting
- `.github/workflows/commitlint.yml` - CI commit message validation

**Related specs:**
- `ci-cd-workflows.md` - Workflows that run linting and validation
- `build-system.md` - Makefile targets for code quality

## Examples

### Example 1: Run Linting

**Input:**
```bash
make lint
```

**Expected Output:**
```
golangci-lint run
# If issues found:
internal/arm/service/service.go:42:2: Error return value is not checked (errcheck)
cmd/arm/main.go:15:1: unused variable 'foo' (unused)
```

**Verification:**
- Exit code 1 if issues found
- Exit code 0 if no issues

### Example 2: Format Code

**Input:**
```bash
make fmt
```

**Expected Output:**
```
gofmt -w .
goimports -w .
```

**Verification:**
- All .go files formatted
- Imports grouped and sorted
- No output if no changes needed

### Example 3: Commit with Invalid Message

**Input:**
```bash
git commit -m "added new feature"
```

**Expected Output:**
```
conventional-pre-commit................................................Failed
- hook id: conventional-pre-commit
- duration: 0.05s
- exit code: 1

Commit message does not follow Conventional Commits format.
Expected: <type>: <description>
Got: added new feature
```

**Verification:**
- Commit rejected
- Working directory unchanged

### Example 4: Commit with Valid Message

**Input:**
```bash
git commit -m "feat: add new registry type"
```

**Expected Output:**
```
Trim Trailing Whitespace.................................................Passed
Fix End of Files.........................................................Passed
Check Yaml...............................................................Passed
Check JSON...............................................................Passed
Check for merge conflicts................................................Passed
Check for added large files..............................................Passed
go fmt...................................................................Passed
go imports...............................................................Passed
golangci-lint............................................................Passed
conventional-pre-commit..................................................Passed
[main abc1234] feat: add new registry type
 1 file changed, 10 insertions(+)
```

**Verification:**
- Commit succeeds
- All hooks pass

### Example 5: Setup Development Environment

**Input:**
```bash
make setup
```

**Expected Output:**
```
go install golang.org/x/tools/cmd/goimports@latest
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
python3 -m venv .venv
.venv/bin/pip install pre-commit
.venv/bin/pre-commit install
.venv/bin/pre-commit install --hook-type commit-msg
go mod tidy
```

**Verification:**
- golangci-lint installed to ~/go/bin/
- goimports installed to ~/go/bin/
- pre-commit installed in .venv/
- Git hooks installed in .git/hooks/

## Notes

- golangci-lint runs all linters in parallel for speed
- Pre-commit hooks only run on staged files (not entire codebase)
- Formatting hooks modify files in place (auto-fix)
- Linting hooks report errors but don't auto-fix
- Conventional commit format enforced by both pre-commit and CI
- gocritic is most opinionated linter (can be noisy but catches subtle issues)
- Timeout set to 5 minutes to handle large codebases

## Known Issues

None - all code quality checks functioning as expected.

## Areas for Improvement

- Add gocyclo linter for cyclomatic complexity
- Add dupl linter for duplicate code detection
- Add gosec linter for security issues
- Add revive linter as alternative to golint
- Add custom linting rules for project-specific patterns
- Add spell checking for comments and documentation
- Add license header checking
- Add import grouping rules (stdlib, external, internal)
- Consider adding gofumpt for stricter formatting
