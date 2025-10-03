# ARM Cache Structure

This directory mimics the structure and contents of the ARM cache located at `~/.arm/cache/`.

## Purpose

The ARM cache stores:
- **Registry metadata and indexes** for faster package discovery
- **Downloaded packages** to avoid re-downloading the same versions
- **Resolved dependencies** and version information
- **Checksums and integrity data** for package verification

## Directory Structure

```
cache/
└── registries/                    # Registry metadata and indexes
    └── <registry-key>/           # Individual registry cache
        ├── index.json            # Registry package index
        ├── repository/           # Git repository cache (if applicable)
        │   └── .git/             # Git repository data
        └── packages/             # Downloaded package files
            └── <package-key>/   # Individual package cache
                └── <version>/   # Version-specific files
                    └── <package-files>  # Actual package content
```

**Note**: The `<registry-key>` and `<package-key>` placeholders shown above are currently in a readable format for demonstration purposes. In the actual ARM cache, these would be generated hashes (e.g., `a1b2c3d4e5f6...`) that uniquely identify registries and packages. The readable format is used here to make the structure more understandable.

## Cache Behavior

- **Registry Indexes**: Cached to avoid repeated API calls to registries
- **Package Files**: Downloaded once and reused across projects
- **Version Resolution**: Cached to speed up dependency resolution
- **Integrity Checks**: SHA256 checksums stored for verification

## Cache Management

ARM automatically manages the cache:
- **Automatic Cleanup**: Old versions and unused packages are cleaned up
- **Size Limits**: Cache size is managed to prevent disk space issues
- **Integrity Verification**: Checksums are verified on cache access
- **Registry Updates**: Indexes are refreshed periodically

## Development Notes

This sample cache structure demonstrates:
- How ARM organizes cached data
- The relationship between registries and packages
- Version-specific caching strategies
- Integration with different registry types (Git, Cloudsmith, etc.)

For production use, the actual cache is located at `~/.arm/cache/` and is managed entirely by the ARM tool.
