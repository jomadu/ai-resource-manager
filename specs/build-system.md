# Build System

## Job to be Done
Compile ARM binaries for multiple platforms with embedded version metadata, enabling developers to build, test, and distribute the CLI tool.

## Activities
1. Build single binary for current platform
2. Build binaries for all supported platforms (Linux, macOS, Windows)
3. Inject version metadata at build time
4. Run tests with coverage reporting
5. Format and lint code
6. Install binary to system PATH
7. Clean build artifacts

## Acceptance Criteria
- [x] Build binary for current platform with `make build`
- [x] Build for all platforms (Linux amd64/arm64, macOS amd64/arm64, Windows amd64) with `make build-all`
- [x] Inject version, commit hash, build timestamp, and platform via LDFLAGS
- [x] Run tests with race detection and coverage with `make test`
- [x] Format code with gofmt and goimports via `make fmt`
- [x] Lint code with golangci-lint via `make lint`
- [x] Install binary to /usr/local/bin with `make install`
- [x] Clean build artifacts with `make clean`
- [x] Install development tools with `make install-tools`
- [x] Setup pre-commit hooks with `make setup-hooks`
- [x] Run all checks (fmt, lint, test) with `make check`
- [x] Complete development setup with `make setup`

## Data Structures

### Build Variables
```makefile
BINARY_NAME := arm
BIN_DIR := bin
DIST_DIR := dist
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S')
LDFLAGS := -ldflags "-X github.com/jomadu/ai-resource-manager/internal/arm/core.buildVersion=$(VERSION) -X github.com/jomadu/ai-resource-manager/internal/arm/core.buildCommit=$(COMMIT) -X github.com/jomadu/ai-resource-manager/internal/arm/core.buildTimestamp=$(BUILD_TIME) -s -w"
```

### Platform Matrix
```
GOOS=linux   GOARCH=amd64  → arm-linux-amd64
GOOS=linux   GOARCH=arm64  → arm-linux-arm64
GOOS=darwin  GOARCH=amd64  → arm-darwin-amd64
GOOS=darwin  GOARCH=arm64  → arm-darwin-arm64
GOOS=windows GOARCH=amd64  → arm-windows-amd64.exe
```

## Algorithm

### Build Single Binary
1. Create bin/ directory if not exists
2. Run `go build` with LDFLAGS to inject version metadata
3. Output binary to bin/arm (or bin/arm.exe on Windows)

### Build All Platforms
1. Create dist/ directory if not exists
2. For each platform in matrix:
   - Set GOOS and GOARCH environment variables
   - Run `go build` with LDFLAGS
   - Output binary to dist/arm-{os}-{arch} (add .exe for Windows)

### Version Injection
1. Extract version from git tags: `git describe --tags --always --dirty`
2. Extract commit hash: `git rev-parse --short HEAD`
3. Generate build timestamp: `date -u '+%Y-%m-%d_%H:%M:%S'`
4. Inject via LDFLAGS into core package variables:
   - `core.buildVersion`
   - `core.buildCommit`
   - `core.buildTimestamp`
5. Add `-s -w` flags to strip debug info and reduce binary size

### Install Binary
1. Build binary with `make build`
2. Check if /usr/local/bin is writable
3. Copy binary to /usr/local/bin/arm (use sudo if needed)

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No git repository | Use "dev" for version, "unknown" for commit |
| Dirty working tree | Append "-dirty" to version string |
| No tags exist | Use commit hash as version |
| /usr/local/bin not writable | Prompt for sudo or fail with error |
| Windows build | Add .exe extension to binary name |

## Dependencies

- Go 1.24.5 or later (specified in go.mod)
- Git (for version extraction)
- golangci-lint (for linting)
- goimports (for formatting)
- pre-commit (for git hooks)

## Implementation Mapping

**Source files:**
- `Makefile` - Build targets and variables
- `internal/arm/core/buildinfo.go` - Version metadata variables and GetBuildInfo()
- `cmd/arm/main.go` - Version command that displays build info
- `go.mod` - Go version and dependencies

**Related specs:**
- `ci-cd-workflows.md` - GitHub Actions that use Makefile targets
- `installation-scripts.md` - Scripts that download and install built binaries

## Examples

### Example 1: Build for Current Platform

**Input:**
```bash
make build
```

**Expected Output:**
```
go build -ldflags "-X github.com/jomadu/ai-resource-manager/internal/arm/core.buildVersion=v3.0.0 -X github.com/jomadu/ai-resource-manager/internal/arm/core.buildCommit=abc1234 -X github.com/jomadu/ai-resource-manager/internal/arm/core.buildTimestamp=2026-01-27_21:39:39 -s -w" -o bin/arm ./cmd/arm
```

**Verification:**
- Binary exists at bin/arm
- `./bin/arm version` displays correct version, commit, and timestamp

### Example 2: Build All Platforms

**Input:**
```bash
make build-all
```

**Expected Output:**
```
Building for all platforms...
GOOS=linux GOARCH=amd64 go build ... -o dist/arm-linux-amd64 ./cmd/arm
GOOS=linux GOARCH=arm64 go build ... -o dist/arm-linux-arm64 ./cmd/arm
GOOS=darwin GOARCH=amd64 go build ... -o dist/arm-darwin-amd64 ./cmd/arm
GOOS=darwin GOARCH=arm64 go build ... -o dist/arm-darwin-arm64 ./cmd/arm
GOOS=windows GOARCH=amd64 go build ... -o dist/arm-windows-amd64.exe ./cmd/arm
```

**Verification:**
- All 5 binaries exist in dist/ directory
- Each binary is executable on its target platform

### Example 3: Run All Checks

**Input:**
```bash
make check
```

**Expected Output:**
```
gofmt -w .
goimports -w .
golangci-lint run
go test -v -race -coverprofile=coverage.out ./...
```

**Verification:**
- Code is formatted
- Linting passes with no errors
- All tests pass
- coverage.out file generated

## Notes

- LDFLAGS includes `-s -w` to strip debug info and reduce binary size (typically 30-40% smaller)
- Version extraction uses `git describe --tags --always --dirty` for human-readable versions
- Build timestamp uses UTC timezone for consistency across environments
- Makefile uses `.PHONY` targets to avoid conflicts with files of the same name

## Known Issues

None - all build targets working as expected.

## Areas for Improvement

- Add `make release` target that combines build-all, checksums, and packaging
- Add `make docker` target for containerized builds
- Consider adding build caching to speed up incremental builds
- Add `make benchmark` target for performance testing
