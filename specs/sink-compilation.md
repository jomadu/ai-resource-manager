# Sink Compilation

## Job to be Done
Compile ARM resources to tool-specific formats and write them to configured output directories, enabling AI tools to consume rules and prompts in their native formats.

## Activities
1. **Compile Resources** - Transform ARM resource definitions (rulesets, promptsets) into tool-specific file formats
2. **Generate Filenames** - Create appropriate filenames based on tool conventions and layout mode
3. **Write Files** - Write compiled resources to disk in hierarchical or flat layout
4. **Manage Index** - Track installed packages and their files in JSON index
5. **Generate Priority Index** - Create human-readable priority documentation for rulesets

## Acceptance Criteria
- [ ] ARM resources compiled to correct tool-specific formats (Cursor, AmazonQ, Copilot, Markdown)
- [ ] Non-resource files copied as-is to preserve original content
- [ ] Filenames generated according to tool conventions with proper extensions
- [ ] Files written to correct paths based on layout mode (hierarchical vs flat)
- [ ] Index file tracks all installed packages and their files
- [ ] Priority index rule file generated for rulesets with sorted priority order
- [ ] Long filenames truncated progressively in flat layout (full path → filename only → truncated filename)
- [ ] Package and path hashes prevent collisions in flat layout
- [ ] Metadata embedded in compiled rules for traceability
- [ ] Empty rulesets remove priority index file
- [ ] Orphaned files cleaned up when packages uninstalled
- [ ] Empty directories removed when last package uninstalled from sink
- [ ] arm-index.json removed when all packages uninstalled
- [ ] arm_index.* files removed when all rulesets/promptsets uninstalled

## Data Structures

### Layout Mode
```go
type Layout string

const (
    LayoutFlat         Layout = "flat"         // Single directory with hash-prefixed names
    LayoutHierarchical Layout = "hierarchical" // Preserves directory structure
)
```

**Layout Selection:**
- Copilot → Flat (required by tool)
- Cursor, AmazonQ, Markdown → Hierarchical (default)

### Index Structure
```json
{
  "version": 1,
  "rulesets": {
    "registry/package@version": {
      "priority": 100,
      "files": ["arm/registry/package/version/rules/file.mdc"]
    }
  },
  "promptsets": {
    "registry/package@version": {
      "files": ["arm/registry/package/version/prompts/file.md"]
    }
  }
}
```

**Fields:**
- `version` - Index schema version (currently 1)
- `rulesets` - Map of package keys to ruleset entries
- `promptsets` - Map of package keys to promptset entries
- `priority` - Installation priority for conflict resolution (rulesets only)
- `files` - List of relative file paths from sink directory

### Tool-Specific Formats

#### Cursor
- **Rules**: `.mdc` extension with YAML frontmatter
- **Prompts**: `.md` extension, plain markdown (no metadata)
- **Frontmatter fields**: `description`, `globs` (from scope), `alwaysApply` (if enforcement=must)

#### AmazonQ
- **Rules**: `.md` extension with embedded metadata
- **Prompts**: `.md` extension, plain markdown (no metadata)
- **Metadata**: YAML block at top of file

#### Copilot
- **Rules**: `.instructions.md` extension with embedded metadata
- **Prompts**: `.md` extension, plain markdown (no metadata)
- **Metadata**: YAML block at top of file

#### Markdown
- **Rules**: `.md` extension with embedded metadata
- **Prompts**: `.md` extension, plain markdown (no metadata)
- **Metadata**: YAML block at top of file

### Embedded Metadata Format
```yaml
---
namespace: registry/package@version
ruleset:
  id: rulesetID
  name: Ruleset Name
  rules:
    - ruleOne
    - ruleTwo
rule:
  id: ruleOne
  name: Rule Name
  enforcement: MUST
  priority: 200
  scope:
    - files: ["*.go", "*.ts"]
---
```

**Fields:**
- `namespace` - Full package identifier with version
- `ruleset.id` - Ruleset resource ID
- `ruleset.name` - Ruleset display name
- `ruleset.rules` - All rule IDs in sorted order
- `rule.id` - Current rule ID
- `rule.name` - Rule display name
- `rule.enforcement` - Uppercase enforcement level (MUST, SHOULD, MAY)
- `rule.priority` - Optional priority value (omitted if 0)
- `rule.scope` - Optional file globs for scoping

## Algorithm

### 1. Compile Ruleset

**Input:** Package, Priority
**Output:** Installed files in sink

```
function CompileRuleset(pkg, priority):
    // Uninstall existing versions first
    Uninstall(pkg.registry, pkg.name)
    
    installedFiles = []
    namespace = "{registry}/{name}@{version}"
    
    for file in pkg.files:
        if IsResourceFile(file):
            if IsRulesetFile(file):
                ruleset = ParseRuleset(file)
                
                // Compile each rule separately
                for ruleID in sorted(ruleset.rules):
                    content = GenerateRule(tool, namespace, ruleset, ruleID)
                    filename = GenerateRuleFilename(tool, ruleset.id, ruleID)
                    
                    fullPath = GetFilePath(pkg, file.dir + filename)
                    WriteFile(fullPath, content)
                    installedFiles.append(RelativePath(fullPath))
        else:
            // Copy non-resource files as-is
            fullPath = GetFilePath(pkg, file.path)
            WriteFile(fullPath, file.content)
            installedFiles.append(RelativePath(fullPath))
    
    // Update index
    index = LoadIndex()
    index.rulesets[PkgKey(pkg)] = {priority, installedFiles}
    SaveIndex(index)
    
    // Generate priority documentation
    GenerateRulesetIndexRuleFile()
```

### 2. Compile Promptset

**Input:** Package
**Output:** Installed files in sink

```
function CompilePromptset(pkg):
    // Uninstall existing versions first
    Uninstall(pkg.registry, pkg.name)
    
    installedFiles = []
    namespace = "{registry}/{name}@{version}"
    
    for file in pkg.files:
        if IsResourceFile(file):
            if IsPromptsetFile(file):
                promptset = ParsePromptset(file)
                
                // Compile each prompt separately
                for promptID in sorted(promptset.prompts):
                    content = GeneratePrompt(tool, namespace, promptset, promptID)
                    filename = GeneratePromptFilename(tool, promptset.id, promptID)
                    
                    fullPath = GetFilePath(pkg, file.dir + filename)
                    WriteFile(fullPath, content)
                    installedFiles.append(RelativePath(fullPath))
        else:
            // Copy non-resource files as-is
            fullPath = GetFilePath(pkg, file.path)
            WriteFile(fullPath, file.content)
            installedFiles.append(RelativePath(fullPath))
    
    // Update index
    index = LoadIndex()
    index.promptsets[PkgKey(pkg)] = {installedFiles}
    SaveIndex(index)
```

### 3. Generate File Path

**Input:** Package metadata, relative path
**Output:** Full filesystem path

```
function GetFilePath(registry, name, version, relativePath):
    if layout == Flat:
        return GetFlatPath(registry, name, version, relativePath)
    else:
        return GetHierarchicalPath(registry, name, version, relativePath)

function GetHierarchicalPath(registry, name, version, relativePath):
    return "{sinkDir}/arm/{registry}/{name}/{version}/{relativePath}"

function GetFlatPath(registry, name, version, relativePath):
    packageHash = Hash(registry + "/" + name + "@" + version)[:4]
    pathHash = Hash(relativePath)[:4]
    pathPart = relativePath.replace("/", "_").replace("\\", "_")
    
    // Truncate if needed (max 100 chars total)
    filenameOverhead = len("arm_xxxx_xxxx_")
    maxPathLen = 100 - filenameOverhead
    
    if len(pathPart) > maxPathLen:
        filename = Basename(relativePath)
        ext = Extension(filename)
        nameWithoutExt = filename.trimSuffix(ext)
        
        if len(filename) <= maxPathLen:
            // Option 2: Use just filename
            pathPart = filename
        else:
            // Option 3: Truncate filename, keep extension
            availableForName = maxPathLen - len(ext)
            pathPart = nameWithoutExt[:availableForName] + ext
    
    return "{sinkDir}/arm_{packageHash}_{pathHash}_{pathPart}"
```

### 4. Generate Rule Content

**Input:** Tool, namespace, ruleset, ruleID
**Output:** Compiled rule content

```
function GenerateRule(tool, namespace, ruleset, ruleID):
    rule = ruleset.rules[ruleID]
    
    if tool == Cursor:
        frontmatter = GenerateCursorFrontmatter(rule)
        metadata = GenerateRuleMetadata(namespace, ruleset, ruleID, rule)
        return frontmatter + "\n\n" + metadata + "\n\n" + rule.body
    else:
        // AmazonQ, Copilot, Markdown use same format
        metadata = GenerateRuleMetadata(namespace, ruleset, ruleID, rule)
        return metadata + "\n\n" + rule.body

function GenerateCursorFrontmatter(rule):
    parts = ["---"]
    
    if rule.description != "":
        parts.append('description: "' + rule.description + '"')
    
    if len(rule.scope) > 0 && len(rule.scope[0].files) > 0:
        globs = join(rule.scope[0].files, ", ")
        parts.append("globs: " + globs)
    
    if rule.enforcement == "must":
        parts.append("alwaysApply: true")
    
    parts.append("---")
    return join(parts, "\n")

function GenerateRuleMetadata(namespace, ruleset, ruleID, rule):
    // See Embedded Metadata Format above
    // Returns YAML block with namespace, ruleset, and rule info
```

### 5. Generate Filename

**Input:** Tool, resource ID, item ID
**Output:** Filename with extension

```
function GenerateRuleFilename(tool, rulesetID, ruleID):
    basename = rulesetID + "_" + ruleID
    
    switch tool:
        case Cursor:   return basename + ".mdc"
        case AmazonQ:  return basename + ".md"
        case Copilot:  return basename + ".instructions.md"
        case Markdown: return basename + ".md"

function GeneratePromptFilename(tool, promptsetID, promptID):
    basename = promptsetID + "_" + promptID
    
    // All tools use .md for prompts
    return basename + ".md"
```

### 6. Generate Priority Index

**Input:** Index with rulesets
**Output:** Human-readable priority documentation

```
function GenerateRulesetIndexRuleFile():
    index = LoadIndex()
    
    if len(index.rulesets) == 0:
        // Remove index file if no rulesets
        DeleteFile(rulesetIndexRulePath)
        return
    
    // Sort rulesets by priority (high to low)
    entries = []
    for key, info in index.rulesets:
        entries.append({key, info.priority, info.files})
    
    sort(entries, by=priority, descending=true)
    
    // Generate markdown content
    content = "# ARM Rulesets\n\n"
    content += "This file defines the installation priorities...\n\n"
    content += "## Priority Rules\n\n"
    content += "1. Higher priority numbers take precedence\n"
    content += "2. Rules from higher priority rulesets override...\n\n"
    content += "## Installed Rulesets\n\n"
    
    for entry in entries:
        content += "### " + entry.key + "\n"
        content += "- **Priority:** " + entry.priority + "\n"
        content += "- **Rules:**\n"
        for file in entry.files:
            content += "  - " + file + "\n"
        content += "\n"
    
    WriteFile(rulesetIndexRulePath, content)
```

### 7. Uninstall Package

**Input:** Registry name, package name
**Output:** Files removed, index updated, empty directories cleaned

```
function Uninstall(registryName, packageName):
    index = LoadIndex()
    
    // Find package key in index
    packageKey = FindPackageKey(index, registryName, packageName)
    if packageKey == "":
        return  // Package not installed, nothing to do
    
    // Determine resource type and get files
    files = []
    if packageKey in index.rulesets:
        files = index.rulesets[packageKey].files
        delete index.rulesets[packageKey]
    else if packageKey in index.promptsets:
        files = index.promptsets[packageKey].files
        delete index.promptsets[packageKey]
    
    // Delete all files
    for filePath in files:
        fullPath = sinkDir + "/" + filePath
        DeleteFile(fullPath)
    
    // Clean up empty directories
    CleanupEmptyDirectories(sinkDir)
    
    // Clean up index files if all packages uninstalled
    if len(index.rulesets) == 0 and len(index.promptsets) == 0:
        DeleteFile(sinkDir + "/arm-index.json")
    else:
        SaveIndex(index)
    
    // Regenerate priority index (removes if no rulesets)
    GenerateRulesetIndexRuleFile()
```

### 8. Cleanup Empty Directories

**Input:** Sink directory
**Output:** Empty directories removed

```
function CleanupEmptyDirectories(sinkDir):
    // Walk directory tree bottom-up
    dirs = WalkDirectoriesBottomUp(sinkDir)
    
    for dir in dirs:
        // Skip sink root directory
        if dir == sinkDir:
            continue
        
        // Check if directory is empty
        entries = ListDirectory(dir)
        if len(entries) == 0:
            RemoveDirectory(dir)
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Empty ruleset (no rules) | Skip compilation, update index with empty files list |
| Empty promptset (no prompts) | Skip compilation, update index with empty files list |
| Non-resource files in package | Copy as-is to preserve original content and structure |
| Long filename in flat layout | Truncate progressively: full path → filename only → truncated filename with extension |
| Package hash collision | Extremely unlikely (SHA256 truncated to 4 chars = 65k combinations), no special handling |
| Path hash collision | Extremely unlikely, no special handling (different packages have different package hashes) |
| Missing rule in ruleset | Return error during compilation |
| Missing prompt in promptset | Return error during compilation |
| Invalid YAML in resource | Return error during parsing (handled by parser layer) |
| Special characters in filename | Replaced with underscores in flat layout |
| No rulesets installed | Remove priority index file to avoid clutter |
| Reinstall same package | Uninstall old version first, then install new version |
| Orphaned files after uninstall | Cleaned up by Clean() operation |
| Last package uninstalled | Remove arm-index.json and all arm_index.* files |
| Empty directories after uninstall | Recursively remove empty directories bottom-up |

## Dependencies

**Prerequisites:**
- Sink configured with directory and tool type
- Package downloaded and extracted to storage
- Resource files parsed successfully

**Required Components:**
- `internal/arm/compiler/` - Tool-specific generators
- `internal/arm/parser/` - Resource parsing
- `internal/arm/filetype/` - File type detection
- `internal/arm/core/` - Package and file structures

**Related Specs:**
- `package-installation.md` - Calls sink compilation during install
- `priority-resolution.md` - Uses priority index for conflict resolution

## Implementation Mapping

**Source files:**
- `internal/arm/sink/manager.go` - Main sink manager with install/uninstall/compile logic
- `internal/arm/compiler/compiler.go` - High-level compilation functions
- `internal/arm/compiler/cursor.go` - Cursor-specific generators
- `internal/arm/compiler/amazonq.go` - AmazonQ-specific generators
- `internal/arm/compiler/copilot.go` - Copilot-specific generators
- `internal/arm/compiler/markdown.go` - Markdown-specific generators
- `internal/arm/compiler/generators.go` - Shared metadata generation
- `internal/arm/compiler/factory.go` - Factory pattern for generator creation
- `internal/arm/compiler/types.go` - Interface definitions

**Related specs:**
- `package-installation.md` - Installation workflow
- `priority-resolution.md` - Priority-based conflict resolution

## Examples

### Example 1: Hierarchical Layout (Cursor)

**Input:**
```bash
arm add sink --tool cursor cursor-rules .cursor/rules
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

**Package Contents:**
- `rules/sample-ruleset.yml` (ARM resource)
  - Resource ID: `cleanCode`
  - Rules: `ruleOne`, `ruleTwo`

**Expected Output:**
```
.cursor/rules/
└── arm/
    ├── ai-rules/
    │   └── clean-code-ruleset/
    │       └── 1.0.0/
    │           └── rules/
    │               ├── cleanCode_ruleOne.mdc
    │               └── cleanCode_ruleTwo.mdc
    ├── arm_index.mdc
    └── arm-index.json
```

**File: cleanCode_ruleOne.mdc**
```markdown
---
description: "First rule description"
globs: *.go, *.ts
alwaysApply: true
---

---
namespace: ai-rules/clean-code-ruleset@1.0.0
ruleset:
  id: cleanCode
  name: Clean Code Rules
  rules:
    - ruleOne
    - ruleTwo
rule:
  id: ruleOne
  name: Rule One
  enforcement: MUST
  priority: 100
  scope:
    - files: ["*.go", "*.ts"]
---

Rule body content here...
```

**Verification:**
- Files created in hierarchical structure under `arm/`
- Cursor frontmatter includes description, globs, alwaysApply
- Metadata block includes namespace and rule info
- Index tracks files and priority
- Priority index file created

### Example 2: Flat Layout (Copilot)

**Input:**
```bash
arm add sink --tool copilot copilot-rules .github/copilot
arm install ruleset ai-rules/clean-code-ruleset copilot-rules
```

**Expected Output:**
```
.github/copilot/
├── arm_1a2b_3c4d_rules_cleanCode_ruleOne.instructions.md
├── arm_1a2b_5e6f_rules_cleanCode_ruleTwo.instructions.md
├── arm_index.instructions.md
└── arm-index.json
```

**Filename Breakdown:**
- `arm_` - Prefix for ARM-managed files
- `1a2b` - Package hash (ai-rules/clean-code-ruleset@1.0.0)
- `3c4d` - Path hash (rules/cleanCode_ruleOne.instructions.md)
- `rules_cleanCode_ruleOne` - Flattened path
- `.instructions.md` - Copilot extension

**File: arm_1a2b_3c4d_rules_cleanCode_ruleOne.instructions.md**
```markdown
---
namespace: ai-rules/clean-code-ruleset@1.0.0
ruleset:
  id: cleanCode
  name: Clean Code Rules
  rules:
    - ruleOne
    - ruleTwo
rule:
  id: ruleOne
  name: Rule One
  enforcement: MUST
  priority: 100
---

Rule body content here...
```

**Verification:**
- Files created in flat structure with hash prefixes
- No Cursor frontmatter (Copilot doesn't use it)
- Metadata block includes namespace and rule info
- Filenames limited to 100 characters
- Package hash groups files from same package

### Example 3: Long Filename Truncation

**Input:**
```
Package: ai-rules/very-long-package-name@1.0.0
File: rules/subdirectory/another-subdirectory/very-long-ruleset-name_very-long-rule-name.instructions.md
```

**Truncation Steps:**

1. **Full path (too long):**
   ```
   arm_1a2b_3c4d_rules_subdirectory_another-subdirectory_very-long-ruleset-name_very-long-rule-name.instructions.md
   (120 characters - exceeds 100 limit)
   ```

2. **Filename only (still too long):**
   ```
   arm_1a2b_3c4d_very-long-ruleset-name_very-long-rule-name.instructions.md
   (75 characters - fits!)
   ```

**Final Output:**
```
arm_1a2b_3c4d_very-long-ruleset-name_very-long-rule-name.instructions.md
```

**Verification:**
- Full path attempted first
- Fell back to filename only
- Package and path hashes preserved
- Extension preserved
- Total length ≤ 100 characters

### Example 4: Non-Resource Files

**Input:**
```
Package: ai-rules/grug-brained-dev@1.0.0
Files:
  - rules/avoid-abstractions.md (plain markdown)
  - rules/complexity-enemy.md (plain markdown)
  - README.md (plain markdown)
```

**Expected Output (Hierarchical):**
```
.cursor/rules/
└── arm/
    └── ai-rules/
        └── grug-brained-dev/
            └── 1.0.0/
                ├── rules/
                │   ├── avoid-abstractions.md
                │   └── complexity-enemy.md
                └── README.md
```

**Verification:**
- Non-resource files copied as-is
- Original directory structure preserved
- No compilation or transformation
- Files tracked in index

### Example 5: Priority Index Generation

**Input:**
```bash
arm install ruleset --priority 200 ai-rules/team-standards cursor-rules
arm install ruleset --priority 100 ai-rules/clean-code-ruleset cursor-rules
arm install ruleset --priority 50 ai-rules/experimental-rules cursor-rules
```

**Expected Output: arm_index.mdc**
```markdown
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### ai-rules/team-standards@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/ai-rules/team-standards/1.0.0/rules/teamStandards_rule1.mdc
  - arm/ai-rules/team-standards/1.0.0/rules/teamStandards_rule2.mdc

### ai-rules/clean-code-ruleset@1.0.0
- **Priority:** 100
- **Rules:**
  - arm/ai-rules/clean-code-ruleset/1.0.0/rules/cleanCode_ruleOne.mdc
  - arm/ai-rules/clean-code-ruleset/1.0.0/rules/cleanCode_ruleTwo.mdc

### ai-rules/experimental-rules@1.0.0
- **Priority:** 50
- **Rules:**
  - arm/ai-rules/experimental-rules/1.0.0/rules/experimental_rule1.mdc
```

**Verification:**
- Rulesets sorted by priority (high to low)
- Clear documentation of priority rules
- All installed files listed
- File removed when no rulesets installed

## Notes

**Why separate files per rule/prompt?**
- Enables fine-grained control and scoping
- Allows AI tools to load only relevant rules
- Simplifies priority resolution (file-level granularity)
- Matches tool expectations (Cursor, Copilot)

**Why embed metadata?**
- Traceability: Know which package/version a rule came from
- Priority resolution: AI agents can understand conflict resolution
- Debugging: Easier to diagnose issues with rule application

**Why two layout modes?**
- Copilot requires flat layout (tool limitation)
- Hierarchical is more human-readable and maintainable
- Both modes tracked in same index format

**Why hash-based naming in flat layout?**
- Prevents filename collisions between packages
- Groups files from same package (same package hash)
- Enables safe uninstall (remove all files with package hash)
- Handles special characters and long paths

**Why progressive truncation?**
- Maximizes information in filename while respecting limits
- Preserves extension for tool recognition
- Keeps hashes for collision prevention
- Balances readability with technical constraints

**Why priority index rule file?**
- AI agents need clear guidance on conflict resolution
- Human-readable documentation of installation priorities
- Single source of truth for priority order
- Automatically updated on install/uninstall
