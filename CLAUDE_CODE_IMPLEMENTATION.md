# Claude Code Implementation Specification

## Overview

This document specifies how ARM will compile rulesets, promptsets, and skillsets to Claude Code format.

## File Format

### Extension
- **Rules**: `.md`
- **Skills/Prompts**: `.md`

### Directory Structure

**Hierarchical Layout** (recommended):
```
.claude/
├── rules/
│   └── arm/
│       └── registry/
│           └── package/
│               └── version/
│                   ├── rulesetId_ruleOne.md
│                   └── rulesetId_ruleTwo.md
└── skills/
    └── arm/
        └── registry/
            └── package/
                └── version/
                    ├── promptsetId_promptOne.md
                    ├── promptsetId_promptTwo.md
                    └── skillKey/                    # Skillset skill directory
                        ├── SKILL.md                 # Generated skill file
                        ├── scripts/
                        │   └── helper.py            # Supporting files
                        └── reference.md
```

**Note**: The structure under `version/` matches whatever directory structure exists in the source package. ARM doesn't create `rules/` or `prompts/` subdirectories - those only exist if the package has them. Skillsets create a directory per skill (using the skill key) containing `SKILL.md` and any supporting files.

## Ruleset Compilation

### File Naming
- **Compiled from ARM YAML**: `{rulesetId}_{ruleKey}.md`
- **Raw .md files in package**: Keep original filename
- Example compiled: `cleanCode_formatting.md`
- Example raw: `my-rule.md` → `my-rule.md`
- Truncation: Same rules as other compilers

### File Content Structure

```markdown
---
paths:
  - "**/*.ts"
  - "**/*.js"
---

---
namespace: registry/package@version
ruleset:
  id: rulesetId
  name: Ruleset Name
  description: Ruleset description
  rules:
    - ruleOne
    - ruleTwo
rule:
  id: ruleOne
  name: Rule One
  description: Rule description
  enforcement: MUST
  scope:
    - files: ["**/*.ts"]
---

# Rule Title

Rule body content here.
```

### Frontmatter Sections

**Section 1: Claude Code Frontmatter (Optional)**
- `paths`: Array of glob patterns for conditional loading
- Only include if rule has `scope.files` defined
- Convert ARM scope patterns to Claude Code paths format

**Section 2: ARM Metadata**
- Namespace: `registry/package@version`
- Ruleset metadata: id, name, description, rules list
- Rule metadata: id, name, description, enforcement, scope

### Scope Mapping

ARM `scope.files` → Claude Code `paths`:
```yaml
# ARM format
scope:
  - files: ["**/*.ts", "**/*.tsx"]

# Becomes Claude Code frontmatter
---
paths:
  - "**/*.ts"
  - "**/*.tsx"
---
```

### Enforcement Mapping

ARM enforcement levels map to documentation only (Claude Code doesn't have built-in enforcement):
- `must` → "MUST" in metadata
- `should` → "SHOULD" in metadata  
- `may` → "MAY" in metadata

## Promptset Compilation

### File Naming
- **Compiled from ARM YAML**: `{promptsetId}_{promptKey}.md`
- **Raw .md files in package**: Keep original filename
- Example compiled: `codeReview_security.md`
- Example raw: `my-prompt.md` → `my-prompt.md`

### File Content Structure

```markdown
Prompt body content here.
```

**No Frontmatter or Metadata**
- Prompts are just the body content
- No ARM metadata section
- No Claude Code frontmatter
- This matches existing behavior for all tools (Cursor, AmazonQ, Copilot, Markdown)

## Skillset Compilation

### Resource Schema

```yaml
apiVersion: v1
kind: Skillset
metadata:
  id: pdf-processor
  name: PDF Processor
  description: Extract and analyze PDF content
spec:
  skills:
    process-pdf:                          # Skill key (becomes directory name)
      name: process-pdf                   # Optional: skill name
      description: Process PDF files      # Optional: skill description
      arguments: "[file] [format]"        # Optional: argument hint for documentation
      visibility: visible                 # Optional: visible | hidden (default: visible)
      body: |                             # Required: SKILL.md content
        Use the bundled script to process PDFs.
        
        Run: `python scripts/process.py input.pdf`
      files:                              # Optional: supporting files
        - path: scripts/process.py        # Relative path within skill directory
          source: ./scripts/process.py    # Path to source file (relative to YAML)
          executable: true                # Optional: chmod +x
        - path: reference.md
          source: ./docs/reference.md     # Can reference files anywhere
        - path: templates/
          source: ./templates/            # Can reference entire directories
```

### Field Definitions

**Generic Fields** (apply across tools):
- `name` - Skill identifier
- `description` - What the skill does
- `arguments` - Hint for expected arguments (e.g., `[file]`, `[issue-number]`, `[component] [from] [to]`)
- `visibility` - Whether skill is user-facing
  - `visible` (default) - Show in menus/lists
  - `hidden` - Background knowledge, not directly invocable
- `body` - Skill instructions/content
- `files` - Supporting files and directories

**Tool-Specific Behavior**:
- Claude Code maps `visibility: hidden` → `user-invocable: false`
- Claude Code maps `arguments` → `argument-hint`
- Other tools may interpret these fields differently or ignore them

**Not Supported** (use raw `.md` files if needed):
- `disable-model-invocation` - Too tool-specific (AI decision-making)
- `allowed-tools` - Too tool-specific (permission models vary)
- `model` - Too tool-specific (not all tools have model selection)
- `context` / `execution` - Too tool-specific (execution models vary)
- `agent` - Too tool-specific (agent types vary by tool)

### Directory Structure

**Source repository:**
```
my-skills/
├── pdf-processor.yml          # Skillset YAML
├── scripts/
│   └── process.py
├── docs/
│   └── reference.md
└── templates/
    └── output.html
```

**Package in registry** (pdf-processor-1.0.0.tar.gz):
```
pdf-processor.yml              # The skillset resource
scripts/
  └── process.py
docs/
  └── reference.md
templates/
  └── output.html
```

**After installation to sink:**
```
.claude/skills/
└── arm/
    └── registry/
        └── package/
            └── 1.0.0/
                └── pdf-processor/         # Archive name becomes directory
                    └── process-pdf/       # Skill key as directory
                        ├── SKILL.md       # Generated during compilation
                        ├── scripts/
                        │   └── process.py # Copied from package
                        ├── reference.md   # Copied from package
                        └── templates/
                            └── output.html
```

**Note**: If the package is delivered as an archive (`.tar.gz`, `.zip`), the archive name (minus extension) becomes a directory in the installation path. This prevents naming collisions between multiple archives and preserves file context for path resolution.

### SKILL.md Generation

```markdown
---
name: process-pdf
description: Process PDF files
argument-hint: "[file] [format]"
user-invocable: false
---

---
namespace: registry/package@1.0.0
skillset:
  id: pdf-processor
  name: PDF Processor
  description: Extract and analyze PDF content
  skills:
    - process-pdf
skill:
  id: process-pdf
  name: process-pdf
  description: Process PDF files
  arguments: "[file] [format]"
  visibility: hidden
  files:
    - scripts/process.py
    - reference.md
    - templates/
---

Use the bundled script to process PDFs.

Run: `python scripts/process.py input.pdf`
```

**Frontmatter Sections:**

1. **Claude Code Frontmatter** (optional)
   - `name`: From skill `name` field (or skill key if not specified)
   - `description`: From skill `description` field (if specified)
   - `argument-hint`: From skill `arguments` field (if specified)
   - `user-invocable`: `false` if `visibility: hidden`, omitted otherwise

2. **ARM Metadata** (always included)
   - Namespace, skillset info, skill info for traceability
   - Includes generic fields: `arguments`, `visibility`
   - Includes `files` list showing supporting files installed by ARM

### File Resolution

**At package time** (when creating the package):
- `source` paths are relative to the YAML file location
- Read file/directory content from `source`
- Package includes the content at the specified `path` location

**At install time** (when installing from registry):
- Extract package (if archive, extract to subdirectory named after archive)
- Discover all `.yml` files (including in extracted archives)
- For each skillset YAML:
  - Resolve `source` paths relative to YAML location
  - For each skill in the skillset:
    - Create skill directory: `.claude/skills/arm/registry/package/version/{archiveName}/{skillKey}/`
    - Generate `SKILL.md` from body and frontmatter fields
    - Copy files from package to skill directory at specified `path`
    - Set executable bit if `executable: true`

**Archive Handling** (v4.0+):
- Archives extract to subdirectory named after archive (minus extension)
- Example: `skills.tar.gz` → `skills/` directory
- This prevents naming collisions and preserves file context
- Path resolution works relative to YAML location within archive

**Example with archive:**
```
Package: pdf-processor-1.0.0.tar.gz
  - pdf-processor.yml
  - scripts/process.py

After extraction:
  - pdf-processor/
    - pdf-processor.yml
    - scripts/process.py

YAML location: pdf-processor/pdf-processor.yml
source: ./scripts/process.py resolves to: pdf-processor/scripts/process.py

Installed to: .claude/skills/arm/registry/package/1.0.0/pdf-processor/process-pdf/scripts/process.py
```

### Supported Claude Code Frontmatter Fields

Generic fields that map to Claude Code frontmatter:
- `name` → `name`
- `description` → `description`
- `arguments` → `argument-hint`
- `visibility: hidden` → `user-invocable: false`

**Not Supported** (use raw `.md` files if needed):
- `disable-model-invocation` - Too tool-specific (AI decision-making)
- `allowed-tools` - Too tool-specific (permission models vary)
- `model` - Too tool-specific (not all tools have model selection)
- `context` / `execution` - Too tool-specific (execution models vary)
- `agent` - Too tool-specific (agent types vary by tool)

For tool-specific control beyond the generic fields, package raw `SKILL.md` files with custom frontmatter instead of using ARM YAML compilation.

### File Path Validation

- `path` must be relative (no leading `/`)
- `path` cannot contain `..` (no directory traversal)
- `source` is relative to YAML file location
- `source` can be a file or directory
- If `source` is a directory, entire tree is copied to `path`
- `executable` flag only applies to files, not directories

## Priority Resolution

### arm_index.md

Generate `arm_index.md` in `.claude/rules/` directory:

```markdown
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### registry/package@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/registry/package/1.0.0/rules/rulesetId_ruleOne.md
  - arm/registry/package/1.0.0/rules/rulesetId_ruleTwo.md

### registry/other-package@2.0.0
- **Priority:** 100
- **Rules:**
  - arm/registry/other-package/2.0.0/rules/otherId_ruleOne.md
```

**Location**: `.claude/rules/arm_index.md`

**When to Generate**:
- After installing any ruleset
- After uninstalling any ruleset
- Remove when no rulesets remain

**Sorting**:
- Descending by priority (highest first)
- Same priority: installation order (later first)

## Sink Configuration

### Default Directories

**Rules Sink**:
```bash
arm add sink --tool claudecode claude-rules .claude/rules
```

**Skills Sink** (for both promptsets and skillsets):
```bash
arm add sink --tool claudecode claude-skills .claude/skills
```

**Note**: Both promptsets and skillsets install to the same sink type (skills), but:
- Promptsets compile to individual `.md` files
- Skillsets compile to directories with `SKILL.md` + supporting files

### Tool Name

Use `claudecode` as the tool identifier (lowercase, no spaces).

## Implementation Checklist

### Parser (`internal/arm/parser/`)

- [ ] Add Skillset kind support to parser
- [ ] Parse `spec.skills` map
- [ ] Parse skill fields: name, description, body, files
- [ ] Parse Claude Code frontmatter fields
- [ ] Parse file entries: path, source, executable
- [ ] Validate file paths (no `..`, no absolute paths)

### Compiler (`internal/arm/compiler/claudecode.go`)

- [ ] Implement `ClaudeCodeCompiler` struct
- [ ] Implement `CompileRuleset()` method
  - [ ] Generate filename: `{rulesetId}_{ruleKey}.md`
  - [ ] Generate Claude Code frontmatter (if scope defined)
  - [ ] Generate ARM metadata section
  - [ ] Append rule body
- [ ] Implement `CompilePromptset()` method
  - [ ] Generate filename: `{promptsetId}_{promptKey}.md`
  - [ ] Generate ARM metadata section
  - [ ] Append prompt body
- [ ] Implement `CompileSkillset()` method
  - [ ] For each skill in skillset:
    - [ ] Create skill directory: `{skillKey}/`
    - [ ] Generate `SKILL.md` with Claude Code frontmatter
    - [ ] Generate ARM metadata section
    - [ ] Append skill body
    - [ ] Copy files from package to skill directory
    - [ ] Set executable bit on files marked executable
- [ ] Implement `GeneratePriorityIndex()` method
  - [ ] Sort rulesets by priority
  - [ ] Generate markdown content
  - [ ] Return as `arm_index.md`

### File Operations

- [ ] Add file copying logic for skillset files
- [ ] Add directory copying logic for skillset directories
- [ ] Add executable bit setting (`chmod +x`)
- [ ] Validate source paths exist during packaging
- [ ] Handle file path resolution relative to YAML location

### Filename Generation

- [ ] Add `GenerateClaudeCodeRuleFilename()` function
- [ ] Add `GenerateClaudeCodePromptFilename()` function
- [ ] Follow existing truncation logic (100 char limit)

### Integration

- [ ] Add `claudecode` case to compiler factory
- [ ] Update sink validation to accept `claudecode` tool
- [ ] Update command help text with `claudecode` option

### Testing (`test/e2e/compile_test.go`)

- [ ] Test ruleset compilation
  - [ ] Verify filename format
  - [ ] Verify frontmatter structure (with and without scope)
  - [ ] Verify ARM metadata section
  - [ ] Verify body content
- [ ] Test promptset compilation
  - [ ] Verify filename format
  - [ ] Verify metadata section
  - [ ] Verify body content
- [ ] Test skillset compilation
  - [ ] Verify skill directory creation
  - [ ] Verify SKILL.md generation
  - [ ] Verify Claude Code frontmatter fields
  - [ ] Verify ARM metadata section
  - [ ] Verify supporting files copied correctly
  - [ ] Verify directory copying
  - [ ] Verify executable bit set correctly
  - [ ] Verify file path validation (reject `..`, absolute paths)
- [ ] Test priority index generation
  - [ ] Verify sorting by priority
  - [ ] Verify markdown format
  - [ ] Verify file paths
- [ ] Test cleanup
  - [ ] Verify arm_index.md removed when no rulesets
  - [ ] Verify empty directories removed
  - [ ] Verify skill directories removed on uninstall

### Documentation

- [ ] Update `docs/sinks.md` with Claude Code examples
- [ ] Update `docs/commands.md` with `claudecode` tool option
- [ ] Update `docs/resource-schemas.md` with Skillset schema
- [ ] Update README.md quick start with Claude Code example
- [ ] Add Claude Code section to `docs/concepts.md`
- [ ] Document skillset file referencing pattern
- [ ] Document executable file support

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Rule without scope | No Claude Code frontmatter, only ARM metadata |
| Rule with scope | Include both Claude Code frontmatter and ARM metadata |
| Empty scope.files | No Claude Code frontmatter |
| Multiple scope entries | Merge all file patterns into single paths array |
| Prompt with scope | Ignore scope (skills don't use paths) |
| Skillset file source doesn't exist | Error during packaging |
| Skillset file path contains `..` | Validation error |
| Skillset file path is absolute | Validation error |
| Skillset source is directory | Copy entire directory tree to path |
| Skillset executable on directory | Ignore (only applies to files) |
| Skillset with no files | Valid (just SKILL.md) |
| Skillset with duplicate paths | Error (path collision) |
| Priority index with no rulesets | Remove arm_index.md |
| Filename collision | Use hash prefix (existing logic) |

## Example Output

### Rule with Scope

**Input** (`clean-code.yml`):
```yaml
apiVersion: v1
kind: Ruleset
metadata:
  id: cleanCode
  name: Clean Code
  description: Clean code practices
spec:
  rules:
    formatting:
      name: Code Formatting
      description: Format code consistently
      enforcement: must
      scope:
        - files: ["**/*.ts", "**/*.tsx"]
      body: |
        Always use 2 spaces for indentation.
```

**Output** (`.claude/rules/arm/registry/package/1.0.0/cleanCode_formatting.md`):
```markdown
---
paths:
  - "**/*.ts"
  - "**/*.tsx"
---

---
namespace: registry/package@1.0.0
ruleset:
  id: cleanCode
  name: Clean Code
  description: Clean code practices
  rules:
    - formatting
rule:
  id: formatting
  name: Code Formatting
  description: Format code consistently
  enforcement: MUST
  scope:
    - files: ["**/*.ts", "**/*.tsx"]
---

# Code Formatting

Always use 2 spaces for indentation.
```

### Rule without Scope

**Input** (`security.yml`):
```yaml
apiVersion: v1
kind: Ruleset
metadata:
  id: security
  name: Security Rules
spec:
  rules:
    auth:
      name: Authentication
      enforcement: must
      body: |
        Always validate user input.
```

**Output** (`.claude/rules/arm/registry/package/1.0.0/security_auth.md`):
```markdown
---
namespace: registry/package@1.0.0
ruleset:
  id: security
  name: Security Rules
  rules:
    - auth
rule:
  id: auth
  name: Authentication
  enforcement: MUST
---

# Authentication

Always validate user input.
```

### Promptset (Simple Skill)

**Input** (`code-review.yml`):
```yaml
apiVersion: v1
kind: Promptset
metadata:
  id: codeReview
  name: Code Review
  description: Review code for quality
spec:
  prompts:
    security:
      name: Security Review
      description: Check for security issues
      body: |
        Review this code for security vulnerabilities.
```

**Output** (`.claude/skills/arm/registry/package/1.0.0/codeReview_security.md`):
```markdown
Review this code for security vulnerabilities.
```

**Note**: Prompts are just the body content with no metadata or frontmatter.

### Skillset (Complex Skill with Files)

**Input** (`pdf-processor.yml`):
```yaml
apiVersion: v1
kind: Skillset
metadata:
  id: pdfProcessor
  name: PDF Processor
  description: Extract and analyze PDF content
spec:
  skills:
    process-pdf:
      name: process-pdf
      description: Process PDF files with Python script
      arguments: "[input-file] [output-format]"
      visibility: visible
      body: |
        Use the bundled script to process PDFs.
        
        Run: `python scripts/process.py input.pdf`
        
        See reference.md for API details.
      files:
        - path: scripts/process.py
          source: ./scripts/process.py
          executable: true
        - path: reference.md
          source: ./docs/api-reference.md
        - path: templates/
          source: ./templates/
```

**Output Structure**:
```
.claude/skills/arm/registry/package/1.0.0/process-pdf/
├── SKILL.md
├── scripts/
│   └── process.py (executable)
├── reference.md
└── templates/
    ├── output.html
    └── report.md
```

**SKILL.md Content**:
```markdown
---
name: process-pdf
description: Process PDF files with Python script
argument-hint: "[input-file] [output-format]"
---

---
namespace: registry/package@1.0.0
skillset:
  id: pdfProcessor
  name: PDF Processor
  description: Extract and analyze PDF content
  skills:
    - process-pdf
skill:
  id: process-pdf
  name: process-pdf
  description: Process PDF files with Python script
  arguments: "[input-file] [output-format]"
  visibility: visible
  files:
    - scripts/process.py
    - reference.md
    - templates/
---

Use the bundled script to process PDFs.

Run: `python scripts/process.py input.pdf`

See reference.md for API details.
```

**Note**: For tool-specific frontmatter like `disable-model-invocation: true`, `allowed-tools`, or `context: fork`, package a raw `SKILL.md` file with that frontmatter instead of using ARM YAML compilation.

## Notes

- Claude Code uses `.md` extension for both rules and skills (no `.mdc`)
- Claude Code frontmatter is optional and only used for conditional loading (rules only)
- ARM metadata is always included for traceability
- Priority index uses same format as other tools but with `.md` extension
- Skills don't support `paths` frontmatter (they're invoked explicitly)
- Skillsets are only supported for `claudecode` tool (not cursor, amazonq, copilot, markdown)
- Skillset files are referenced by `source` path, not embedded inline
- Skillset `source` paths are relative to the YAML file location
- Skillset `path` is relative to the skill directory root
- Executable bit is preserved for files marked with `executable: true`
- Directory sources copy entire tree to destination path
