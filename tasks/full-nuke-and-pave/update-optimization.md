# Update Command Optimization

## Problem

The current `UpdateRuleset` method always performs a full uninstall/reinstall cycle, even when the resolved version hasn't changed. This is inefficient and causes unnecessary file system operations.

## Current Flow

```go
func (a *ArmService) UpdateRuleset(ctx context.Context, registry, ruleset string) error {
    manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
    if err != nil {
        return fmt.Errorf("failed to get manifest entry: %w", err)
    }

    slog.InfoContext(ctx, "Updating ruleset", "registry", registry, "ruleset", ruleset)
    return a.InstallRuleset(ctx, &InstallRequest{
        Registry: registry,
        Ruleset:  ruleset,
        Version:  manifestEntry.Version,
        Include:  manifestEntry.Include,
        Exclude:  manifestEntry.Exclude,
    })
}
```

## Optimized Flow

1. Resolve the version constraint from manifest
2. Check filesystem via installer.ListInstalled() to see what's actually installed
3. Compare resolved version with filesystem version
4. If versions match, verify checksum from lockfile (if available)
5. Only reinstall if versions differ or checksum verification fails
6. Log appropriate messages for "no update needed" vs "updating"

## Implementation

```go
func (a *ArmService) UpdateRuleset(ctx context.Context, registry, ruleset string) error {
    // Get manifest entry for version constraint
    manifestEntry, err := a.manifestManager.GetEntry(ctx, registry, ruleset)
    if err != nil {
        return fmt.Errorf("failed to get manifest entry: %w", err)
    }

    // Get current lockfile entry
    currentLockEntry, err := a.lockFileManager.GetEntry(ctx, registry, ruleset)
    if err != nil {
        // No current installation, proceed with install
        slog.InfoContext(ctx, "Installing ruleset (not currently installed)", "registry", registry, "ruleset", ruleset)
        return a.InstallRuleset(ctx, &InstallRequest{
            Registry: registry,
            Ruleset:  ruleset,
            Version:  manifestEntry.Version,
            Include:  manifestEntry.Include,
            Exclude:  manifestEntry.Exclude,
        })
    }

    // Resolve what version we should have
    registries, err := a.manifestManager.GetRawRegistries(ctx)
    if err != nil {
        return fmt.Errorf("failed to get registries: %w", err)
    }
    registryConfig, exists := registries[registry]
    if !exists {
        return fmt.Errorf("registry %s not configured", registry)
    }

    registryClient, err := registry.NewRegistry(registry, registryConfig)
    if err != nil {
        return fmt.Errorf("failed to create registry: %w", err)
    }

    versionStr := manifestEntry.Version
    if versionStr == "" {
        versionStr = "latest"
    }
    versionStr = expandVersionShorthand(versionStr)

    resolvedVersionResult, err := registryClient.ResolveVersion(ctx, versionStr)
    if err != nil {
        return fmt.Errorf("failed to resolve version: %w", err)
    }

    // Check if we're already at the resolved version
    if currentLockEntry.Version == resolvedVersionResult.Version.Version {
        // Verify checksum to ensure integrity
        selector := types.ContentSelector{Include: manifestEntry.Include, Exclude: manifestEntry.Exclude}
        files, err := registryClient.GetContent(ctx, resolvedVersionResult.Version, selector)
        if err != nil {
            return fmt.Errorf("failed to get content for verification: %w", err)
        }

        if lockfile.VerifyChecksum(files, currentLockEntry.Checksum) {
            slog.InfoContext(ctx, "Ruleset already up to date", "registry", registry, "ruleset", ruleset, "version", currentLockEntry.Display)
            return nil
        } else {
            slog.InfoContext(ctx, "Checksum mismatch, reinstalling", "registry", registry, "ruleset", ruleset, "version", currentLockEntry.Display)
        }
    } else {
        slog.InfoContext(ctx, "Updating ruleset", "registry", registry, "ruleset", ruleset, "from", currentLockEntry.Display, "to", resolvedVersionResult.Version.Display)
    }

    // Version changed or checksum failed, proceed with update
    return a.InstallRuleset(ctx, &InstallRequest{
        Registry: registry,
        Ruleset:  ruleset,
        Version:  manifestEntry.Version,
        Include:  manifestEntry.Include,
        Exclude:  manifestEntry.Exclude,
    })
}
```

## Benefits

1. **Performance**: Avoids unnecessary file operations when no update is needed
2. **User Experience**: Clear logging about whether updates are needed
3. **Integrity**: Still verifies checksums to catch corruption
4. **Backwards Compatible**: Same external API, just optimized internally

## Edge Cases Handled

1. **No current installation**: Falls back to install
2. **Checksum mismatch**: Reinstalls even if version matches
3. **Version resolution failure**: Returns appropriate error
4. **Registry not found**: Returns appropriate error
