# Phase 1: CLI Layer Cleanup - COMPLETED ✅

## Task 1.1: Split CLI Commands ✅ DONE

### ✅ Achieved State
- `cmd/arm/main.go` is 30 lines with only root command setup
- All commands extracted into separate focused files
- Business logic separated from CLI handling

### ✅ Final Structure
```
cmd/arm/
├── main.go           # Root command setup only (~50 lines)
├── install.go        # Install command (~80 lines)
├── info.go          # Info/list/outdated commands (~100 lines)
├── config.go        # Config commands (cleanup existing)
├── cache.go         # Cache commands (cleanup existing)
├── version.go       # Version command (~20 lines)
├── parser.go        # Argument parsing utilities (~60 lines)
└── output.go        # Output formatting utilities (~80 lines)
```

### Implementation Steps

#### Step 1.1.1: Create parser.go
```go
package main

type RulesetRef struct {
    Registry string
    Name     string
    Version  string
}

func ParseRulesetArg(arg string) (RulesetRef, error)
func ParseVersionConstraint(version string) string
func ValidateRulesetName(name string) error
```

#### Step 1.1.2: Create output.go
```go
package main

type OutputFormatter interface {
    FormatRulesetInfo(*arm.RulesetInfo, bool) string
    FormatOutdatedTable([]arm.OutdatedRuleset) error
    FormatJSON(interface{}) error
}

func NewTableFormatter() OutputFormatter
func NewJSONFormatter() OutputFormatter
```

#### Step 1.1.3: Extract install.go
- Move install command definition
- Use parser utilities
- Keep only CLI concerns

#### Step 1.1.4: Extract info.go
- Move list, info, outdated commands
- Use output formatters
- Consolidate similar commands

#### Step 1.1.5: Extract version.go
- Simple version command
- Use output formatter

#### Step 1.1.6: Clean up main.go
- Only root command and subcommand registration
- Move global service initialization

## Task 1.2: Extract CLI Utilities

### Current Issues
- Manual string parsing with custom functions
- Inconsistent error handling
- Mixed formatting logic

### Target Improvements
- Proper argument validation
- Consistent error messages
- Reusable formatting components

### Implementation Steps

#### Step 1.2.1: Improve Argument Parsing
```go
// Replace manual parsing with structured approach
func ParseRulesetArg(arg string) (RulesetRef, error) {
    if arg == "" {
        return RulesetRef{}, fmt.Errorf("ruleset argument cannot be empty")
    }

    // Handle registry/ruleset[@version] format
    parts := strings.SplitN(arg, "/", 2)
    if len(parts) != 2 {
        return RulesetRef{}, fmt.Errorf("invalid ruleset format: %s (expected registry/ruleset[@version])", arg)
    }

    registry := parts[0]
    rulesetAndVersion := parts[1]

    versionParts := strings.SplitN(rulesetAndVersion, "@", 2)
    ruleset := versionParts[0]
    version := ""
    if len(versionParts) > 1 {
        version = versionParts[1]
    }

    if err := ValidateRulesetName(ruleset); err != nil {
        return RulesetRef{}, fmt.Errorf("invalid ruleset name: %w", err)
    }

    return RulesetRef{
        Registry: registry,
        Name:     ruleset,
        Version:  version,
    }, nil
}
```

#### Step 1.2.2: Standardize Output Formatting
```go
type TableFormatter struct {
    writer io.Writer
}

func (f *TableFormatter) FormatRulesetInfo(info *arm.RulesetInfo, detailed bool) string {
    var buf strings.Builder

    if detailed {
        f.formatDetailedInfo(&buf, info)
    } else {
        f.formatSummaryInfo(&buf, info)
    }

    return buf.String()
}
```

## Task 1.3: Improve Error Handling

### Current Issues
- Inconsistent error message format
- Some errors logged at CLI level
- Missing context in error messages

### Target Improvements
- Consistent error format: "failed to <action>: <reason>"
- All logging handled by service layer
- Rich error context for debugging

### Implementation Steps

#### Step 1.3.1: Standardize Error Messages
```go
// Before
return fmt.Errorf("registry %s not found", registryName)

// After
return fmt.Errorf("failed to install ruleset: registry %s not configured", registryName)
```

#### Step 1.3.2: Remove CLI Logging
- Remove all `slog` calls from CLI layer
- Let service layer handle logging
- CLI only handles user-facing messages

#### Step 1.3.3: Add Error Context
```go
func (cmd *installCommand) RunE(cmd *cobra.Command, args []string) error {
    rulesets, err := parseRulesetArgs(args)
    if err != nil {
        return fmt.Errorf("failed to parse arguments: %w", err)
    }

    for _, ruleset := range rulesets {
        if err := cmd.service.InstallRuleset(ctx, ruleset); err != nil {
            return fmt.Errorf("failed to install %s/%s: %w",
                ruleset.Registry, ruleset.Name, err)
        }
    }

    return nil
}
```

## Acceptance Criteria - ALL COMPLETED ✅

### Task 1.1 Complete ✅:
- [x] `main.go` is 30 lines with only root command setup
- [x] Each command file is focused on single responsibility
- [x] All commands use shared utilities from `parser.go` and `output.go`
- [x] No business logic in CLI command handlers

### Task 1.2 Complete ✅:
- [x] All argument parsing uses structured approach
- [x] Input validation happens at CLI layer with clear error messages
- [x] Output formatting is consistent across all commands
- [x] Proper parser functions replace manual string manipulation

### Task 1.3 Complete ✅:
- [x] All error messages follow consistent format
- [x] No logging calls in CLI layer
- [x] Error context includes operation and affected resources
- [x] User-friendly error messages for common failure cases

## Testing Strategy

### Unit Tests
- Test argument parsing with valid/invalid inputs
- Test output formatting with various data structures
- Test error message formatting

### Integration Tests
- Test command execution with mocked service
- Verify error propagation from service to CLI
- Test output format consistency

## Time Estimate: 6-8 hours total
- Task 1.1: 3-4 hours
- Task 1.2: 2-3 hours
- Task 1.3: 1 hour
