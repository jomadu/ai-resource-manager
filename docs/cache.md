# Cache System

## Purpose

The cache stores previously downloaded package versions to avoid redundant network requests. When a project requests a version that exists in cache, ARM serves it locally instead of downloading from the registry.

## Location

Cache directory: `~/.arm/cache`

## Directory Structure

```txt
~/.arm/cache/registries/
    <registry-key>/
        index.json              # Registry metadata and package index
        packages/               # Cached package versions
            <package-key>/
                <version>/
                    # Package files
        repository/             # Git repositories only
            .git/
            # Repository files
```

## Key Generation

### Registry Key
Unique hash identifying a registry based on:
- **Git registries**: URL + type
- **GitLab registries**: URL + type + (group_id OR project_id)
- **Cloudsmith registries**: URL + type

### Package Key
Unique hash identifying a package based on:
- **Git registries**: Normalized includes/excludes patterns
- **Non-git registries**: Package name

## Index Structure

Each registry maintains an `index.json` file containing:
- Registry metadata for key generation
- Package metadata and version tracking
- Timestamp information for cache management

### Example: Cloudsmith Registry Index

```json
{
    "registry_metadata": {
        "url": "https://api.cloudsmith.io",
        "type": "cloudsmith",
        "owner": "sample-org",
        "repository": "arm-registry"
    },
    "created_on": "2025-01-08T23:10:43.984784Z",
    "last_updated_on": "2025-01-08T23:10:43.984784Z",
    "last_accessed_on": "2025-01-08T23:10:43.984784Z",
    "packages": {
        "<package-key>": {
            "package_metadata": {
                "name": "clean-code-ruleset"
            },
            "created_on": "2025-01-08T23:10:43.984784Z",
            "last_updated_on": "2025-01-08T23:10:43.984784Z",
            "last_accessed_on": "2025-01-08T23:10:43.984784Z",
            "versions": {
                "1.0.0": {
                    "created_on": "2025-01-08T23:10:43.984784Z",
                    "last_updated_on": "2025-01-08T23:10:43.984784Z",
                    "last_accessed_on": "2025-01-08T23:10:43.984784Z"
                }
            }
        }
    }
}
```

## Git Repository Caching

Git-based registries include a `repository/` directory containing a local clone of the remote repository. This enables:
- Efficient file access without API rate limits
- Fast update checks using Git operations
- Offline access to previously cached content

## Cache Benefits

- **Performance**: Eliminates redundant downloads
- **Reliability**: Reduces dependency on network availability
- **API Efficiency**: Minimizes registry API calls
- **Offline Support**: Enables work with cached packages when offline
