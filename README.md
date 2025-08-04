![arm-header](./assets/header.png)

A package manager for AI coding assistant rulesets. Install, update, and manage coding rules across different AI tools like Cursor and Amazon Q Developer.

## What is ARM?

ARM solves the problem of managing and sharing AI coding rulesets across teams and projects. Instead of manually copying `.cursorrules` files or `.amazonq/rules` directories, ARM provides a centralized way to distribute, version, and update coding rules from multiple registries.

## Quick Start

### Installation

```bash
# Install via curl (coming soon)
curl -sSL https://install.arm.dev | sh

# Or download binary from releases
wget https://github.com/user/arm/releases/latest/download/arm-linux-amd64
chmod +x arm-linux-amd64
sudo mv arm-linux-amd64 /usr/local/bin/arm
```

### Basic Usage

```bash
# Install from GitLab/S3/Filesystem registries (version discovery supported)
arm install company@typescript-rules

# Install from HTTP registry (exact version required)
arm install company@typescript-rules@1.0.0

# Install from manifest
arm install

# List installed rulesets
arm list
arm list --format=json  # JSON output

# Update all rulesets
arm update

# Check for outdated rulesets
arm outdated

# Manage configuration
arm config list
arm config get sources.default
arm config set sources.company https://internal.company.local/
```

## Configuration

### Project Configuration (rules.json)

```json
{
  "targets": [".cursorrules", ".amazonq/rules"],
  "dependencies": {
    "typescript-rules": "^1.0.0",
    "company@security-rules": "^2.1.0"
  }
}
```

### Registry Configuration (.armrc)

```ini
[sources]
default = https://registry.armjs.org/
company = https://gitlab.company.com

[sources.company]
type = gitlab
projectID = 12345
authToken = $COMPANY_REGISTRY_TOKEN
```

## Commands

| Command | Description | Status |
|---------|-------------|--------|
| `arm install [ruleset]` | Install rulesets | ✅ |
| `arm uninstall <ruleset>` | Remove a ruleset | ✅ |
| `arm update [ruleset]` | Update rulesets | ✅ |
| `arm list [--format=table\|json]` | List installed rulesets | ✅ |
| `arm outdated` | Show outdated rulesets | ✅ |
| `arm config [list\|get\|set] [key] [value]` | Manage configuration | ✅ |
| `arm clean` | Clean cache and unused files | 🚧 |
| `arm help` | Show help | ✅ |
| `arm version` | Show version | ✅ |

## Supported Targets

- **Cursor IDE**: `.cursorrules`
- **Amazon Q Developer**: `.amazonq/rules/`
- Extensible for future AI coding tools

## Supported Registries

- **GitLab Package Registry** - Project and group-level registries with full metadata and version discovery
- **AWS S3** - S3 bucket-based registries with prefix support and S3 prefix-based version discovery
- **Generic HTTP** - Simple file server registries (exact versions required)
- **Local File System** - Local directory registries with filesystem-based version discovery

## File Structure

After installation, your project will look like:

```
# Project files
.cursorrules/
  arm/
    company/
      typescript-rules/
        1.0.1/
          rule-1.md
          rule-2.md
.amazonq/
  rules/
    arm/
      company/
        typescript-rules/
          1.0.1/
            rule-1.md
            rule-2.md
rules.json
rules.lock
.armrc

# Global cache (shared across projects)
~/.arm/
  cache/
    packages/
      registry.armjs.org/
        typescript-rules/
          1.0.1/
            package.tar.gz
    registry/
      registry.armjs.org/
        metadata.json
        versions.json
```

## Development Status

🚧 **Phase 4 In Progress** - Cache system implemented, clean command in development. See our [development roadmap](docs/project/roadmap.md):

- **Phase 1**: Core commands (install, uninstall, list) - ✅ **COMPLETED**
- **Phase 2**: Configuration and registry support - ✅ **COMPLETED**
- **Phase 3**: Update/outdated functionality - ✅ **COMPLETED**
- **Phase 4**: Cache management and cleanup - 🚧 **IN PROGRESS**
- **Phase 5**: Testing and documentation - 📋 **PLANNED**

✅ **Phase 3 Complete**: Update/Outdated functionality implemented with version constraints and progress reporting.
🚧 **Phase 4 Progress**: Global cache system implemented with 60%+ performance improvements. Clean command in development.

📈 **Technical Tasks**: See [docs/project/tasks.md](docs/project/tasks.md) for detailed implementation tracking.

## Documentation

📚 **[Complete Documentation](docs/)** - Organized by audience:
- **[Product Requirements](docs/product/)** - Specifications and business requirements
- **[Project Planning](docs/project/)** - Roadmaps, status, and milestones
- **[Technical Implementation](docs/technical/)** - Development tasks and guides

## Contributing

This project is implemented in Go for fast, dependency-free distribution. See [docs/product/requirements.md](docs/product/requirements.md) for detailed requirements and architecture decisions.

For development work, start with [docs/technical/tasks.md](docs/technical/tasks.md) for current implementation status.

## License

MIT License - see [LICENSE](LICENSE) for details.
