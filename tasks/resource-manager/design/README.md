![ARM Header](assets/header.png)

# AI Rules Manager (ARM)

## What is ARM?

A package manager for AI rules that treats rulesets and promptsets like code dependencies - with semantic versioning, reproducible installs, and automatic distribution to your AI tools.

Connect to Git repositories like awesome-cursorrules or your team's rule collections, install versioned resources across projects, and keep them automatically synced with their source of truth.

## Why ARM?

AI coding assistants like Cursor and Amazon Q rely on rules and prompts to guide their behavior, but managing these resources is broken:

- **Manual copying** severs the connection to the source of truth - once copied, resources are orphaned with no way to get updates
- **Breaking changes blindness** - when you pull latest resources, you have no idea if they'll break your AI's behavior
- **Doesn't scale** - managing resources across even 3 projects becomes unmanageable overhead
- **Format fragmentation** - each AI tool uses different formats, requiring manual conversion

ARM solves these problems with a **package manager approach** - semantic versioning, reproducible installs, and automatic distribution keep resources connected to their source of truth. ARM resource definitions provide a unified format that compiles to platform-specific outputs.

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash
```

### Install Specific Version

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/install.sh | bash -s v1.0.0
```

### Manual Installation

1. Download the latest release from [GitHub](https://github.com/jomadu/ai-rules-manager/releases)
2. Extract and move the binary to your PATH
3. Run `arm help` to verify installation

### Verify Installation

```bash
arm version
arm help
```

> **Upgrading from v2?** See the [migration guide](docs/migration-v2-to-v3.md) for breaking changes and upgrade steps.

## Uninstall

```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-rules-manager/main/scripts/uninstall.sh | bash
```

## Quick Start

Add Git registry:
```bash
arm add registry --type git ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry
```

Add GitLab registry:
```bash
arm add registry --type gitlab --gitlab-project-id 123 my-gitlab https://gitlab.example.com
```

Add Cloudsmith registry:
```bash
arm add registry --type cloudsmith my-cloudsmith https://app.cloudsmith.com/myorg/ai-rules
```

Configure sinks:

**Cursor:**
```bash
arm add sink --type cursor cursor-rules .cursor/rules
arm add sink --type cursor cursor-prompts .cursor/prompts
```

**GitHub Copilot:**
```bash
arm add sink --type copilot copilot-rules .github/copilot
```

**Amazon Q:**
```bash
arm add sink --type amazonq q-rules .amazonq/rules
arm add sink --type amazonq q-prompts .amazonq/prompts
```

Install ruleset:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install promptset:
```bash
arm install promptset ai-rules/code-review-promptset cursor-prompts
```

Install to multiple sinks:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules
arm install promptset ai-rules/code-review-promptset cursor-prompts q-prompts
```

## Concepts

### Files

ARM uses four key files to manage your AI resources:

- **`arm.json`** - Team-shared project manifest with registries, rulesets, promptsets, and sinks
- **`arm-lock.json`** - Team-shared locked versions for reproducible installs
- **`arm-index.json`** - Local sink inventory tracking what ARM has installed in that sink, used to generate the `arm_index.*` file
- **`arm_index.*`** - Generated priority rule that helps AI tools resolve conflicts between rulesets

### Resources

ARM manages two types of AI resources:

- **Rulesets**: Collections of rules that provide instructions, guidelines, and context to AI coding assistants
- **Promptsets**: Collections of prompts that provide reusable prompt templates for AI interactions

Different AI tools use different formats:

- **Cursor**: `.cursorrules` files, `.mdc` files with YAML frontmatter, or `.md` files for prompts
- **Amazon Q**: `.md` files in `.amazonq/rules/` or `.amazonq/prompts/` directories
- **GitHub Copilot**: `.instructions.md` files in `.github/copilot/` directory

ARM consumes these resources as versioned packages, ensuring consistency across projects while respecting each tool's format requirements.

### Rulesets

Rulesets are collections of AI rules packaged as versioned units, identified by names like `ai-rules/clean-code-ruleset` where `ai-rules` is the registry and `clean-code-ruleset` is the ruleset name.

**Key Commands:**
- `arm install ruleset <registry>/<ruleset>[@version] <sink>...` - Install ruleset to specific sinks
- `arm update ruleset [ruleset...]` - Update to latest compatible versions
- `arm uninstall ruleset <registry>/<ruleset>` - Remove ruleset
- `arm list ruleset` - Show installed rulesets
- `arm info ruleset [ruleset...]` - Show detailed information
- `arm config ruleset set <registry>/<ruleset> <key> <value>` - Update ruleset configuration

#### ARM Resource Format

ARM uses Kubernetes-style resource definitions for rulesets:

```yaml
apiVersion: v1                    # Required
kind: Ruleset                     # Required
metadata:                         # Required
  id: "cleanCode"                 # Required
  name: "Clean Code"              # Optional
  description: "Make clean code"  # Optional
spec:                             # Required
  rules:                          # Required
    ruleOne:                      # Required (rule key)
      name: "Rule 1"              # Optional
      description: "Rule description"  # Optional
      priority: 100               # Optional (default: 100)
      enforcement: must           # Optional (must|should|may)
      scope:                      # Optional
        - files: ["**/*.py"]      # Optional
      body: |                     # Required
        This is the body of the rule.
```

ARM resource definitions compile to tool-specific formats with embedded metadata for priority resolution and conflict management.

### Promptsets

Promptsets are collections of AI prompts packaged as versioned units, identified by names like `ai-rules/code-review-promptset` where `ai-rules` is the registry and `code-review-promptset` is the promptset name.

**Key Commands:**
- `arm install promptset <registry>/<promptset>[@version] <sink>...` - Install promptset to specific sinks
- `arm update promptset [promptset...]` - Update to latest compatible versions
- `arm uninstall promptset <registry>/<promptset>` - Remove promptset
- `arm list promptset` - Show installed promptsets
- `arm info promptset [promptset...]` - Show detailed information
- `arm config promptset set <registry>/<promptset> <key> <value>` - Update promptset configuration

#### ARM Promptset Format

ARM uses Kubernetes-style resource definitions for promptsets:

```yaml
apiVersion: v1                    # Required
kind: Promptset                   # Required
metadata:                         # Required
  id: "code-review"               # Required
  name: "Code Review"             # Optional
  description: "Review the code!" # Optional
spec:                             # Required
  prompts:                        # Required
    promptOne:                    # Required (prompt key)
      name: "Prompt One"          # Optional
      description: "Prompt description"  # Optional
      body: |                     # Required
        This is how you review the code.
```

### Registries

Registries are remote sources where rulesets and promptsets are stored and versioned, similar to npm registries. ARM supports:

- **Git registries**: GitHub repositories, GitLab projects, or any Git remote
- **GitLab Package registries**: GitLab's Generic Package Registry for versioned packages
- **Cloudsmith registries**: Cloudsmith's package repository service for single-file artifacts

#### Registry Structure

**Recommended structure:**
```
clean-code-ruleset.yml      # ARM resource definitions
security-ruleset.yml
code-review-promptset.yml   # ARM resource definitions
build/                      # Pre-compiled rules (optional)
├── cursor/
│   ├── clean-code.mdc
│   └── security.mdc
└── amazonq/
    ├── clean-code.md
    └── security.md
```

This structure works for all registry types, with ARM resource files at the root level and pre-compiled rules organized under `build/` by AI tool. ARM defaults to installing resource files (`*.yml, *.yaml`) when no `--include` patterns are specified.

#### Archive Support

ARM automatically extracts and processes **zip** and **tar.gz** archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic extraction**: Archives are detected by extension and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence in case of path conflicts
- **Security**: Path sanitization prevents directory traversal attacks
- **Pattern filtering**: `--include` patterns are applied to the merged content after extraction

**Examples:**
- [PatrickJS/awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules) - Community collection of Cursor rules
- [snarktank/ai-dev-tasks](https://github.com/snarktank/ai-dev-tasks) - AI development task templates
- [steipete/agent-rules](https://github.com/steipete/agent-rules) - Agent configuration rules

**Key Commands:**
- `arm add registry --type git <name> <url>` - Add Git registry
- `arm add registry --type gitlab --gitlab-project-id <id> <name> <url>` - Add GitLab registry
- `arm add registry --type cloudsmith <name> <url>` - Add Cloudsmith registry
- `arm list registry` - List registries
- `arm remove registry <name>` - Remove registry

### Sinks

Sinks define where installed resources should be placed in your local filesystem. Each sink targets a specific directory for a particular AI tool.

**Key Commands:**
- `arm add sink --type <cursor|amazonq|copilot> <name> <directory>` - Add sink
- `arm list sink` - List sinks
- `arm remove sink <name>` - Remove sink

#### Layout Modes

**Hierarchical Layout** (default): Preserves directory structure from resources
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

**Flat Layout**: Places all files in single directory with hash-prefixed names
```
.github/copilot/
├── 183791a9_cleanCode_ruleOne.instructions.md
├── arm_index.instructions.md
└── arm-index.json
```

#### Compilation

ARM resource definitions are automatically compiled to tool-specific formats:

- **Cursor**: Markdown with YAML frontmatter (`.mdc`) for rules, plain markdown (`.md`) for prompts
- **Amazon Q**: Pure markdown (`.md`) for both rules and prompts
- **Copilot**: Instructions format (`.instructions.md`) for both rules and prompts

Each compiled resource includes embedded metadata for priority resolution and resource tracking.

## Package Management

Install all configured packages:
```bash
arm install
```

Update packages:
```bash
arm update
arm update ruleset
arm update promptset
```

Check for outdated packages:
```bash
arm outdated
arm outdated ruleset
arm outdated promptset
```

List installed packages:
```bash
arm list
arm list ruleset
arm list promptset
```

## Utilities

Clean cache:
```bash
arm clean cache
arm clean cache --nuke
```

Clean sinks:
```bash
arm clean sinks
arm clean sinks --nuke
```

Compile resources:
```bash
arm compile --target cursor ruleset.yml ./output/
arm compile --target amazonq --recursive ./src/ ./build/
```
