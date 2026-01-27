# Standalone Compilation

## Job to be Done
Compile local ARM resource files (YAML) to tool-specific formats without installing from registries, enabling local development and testing of rulesets and promptsets.

## Activities
1. Discover ARM resource files from paths (files or directories)
2. Validate YAML structure and schema
3. Compile to tool-specific formats (Cursor, Amazon Q, Copilot, Markdown)
4. Write compiled files to output directory
5. Support pattern filtering and recursive discovery

## Acceptance Criteria
- [x] Compile single ARM resource file
- [x] Compile all YAML files in directory
- [x] Recursive directory traversal with --recursive
- [x] Validate-only mode without writing files
- [x] Pattern filtering with --include and --exclude
- [x] Custom namespace with --namespace
- [x] Force overwrite with --force
- [x] Fail-fast mode with --fail-fast
- [x] Support all four tool formats (cursor, amazonq, copilot, markdown)
- [x] Auto-detect file type (ruleset vs promptset)

## Data Structures

### CompileRequest
```go
type CompileRequest struct {
    Paths        []string  // Input file/directory paths
    Tool         string    // Target tool: cursor, amazonq, copilot, markdown
    OutputDir    string    // Output directory path
    Namespace    string    // Optional namespace (defaults to resource ID)
    Force        bool      // Overwrite existing files
    Recursive    bool      // Recursively search directories
    Verbose      bool      // Verbose output
    ValidateOnly bool      // Only validate, don't compile
    Include      []string  // Include patterns (default: *.yml, *.yaml)
    Exclude      []string  // Exclude patterns
    FailFast     bool      // Stop on first error
}
```

## Algorithm

### Compile Workflow
1. Parse command-line arguments
2. Discover files from input paths:
   - If path is file, add to list
   - If path is directory, walk directory (recursive if flag set)
   - Apply include/exclude patterns (default: *.yml, *.yaml)
3. For each discovered file:
   - Detect file type (ruleset or promptset)
   - Parse YAML and validate schema
   - If validate-only, skip compilation
   - Determine namespace (use --namespace or resource ID)
   - Compile to tool-specific format
   - Write to output directory
4. Report success or errors

### File Discovery
1. Default include patterns: `*.yml`, `*.yaml`
2. Walk directories (recursive if flag set)
3. Check exclude patterns first (exclude overrides include)
4. Check include patterns (if specified)
5. Return list of matching files

### Tool-Specific Compilation
- **Cursor**: `.mdc` with frontmatter for rules, `.md` for prompts
- **Amazon Q**: `.md` for both rules and prompts
- **Copilot**: `.instructions.md` for rules
- **Markdown**: `.md` for both rules and prompts

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| No files found | Error: "no files found matching criteria" |
| Invalid tool name | Error: "invalid tool: X (must be cursor, copilot, amazonq, or markdown)" |
| File not ruleset/promptset | Error: "file X is not a valid ruleset or promptset" |
| Invalid YAML | Error: "failed to parse ruleset/promptset X: <details>" |
| Output file exists | Error unless --force specified |
| Multiple errors | Report first error unless --fail-fast disabled |
| No output path in validate mode | Skip output path requirement |
| Empty namespace | Use resource metadata ID as namespace |

## Dependencies

- ARM resource parsing (parser)
- File type detection (filetype)
- Tool-specific compilers (compiler)

## Implementation Mapping

**Source files:**
- `cmd/arm/main.go` - handleCompile() CLI handler
- `internal/arm/service/service.go` - CompileFiles(), compileFile(), discoverFiles()
- `internal/arm/compiler/compiler.go` - CompileRuleset(), CompilePromptset()
- `internal/arm/filetype/filetype.go` - IsRulesetFile(), IsPromptsetFile()
- `internal/arm/parser/parser.go` - ParseRuleset(), ParsePromptset()
- `test/e2e/compile_test.go` - E2E compilation tests

**Related specs:**
- `sink-compilation.md` - Compilation as part of package installation

## Examples

### Compile Single File
```bash
arm compile my-rules.yml --tool cursor --output .cursor/rules/
```

**Input:** `my-rules.yml` (ARM ruleset)
**Output:** `.cursor/rules/myRules_ruleOne.mdc`, `.cursor/rules/myRules_ruleTwo.mdc`, etc.

### Compile Directory Recursively
```bash
arm compile ./rules/ --tool copilot --output .github/copilot/ --recursive
```

**Input:** All `*.yml` and `*.yaml` files in `./rules/` and subdirectories
**Output:** Compiled `.instructions.md` files in `.github/copilot/`

### Validate Without Compiling
```bash
arm compile ./rules/ --validate-only
```

**Output:** "Validation successful" or error messages

### Compile with Custom Namespace
```bash
arm compile my-rules.yml --tool cursor --output .cursor/rules/ --namespace my-team
```

**Output:** Files prefixed with `my-team` instead of resource ID

### Compile with Pattern Filtering
```bash
arm compile ./rules/ --tool amazonq --output .amazonq/rules/ \
  --include "security-*.yml" --exclude "*-experimental.yml" --recursive
```

**Input:** Only files matching `security-*.yml` but not `*-experimental.yml`
**Output:** Compiled `.md` files in `.amazonq/rules/`

### Force Overwrite
```bash
arm compile my-rules.yml --tool cursor --output .cursor/rules/ --force
```

**Output:** Overwrites existing files without error

## Notes

- Standalone compilation is independent of registries, manifests, and lock files
- Useful for local development and testing before publishing to registries
- Does not generate arm_index.* priority files (use sink installation for that)
- Does not track installations in arm-index.json (standalone operation)
- Namespace defaults to resource metadata ID if not specified
- File discovery uses simple pattern matching (not glob ** syntax like package installation)
