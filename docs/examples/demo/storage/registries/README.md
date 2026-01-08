# ARM Registry Structure

This directory mimics the structure and contents of the ARM storage located at `~/.arm/storage/registries/`.

## Purpose

The ARM registries directory stores:
- **Registry metadata and git repositories** for package discovery
- **Downloaded packages** to avoid re-downloading the same versions
- **Resolved dependencies** and version information
- **Checksums and integrity data** for package verification

## Directory Structure

```
registries/
└── <registry-key>/               # Individual registry
    ├── metadata.json             # Registry metadata and timestamps
    ├── repo/                     # Git repository (if applicable)
    │   └── .git/                 # Git repository data
    └── packages/                 # Downloaded package files
        └── <package-key>/       # Individual package cache
            ├── metadata.json     # Package metadata and timestamps
            └── <version>/       # Version-specific files
                ├── metadata.json # Version timestamps
                └── files/       # Actual package content
                    └── <package-files>
```

**Note**: The `<registry-key>` and `<package-key>` placeholders shown above are currently in a readable format for demonstration purposes. In the actual ARM cache, these would be generated hashes (e.g., `a1b2c3d4e5f6...`) that uniquely identify registries and packages. The readable format is used here to make the structure more understandable.

## Cache Behavior

- **Registry Metadata**: Cached to avoid repeated API calls to registries
- **Package Files**: Downloaded once and reused across projects
- **Version Resolution**: Cached to speed up dependency resolution
- **Integrity Checks**: SHA256 checksums stored for verification

## Cache Management

ARM automatically manages the cache:
- **Automatic Cleanup**: Old versions and unused packages are cleaned up
- **Size Limits**: Cache size is managed to prevent disk space issues
- **Integrity Verification**: Checksums are verified on cache access
- **Registry Updates**: Metadata is refreshed periodically

## Development Notes

This sample registry structure demonstrates:
- How ARM organizes cached data
- The relationship between registries and packages
- Version-specific caching strategies
- Integration with different registry types (Git, Cloudsmith, etc.)

For production use, the actual storage is located at `~/.arm/storage/` and is managed entirely by the ARM tool.
