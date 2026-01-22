# Known Issues

## ARM Resource Compilation Not Happening During Installation

**Status:** Bug  
**Severity:** High  
**Discovered:** 2026-01-22

### Description

When installing rulesets or promptsets from ARM resource YAML files (e.g., `rulesets/*.yml`, `promptsets/*.yml`), the files are being copied as-is to the sink directories instead of being compiled into individual rule/prompt files.

### Expected Behavior

When installing a ruleset like `rulesets/grug-brained-dev.yml` containing:
```yaml
spec:
  rules:
    grug-simplicity:
      name: "Grug Simplicity Rule"
      ...
    grug-testing:
      name: "Grug Testing Rule"
      ...
```

ARM should compile this into individual files:
- `.cursor/rules/arm/local-test/grug-rules/v2.0.0/grugBrainedDev_grug-simplicity.mdc`
- `.cursor/rules/arm/local-test/grug-rules/v2.0.0/grugBrainedDev_grug-testing.mdc`

### Actual Behavior

ARM copies the raw YAML file without compilation:
- `.cursor/rules/arm/local-test/grug-rules/v2.0.0/rulesets/grug-brained-dev.yml`

### Reproduction Steps

1. Run the sample Git workflow:
   ```bash
   cd scripts/workflows/git
   ./sample-git-workflow.sh
   ```

2. Check the installed files:
   ```bash
   cd sandbox
   find .cursor/rules/arm -type f -name "*.mdc"
   # Expected: Multiple .mdc files (one per rule)
   # Actual: Only arm_index.mdc
   ```

3. Verify the raw YAML is present:
   ```bash
   find .cursor/rules/arm -name "*.yml"
   # Shows: .cursor/rules/arm/local-test/grug-rules/v2.0.0/rulesets/grug-brained-dev.yml
   ```

### Impact

- AI tools (Cursor, Amazon Q, GitHub Copilot) cannot use the rules because they're in YAML format instead of compiled format
- The `arm compile` command works correctly when run manually, so the compiler itself is functional
- The issue is in the installation pipeline where compilation should happen but doesn't

### Related Code

- Installation logic: `cmd/arm/install*.go`, `internal/arm/service/`
- Compilation logic: `internal/arm/compiler/compiler.go` (works correctly)
- Test workflows: `scripts/workflows/git/sample-git-workflow.sh`

### Workaround

None currently. Manual compilation with `arm compile` doesn't integrate with the installation system.

### Notes

- The `arm_index.mdc` file is correctly generated and references the YAML file
- The `arm.json` manifest correctly tracks the installation
- Legacy format files (pre-compiled `.mdc`, `.md`, `.instructions.md`) work fine when installed
