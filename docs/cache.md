# Cache System

## Purpose

The cache stores previously downloaded package versions to avoid redundant network requests. When a project requests a version that exists in cache, ARM serves it locally instead of downloading from the registry.

## Location

Cache directory: `~/.arm/registries`

## Directory Structure

```txt
~/.arm/registries/
    <registry-key>/
        metadata.json           # Registry metadata and timestamps
        repo/                   # Git repositories only
            .git/
            # Repository files
        packages/               # Cached package versions
            <package-key>/
                metadata.json   # Package metadata and timestamps
                <version>/
                    metadata.json # Version timestamps
                    files/      # Package files
                        # Actual package content
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

## Metadata Structure

The cache uses a three-level metadata structure for efficient management:
- Registry metadata for key generation and registry-level timestamps
- Package metadata for key generation and package-level timestamps  
- Version metadata for version-level timestamps

### Registry Metadata (`metadata.json`)

```json
{
    "metadata": {
        "url": "https://api.cloudsmith.io",
        "type": "cloudsmith",
        "owner": "sample-org",
        "repository": "arm-registry"
    },
    "created_on": "2025-01-08T23:10:43.984784Z",
    "last_updated_on": "2025-01-08T23:10:43.984784Z",
    "last_accessed_on": "2025-01-08T23:10:43.984784Z"
}
```

### Package Metadata (`packages/<package-key>/metadata.json`)

```json
{
    "metadata": {
        "name": "clean-code-ruleset",
        "description": "Clean code best practices",
        "includes": ["**/*.yml"],
        "excludes": ["**/test/**"]
    },
    "created_on": "2025-01-08T23:10:43.984784Z",
    "last_updated_on": "2025-01-08T23:10:43.984784Z",
    "last_accessed_on": "2025-01-08T23:10:43.984784Z"
}
```

### Version Metadata (`packages/<package-key>/<version>/metadata.json`)

```json
{
    "created_on": "2025-01-08T23:10:43.984784Z",
    "last_updated_on": "2025-01-08T23:10:43.984784Z",
    "last_accessed_on": "2025-01-08T23:10:43.984784Z"
}
```

## Git Repository Caching

Git-based registries include a `repo/` directory containing a local clone of the remote repository. This enables:
- Efficient file access without API rate limits
- Fast update checks using Git operations
- Offline access to previously cached content

## Cache Benefits

- **Performance**: Eliminates redundant downloads
- **Reliability**: Reduces dependency on network availability
- **API Efficiency**: Minimizes registry API calls
- **Offline Support**: Enables work with cached packages when offline
