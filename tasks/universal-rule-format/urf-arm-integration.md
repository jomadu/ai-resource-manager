# ARM Integration with Universal Rule Format

This document describes how ARM (AI Rules Manager) consumes Universal Rule Format (URF) rulesets and manages installation priorities across multiple rulesets.

## Installation Priority System

ARM assigns installation priorities to rulesets during the `arm install` command:

```bash
arm install ai-rules/secure-coding-practices --sinks cursor --priority 100
arm install ai-rules/python-rules --sinks cursor --priority 50
```

### Priority Resolution

When multiple rulesets are installed, ARM resolves conflicts using:

1. **Enforcement level** (must > should > may)
2. **Ruleset installation priority** (higher number = higher priority)
3. **Rule priority within ruleset** (higher number = higher priority)

### Compilation Process

ARM compiles URF rulesets into tool-specific formats while preserving priority information:

1. **Parse URF YAML** - Extract rules, metadata, and priorities
2. **Apply installation context** - Add ruleset installation priority
3. **Generate cross-references** - Create "Other Rules in Ruleset" sections
4. **Compile to target format** - Output tool-specific files (Cursor .mdc, Amazon Q .md)
5. **Generate priority manifest** - Create `arm-priorities.md` for AI agents

### Directory Structure

```
.cursor/rules/
├── arm-priorities.md                    # Installation priority manifest
└── arm/
    ├── ai-rules/
    │   ├── secure-coding-practices/
    │   │   └── 1.5.0/
    │   │       └── secure-coding-practices/
    │   │           ├── input-validation.mdc
    │   │           └── sql-injection-prevention.mdc
    │   └── python-rules/
    │       └── 2.1.0/
    │           └── python-best-practices/
    │               ├── type-hints-required.mdc
    │               └── pep8-naming.mdc
```

### Priority Manifest

ARM generates `arm-priorities.md` at each sink to inform AI agents about:
- Installed rulesets and their priorities
- Relative paths to all rule files
- Priority hierarchy for conflict resolution

This ensures AI agents understand the complete context when applying rules across multiple rulesets.

## Tool-Specific Compilation

### Cursor (.mdc files)
- YAML frontmatter with `globs` and `alwaysApply`
- Metadata in human-readable format
- Relative file references for cross-linking

### Amazon Q (.md files)
- Pure markdown format
- All metadata in body text
- Enforcement level in headers

Both formats include:
- Ruleset identification
- Enforcement level and priority
- File scope patterns
- Cross-references to other rules in the same ruleset
