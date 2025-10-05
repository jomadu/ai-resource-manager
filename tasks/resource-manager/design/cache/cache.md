The cache saves previously downloaded versions of packages, so that if a project requests a previously downloaded version, ARM doesn't have to go download the requested version from the registry.

The cache lives in the users home directory under   `~/.arm/cache`.

Here's the structure:

```txt
~/.arm/cache/registries/
    <registry-key>/ # a non-git based registry
        index.json
        packages/
            <package-key>/
                <version>/
                    # package files
                    ...
    <registry-key>/ # a git based registry
        index.json
        repository/
            .git/
            # repository files
            ...
        packages/
            <package-key>/
                <version>/
                    # package files
                    ...
```

you'll notice two types of keys:

1. `<registry-key>` this key is a hash of an object that represents the a unique registry. For git based registries, this is the url and the registry type. For gitlab based registries, this is the url, type, and either the group id or the project id. For cloudsmith, this is the url and type.
2. `<package-key>` this key is a hash of properties that make a package unique. for rulesets and promptsets from git based registries, this is the normalized includes and excludes. for non-git based registries, this is just the package name.

you'll also notice that each registry contains an `index.json`. this is there to store package and registry metadata, like the metadata used to build the registry and package keys, and created/updated/accessed timestamps.

see the following example for a cloudsmith registry `index.json`:

```json
{
    "registry_metadata": {
        "url": "https://app.cloudsmith.com/sample-org/arm-registry",
        "type": "cloudsmith"
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
        },
        "<package-key>": {
            "package_metadata": {
                "name": "code-review-promptset"
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

lastly, you'll notice that for git-based registries, there is a folder named repository, in which the git repository is locally cloned. this is useful for downloading files and checking for updates, without having to hammer an api.
