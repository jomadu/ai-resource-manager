# Sink Lifecycle Sync Design

## Problem

When users add, update, or remove sinks, the installed rulesets don't automatically sync to match the new configuration. This creates inconsistency between sink config and actual installed files.

## Architecture Principle

**Lockfile is NOT the source of truth for what's installed** - it only records resolved versions for reproducible builds. The actual installation state must be determined by scanning sink directories.

## Solution

Implement automatic ruleset sync during sink lifecycle operations:

- **Sink added**: Install rulesets from manifest that match new sink patterns
- **Sink updated**: Sync actual installations with updated patterns/directories
- **Sink removed**: Uninstall rulesets from removed sink directories

## Implementation

### 1. Add Sync Methods to ArmService

```go
// SyncSink syncs rulesets for a specific sink
func (a *ArmService) SyncSink(ctx context.Context, sinkName string, sink *config.SinkConfig) error

// SyncRemovedSink uninstalls rulesets from removed sink
func (a *ArmService) SyncRemovedSink(ctx context.Context, removedSink *config.SinkConfig) error
```

### 2. Modify Command Handlers

Update `cmd/arm/config.go` to call sync after config changes:

```go
var sinkAddCmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        // Add sink to config
        err := configManager.AddSink(...)
        if err != nil {
            return err
        }
        // Sync new sink
        sink, _ := configManager.GetSink(ctx, name)
        return armService.SyncSink(ctx, name, sink)
    },
}

var sinkUpdateCmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        // Update sink config
        err := configManager.UpdateSink(...)
        if err != nil {
            return err
        }
        // Sync updated sink
        sink, _ := configManager.GetSink(ctx, name)
        return armService.SyncSink(ctx, name, sink)
    },
}

var sinkRemoveCmd = &cobra.Command{
    RunE: func(cmd *cobra.Command, args []string) error {
        // Get sink before removal
        sink, _ := configManager.GetSink(ctx, name)
        // Remove from config
        err := configManager.RemoveSink(...)
        if err != nil {
            return err
        }
        // Sync removed sink
        return armService.SyncRemovedSink(ctx, sink)
    },
}
```

### 3. Sync Logic

**SyncSink**:
- Scan sink directories to discover what's actually installed
- Get manifest entries (what should be installed)
- Filter manifest entries matching sink patterns
- Install missing rulesets, remove rulesets that no longer match

**SyncRemovedSink**:
- Scan removed sink directories to find installed rulesets
- Uninstall all found rulesets from those directories

### 4. Directory Scanning Approach

Use existing `installer.ListInstalled()` method to scan directories:
- Provides true source of truth for installation state
- Handles cases where files were manually modified
- Slower than lockfile lookup but ensures correctness

## Error Handling

- Log sync errors but don't rollback config changes
- Users can run `arm install` to fix inconsistencies
- Partial failures are acceptable (some rulesets may install while others fail)

## Dependencies

- Commands need access to both configManager and armService
- Config manager stays pure (no business logic)
- Commands orchestrate config changes + sync operations

## Benefits

- Immediate consistency between config and installed files
- No manual sync step required
- Fixes the "adding sink after install" bug
- True source of truth from actual filesystem state
- Handles manual file modifications gracefully
