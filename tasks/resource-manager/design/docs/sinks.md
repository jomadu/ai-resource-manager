# Sinks

Sinks define where installed resources should be placed in your local filesystem. Each sink targets a specific directory for a particular AI tool.

## Commands

- `arm add sink --type <cursor|amazonq|copilot> <name> <directory>` - Add sink
- `arm list sink` - List sinks
- `arm info sink [name...]` - Show sink information
- `arm config sink set <name> <key> <value>` - Update sink configuration
- `arm remove sink <name>` - Remove sink

## Layout Modes

### Hierarchical Layout (default)
Preserves directory structure from resources:

```
.cursor/rules/
└── arm/
    ├── ai-rules/
    │   └── clean-code-ruleset/
    │       └── 1.0.0/
    │           └── cleanCode_ruleOne.mdc
    ├── arm_index.mdc
    └── arm-index.json
```

### Flat Layout
Places all files in single directory with hash-prefixed names:

```
.github/copilot/
├── 183791a9_cleanCode_ruleOne.instructions.md
├── arm_index.instructions.md
└── arm-index.json
```

## Compilation

ARM resource definitions are automatically compiled to tool-specific formats:

- **Cursor**: Markdown with YAML frontmatter (`.mdc`) for rules, plain markdown (`.md`) for prompts
- **Amazon Q**: Pure markdown (`.md`) for both rules and prompts
- **Copilot**: Instructions format (`.instructions.md`) for both rules and prompts

Each compiled resource includes embedded metadata for priority resolution and resource tracking.

## Configuration

Available sink configuration keys:
- `layout` - Set to `hierarchical` or `flat`
- `directory` - Output path for the sink
- `compile_target` - Set to `md`, `cursor`, `amazonq`, or `copilot`

## Examples

Add Cursor sinks:
```bash
arm add sink --type cursor cursor-rules .cursor/rules
arm add sink --type cursor cursor-prompts .cursor/prompts
```

Add Amazon Q sinks:
```bash
arm add sink --type amazonq q-rules .amazonq/rules
arm add sink --type amazonq q-prompts .amazonq/prompts
```

Add GitHub Copilot sink:
```bash
arm add sink --type copilot copilot-rules .github/copilot
```

Change sink layout:
```bash
arm config sink set cursor-rules layout flat
```
