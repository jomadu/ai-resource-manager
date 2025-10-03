# Sample Ruleset Compilation

This directory contains a sample ruleset (`sample-ruleset.yml`) compiled to multiple build targets, demonstrating how the AI Rules Manager transforms a single ruleset definition into platform-specific formats.

## Compilation Process

The compilation process takes a source ruleset YAML file and generates platform-specific output files for each rule. Each build target has its own subdirectory with files that follow the target's specific format and requirements.

### Source Ruleset Structure

The source `sample-ruleset.yml` contains:
- **Ruleset metadata**: ID, name, description
- **Rules**: Each with optional properties like name, description, priority, enforcement level, scope, and body content
- **Enforcement levels**: `must`, `should`, `may` (or none)
- **Scope**: Optional file patterns that define which files the rule applies to

## Build Targets

### `/md/` - Markdown Format

**File Extension**: `.md`

**Purpose**: General-purpose markdown documentation format with comprehensive metadata.

**Structure**:
- **Frontmatter**: Full metadata block containing namespace, ruleset info, and rule details
- **Header**: Rule name with enforcement hint (e.g., "Rule 1 (MUST)")
- **Body**: Rule content

**Key Features**:
- Complete metadata preservation
- Enforcement level indicators in headers
- Rich frontmatter for documentation systems
- Consistent with other markdown-based tools

**Example**:
```yaml
---
namespace: someNamespace
ruleset:
    id: "cleanCode"
    name: "Clean Code"
    description: "Make clean code"
    rules: [ruleOne, ruleTwo, ...]
rule:
    id: ruleOne
    name: "Rule 1"
    description: "..."
    priority: 100
    enforcement: must
    scope:
        - files: ["**/*.py"]
---

# Rule 1 (MUST)

This is the body of the rule.
```

### `/cursor/` - Cursor IDE Format

**File Extension**: `.mdc`

**Purpose**: Optimized for Cursor IDE with minimal, functional frontmatter.

**Structure**:
- **Frontmatter**: Only essential Cursor-specific properties
- **Header**: Rule name with enforcement hint
- **Body**: Rule content

**Key Features**:
- **`description`**: Rule description (only if non-empty)
- **`globs`**: File patterns from scope as comma-separated string (only if non-empty)
- **`alwaysApply`**: Boolean derived from enforcement level (`true` if `must`, omitted if `false`)
- **Minimal frontmatter**: Empty properties are omitted entirely
- **No frontmatter**: If all properties would be empty, no frontmatter block is included

**Example**:
```yaml
---
description: "Rule with optional description, priority, enforcement..."
globs: "**/*.py"
alwaysApply: true
---

# Rule 1 (MUST)

This is the body of the rule.
```

### `/copilot/` - GitHub Copilot Format

**File Extension**: `.instructions.md`

**Purpose**: GitHub Copilot instructions with separate metadata and frontmatter blocks.

**Structure**:
- **Frontmatter**: Copilot-specific `applyTo` property
- **Metadata Block**: Separate YAML block with full rule metadata
- **Header**: Rule name with enforcement hint
- **Body**: Rule content

**Key Features**:
- **`applyTo`**: File patterns as comma-separated string (defaults to `"**/*"` if no scope)
- **Separate metadata**: Full metadata in its own YAML block below frontmatter
- **File pattern handling**: Converts scope arrays to comma-separated strings
- **Default scope**: Uses `"**/*"` when no scope is defined

**Example**:
```yaml
---
applyTo: "**/*.py"
---

---
namespace: someNamespace
ruleset:
    id: "cleanCode"
    name: "Clean Code"
    description: "Make clean code"
    rules: [ruleOne, ruleTwo, ...]
rule:
    id: ruleOne
    name: "Rule 1"
    description: "..."
    priority: 100
    enforcement: must
    scope:
        - files: ["**/*.py"]
---

# Rule 1 (MUST)

This is the body of the rule.
```

## Usage

Each build target is designed for specific use cases:

- **Markdown**: Documentation, wikis, general-purpose rule display
- **Cursor**: IDE integration with minimal overhead
- **Copilot**: AI assistant integration with clear file targeting

The compilation process ensures that each target receives the appropriate level of detail and formatting for its intended use case while maintaining consistency in the core rule content and structure.
