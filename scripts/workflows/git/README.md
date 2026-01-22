# Git Registry Testing Workflows

This directory contains scripts for testing ARM with Git registries.

## Quick Start

Run a complete sample workflow with a local Git registry:

```bash
./sample-git-workflow.sh
```

This creates a sandbox, sets up a local Git registry, and demonstrates installing rulesets and promptsets.

## Scripts

### Core Setup Scripts

**`init-git-sandbox.sh`**
- Builds ARM binary from source
- Creates clean sandbox environment at `./sandbox/`
- Foundation for all workflows

**`init-local-git-registry.sh <repo-path>`**
- Creates a local Git repository with sample content
- Includes multiple versions (v1.0.0, v1.0.1, v1.1.0, v2.0.0)
- Contains both legacy format and ARM resource YAML files
- Requires: repository path argument

### Sandbox Setup Scripts

**`init-sample-local-git-registry-sandbox.sh`**
- Initializes sandbox with ARM binary
- Creates local Git registry at `./sandbox/local-registry`
- Adds registry as `local-test`
- Ready for immediate testing

**`init-git-registry-sandbox.sh`**
- Initializes sandbox with ARM binary
- Adds user-specified Git registry
- Configured via `.env` file (see `.env.example`)
- Environment variables: `GIT_REGISTRY_URL`, `GIT_REGISTRY_NAME`

### Workflow Scripts

**`sample-git-workflow.sh`**
- Complete workflow demonstration
- Uses local Git registry
- Installs rulesets and promptsets
- Shows common commands

## Local Testing (No GitHub Required)

Use local git repositories for quick testing without network access:

```bash
# Run complete workflow with local registry
./sample-git-workflow.sh

# Or set up components manually:
./init-git-sandbox.sh
./init-local-git-registry.sh /tmp/my-test-repo
cd sandbox
./arm add registry git --url file:///tmp/my-test-repo my-registry
```

**Advantages:**
- No GitHub account or authentication needed
- Faster (no network calls)
- Works offline
- Easy to modify test data

**Use cases:**
- Local development and testing
- CI/CD pipelines
- Quick iteration on features
- Testing version resolution logic

## Remote Testing (GitHub)

Use actual GitHub repositories for integration testing:

```bash
# Configure .env file
cp .env.example .env
# Edit .env:
#   GIT_REGISTRY_URL=https://github.com/user/repo
#   GIT_REGISTRY_NAME=my-registry

# Set up sandbox with GitHub registry
./init-git-registry-sandbox.sh

# Or use existing sample repository
# Edit .env:
#   GIT_REGISTRY_URL=https://github.com/jomadu/ai-rules-manager-sample-git-registry
#   GIT_REGISTRY_NAME=sample-repo
./init-git-registry-sandbox.sh
```

**Advantages:**
- Tests real-world scenarios
- Validates authentication
- Tests network error handling
- Demonstrates actual usage

**Use cases:**
- Integration testing
- Documentation examples
- Demonstrating to users
- Testing GitHub-specific features

## Sample Repository Structure

The local Git registry created by `init-local-git-registry.sh` contains:

### Version History
- **v1.0.0** - Initial release with basic grug rules (legacy format)
- **v1.0.1** - Bug fixes (patch release)
- **v1.1.0** - Added task management rules (minor release)
- **v2.0.0** - Breaking changes with ARM resource YAML format

### File Formats

**Legacy Format** (v1.x):
- `rules/cursor/*.mdc` - Cursor rules
- `rules/amazonq/*.md` - Amazon Q rules
- `rules/copilot/*.instructions.md` - GitHub Copilot instructions

**ARM Resource Format** (v2.x+):
- `rulesets/grug-brained-dev.yml` - Ruleset specification
- `promptsets/code-review.yml` - Promptset specification

## Examples

```bash
# Basic workflow
./sample-git-workflow.sh

# Custom local registry path
./init-local-git-registry.sh /tmp/custom-repo
# Edit .env:
#   GIT_REGISTRY_URL=file:///tmp/custom-repo
#   GIT_REGISTRY_NAME=custom-test
./init-git-registry-sandbox.sh

# GitHub registry
# Edit .env:
#   GIT_REGISTRY_URL=https://github.com/PatrickJS/awesome-cursorrules
#   GIT_REGISTRY_NAME=awesome-cursorrules
./init-git-registry-sandbox.sh

# Manual setup
./init-git-sandbox.sh
cd sandbox
./arm add registry git --url https://github.com/user/repo my-registry
./arm add sink --tool cursor cursor-rules .cursor/rules
./arm install ruleset my-registry/package cursor-rules
```
