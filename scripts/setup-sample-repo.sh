#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Default repository name
DEFAULT_REPO="ai-rules-manager-sample-git-registry"

usage() {
    echo "Usage: $0 [repo-name]"
    echo ""
    echo "Creates a comprehensive sample repository for ARM testing."
    echo ""
    echo "Arguments:"
    echo "  repo-name    - Name for sample repository (default: $DEFAULT_REPO)"
    echo ""
    echo "Examples:"
    echo "  $0                    # Use default name"
    echo "  $0 my-sample-repo    # Custom name"
    echo ""
    echo "Requirements:"
    echo "  - GitHub CLI (gh) must be installed and authenticated"
    echo "  - Git must be installed"
}

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}



create_version_1_0_0() {
    mkdir -p rules/cursor rules/amazonq rules/copilot

    cat > README.md << 'EOF'
# AI Rules Test Repository

Test repository for ARM (AI Rules Manager) with grug-brained-dev rules.

## Repository Structure

### Legacy Format (v1.x)
- `rules/cursor/grug-brained-dev.mdc` - Cursor rules for grug-brained development
- `rules/amazonq/grug-brained-dev.md` - Amazon Q rules for grug-brained development
- `rules/copilot/grug-brained-dev.instructions.md` - GitHub Copilot instructions for grug-brained development

### URF Format (v2.x+)
- `rulesets/grug-brained-dev.yml` - Universal Rule Format specification with grug-brained principles

## Version History

- **v1.0.0** - Initial release with basic grug rules
- **v1.1.0** - Added task management rules
- **v1.0.1** - Bug fixes in grug rules
- **v2.0.0** - **BREAKING**: Introduced URF format with structured rules
- **v2.1.0** - Added clean code rules

## URF Format Benefits

- **Structured Rules**: Consistent metadata and priority system
- **Multi-tool Support**: Single source compiles to cursor, amazonq, copilot formats
- **Priority Management**: Rules have explicit priority and enforcement levels
- **Scope Definition**: Rules specify which files they apply to
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


}

create_version_1_1_0() {
    # Add task management rules
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

    cat > rules/cursor/process-tasks.mdc << 'EOF'
# Process Tasks Rules (Cursor)

*Grug process tasks efficiently.*

## Task Processing
- One task at time
- Finish before start new
- Mark done when complete

## Task Review
- Check work before mark done
- Ask for help if stuck
- Learn from mistakes
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

    cat > rules/amazonq/process-tasks.md << 'EOF'
# Process Tasks Rules (Amazon Q)

*Grug process tasks efficiently.*

## Task Processing
- One task at time
- Finish before start new
- Mark done when complete

## Task Review
- Check work before mark done
- Ask for help if stuck
- Learn from mistakes
EOF

    cat > rules/copilot/generate-tasks.instructions.md << 'EOF'
---
description: 'Task generation instructions for GitHub Copilot'
---

# Generate Tasks Instructions

*Grug generate tasks for work.*

## Instructions

- Break big work into small tasks
- Small tasks easier for grug brain
- Write tasks down so grug no forget
- Important tasks first
- Easy tasks when grug tired
- Hard tasks when grug fresh

## Additional Guidelines

- Use clear task descriptions
- Include acceptance criteria
- Estimate task complexity
EOF

    cat > rules/copilot/process-tasks.instructions.md << 'EOF'
---
description: 'Task processing instructions for GitHub Copilot'
---

# Process Tasks Instructions

*Grug process tasks efficiently.*

## Instructions

- One task at time
- Finish before start new
- Mark done when complete
- Check work before mark done
- Ask for help if stuck
- Learn from mistakes

## Additional Guidelines

- Update task status regularly
- Document blockers and solutions
- Review completed work
EOF


}

create_version_1_2_0() {
    # Bug fix in grug-brained-dev files (patch release)
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
- FIXED: Grug remember to save work often
- Grug test before ship
- Broken code make grug sad
- Test save grug from angry users
- FIXED: Grug test edge cases too

## Additional Guidelines

- Write code that tells story
- Use names that make sense
- Small functions better than big functions
EOF


}

create_version_2_0_0() {
    # Add URF format file
    mkdir -p rulesets
    cat > rulesets/grug-brained-dev.yml << 'EOF'
version: "1.0"
metadata:
  id: "ai-rules-sample"
  name: "AI Rules Sample Ruleset"
  version: "2.0.0"
  description: "Sample URF ruleset for ARM testing with grug-brained development principles"
rules:
  - id: "grug-simplicity"
    name: "Grug Simplicity Rule"
    description: "Keep code simple for grug brain to understand"
    priority: 100
    enforcement: "must"
    scope:
      - files: ["**/*.js", "**/*.ts", "**/*.py", "**/*.go"]
    body: |
      Grug no like complex code. Simple code good, complex code bad.

      ## Rules
      - If grug no understand, too complex
      - Use simple names that make sense
      - Small functions better than big functions
      - One thing per function

      ## Examples
      ```javascript
      // Good - grug understand
      function addNumbers(a, b) {
        return a + b;
      }

      // Bad - grug confused
      const performArithmeticOperation = (operandA, operandB, operationType) => {
        return operationType === 'addition' ? operandA + operandB : null;
      }
      ```
  - id: "grug-testing"
    name: "Grug Testing Rule"
    description: "Test everything before ship to avoid angry users"
    priority: 90
    enforcement: "should"
    scope:
      - files: ["**/*.test.js", "**/*.spec.ts", "**/*_test.py"]
    body: |
      Grug test before ship. Broken code make grug sad.

      ## Rules
      - Test save grug from angry users
      - Write test for new code
      - Test edge cases too
      - Run tests before commit

      ## Test Types
      - Unit tests for small parts
      - Integration tests for big parts
      - End-to-end tests for whole thing
  - id: "grug-documentation"
    name: "Grug Documentation Rule"
    description: "Document why, not what - code should explain what"
    priority: 70
    enforcement: "may"
    scope:
      - files: ["**/*.md", "**/*.rst", "**/*.txt"]
    body: |
      Grug write docs when needed. Code should tell story.

      ## Rules
      - Comments explain why, not what
      - README tell how to use
      - Document complex business logic
      - Keep docs up to date

      ## Good Comments
      ```javascript
      // Retry 3 times because API sometimes flaky
      for (let i = 0; i < 3; i++) {
        try {
          return await apiCall();
        } catch (error) {
          if (i === 2) throw error;
        }
      }
      ```
EOF

    # Breaking changes to task rules
    cat > rules/cursor/generate-tasks.mdc << 'EOF'
# Generate Tasks Rules v2 (Cursor)

*BREAKING CHANGE: New task generation system.*

## Advanced Task Generation
- Use AI to break down complex work
- Estimate time for each task
- Assign priority scores automatically

## Task Templates
- Predefined templates for common work
- Customizable task structures
- Integration with project management tools
EOF

    cat > rules/cursor/process-tasks.mdc << 'EOF'
# Process Tasks Rules v2 (Cursor)

*BREAKING CHANGE: New task processing workflow.*

## Automated Processing
- Smart task scheduling
- Dependency tracking
- Progress monitoring

## Quality Gates
- Automated quality checks
- Peer review requirements
- Documentation standards
EOF

    cat > rules/amazonq/generate-tasks.md << 'EOF'
# Generate Tasks Rules v2 (Amazon Q)

*BREAKING CHANGE: New task generation system.*

## Advanced Task Generation
- Use AI to break down complex work
- Estimate time for each task
- Assign priority scores automatically

## Task Templates
- Predefined templates for common work
- Customizable task structures
- Integration with project management tools
EOF

    cat > rules/amazonq/process-tasks.md << 'EOF'
# Process Tasks Rules v2 (Amazon Q)

*BREAKING CHANGE: New task processing workflow.*

## Automated Processing
- Smart task scheduling
- Dependency tracking
- Progress monitoring

## Quality Gates
- Automated quality checks
- Peer review requirements
- Documentation standards
EOF

    cat > rules/copilot/generate-tasks.instructions.md << 'EOF'
---
description: 'Advanced task generation instructions for GitHub Copilot v2'
---

# Generate Tasks Instructions v2

*BREAKING CHANGE: New task generation system.*

## Instructions

- Use AI to break down complex work
- Estimate time for each task
- Assign priority scores automatically
- Use predefined templates for common work
- Customize task structures
- Integrate with project management tools

## Additional Guidelines

- Leverage automation for task creation
- Apply consistent task categorization
- Include dependency mapping
EOF

    cat > rules/copilot/process-tasks.instructions.md << 'EOF'
---
description: 'Advanced task processing instructions for GitHub Copilot v2'
---

# Process Tasks Instructions v2

*BREAKING CHANGE: New task processing workflow.*

## Instructions

- Use smart task scheduling
- Track dependencies automatically
- Monitor progress continuously
- Apply automated quality checks
- Require peer review
- Follow documentation standards

## Additional Guidelines

- Implement quality gates
- Use automated validation
- Maintain audit trails
EOF
}

check_dependencies() {
    log "Checking dependencies..."

    if ! command -v gh &> /dev/null; then
        error "GitHub CLI (gh) not found!"
        echo "Please install it from: https://cli.github.com/"
        echo "Then run: gh auth login"
        return 1
    fi

    if ! command -v git &> /dev/null; then
        error "Git not found!"
        echo "Please install Git first."
        return 1
    fi

    # Check if authenticated
    if ! gh auth status &> /dev/null; then
        error "GitHub CLI not authenticated!"
        echo "Please run: gh auth login"
        return 1
    fi

    success "Dependencies check passed"
}

create_sample_repo() {
    local repo_name="$1"
    local temp_dir="/tmp/arm-setup-$$"

    log "Checking if repository exists: $repo_name"

    # Check if repo already exists
    if gh repo view "$repo_name" &> /dev/null; then
        error "Repository $repo_name already exists!"
        echo "Please choose a different name or delete the existing repository."
        echo "To delete: gh repo delete $repo_name"
        return 1
    fi

    log "Creating sample repository: $repo_name"

    mkdir -p "$temp_dir"
    cd "$temp_dir"

    git init

    # Create v1.0.0 - Basic content
    create_version_1_0_0
    git add .
    git commit -m "feat: initial ARM sample repository with grug-brained-dev rules"
    git tag v1.0.0

    # Create v1.1.0 - Enhanced content
    create_version_1_1_0
    git add .
    git commit -m "feat: add task management rules"
    git tag v1.1.0

    # Create v1.0.1 - Bug fix
    create_version_1_2_0
    git add .
    git commit -m "fix: bug fix in grug-brained-dev.mdc rule"
    git tag v1.0.1

    # Create v2.0.0-rc.1 - Pre-release with breaking changes
    git checkout -b rc
    create_version_2_0_0
    git add .
    git commit -m "feat!: breaking changes to task rules (testing phase)"
    git tag v2.0.0-rc.1

    # Create v2.0.0 - Merge breaking changes to main
    git checkout main
    git merge rc --no-ff -m "feat!: breaking changes merged to main (stable release)"
    git tag v2.0.0

    # Create v2.1.0 - Add clean code rules
    cat > rules/cursor/clean-code.mdc << 'EOF'
# Clean Code Rules (Cursor)

*Grug write clean code for happy team.*

## Clean Code Principles
- Code should tell story
- Names should make sense
- Functions should be small
- Comments explain why, not what

## Refactoring
- Clean code little bit every day
- Remove dead code
- Fix bad names when see them
EOF

    cat > rules/amazonq/clean-code.md << 'EOF'
# Clean Code Rules (Amazon Q)

*Grug write clean code for happy team.*

## Clean Code Principles
- Code should tell story
- Names should make sense
- Functions should be small
- Comments explain why, not what

## Refactoring
- Clean code little bit every day
- Remove dead code
- Fix bad names when see them
EOF

    cat > rules/copilot/clean-code.instructions.md << 'EOF'
---
description: 'Clean code instructions for GitHub Copilot'
---

# Clean Code Instructions

*Grug write clean code for happy team.*

## Instructions

- Code should tell story
- Names should make sense
- Functions should be small
- Comments explain why, not what
- Clean code little bit every day
- Remove dead code
- Fix bad names when see them

## Additional Guidelines

- Apply consistent formatting
- Use meaningful variable names
- Keep functions focused on single responsibility
- Refactor continuously
EOF

    git add .
    git commit -m "feat: add clean code rules (new features)"
    git tag v2.1.0

    # Create and push repository
    gh repo create "$repo_name" --public --source=. --remote=origin --push
    git push origin --tags

    cd /
    rm -rf "$temp_dir"

    success "Sample repository created: https://github.com/$(gh api user --jq .login)/$repo_name"
}



main() {
    local repo_name="${1:-$DEFAULT_REPO}"

    # Check for help
    if [[ "$1" == "--help" || "$1" == "-h" ]]; then
        usage
        exit 0
    fi

    log "Setting up ARM sample repository..."
    log "Repository: $repo_name"

    if ! check_dependencies; then
        exit 1
    fi

    create_sample_repo "$repo_name"

    success "🎉 Sample repository created successfully!"
    echo ""
    echo "Next steps:"
    echo "1. Test your setup:"
    echo "   ./scripts/sample-workflow.sh all \"https://github.com/\$(gh api user --jq .login)/$repo_name\""
    echo ""
    echo "2. Or run interactively:"
    echo "   ./scripts/sample-workflow.sh all"
    echo ""
    echo "Your sample repository is ready for action!"
}

main "$@"
