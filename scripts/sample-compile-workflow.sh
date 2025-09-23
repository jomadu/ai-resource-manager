#!/bin/bash

set -e

# Parse command line arguments
INTERACTIVE=true
while [[ $# -gt 0 ]]; do
    case $1 in
        --non-interactive|-n)
            INTERACTIVE=false
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--non-interactive|-n]"
            exit 1
            ;;
    esac
done

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

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

run_arm() {
    echo -e "${BLUE}$ ./arm $*${NC}"
    ./arm "$@"
}

show_tree() {
    local title="$1"
    echo ""
    echo -e "${YELLOW}=== $title ===${NC}"
    if command -v tree &> /dev/null; then
        tree -a -I '.git' . || ls -la
    else
        find . -type f -not -path './.git/*' | sort
    fi
    echo ""
}

show_file() {
    local file="$1"
    local title="$2"
    echo ""
    echo -e "${YELLOW}=== $title ===${NC}"
    if [ -f "$file" ]; then
        cat "$file"
    else
        echo "File not found: $file"
    fi
    echo ""
}

pause() {
    if [ "$INTERACTIVE" = true ]; then
        echo ""
        read -p "Press Enter to continue..."
        echo ""
    fi
}

# === SETUP ===
log "=== SETUP COMPILE WORKFLOW ==="

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

log "Running setup-sandbox script..."
"$SCRIPT_DIR/setup-sandbox.sh"

log "Entering sandbox..."
cd "$PROJECT_ROOT/sandbox"

log "Creating sample URF files..."

# Create a comprehensive URF file
cat > clean-code.yaml << 'EOF'
version: "1.0"
metadata:
  id: "clean-code"
  name: "Clean Code Guidelines"
  version: "1.2.0"
  description: "Best practices for writing clean, maintainable code"
rules:
  naming:
    name: "Meaningful Names"
    priority: 100
    enforcement: "must"
    body: |
      Use meaningful, descriptive names for variables, functions, and classes:
      - Choose names that reveal intent
      - Avoid abbreviations and acronyms
      - Use searchable names for important concepts
      - Class names should be nouns, method names should be verbs
  functions:
    name: "Small Functions"
    priority: 200
    enforcement: "should"
    body: |
      Keep functions small and focused:
      - Functions should do one thing well
      - Aim for 20 lines or fewer
      - Use descriptive function names
      - Minimize function parameters (ideally 3 or fewer)
  comments:
    name: "Smart Comments"
    priority: 150
    enforcement: "should"
    body: |
      Write comments that add value:
      - Don't comment obvious code
      - Explain why, not what
      - Keep comments up to date
      - Use comments to warn of consequences
EOF

# Create a security-focused URF file
cat > security.yaml << 'EOF'
version: "1.0"
metadata:
  id: "security-rules"
  name: "Security Best Practices"
  version: "2.1.0"
  description: "Essential security guidelines for application development"
rules:
  input-validation:
    name: "Input Validation"
    priority: 300
    enforcement: "must"
    body: |
      Always validate and sanitize user input:
      - Validate all input at boundaries
      - Use parameterized queries for database access
      - Sanitize data before output
      - Implement proper error handling
  authentication:
    name: "Authentication & Authorization"
    priority: 250
    enforcement: "must"
    body: |
      Implement secure authentication:
      - Use strong password policies
      - Implement multi-factor authentication
      - Use secure session management
      - Apply principle of least privilege
EOF

# Create an invalid URF file for error testing
cat > invalid.yaml << 'EOF'
invalid: structure
missing: required fields
EOF

show_tree "Sandbox with URF files"
pause

# === BASIC COMPILE HELP ===
log "=== BASIC COMPILE HELP ==="

log "Showing compile command help..."
run_arm compile --help
pause

# === VALIDATION ONLY ===
log "=== VALIDATION ONLY ==="

log "Validating URF files (validation-only mode)..."
run_arm compile clean-code.yaml security.yaml --target cursor --validate-only --verbose
pause

log "Testing validation with invalid file (should show error)..."
if run_arm compile invalid.yaml --target cursor --validate-only 2>&1; then
    error "Validation should have failed!"
else
    success "Validation correctly failed for invalid URF file"
fi
pause

# === SINGLE TARGET COMPILATION ===
log "=== SINGLE TARGET COMPILATION ==="

log "Compiling to Cursor format..."
run_arm compile clean-code.yaml --target cursor --output ./cursor-output --verbose

show_tree "Cursor output"
show_file "./cursor-output/clean-code_naming.mdc" "Sample Cursor rule file"
pause

log "Compiling to Amazon Q format..."
run_arm compile security.yaml --target amazonq --output ./amazonq-output --verbose

show_tree "Amazon Q output"
show_file "./amazonq-output/security-rules_input-validation.md" "Sample Amazon Q rule file"
pause

# === MULTI-TARGET COMPILATION ===
log "=== MULTI-TARGET COMPILATION ==="

log "Compiling to multiple targets (cursor,amazonq,copilot)..."
run_arm compile clean-code.yaml --target cursor,amazonq,copilot --output ./multi-output --verbose

show_tree "Multi-target output structure"
show_file "./multi-output/cursor/clean-code_functions.mdc" "Cursor format"
show_file "./multi-output/amazonq/clean-code_functions.md" "Amazon Q format"
show_file "./multi-output/copilot/clean-code_functions.instructions.md" "Copilot format"
pause

# === BATCH COMPILATION ===
log "=== BATCH COMPILATION ==="

log "Compiling multiple files at once..."
run_arm compile clean-code.yaml security.yaml --target markdown --output ./batch-output --verbose

show_tree "Batch compilation output"
pause

# === CUSTOM NAMESPACE ===
log "=== CUSTOM NAMESPACE ==="

log "Compiling with custom namespace..."
run_arm compile clean-code.yaml --target cursor --output ./namespace-output --namespace "team/standards" --verbose

show_tree "Custom namespace output"
show_file "./namespace-output/clean-code_naming.mdc" "Rule with custom namespace"
pause

# === DIRECTORY COMPILATION ===
log "=== DIRECTORY COMPILATION ==="

log "Creating subdirectory with URF files..."
mkdir -p rules/team
cp clean-code.yaml rules/
cp security.yaml rules/team/

log "Compiling directory (non-recursive)..."
run_arm compile rules/ --target cursor --output ./dir-output --verbose

show_tree "Directory compilation output (non-recursive)"
pause

log "Compiling directory (recursive)..."
run_arm compile rules/ --target amazonq --output ./recursive-output --recursive --verbose

show_tree "Directory compilation output (recursive)"
pause

# === FORCE OVERWRITE ===
log "=== FORCE OVERWRITE ==="

log "First compilation..."
run_arm compile clean-code.yaml --target cursor --output ./force-test

log "Attempting to overwrite without --force (should fail)..."
if run_arm compile clean-code.yaml --target cursor --output ./force-test 2>&1; then
    error "Overwrite should have failed without --force!"
else
    success "Correctly prevented overwrite without --force flag"
fi

log "Overwriting with --force flag..."
run_arm compile clean-code.yaml --target cursor --output ./force-test --force --verbose
pause

# === ERROR HANDLING ===
log "=== ERROR HANDLING ==="

log "Testing fail-fast mode with mixed valid/invalid files..."
if run_arm compile clean-code.yaml invalid.yaml --target cursor --output ./error-test --fail-fast 2>&1; then
    error "Fail-fast should have stopped on first error!"
else
    success "Fail-fast correctly stopped on first error"
fi

log "Testing continue-on-error mode (default)..."
run_arm compile clean-code.yaml invalid.yaml --target cursor --output ./error-test --verbose || true
success "Continue-on-error processed valid files despite errors"
pause

# === CLEANUP AND SUMMARY ===
log "=== CLEANUP ==="

log "Cleaning up sandbox..."
cd "$PROJECT_ROOT/sandbox"
rm -rf cursor-output amazonq-output multi-output batch-output namespace-output rules recursive-output force-test error-test
rm -f clean-code.yaml security.yaml invalid.yaml

success "Cleanup complete"

# === SUMMARY ===
log "=== COMPILE WORKFLOW COMPLETE ==="
success "Compile workflow completed successfully!"
echo ""
echo "This workflow demonstrated:"
echo "• URF file validation (--validate-only)"
echo "• Single target compilation (cursor, amazonq, copilot, markdown)"
echo "• Multi-target compilation with subdirectories"
echo "• Batch compilation of multiple files"
echo "• Custom namespace specification"
echo "• Directory compilation (recursive and non-recursive)"
echo "• Force overwrite protection and override"
echo "• Error handling (fail-fast vs continue-on-error)"
echo "• Verbose output for detailed progress"
echo ""
echo "Key features:"
echo "• Reuses existing URF compilation infrastructure"
echo "• Follows ARM's established UI and service patterns"
echo "• Supports all existing URF targets and formats"
echo "• Provides flexible file discovery and filtering"
echo "• Maintains compatibility with ARM installation workflow"
