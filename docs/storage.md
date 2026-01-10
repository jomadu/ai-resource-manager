# Storage System

## Purpose

The storage system stores previously downloaded package versions to avoid redundant network requests. When a project requests a version that exists in storage, ARM serves it locally instead of downloading from the registry.

## Location

Storage directory: `~/.arm/storage`

## Directory Structure

```txt
~/.arm/storage/
    registries/
        <registry-key>/
            metadata.json           # Registry metadata
            repo/                   # Git repositories only
                .git/
                # Repository files
            packages/               # Cached package versions
                <package-key>/
                    metadata.json   # Package metadata
                    <version>/
                        metadata.json # Version metadata + timestamps
                        files/      # Package files
                            # Actual package content
```

## Key Generation

### Registry Key
Unique hash identifying a registry based on:
- **Git registries**: URL + type
- **GitLab registries**: URL + type + (group_id OR project_id)
- **Cloudsmith registries**: URL + type + owner + repository

### Package Key
Unique hash identifying a package based on:
- **Git registries**: Normalized includes/excludes patterns (no name) - the entire repository is the source, "packages" are defined by file patterns
- **Non-git registries**: Package name + includes/excludes patterns

## Metadata Structure

The cache uses a three-level metadata structure for efficient management:
- Registry metadata for registry configuration
- Package metadata for package identification  
- Version metadata for version information and timestamps

### Registry Metadata (`metadata.json`)

**Git Registry:**
```json
{
    "url": "https://github.com/PatrickJS/awesome-cursorrules",
    "type": "git"
}
```

**GitLab Group Registry:**
```json
{
    "url": "https://gitlab.example.com",
    "type": "gitlab",
    "group_id": "123"
}
```

**GitLab Project Registry:**
```json
{
    "url": "https://gitlab.example.com",
    "type": "gitlab",
    "project_id": "456"
}
```

**Cloudsmith Registry:**
```json
{
    "url": "https://api.cloudsmith.io",
    "type": "cloudsmith",
    "owner": "sample-org",
    "repository": "arm-registry"
}
```

### Package Metadata (`packages/<package-key>/metadata.json`)

**Git Registry Package:**
```json
{
    "includes": ["**/*.yml"],
    "excludes": ["**/test/**"]
}
```

**Non-Git Registry Package:**
```json
{
    "name": "clean-code-ruleset",
    "includes": ["**/*.yml"],
    "excludes": ["**/test/**"]
}
```

### Version Metadata (`packages/<package-key>/<version>/metadata.json`)

```json
{
    "version": {
        "major": 1,
        "minor": 0,
        "patch": 0
    },
    "createdAt": "2025-01-08T23:10:43.984784Z",
    "updatedAt": "2025-01-08T23:10:43.984784Z",
    "accessedAt": "2025-01-08T23:10:43.984784Z"
}
```

## Git Repository Caching

Git-based registries include a `repo/` directory containing a local clone of the remote repository. This enables:
- Efficient file access without API rate limits
- Fast update checks using Git operations
- Offline access to previously cached content
