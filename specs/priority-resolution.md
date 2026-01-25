# Priority Resolution

## Job to be Done
Resolve conflicts between overlapping rules from multiple installed rulesets, enabling users to layer rules with clear precedence (e.g., team standards override community best practices).

## Activities
1. **Assign Priority** - Set priority value during ruleset installation or update
2. **Track Priority** - Store priority in sink index alongside installed files
3. **Generate Priority Index** - Create human-readable documentation explaining priority order
4. **Embed Metadata** - Include priority in compiled rule metadata for AI tool consumption
5. **Sort by Priority** - Order rulesets from highest to lowest priority in index file

## Acceptance Criteria
- [ ] Priority assigned during installation (default: 100)
- [ ] Priority stored in arm-index.json per ruleset
- [ ] Priority updatable via `arm set ruleset` command
- [ ] arm_index.* file generated with rulesets sorted by priority (high to low)
- [ ] Priority embedded in compiled rule metadata
- [ ] Index file explains resolution rules clearly for AI tools
- [ ] Empty rulesets remove priority index file
- [ ] Priority only applies to rulesets (not promptsets)
- [ ] Higher priority numbers take precedence over lower numbers
- [ ] Index file regenerated after any ruleset install/uninstall/update

## Data Structures

### Priority Value
```go
type Priority int

const DefaultPriority = 100
```

**Constraints:**
- Integer value (can be negative, zero, or positive)
- Default: 100
- Higher values = higher precedence
- No upper/lower bounds enforced

### Index Entry
```json
{
  "rulesets": {
    "registry/package@version": {
      "priority": 100,
      "files": ["path/to/rule.mdc"]
    }
  }
}
```

**Fields:**
- `priority` - Integer priority value for conflict resolution
- `files` - List of installed file paths (relative to sink directory)

### Priority Index File Content
```markdown
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### registry/package@version
- **Priority:** 200
- **Rules:**
  - path/to/rule1.mdc
  - path/to/rule2.mdc

### registry/other@version
- **Priority:** 100
- **Rules:**
  - path/to/rule3.mdc
```

**Purpose:**
- Provides AI tools with explicit priority ordering
- Documents which rules take precedence
- Enables manual inspection of priority configuration

### Embedded Rule Metadata
```yaml
---
namespace: registry/package@version
ruleset:
  id: rulesetID
  name: Ruleset Name
rule:
  id: ruleOne
  description: Rule description
  priority: 100
---
```

**Fields:**
- `priority` - Copied from ruleset installation priority (only if > 0)
- Embedded in compiled rule files for traceability

## Algorithm

### Install Ruleset with Priority

1. **Accept priority parameter** (default: 100)
2. **Uninstall existing versions** of same package
3. **Compile and write rules** to sink directory
4. **Update index** with priority and file list
5. **Regenerate priority index file** with sorted rulesets

**Pseudocode:**
```
function InstallRuleset(pkg, priority):
    if priority == 0:
        priority = DefaultPriority
    
    Uninstall(pkg.registry, pkg.name)
    
    installedFiles = []
    for file in pkg.files:
        if IsRulesetFile(file):
            for ruleID in file.rules:
                content = GenerateRule(file, ruleID, priority)
                path = WriteFile(content)
                installedFiles.append(path)
    
    index.rulesets[pkg.key] = {
        priority: priority,
        files: installedFiles
    }
    SaveIndex(index)
    GeneratePriorityIndexFile(index)
```

### Update Ruleset Priority

1. **Validate package exists** in manifest
2. **Load current index**
3. **Update priority value** in index entry
4. **Save index**
5. **Regenerate priority index file**
6. **Recompile rules** with new priority metadata

**Pseudocode:**
```
function SetRulesetPriority(registry, name, priority):
    if not ExistsInManifest(registry, name):
        return error("package not installed")
    
    index = LoadIndex()
    key = FindPackageKey(index, registry, name)
    
    if key not found:
        return error("package not in index")
    
    index.rulesets[key].priority = priority
    SaveIndex(index)
    GeneratePriorityIndexFile(index)
    
    // Recompile to update embedded metadata
    pkg = LoadPackage(registry, name)
    RecompileRules(pkg, priority)
```

### Generate Priority Index File

1. **Load index** from disk
2. **Check if rulesets exist** (if empty, delete index file and return)
3. **Build header** with priority rules explanation
4. **Collect entries** (key, priority, files)
5. **Sort entries** by priority (high to low)
6. **Format markdown** with sections per ruleset
7. **Write file** to sink directory

**Pseudocode:**
```
function GeneratePriorityIndexFile(index):
    if len(index.rulesets) == 0:
        DeleteFile(indexRulePath)
        return
    
    entries = []
    for key, info in index.rulesets:
        entries.append({
            key: key,
            priority: info.priority,
            files: info.files
        })
    
    // Bubble sort by priority (high to low)
    for i in 0..len(entries):
        for j in i+1..len(entries):
            if entries[j].priority > entries[i].priority:
                swap(entries[i], entries[j])
    
    content = BuildMarkdown(entries)
    WriteFile(indexRulePath, content)
```

### Embed Priority in Rule Metadata

1. **Check if priority > 0** (skip if default/zero)
2. **Add priority field** to metadata YAML
3. **Include in frontmatter** or metadata block

**Pseudocode:**
```
function GenerateRuleMetadata(rule):
    metadata = BuildBaseMetadata(rule)
    
    if rule.priority > 0:
        metadata.append("  priority: " + rule.priority)
    
    return metadata
```

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Same priority for multiple rulesets | Undefined precedence order. Users should avoid this. Index file lists in arbitrary order. |
| No rulesets installed | Priority index file deleted. No conflicts possible. |
| Single ruleset installed | Priority stored but irrelevant. Index file still generated for consistency. |
| Priority = 0 | Treated as default (100). Not embedded in metadata. |
| Negative priority | Allowed. Lower precedence than positive priorities. |
| Very large priority | No upper bound. Sorted correctly. |
| Missing priority in index | Should not occur. If corrupted, treat as default (100). |
| Promptset with priority | Ignored. Priority only applies to rulesets. |
| Priority update without recompile | Metadata in files becomes stale. SetRulesetPriority must recompile. |
| Index file deleted manually | Regenerated on next install/uninstall. No data loss (source is arm-index.json). |

## Dependencies

- **Sink Manager** - Manages index file and installation
- **Compiler** - Embeds priority in rule metadata
- **Manifest** - Tracks installed packages
- **Parser** - Reads ruleset YAML files

## Implementation Mapping

**Source files:**
- `internal/arm/sink/manager.go` - InstallRuleset, generateRulesetIndexRuleFile, loadIndex, saveIndex
- `internal/arm/compiler/generators.go` - GenerateRuleMetadata (embeds priority)
- `internal/arm/service/service.go` - SetRulesetPriority method
- `internal/arm/service/ruleset_setters.go` - Priority update logic

**Related specs:**
- `sink-compilation.md` - Compilation process and metadata embedding
- `package-installation.md` - Installation workflow
- `registry-management.md` - Package storage and retrieval

## Examples

### Example 1: Install with Custom Priority

**Input:**
```bash
arm install ruleset --priority 200 team-rules/standards cursor-rules
arm install ruleset --priority 100 community-rules/best-practices cursor-rules
```

**Expected Output:**

arm-index.json:
```json
{
  "version": 1,
  "rulesets": {
    "team-rules/standards@1.0.0": {
      "priority": 200,
      "files": ["arm/team-rules/standards/1.0.0/rule1.mdc"]
    },
    "community-rules/best-practices@2.1.0": {
      "priority": 100,
      "files": ["arm/community-rules/best-practices/2.1.0/rule2.mdc"]
    }
  }
}
```

arm_index.mdc:
```markdown
# ARM Rulesets

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.**

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules

## Installed Rulesets

### team-rules/standards@1.0.0
- **Priority:** 200
- **Rules:**
  - arm/team-rules/standards/1.0.0/rule1.mdc

### community-rules/best-practices@2.1.0
- **Priority:** 100
- **Rules:**
  - arm/community-rules/best-practices/2.1.0/rule2.mdc
```

**Verification:**
- Team rules listed first (higher priority)
- Community rules listed second (lower priority)
- Priority values match installation flags

### Example 2: Update Priority

**Input:**
```bash
arm set ruleset community-rules/best-practices priority 150
```

**Expected Output:**

arm-index.json updated:
```json
{
  "rulesets": {
    "community-rules/best-practices@2.1.0": {
      "priority": 150,
      "files": ["arm/community-rules/best-practices/2.1.0/rule2.mdc"]
    }
  }
}
```

arm_index.mdc updated with new sort order (if 150 > other priorities).

Rule metadata recompiled:
```yaml
---
priority: 150
---
```

**Verification:**
- Priority updated in index
- Index file regenerated with new sort order
- Rule files recompiled with new priority metadata

### Example 3: Default Priority

**Input:**
```bash
arm install ruleset my-rules/general cursor-rules
```

**Expected Output:**

arm-index.json:
```json
{
  "rulesets": {
    "my-rules/general@1.0.0": {
      "priority": 100,
      "files": ["arm/my-rules/general/1.0.0/rule.mdc"]
    }
  }
}
```

Rule metadata (no priority field):
```yaml
---
namespace: my-rules/general@1.0.0
ruleset:
  id: general
---
```

**Verification:**
- Default priority (100) stored in index
- Priority not embedded in metadata (only if > 0)

### Example 4: Empty Rulesets

**Input:**
```bash
arm uninstall my-rules/general
# (last ruleset removed)
```

**Expected Output:**

arm-index.json:
```json
{
  "version": 1,
  "rulesets": {}
}
```

arm_index.mdc deleted.

**Verification:**
- Index file removed when no rulesets remain
- No orphaned index file

### Example 5: Same Priority Conflict

**Input:**
```bash
arm install ruleset --priority 100 rules-a/set1 cursor-rules
arm install ruleset --priority 100 rules-b/set2 cursor-rules
```

**Expected Output:**

arm_index.mdc (order undefined):
```markdown
### rules-a/set1@1.0.0
- **Priority:** 100

### rules-b/set2@1.0.0
- **Priority:** 100
```

**Verification:**
- Both rulesets listed
- Order between same-priority rulesets is arbitrary
- Users should avoid this scenario

## Notes

### Why Priority Matters

Priority enables layering of rulesets with clear precedence:
- **Team standards** (priority 200) override **community best practices** (priority 100)
- **Project-specific rules** (priority 300) override **team standards** (priority 200)
- **Experimental rules** (priority 50) have lowest precedence

### AI Tool Consumption

The arm_index.* file serves as documentation for AI tools:
- Explains priority system in natural language
- Lists rulesets in precedence order
- Provides explicit conflict resolution guidance

AI tools should:
1. Read arm_index.* to understand priority order
2. Check embedded metadata in individual rules
3. Apply higher priority rules when conflicts occur

### Design Decisions

**Why integer priority?**
- Simple to understand and compare
- No upper/lower bounds (flexibility)
- Default of 100 allows room above and below

**Why regenerate index file on every change?**
- Ensures consistency between arm-index.json and arm_index.*
- Prevents stale documentation
- Low cost (small file, fast generation)

**Why embed priority in rule metadata?**
- Enables traceability (which priority was this rule installed with?)
- Allows AI tools to verify priority without reading index
- Useful for debugging priority conflicts

**Why allow same priority?**
- Simplifies implementation (no validation needed)
- Users may intentionally want equal priority
- Undefined order is acceptable (users should avoid if order matters)

**Why delete index file when empty?**
- Avoids confusing AI tools with empty priority list
- Cleaner sink directory
- Regenerated automatically when rulesets installed

### Testing Considerations

Tests should verify:
- Default priority (100) applied when flag omitted
- Custom priority stored correctly
- Index file sorted by priority (high to low)
- Priority embedded in metadata only if > 0
- Index file deleted when last ruleset removed
- SetRulesetPriority updates index and recompiles rules
- Same priority handled gracefully (no errors)
- Negative priorities sorted correctly
- Priority only applies to rulesets (not promptsets)
