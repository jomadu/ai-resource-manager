# Amazon Q Prompts

Saved prompts for Ralph workflow automation.

## Available Prompts

### @prd-generator
Generate a Product Requirements Document for a new feature.

**Usage:**
```
@prd-generator I want to add version constraint validation
```

**Output:** `specs/[feature-name]/spec.md`

### @prd-to-ralph
Convert a markdown PRD to Ralph's JSON format.

**Usage:**
```
@prd-to-ralph Convert specs/my-feature/spec.md to prd.json
```

**Output:** `specs/[feature-name]/prd.json`

## Workflow

1. **Generate PRD:**
   ```
   @prd-generator Add status tracking to versions
   ```
   Answer clarifying questions, get `specs/add-status/spec.md`

2. **Convert to Ralph format:**
   ```
   @prd-to-ralph Convert specs/add-status/spec.md
   ```
   Get `specs/add-status/prd.json`

3. **Run Ralph:**
   ```bash
   cp specs/add-status/prd.json ./prd.json
   ./ralph-kiro.sh
   ```

## Notes

- Stories sized for one context window
- Verifiable acceptance criteria
- Dependency-ordered execution
