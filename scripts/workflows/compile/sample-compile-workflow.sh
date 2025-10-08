#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
pause() { echo ""; read -p "Press Enter to continue..."; echo ""; }

run_arm() {
    echo -e "${BLUE}$ ./arm $*${NC}"
    ./arm "$@"
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log "=== Compile Workflow ==="

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-compile-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"
pause

# Basic compile examples
log "Compiling rulesets to Cursor format..."
run_arm compile --target cursor example-rulesets/*.yml ./cursor-output
pause

log "Compiling rulesets to multiple targets..."
run_arm compile --target cursor example-rulesets/clean-code.yml ./multi-output/cursor
run_arm compile --target amazonq example-rulesets/clean-code.yml ./multi-output/amazonq
run_arm compile --target copilot example-rulesets/clean-code.yml ./multi-output/copilot
pause

log "Compiling with validation only..."
run_arm compile --validate-only example-rulesets/*.yml
pause

log "Compiling promptsets..."
run_arm compile --target cursor example-promptsets/*.yml ./promptset-output
pause

log "Demonstrating resource-specific compilation..."
run_arm compile --target cursor example-rulesets/clean-code.yml ./ruleset-specific-output
run_arm compile --target cursor example-promptsets/code-review.yml ./promptset-specific-output
pause

success "Compile workflow complete! Check outputs:"
echo ""
echo "Generated files:"
echo "  cursor-output/           - Cursor format (.mdc)"
echo "  multi-output/            - Multiple formats"
echo "  promptset-output/        - Promptset compilation"
echo "  ruleset-specific-output/ - Ruleset-specific compilation"
echo "  promptset-specific-output/ - Promptset-specific compilation"
echo ""
echo "Try more commands:"
echo "  ./arm compile --help"
echo "  ./arm compile --target amazonq example-rulesets/*.yml ./amazonq-output"
echo "  ./arm compile --target cursor example-promptsets/*.yml ./promptset-output"
echo ""
echo "Resource-specific compilation:"
echo "  ./arm compile --target cursor example-rulesets/*.yml ./output"
echo "  ./arm compile --target cursor example-promptsets/*.yml ./output"
