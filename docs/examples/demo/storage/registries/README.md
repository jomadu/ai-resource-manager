# ARM Registry Structure

This directory mimics the structure and contents of the ARM storage located at `~/.arm/storage/registries/`.

## Purpose

The ARM registries directory stores:
- **Registry metadata** for package discovery
- **Downloaded packages** to avoid re-downloading the same versions
- **Git repositories** (for git-based registries) for efficient file access
- **Version timestamps** for cache management

## Directory Structure

```
registries/
└── <registry-key>/               # Individual registry
    ├── metadata.json             # Registry metadata
    ├── repo/                     # Git repository (if applicable)
    │   └── .git/                 # Git repository data
    └── packages/                 # Downloaded package files
        └── <package-key>/       # Individual package cache
            ├── metadata.json     # Package metadata
            └── <version>/       # Version-specific files
                ├── metadata.json # Version metadata + timestamps
                └── files/       # Actual package content
                    └── <package-files>
```

**Note**: The `<registry-key>` and `<package-key>` placeholders shown above are currently in a readable format for demonstration purposes. In the actual ARM cache, these would be generated hashes (e.g., `a1b2c3d4e5f6...`) that uniquely identify registries and packages. The readable format is used here to make the structure more understandable.

## Cache Behavior

- **Registry Metadata**: Stores registry configuration for key generation
- **Package Files**: Downloaded once and reused across projects
- **Version Timestamps**: Track creation, updates, and access times for cache management
- **Git Repository Caching**: Local clones for efficient file access without API limits

## Cache Management

ARM automatically manages the cache:
- **Automatic Cleanup**: Old versions removed based on `updatedAt` timestamps
- **Unused Version Removal**: Versions removed based on `accessedAt` timestamps
- **Package Removal**: Empty package directories removed when all versions are cleaned
- **Git Repository Updates**: Repository clones updated for git-based registries

## Development Notes

This sample registry structure demonstrates:
- How ARM organizes cached data
- The relationship between registries and packages
- Version-specific caching strategies
- Integration with different registry types (Git, Cloudsmith, etc.)

For production use, the actual storage is located at `~/.arm/storage/` and is managed entirely by the ARM tool.
