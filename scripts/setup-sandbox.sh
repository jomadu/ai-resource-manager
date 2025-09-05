#!/bin/bash

set -e

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

pause() {
    echo ""
    read -p "Press Enter to continue..."
    echo ""
}

# Build ARM
log "Building ARM..."
cd ..
make build
success "ARM built successfully"

# Setup sandbox
log "Setting up sandbox environment..."
rm -rf sandbox/
mkdir sandbox
cp ./bin/arm ./sandbox
cd sandbox
run_arm cache nuke

show_tree "Initial sandbox structure"

# === SETUP ===
log "=== SETUP PHASE ==="

log "Adding registry configuration..."
run_arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git

log "Adding sink configurations..."
run_arm config sink add q --directories .amazonq/rules --include "ai-rules/amazonq-*"
run_arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
run_arm config sink add copilot --directories .github/instructions --include "ai-rules/copilot-*" --layout flat

success "Configuration complete"

log "Showing arm.json manifest:"
cat arm.json | jq .

log "Showing .armrc.json configuration:"
cat .armrc.json | jq .
pause
