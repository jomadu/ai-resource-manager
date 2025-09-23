# Universal Rule Format (URF)

## Overview

Universal Rule Format (URF) is ARM's solution to AI rule format fragmentation. Instead of writing separate rules for each AI tool, URF lets you write rules once in a standardized YAML format, then automatically compile them for any supported AI tool.

## The Problem

Each AI tool uses different rule formats:
- **Cursor**: `.mdc` files with YAML frontmatter
- **Amazon Q**: `.md` files in specific directories
- **GitHub Copilot**: `.instructions.md` files with special formatting

This forces teams to maintain multiple versions of the same rules, leading to inconsistency and maintenance overhead.

## The Solution

URF provides a single source of truth for AI rules with automatic compilation to tool-specific formats.

## URF File Structure

```yaml
version: "1.0"
metadata:
  id: "grug-brained-dev"
  name: "Grug-Brained Developer Rules"
  version: "1.0.0"  # required
  description: "Simple, obvious coding practices"  # optional
rules:
  simple-code:  # Rule ID as map key
    name: "Keep Code Simple"
    description: "Write code that grug brain can understand"  # optional
    priority: 100  # optional
    enforcement: "must"  # optional
    scope:  # optional
      - files: ["**/*.py", "**/*.js"]
    body: |
      Write simple, obvious code that grug brain can understand.
      Avoid clever tricks and complex abstractions.

      ## Examples

      **Good:**
      ```python
      def calculate_total(items):
          total = 0
          for item in items:
              total += item.price
          return total
      ```

      **Bad:**
      ```python
      calc_tot = lambda i: reduce(add, map(attrgetter('price'), i), 0)
      ```
  meaningful-names:  # Another rule ID as map key
    name: "Use Meaningful Names"
    # Minimal rule - only required fields
    body: |
      Choose names that clearly express what the code does.
      Avoid abbreviations and single-letter variables.
```

## Field Reference

### Metadata Section
- **`id`** (required): Unique identifier for the ruleset
- **`name`** (required): Human-readable ruleset name
- **`version`** (required): Semantic version (e.g., "1.0.0")
- **`description`** (optional): Brief description of the ruleset's purpose

### Rules Section
Rules are defined as a map where each key is the rule ID and the value contains the rule definition. This structure automatically prevents duplicate rule IDs.

### Rule Fields
- **Rule ID** (map key): Unique identifier within the ruleset
- **`name`** (required): Human-readable rule name
- **`description`** (optional): Brief explanation of the rule's purpose
- **`priority`** (optional): Numeric priority (higher = more important)
- **`enforcement`** (optional): Rule strictness level
  - `"must"`: Critical rules that should always be followed
  - `"should"`: Important best practices
  - `"may"`: Optional suggestions
- **`scope`** (optional): Array of scope objects defining where the rule applies
  - Each scope object currently supports: `files: ["pattern1", "pattern2"]`
- **`body`** (required): Rule content in markdown format

## Compilation Process

When ARM installs a URF ruleset, it:

1. **Parses** the YAML file and validates structure
2. **Generates metadata** for each rule including priority and enforcement
3. **Compiles** to tool-specific formats
4. **Embeds metadata** in each compiled rule for priority resolution

## Tool-Specific Output

ARM automatically compiles URF files to the appropriate format for each AI tool:

- **Cursor**: `.mdc` files with YAML frontmatter and embedded metadata
- **Amazon Q**: `.md` files with embedded metadata blocks
- **GitHub Copilot**: `.instructions.md` files with embedded metadata

All compiled formats include the original rule content plus ARM metadata for priority resolution and conflict management.

### Compiler Defaults

When optional fields are omitted, compilers apply tool-specific defaults:

**Scope handling:**
- **Cursor**: Defaults to `globs: ["**/*"]` when scope is empty
- **Copilot**: Defaults to `applyTo: "**"` when scope is empty
- **Amazon Q**: Scope only appears in metadata (no tool-specific frontmatter)

**Enforcement handling:**
- **Cursor**: Only sets `alwaysApply: true` when enforcement is "must"
- **Copilot/Amazon Q**: Enforcement only appears in metadata and rule titles
- **All tools**: Omit enforcement from rule titles when not specified

## Priority Resolution

When multiple rulesets are installed, ARM generates an `arm_index.*` file that helps AI tools resolve conflicts:

1. **Enforcement level** takes precedence (`must` > `should` > `may`)
2. **Rule priority** breaks ties within the same enforcement level
3. **Ruleset priority** (set during installation) resolves conflicts between rulesets

## Best Practices

### Writing URF Rules

1. **Use clear, descriptive rule IDs** - `simple-code` not `rule1` (rule IDs are map keys)
2. **Set appropriate priorities** - Reserve 90+ for critical rules
3. **Choose enforcement levels carefully** - Use `must` sparingly for truly critical rules
4. **Include examples** - Show both good and bad code patterns
5. **Scope rules appropriately** - Target specific file types when relevant
6. **Avoid duplicate rule IDs** - Map structure automatically prevents duplicates

### Organizing Rulesets

1. **Group related rules** - Keep similar concepts in the same ruleset
2. **Use semantic versioning** - Increment versions when making breaking changes
3. **Document rule interactions** - Explain how rules work together
4. **Test across tools** - Verify compiled output works as expected

## Migration from Legacy Formats

ARM automatically detects URF vs legacy formats, allowing gradual migration:

1. **Start with new rules** - Write new rules in URF format
2. **Convert high-value rules** - Migrate frequently-used rules first
3. **Maintain compatibility** - Keep legacy rules until migration is complete
4. **Validate output** - Test compiled rules in each target tool

## Examples

See the [sample registry](https://github.com/jomadu/ai-rules-manager-sample-git-registry) for a complete URF example:

- **grug-brained-dev.yml** - Simple, practical coding rules
