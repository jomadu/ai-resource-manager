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

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log "=== Compile Workflow ==="

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-compile-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"
pause

# Basic compile examples
log "Compiling to Cursor format..."
./arm compile example-rulesets/*.yml --target cursor --output ./cursor-output
pause

log "Compiling to multiple targets..."
./arm compile example-rulesets/clean-code.yml --target cursor,amazonq,copilot --output ./multi-output
pause

log "Compiling with validation only..."
./arm compile example-rulesets/*.yml --target cursor --validate-only
pause

success "Compile workflow complete! Check outputs:"
echo ""
echo "Generated files:"
echo "  cursor-output/     - Cursor format (.mdc)"
echo "  multi-output/      - Multiple formats"
echo ""
echo "Try more commands:"
echo "  ./arm compile --help"
echo "  ./arm compile example-rulesets/*.yml --target amazonq --output ./amazonq-output"
