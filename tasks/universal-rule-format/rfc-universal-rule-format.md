# RFC: Universal Rule Format for AI Development Tools

## Abstract

This RFC proposes a Universal Rule Format (URF) that enables tool-agnostic rule authoring for AI development assistants. The format addresses the current fragmentation of AI rule formats across different tools while maintaining compatibility and enabling advanced features like priority-based rule ordering and enforcement levels.

## Problem Statement

Current AI rule management suffers from several critical issues:

1. **Format Fragmentation**: Each AI tool (Cursor, Amazon Q, GitHub Copilot) uses different rule formats and conventions
2. **Manual Conversion**: Rule authors must manually adapt rules for each target tool
3. **No Priority System**: Rules cannot be ordered by importance or enforcement level
4. **Limited Metadata**: Existing formats lack structured metadata for rule management
5. **Poor Composition**: Rules from different sources cannot be intelligently combined

## Proposed Solution

### Universal Rule Format (URF)

A YAML-based format that captures rule intent, priority, and metadata in a tool-agnostic way:

```yaml
version: "1.0"
metadata:
  id: "ruleset-id"
  name: "Ruleset Name"
  version: "1.0.0"
  description: "Description"
rules:
  - id: "critical-security-check"
    name: "Critical Security Check"
    description: "Validate all user inputs to prevent security vulnerabilities"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
  - id: "recommended-best-practice"
    name: "Recommended Best Practice"
    description: "Follow established coding conventions for maintainability"
    priority: 80
    enforcement: "should"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
  - id: "optional-optimization"
    name: "Optional Optimization"
    description: "Consider performance improvements when feasible"
    priority: 60
    enforcement: "may"
    scope:
      - files: ["**/*.ext"]
    body: |
      Rule content in markdown format.
```

compiled meta-data for the critical security check rule
```
---
namespace: namespace-id
ruleset:
  id: ruleset-id
  name: Ruleset Name
  version: 1.0.0
  rules:
    - critical-security-check
    - recommended-best-practice
    - additional-optimization
rule:
  id: critical-security-check
  name: Rule Name
  enforcement: MUST
  priority: 100
  scope:
    - files: "**/*.py"
---
```


### Compilation Strategy

URF files are compiled to tool-specific files, formats, and directory structures by implementation tools.

### Priority System

Rules within a ruleset are ordered by:

1. **Enforcement level** (must > should > may)
2. **Rule priority** (higher number = higher priority)

Implementation tools may layer additional priority systems for managing multiple rulesets.

## Implementation Details

### Compilation Process

1. **Parse URF**: Validate YAML structure and required fields
2. **Generate Metadata**: Create cross-references and embed all required metadata
3. **Compile Format**: Transform to tool-specific format with priority metadata
4. **Write Files**: Individual files per rule with embedded priority information

### Required Metadata in Compiled Files

Each compiled rule file must include the following metadata:

#### Namespace Context
- **Namespace**: Deployment namespace identifier

#### Ruleset Context
- **Ruleset ID**: Unique ruleset identifier
- **Ruleset Name**: Human-readable ruleset name
- **Ruleset Version**: Semantic version of the source ruleset
- **Rules List**: All rule IDs in the ruleset

#### Core Rule Metadata
- **Rule ID**: Unique identifier within the ruleset
- **Rule Name**: Human-readable rule name
- **Enforcement Level**: MUST, SHOULD, or MAY
- **Priority**: Numeric priority within ruleset
- **Scope**: File patterns where rule applies (if specified)

### Tool-Specific Formats

#### Cursor Target Format
- **Format**: Markdown with YAML frontmatter (`.mdc` extension)
- **Frontmatter Fields**:
  - `description`: Rule description
  - `globs`: File patterns from scope (array format)
  - `alwaysApply`: true for MUST, false for SHOULD/MAY
- **Body Structure**:
  - Compact metadata block
  - Rule title with enforcement suffix (H1)
  - Rule body content

**Example**:
```markdown
---
description: Validate all user inputs to prevent security vulnerabilities
globs: ["**/*.ext"]
alwaysApply: true
---

---
namespace: namespace-id
ruleset:
  id: ruleset-id
  name: Ruleset Name
  version: 1.0.0
  rules:
    - critical-security-check
    - recommended-best-practice
    - optional-optimization
rule:
  id: critical-security-check
  name: Critical Security Check
  enforcement: MUST
  priority: 100
  scope:
    - files: "**/*.ext"
---

# Critical Security Check (MUST)

Rule content in markdown format.
```

#### Amazon Q Target Format
- **Format**: Pure markdown (`.md` extension)
- **Body Structure**:
  - Compact metadata block
  - Rule title with enforcement suffix (H1)
  - Rule body content

**Example**:
```markdown
---
namespace: namespace-id
ruleset:
  id: ruleset-id
  name: Ruleset Name
  version: 1.0.0
  rules:
    - critical-security-check
    - recommended-best-practice
    - optional-optimization
rule:
  id: critical-security-check
  name: Critical Security Check
  enforcement: MUST
  priority: 100
  scope:
    - files: "**/*.ext"
---

# Critical Security Check (MUST)

Rule content in markdown format.
```

### File Organization

Implementation tools should organize compiled rules to preserve ruleset identity and version information.

## Benefits

1. **Single Source of Truth**: Write rules once, deploy everywhere
2. **Intelligent Prioritization**: Rules are ordered by importance and enforcement
3. **Tool Compatibility**: Automatic adaptation to tool-specific formats
4. **Rich Metadata**: Structured information for rule management
5. **Extensibility**: Easy to add support for new AI tools

## Backward Compatibility

- **Auto-detection**: Implementation tools should detect URF vs legacy formats automatically
- **Mixed Support**: Projects should support both formats during transition
- **Migration Path**: Legacy rules can be converted to URF format

## Security Considerations

- **Input Validation**: URF files are validated during parsing
- **Scope Restrictions**: File patterns are sanitized to prevent path traversal
- **Content Safety**: Rule bodies are treated as untrusted content

## Future Extensions

- **Rule Dependencies**: Rules that depend on other rules
- **Conditional Rules**: Rules that apply based on project context
- **Rule Templates**: Parameterized rules for common patterns
- **Validation Rules**: Rules that validate code against other rules

## Reference Implementation

A reference implementation should include:

1. **URF Parser**: YAML validation and structure verification
2. **Compiler Framework**: Extensible system for adding new target formats
3. **Metadata Generator**: Logic for embedding priority and cross-reference information
4. **Validation Tools**: Format checking and semantic validation

## Conclusion

The Universal Rule Format provides a foundation for scalable, maintainable AI rule management. By separating rule intent from tool-specific formatting, URF enables rule authors to focus on content while implementation tools handle the complexity of multi-tool deployment.
