#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }

usage() {
    echo "Usage: $0 [repo-path]"
    echo ""
    echo "Creates a local Git repository for ARM testing."
    echo ""
    echo "Arguments:"
    echo "  repo-path    - Path for local repository (default: ./)"
}

if [[ "$1" == "--help" || "$1" == "-h" ]]; then
    usage
    exit 0
fi

create_version_1_0_0() {
    mkdir -p rules/cursor rules/amazonq rules/copilot promptsets

    cat > README.md << 'EOF'
# AI Rules Test Repository (Local)

Local test repository for ARM (AI Rules Manager) with grug-brained-dev rules.

## Repository Structure

### Legacy Format (v1.x)
- `rules/cursor/grug-brained-dev.mdc` - Cursor rules for grug-brained development
- `rules/amazonq/grug-brained-dev.md` - Amazon Q rules for grug-brained development
- `rules/copilot/grug-brained-dev.instructions.md` - GitHub Copilot instructions for grug-brained development

### Resource Format (v2.x+)
- `rulesets/grug-brained-dev.yml` - Ruleset specification with grug-brained principles
- `promptsets/code-review.yml` - Promptset specification for code review prompts

## Version History

- **v1.0.0** - Initial release with basic grug rules
- **v1.1.0** - Added task management rules
- **v1.0.1** - Bug fixes in grug rules
- **v2.0.0** - **BREAKING**: Introduced URF format with structured rules
- **v2.1.0** - Added clean code rules
EOF

    cat > rules/cursor/grug-brained-dev.mdc << 'EOF'
# Grug Brained Dev Rules (Cursor)

*Simple rules for simple grug brain.*

## Grug Rule 1: Keep Simple
- Grug no like complex code
- Simple code good, complex code bad
- If grug no understand, too complex

## Grug Rule 2: Test Everything
- Grug test before ship
- Broken code make grug sad
- Test save grug from angry users
EOF

    cat > rules/amazonq/grug-brained-dev.md << 'EOF'
# Grug Brained Dev Rules (Amazon Q)

*Simple rules for simple grug brain.*

## Grug Rule 1: Keep Simple
- Grug no like complex code
- Simple code good, complex code bad
- If grug no understand, too complex

## Grug Rule 2: Test Everything
- Grug test before ship
- Broken code make grug sad
- Test save grug from angry users
EOF

    cat > rules/copilot/grug-brained-dev.instructions.md << 'EOF'
---
description: 'Grug-brained development instructions for GitHub Copilot'
---

# Grug Brained Dev Instructions

*Simple rules for simple grug brain.*

## Instructions

- Keep code simple - grug no like complex code
- Simple code good, complex code bad
- If grug no understand, too complex
- Grug test before ship
- Broken code make grug sad
- Test save grug from angry users

## Additional Guidelines

- Write code that tells story
- Use names that make sense
- Small functions better than big functions
EOF

    cat > promptsets/code-review.yml << 'EOF'
apiVersion: v1
kind: Promptset
metadata:
  id: "codeReview"
  name: "Code Review Assistant"
  description: "Prompts for comprehensive code review analysis"
spec:
  prompts:
    review-analysis:
      name: "Code Review Analysis"
      description: "Analyze code for quality, security, and best practices"
      body: |
        Please review this code for:
        1. Code quality and readability
        2. Security vulnerabilities
        3. Performance issues
        4. Best practices adherence
        5. Potential bugs or edge cases

        Provide specific feedback with line numbers and suggestions for improvement.
    architecture-review:
      name: "Architecture Review"
      description: "Review code architecture and design patterns"
      body: |
        Analyze the architecture and design of this code:
        1. Design patterns used
        2. Separation of concerns
        3. Coupling and cohesion
        4. Scalability considerations
        5. Maintainability factors

        Suggest architectural improvements if needed.
EOF
}

create_version_1_1_0() {
    cat > rules/cursor/generate-tasks.mdc << 'EOF'
# Generate Tasks Rules (Cursor)

*Grug generate tasks for work.*

## Task Generation
- Break big work into small tasks
- Small tasks easier for grug brain
- Write tasks down so grug no forget

## Task Priority
- Important tasks first
- Easy tasks when grug tired
- Hard tasks when grug fresh
EOF

    cat > rules/amazonq/generate-tasks.md << 'EOF'
# Generate Tasks Rules (Amazon Q)

*Grug generate tasks for work.*

## Task Generation
- Break big work into small tasks
- Small tasks easier for grug brain
- Write tasks down so grug no forget

## Task Priority
- Important tasks first
- Easy tasks when grug tired
- Hard tasks when grug fresh
EOF

    cat > promptsets/testing.yml << 'EOF'
apiVersion: v1
kind: Promptset
metadata:
  id: "testing"
  name: "Testing Assistant"
  description: "Prompts for generating and improving test code"
spec:
  prompts:
    test-generation:
      name: "Test Generation"
      description: "Generate comprehensive test cases for code"
      body: |
        Generate comprehensive test cases for this code:
        1. Unit tests for all public methods
        2. Edge cases and boundary conditions
        3. Error handling scenarios
        4. Integration test suggestions
        5. Performance test considerations

        Include test data setup and expected outcomes.
EOF
}

create_version_1_0_1() {
    cat > rules/cursor/grug-brained-dev.mdc << 'EOF'
# Grug Brained Dev Rules (Cursor)

*Simple rules for simple grug brain.*

## Grug Rule 1: Keep Simple
- Grug no like complex code
- Simple code good, complex code bad
- If grug no understand, too complex
- FIXED: Grug remember to save work often

## Grug Rule 2: Test Everything
- Grug test before ship
- Broken code make grug sad
- Test save grug from angry users
- FIXED: Grug test edge cases too
EOF
}

create_version_2_0_0() {
    mkdir -p rulesets
    cat > rulesets/grug-brained-dev.yml << 'EOF'
apiVersion: v1
kind: Ruleset
metadata:
  id: "grugBrainedDev"
  name: "Grug Brained Development"
  description: "Sample ruleset for ARM testing with grug-brained development principles"
spec:
  rules:
    grug-simplicity:
      name: "Grug Simplicity Rule"
      description: "Keep code simple for grug brain to understand"
      priority: 100
      enforcement: must
      scope:
        - files: ["**/*.js", "**/*.ts", "**/*.py", "**/*.go"]
      body: |
        Grug no like complex code. Simple code good, complex code bad.

        ## Rules
        - If grug no understand, too complex
        - Use simple names that make sense
        - Small functions better than big functions
        - One thing per function
    grug-testing:
      name: "Grug Testing Rule"
      description: "Test everything before ship to avoid angry users"
      priority: 90
      enforcement: should
      body: |
        Grug test before ship. Broken code make grug sad.

        ## Rules
        - Test save grug from angry users
        - Write test for new code
        - Test edge cases too
        - Run tests before commit
EOF
}

create_local_repo() {
    local repo_path="$1"

    if [ -d "$repo_path" ]; then
        error "Repository path already exists: $repo_path"
        echo "Remove it first: rm -rf $repo_path"
        return 1
    fi

    log "Creating local repository at: $repo_path"

    mkdir -p "$repo_path"
    cd "$repo_path"
    git init
    git config user.email "test@example.com"
    git config user.name "ARM Test"

    create_version_1_0_0
    git add .
    git commit -m "feat: initial ARM sample repository with grug-brained-dev rules"
    git tag v1.0.0

    create_version_1_0_1
    git add .
    git commit -m "fix: bug fix in grug-brained-dev.mdc rule"
    git tag v1.0.1

    create_version_1_1_0
    git add .
    git commit -m "feat: add task management rules"
    git tag v1.1.0

    create_version_2_0_0
    git add .
    git commit -m "feat!: breaking changes with URF format"
    git tag v2.0.0

    git checkout -b develop
    echo "# Development Branch" >> README.md
    git add .
    git commit -m "chore: create develop branch"
    git checkout main

    success "Local repository created at: $repo_path"
}

REPO_PATH="${1:-.}"

log "Setting up local ARM test repository..."
log "Repository path: $REPO_PATH"

create_local_repo "$REPO_PATH"

success "Local repository created at: $REPO_PATH"
