# Universal Rule Format Specification

## Format Specification

### YAML Structure

```yaml
version: "1.0"

metadata:
  name: "ruleset-name"           # Unique identifier for the ruleset
  version: "1.0.0"               # Semantic version of the ruleset
  description: "Description"     # Human-readable description

rules:
  - id: "rule-id"                # Unique identifier within ruleset
    description: "Rule purpose"   # Brief description of the rule
    priority: 100                 # Numeric priority (higher = more important)
    enforcement: "must"           # Enforcement level: must|should|may
    scope:                        # Optional: file patterns where rule applies
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
      This is the actual rule text that will be compiled
      into tool-specific formats.
```

### Field Definitions

- **version**: Format version (currently "1.0")
- **metadata.name**: Unique ruleset identifier
- **metadata.version**: Semantic version for the ruleset
- **metadata.description**: Human-readable ruleset description
- **rules[].id**: Unique rule identifier within the ruleset
- **rules[].description**: Brief rule description
- **rules[].priority**: Numeric priority (1-1000+, higher = more important)
- **rules[].enforcement**: Enforcement level
  - `must`: Critical rules that should be strictly followed
  - `should`: Important rules that are strongly recommended
  - `may`: Optional rules that provide guidance
- **rules[].scope**: Optional file pattern restrictions
- **rules[].body**: Rule content in markdown format

## Design Principles

### Priority System
- **Rule priority**: Assigns priority to rules within the ruleset
- **Installation priority**: Assigns priority to the ruleset against other rulesets

### Enforcement Levels
- **must**: Critical rules that should be strictly followed
- **should**: Important rules that are strongly recommended
- **may**: Optional rules that provide guidance

### Rule Composition
- Rulesets are self-contained and don't depend on other rulesets
- Rules within a ruleset reference each other through relative file paths in compiled output
- Rule relationships are expressed in the body text and compiled markdown

### Conflict Resolution Hierarchy
1. **Enforcement level** (must > should > may)
2. **Ruleset installation priority** (higher priority wins)
3. **Individual rule priority** (within each ruleset)

## Format Examples

### Python Development Ruleset

```yaml
version: "1.0"

metadata:
  name: "python-best-practices"
  version: "2.1.0"
  description: "Python development best practices and conventions"

rules:
  - id: "type-hints-required"
    description: "Use type hints for function parameters and return values"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.py"]
    body: |
      All function parameters and return values must include type hints.
      Use typing module for complex types. This improves code readability
      and enables better IDE support and static analysis.

  - id: "pep8-naming"
    description: "Follow PEP 8 naming conventions"
    priority: 90
    enforcement: "must"
    scope:
      - files: ["**/*.py"]
    body: |
      Use snake_case for variables and functions, PascalCase for classes,
      UPPER_CASE for constants. Follow PEP 8 naming guidelines strictly.

  - id: "docstrings-public"
    description: "Document all public functions and classes"
    priority: 80
    enforcement: "should"
    scope:
      - files: ["**/*.py"]
    body: |
      All public functions and classes should have docstrings following
      Google or NumPy style. Include parameter descriptions and return values.

  - id: "list-comprehensions"
    description: "Prefer list comprehensions over loops when appropriate"
    priority: 70
    enforcement: "should"
    scope:
      - files: ["**/*.py"]
    body: |
      Use list comprehensions for simple transformations and filtering.
      They are more Pythonic and often more readable than equivalent loops.

  - id: "f-strings"
    description: "Use f-strings for string formatting"
    priority: 60
    enforcement: "may"
    scope:
      - files: ["**/*.py"]
    body: |
      Prefer f-strings over .format() or % formatting for better readability
      and performance in Python 3.6+.

  - id: "context-managers"
    description: "Use context managers for resource management"
    priority: 50
    enforcement: "may"
    scope:
      - files: ["**/*.py"]
    body: |
      Use 'with' statements for file operations, database connections,
      and other resources that need proper cleanup.
```

### Secure Coding Ruleset

```yaml
version: "1.0"

metadata:
  name: "secure-coding-practices"
  version: "1.5.0"
  description: "Security-focused coding guidelines and vulnerability prevention"

rules:
  - id: "input-validation"
    description: "Validate and sanitize all user inputs"
    priority: 1000
    enforcement: "must"
    body: |
      All user inputs must be validated, sanitized, and escaped before use.
      Never trust user input. Use allowlists over denylists when possible.
      Implement proper input length limits and type checking.

  - id: "sql-injection-prevention"
    description: "Use parameterized queries to prevent SQL injection"
    priority: 950
    enforcement: "must"
    scope:
      - files: ["**/*.py", "**/*.js", "**/*.java", "**/*.cs"]
    body: |
      Never concatenate user input directly into SQL queries.
      Always use parameterized queries or prepared statements.
      Use ORM query builders when available.

  - id: "secrets-management"
    description: "Never hardcode secrets in source code"
    priority: 900
    enforcement: "should"
    body: |
      Store secrets in environment variables, secure vaults, or configuration
      files that are excluded from version control. Use secret management
      services in production environments.

  - id: "https-enforcement"
    description: "Enforce HTTPS for all external communications"
    priority: 850
    enforcement: "should"
    body: |
      All external API calls and web requests should use HTTPS.
      Validate SSL certificates and implement certificate pinning
      for critical connections.

  - id: "error-information-disclosure"
    description: "Avoid exposing sensitive information in error messages"
    priority: 800
    enforcement: "may"
    body: |
      Error messages should not reveal system internals, file paths,
      database schemas, or other sensitive information that could
      aid attackers.

  - id: "dependency-scanning"
    description: "Regularly scan dependencies for known vulnerabilities"
    priority: 750
    enforcement: "may"
    body: |
      Use automated tools to scan third-party dependencies for known
      security vulnerabilities. Keep dependencies updated and remove
      unused packages.
```

## Installation & Compilation

```bash
arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
arm config sink add cursor .cursor/rules --compile-to cursor
arm install ai-rules/secure-coding-practices --sinks cursor --priority 100
arm install ai-rules/python-rules --sinks cursor --priority 50
```

## Priority Resolution Example

With the above installation priorities, rules are ordered as follows:

```
MUST (Critical - Always Applied First)
├── secure-coding-practices (installation priority: 100)
│   ├── input-validation (rule priority: 1000)
│   └── sql-injection-prevention (rule priority: 950)
└── python-rules (installation priority: 50)
    ├── type-hints-required (rule priority: 100)
    └── pep8-naming (rule priority: 90)

SHOULD (Important - Applied Second)
├── secure-coding-practices (installation priority: 100)
│   ├── secrets-management (rule priority: 900)
│   └── https-enforcement (rule priority: 850)
└── python-rules (installation priority: 50)
    ├── docstrings-public (rule priority: 80)
    └── list-comprehensions (rule priority: 70)

MAY (Optional - Applied Last)
├── secure-coding-practices (installation priority: 100)
│   ├── error-information-disclosure (rule priority: 800)
│   └── dependency-scanning (rule priority: 750)
└── python-rules (installation priority: 50)
    ├── f-strings (rule priority: 60)
    └── context-managers (rule priority: 50)
```

This hierarchy ensures security rules always take precedence over development style rules, while maintaining logical ordering within each category.

## Compiled Output Examples

ARM compiles universal format rules into tool-specific formats. Here's how the `type-hints-required` rule would be compiled for different tools:

### Cursor Compilation

**File**: `.cursor/rules/arm/ai-rules/python-rules/2.1.0/python-best-practices/type-hints-required.mdc`

```markdown
---
description: Use type hints for function parameters and return values
globs: ["**/*.py"]
alwaysApply: true
---

**Ruleset:** ai-rules/python-rules/2.1.0/python-best-practices
**Enforcement:** MUST
**Priority Within Ruleset:** 100
**Scope:**
- files: ["**/*.py"]
**Other Rules in Ruleset:**
1. ./pep8-naming.mdc (enforcement MUST, Priority 90)
2. ./docstrings-public.mdc (enforcement SHOULD, Priority 80)
3. ./list-comprehensions.mdc (enforcement SHOULD, Priority 70)
4. ./f-strings.mdc (enforcement MAY, Priority 60)
5. ./context-managers.mdc (enforcement MAY, Priority 50)

# Type Hints Required

All function parameters and return values must include type hints.
Use typing module for complex types. This improves code readability
and enables better IDE support and static analysis.
```

### Amazon Q Compilation

**File**: `.amazonq/rules/arm/ai-rules/python-rules/2.1.0/python-best-practices/type-hints-required.md`

```markdown
**Ruleset:** ai-rules/python-rules/2.1.0/python-best-practices
**Enforcement:** MUST
**Priority Within Ruleset:** 100
**Scope:**
- files: ["**/*.py"]
**Other Rules in Ruleset:**
1. ./pep8-naming.md (enforcement MUST, Priority 90)
2. ./docstrings-public.md (enforcement SHOULD, Priority 80)
3. ./list-comprehensions.md (enforcement SHOULD, Priority 70)
4. ./f-strings.md (enforcement MAY, Priority 60)
5. ./context-managers.md (enforcement MAY, Priority 50)

# Type Hints Required (MUST)

All function parameters and return values must include type hints.
Use typing module for complex types. This improves code readability
and enables better IDE support and static analysis.
```

### Key Differences

**Cursor Format (MDC)**:
- Uses `globs` frontmatter field for file patterns
- `alwaysApply: true` for MUST rules, `false` for SHOULD/MAY
- `description` includes enforcement level and priority
- Uses relative paths for rule references (`.mdc` extension)
- Simpler, more concise presentation

**Amazon Q Format**:
- No frontmatter - pure markdown format
- Metadata included in human-readable form in the body
- Explicit priority context and enforcement explanations
- Better suited for complex rule hierarchies

## Tool Compatibility

- **Implementation agnostic**: Format doesn't specify tool-specific behavior
- **Best effort encoding**: Unsupported features are encoded in rule body text
- **Frontmatter support**: Tools that support YAML frontmatter get structured metadata
- **Fallback text**: All metadata is also included in human-readable form in the body

## Supported Tools

- **Cursor**: Compiles to markdown files in `.cursor/rules`
- **Amazon Q**: Compiles to markdown files in `.amazonq/rules`
- **Extensible**: New compilers can be added for other AI tools
