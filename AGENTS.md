# ARM Agent Operations Guide

## Work Tracking System

**Status**: Beads issue tracking configured.

**Rationale**: This project uses Beads for AI-native issue tracking. Agents should:
- Use `bd` CLI commands to query and manage issues
- Issues live in `.beads/issues.jsonl` and sync with git
- Reference issue IDs in commit messages and PRs

**Commands**: 
```bash
bd list              # View all issues
bd show <issue-id>   # View issue details
bd create "title"    # Create new issue
bd update <id>       # Update issue status
bd sync              # Sync with git remote
```

## Story/Bug Input

**Location**: `TASK.md` in repository root

**Format**: Markdown file containing current work items, bugs, and feature requests

**Reading**: Agents should read `TASK.md` for current priorities, acceptance criteria, and linked specifications from `specs/` directory.

## Planning System

**Draft Plans**: Not formally defined. Plans should be documented in:
- Issue comments for small changes
- New markdown files in `specs/` for major features (following TEMPLATE.md)

**Publishing**: Commit plans to repository and reference in PRs.

## Build/Test/Lint Commands

**Test**:
```bash
go test ./...                    # Run all tests
go test ./... -v                 # Verbose output
go test -race -coverprofile=coverage.out ./...  # With race detection and coverage
make test                        # Makefile target
```

**Build**:
```bash
go build -o arm cmd/arm/main.go  # Build binary
make build                       # Makefile target with version injection
make build-all                   # Cross-platform builds
```

**Lint**:
```bash
make lint                        # Run golangci-lint (requires installation)
make fmt                         # Format code with gofmt and goimports
make check                       # Run fmt, lint, and test
```

**Setup**:
```bash
make install-tools               # Install goimports and golangci-lint
make setup-hooks                 # Install pre-commit hooks
make setup                       # Full development setup
```

**Rationale**: Commands derived from Makefile and .github/workflows/build.yml. The Makefile provides consistent interface across development and CI environments.

## Specification Definition

**Paths**: `specs/*.md`

**Format**: Markdown files following JTBD (Jobs to be Done) structure with:
- Job description
- Acceptance criteria
- Algorithm/design decisions
- Examples

**Template**: `specs/TEMPLATE.md`

**Rationale**: Specifications define "what should exist" - the requirements, algorithms, and design decisions. They are implementation-agnostic and focus on the problem space.

## Implementation Definition

**Paths**:
- `cmd/arm/*.go` - CLI commands and handlers
- `internal/arm/**/*.go` - Core business logic
- `test/e2e/*.go` - End-to-end integration tests
- `*_test.go` - Unit tests

**Excluded**:
- `docs/*.md` - User-facing documentation (not implementation)
- `specs/*.md` - Specifications (not implementation)
- `vendor/` - Third-party dependencies
- `bin/`, `dist/` - Build artifacts
- `.github/`, `scripts/` - Infrastructure (separate concern)

**Rationale**: Implementation is "what actually exists" - the Go code that realizes the specifications. Tests are part of implementation as they verify behavior.

## Quality Criteria

**Specifications**:
- [ ] Follows TEMPLATE.md structure
- [ ] Includes clear acceptance criteria
- [ ] Provides examples where applicable
- [ ] References related specs

**Implementation**:
- [ ] All tests pass: `go test ./...`
- [ ] All linters pass: `make lint`
- [ ] Code formatted: `make fmt`
- [ ] Conventional commit format used
- [ ] No race conditions: `go test -race ./...`
- [ ] Test coverage maintained (tracked via Codecov)

**Pre-Commit Checklist**:
```bash
make lint      # Must pass with no errors
go test ./...  # Must pass all tests
```

**Rationale**: Quality criteria derived from CI/CD workflows (.github/workflows/build.yml), Makefile targets, and CONTRIBUTING.md. These are boolean checks that must pass before merging.

## Project Structure

- `cmd/arm/` - CLI entry point and command handlers
- `internal/arm/service/` - Business logic layer
- `internal/arm/compiler/` - Tool-specific compilers (Cursor, AmazonQ, Copilot, Markdown)
- `internal/arm/parser/` - YAML resource parsing
- `internal/arm/registry/` - Registry implementations (Git, GitLab, Cloudsmith)
- `internal/arm/sink/` - Sink management and compilation
- `internal/arm/manifest/` - Manifest file handling
- `internal/arm/storage/` - Package storage and caching
- `internal/arm/core/` - Version resolution, pattern matching, archive extraction
- `internal/arm/config/` - .armrc configuration management
- `internal/arm/packagelockfile/` - Lock file management
- `internal/arm/filetype/` - File type detection
- `internal/arm/resource/` - Resource type definitions
- `docs/` - User documentation (12 files, 2686 lines)
- `specs/` - Technical specifications (20 files, JTBD-based)
- `test/e2e/` - End-to-end integration tests (14 test suites)
- `scripts/` - Installation and workflow scripts
- `.github/workflows/` - CI/CD automation

## Tech Stack

- **Language**: Go 1.24.5
- **Build**: Makefile + go.mod
- **Testing**: go test with race detection
- **Linting**: golangci-lint (13 linters enabled)
- **Formatting**: gofmt + goimports
- **CI/CD**: GitHub Actions (build, test, lint, security, release)
- **Pre-commit**: Python-based hooks with conventional commit validation
- **Release**: semantic-release with conventional commits

## Git Workflow

```bash
# Check status
git status

# Stage changes
git add -A

# Commit with conventional format
git commit -m "type: description"
# Types: feat, fix, docs, refactor, test, chore
# Breaking changes: feat!, fix!

# Push changes
git push
```

## Operational Learnings

**Bootstrap Findings**:
- Project has comprehensive specifications in `specs/` following JTBD methodology
- Build system is well-defined with Makefile and GitHub Actions
- Beads issue tracking configured for AI-native workflow
- Quality gates enforced via CI/CD and pre-commit hooks
- All 14 E2E test suites passing with 100% success rate
- Documentation is extensive (12 docs files, 20 spec files)

**Warnings**:
- Planning system not formally defined (should document in specs/ for major features)
- Agents should verify golangci-lint is installed before running `make lint`

**Current Status**: All functionality implemented and tested. Project is in maintenance mode with focus on bug fixes and minor enhancements.

## Landing the Plane (Session Completion)

**When ending a work session**, you MUST complete ALL steps below. Work is NOT complete until `git push` succeeds.

**MANDATORY WORKFLOW:**

1. **File issues for remaining work** - Create issues for anything that needs follow-up
2. **Run quality gates** (if code changed) - Tests, linters, builds
3. **Update issue status** - Close finished work, update in-progress items
4. **PUSH TO REMOTE** - This is MANDATORY:
   ```bash
   git pull --rebase
   bd sync
   git push
   git status  # MUST show "up to date with origin"
   ```
5. **Clean up** - Clear stashes, prune remote branches
6. **Verify** - All changes committed AND pushed
7. **Hand off** - Provide context for next session

**CRITICAL RULES:**
- Work is NOT complete until `git push` succeeds
- NEVER stop before pushing - that leaves work stranded locally
- NEVER say "ready to push when you are" - YOU must push
- If push fails, resolve and retry until it succeeds
