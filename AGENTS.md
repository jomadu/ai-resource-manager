# ARM Agent Operations Guide

## Spec vs Implementation

**Specification** (what should exist):
- `specs/*` - JTBDs, acceptance criteria, algorithms, design decisions

**Implementation** (what actually exists):
- `internal/*` - Core business logic
- `cmd/*` - CLI commands and handlers
- `test/*` - Unit and integration tests
- `docs/*` - User-facing documentation
- Root files - README.md, Makefile, go.mod, etc.
- Infrastructure - `.github/*` (workflows, CI/CD), scripts/*, etc.

## Build & Test

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run specific package tests
go test ./internal/arm/service
go test ./cmd/arm

# Run linting
make lint
```

## Pre-Commit Checklist

**ALWAYS run before committing:**
```bash
make lint    # Fix all linting errors
go test ./...  # Ensure all tests pass
```

## Development

```bash
# Build the binary
go build -o arm cmd/arm/main.go

# Run the binary
./arm help
./arm version
```

## Git Workflow

```bash
# Check status
git status

# Stage all changes
git add -A

# Commit with conventional commit format
git commit -m "type: description"
# Types: feat, fix, docs, refactor, test, chore
# Breaking changes: feat!, fix!

# Push changes
git push

# Create tag (ralph-* prefix for agent work)
git tag ralph-0.0.X
git push origin ralph-0.0.X
```

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
- `docs/` - User documentation
- `specs/` - Technical specifications (JTBDs, acceptance criteria, algorithms)
- `test/e2e/` - End-to-end integration tests

## Current Status

All functionality implemented and tested. See IMPLEMENTATION_PLAN.md and specs/ for details.
