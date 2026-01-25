# Sinks

## What are Sinks?

Sinks are output destinations where ARM compiles and writes installed resources. Think of them as build targets - the same source content can be compiled to different formats for different AI tools.

- The `--directory` parameter defines WHERE files go, usually where an AI tool expects the resources to already be (e.g. `.cursor/rules`, etc.)
- The `--tool` parameter specifies HOW any ARM resource files are compiled (e.g. to amazonq, cursor, copilot, etc.)

## Commands

For detailed command usage and examples, see [Sink Management](commands.md#sink-management) in the commands reference.

## Layout Modes

### Example Packages

To illustrate how layouts work, we'll use two example packages:

**Package 1: `ai-rules/clean-code-ruleset@1.0.0`**
- Contains: `rules/sample-ruleset.yml` (ARM resource file)
- Resource ID: `cleanCode`
- Rules: `ruleOne`, `ruleTwo`, `ruleThree`, `ruleFour`, `ruleFive`

**Package 2: `ai-rules/grug-brained-dev@1.0.0`**
- Contains: Non-resource files (plain markdown rules)
- Files: `rules/avoid-abstractions.md`, `rules/complexity-enemy.md`, `rules/readable-over-clever.md`

### Hierarchical Layout (default)

Preserves directory structure from packages.

**Path Schema:**
```
{sink_directory}/arm/{registry}/{package}/{version}/{original_path}
```

**Example:**
```
.cursor/rules/
└── arm/
    ├── ai-rules/
    │   ├── clean-code-ruleset/
    │   │   └── 1.0.0/
    │   │       └── rules/
    │   │           ├── cleanCode_ruleOne.mdc
    │   │           ├── cleanCode_ruleTwo.mdc
    │   │           ├── cleanCode_ruleThree.mdc
    │   │           ├── cleanCode_ruleFour.mdc
    │   │           └── cleanCode_ruleFive.mdc
    │   └── grug-brained-dev/
    │       └── 1.0.0/
    │           └── rules/
    │               ├── avoid-abstractions.md
    │               ├── complexity-enemy.md
    │               └── readable-over-clever.md
    ├── arm_index.mdc
    └── arm-index.json
```

### Flat Layout

Places all files in single directory with hash-prefixed names including relative path.

**Path Schema:**
```
{sink_directory}/arm_{package_hash}_{path_hash}_{flattened_path}
```

**Example:**

```
.github/copilot/
├── arm_1a2b_3c4d_rules_cleanCode_ruleOne.instructions.md
├── arm_1a2b_5e6f_rules_cleanCode_ruleTwo.instructions.md
├── arm_1a2b_7g8h_rules_cleanCode_ruleThree.instructions.md
├── arm_1a2b_9i0j_rules_cleanCode_ruleFour.instructions.md
├── arm_1a2b_k1l2_rules_cleanCode_ruleFive.instructions.md
├── arm_m3n4_o5p6_rules_avoid-abstractions.md
├── arm_m3n4_q7r8_rules_complexity-enemy.md
├── arm_m3n4_s9t0_rules_readable-over-clever.md
├── arm_index.instructions.md
└── arm-index.json
```

Files from the same package share the same package hash (e.g., `1a2b` for clean-code-ruleset, `m3n4` for grug-brained-dev), making them group together when sorted.

**Filename Truncation:**
Filenames are limited to 100 characters total. When truncation is needed, ARM uses a progressive fallback approach:

1. **Try full path**: `arm_1a2b_3c4d_rules_cleanCode_ruleOne.instructions.md`
2. **Try filename only**: `arm_1a2b_3c4d_cleanCode_ruleOne.instructions.md` (drops directory)
3. **Try truncated filename**: `arm_1a2b_3c4d_cleanCode_ruleO.instructions.md` (truncates name, keeps extension)

The package and path hashes are always preserved to prevent collisions.

## Compilation

ARM resource definitions are automatically compiled to tool-specific formats:

- **Cursor**: Markdown with YAML frontmatter (`.mdc`) for rules, plain markdown (`.md`) for prompts
- **Amazon Q**: Pure markdown (`.md`) for both rules and prompts
- **Copilot**: Instructions format (`.instructions.md`) for both rules (copilot doesn't have a "prompts" resource)

Each compiled resource includes embedded metadata for priority resolution and resource tracking.
