# Concepts

## Files

ARM uses four key files to manage your AI resources:

- **`arm.json`** - Team-shared project manifest with registries, rulesets, promptsets, and sinks
- **`arm-lock.json`** - Team-shared locked versions for reproducible installs
- **`arm-index.json`** - Local sink inventory tracking what ARM has installed in that sink, used to generate the `arm_index.*` file
- **`arm_index.*`** - Generated priority rule that helps AI tools resolve conflicts between rulesets

## Resources

ARM consumes two types of AI resources:

- **Rulesets**: Collections of rules that provide instructions, guidelines, and context to AI coding assistants
- **Promptsets**: Collections of prompts that provide reusable prompt templates for AI interactions

Different AI tools use different formats:

- **Cursor**: `.cursorrules` files, `.mdc` files with YAML frontmatter, or `.md` files for prompts
- **Amazon Q**: `.md` files in `.amazonq/rules/` or `.amazonq/prompts/` directories
- **GitHub Copilot**: `.instructions.md` files in `.github/copilot/` directory

ARM consumes these resources as versioned packages, ensuring consistency across projects while respecting each tool's format requirements.

## Resource Types

### Rulesets vs Promptsets

**Rulesets** contain rules with:
- Priority-based conflict resolution
- Enforcement levels (must/should/may)
- File scope targeting
- Compilation with metadata

**Promptsets** contain prompts with:
- Simple content-only compilation
- No priority conflicts
- Universal compatibility
- Minimal metadata

Both use ARM's Kubernetes-style resource definitions but serve different purposes in AI assistant workflows.
