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
run_arm config sink add q-generic --directories .amazonq/rules --include "ai-rules/q-generic"
run_arm config sink add cursor --directories .cursor/rules --include "ai-rules/cursor-*"
run_arm config sink add copilot --directories .github/instructions --include "ai-rules/copilot-*" --layout flat

success "Configuration complete"

log "Showing arm.json manifest:"
cat arm.json | jq .

log "Showing .armrc.json configuration:"
cat .armrc.json | jq .
pause

# === VERSION ===
log "=== VERSION COMMAND ==="
run_arm version
pause

# === HELP ===
log "=== HELP COMMAND ==="
run_arm help
pause

# === INSTALL - Latest Version ===
log "=== INSTALL - Latest Version ==="
log "Installing latest version (should resolve to ^2.1.0)..."
run_arm install ai-rules/amazonq-rules --include "rules/amazonq/*.md"
run_arm install ai-rules/cursor-rules --include "rules/cursor/*.mdc"
run_arm install ai-rules/copilot-rules --include "rules/copilot/*.instructions.md"
run_arm install ai-rules/q-generic --include "rules/**/generic.md"

success "Installation complete"

log "Generated arm.json:"
cat arm.json | jq .

log "Generated arm-lock.json:"
cat arm-lock.json | jq .

show_tree "Project structure after latest install"
pause

# === LIST ===
log "=== LIST COMMAND ==="
run_arm list
pause

# === INFO - Single Ruleset ===
log "=== INFO - Single Ruleset ==="
run_arm info ai-rules/amazonq-rules
pause

# === INFO - All Rulesets ===
log "=== INFO - All Rulesets ==="
run_arm info
pause

# === UNINSTALL ===
log "=== UNINSTALL ==="
log "Uninstalling cursor rules..."
run_arm uninstall ai-rules/cursor-rules

success "Uninstall complete"

log "Updated arm.json:"
cat arm.json | jq .

log "Updated arm-lock.json:"
cat arm-lock.json | jq .

show_tree "Project structure after uninstall"
pause

# === INSTALL - Specific Version ===
log "=== INSTALL - Specific Version ==="
log "Installing specific version 1.0.0..."
run_arm install ai-rules/cursor-rules@1.0.0 --include "rules/cursor/*.mdc"
run_arm install ai-rules/copilot-rules@1.0.0 --include "rules/copilot/*.instructions.md"

success "Version-specific installation complete"

log "Updated arm.json:"
cat arm.json | jq .

log "Updated arm-lock.json:"
cat arm-lock.json | jq .

show_tree "Project structure after version-specific install"
pause

# === OUTDATED ===
log "=== OUTDATED COMMAND ==="
run_arm outdated
pause

# === UPDATE ===
log "=== UPDATE COMMAND ==="
log "Updating cursor and copilot rules to latest compatible version..."
run_arm update ai-rules/cursor-rules
run_arm update ai-rules/copilot-rules

success "Update complete"

log "Updated arm-lock.json after update:"
cat arm-lock.json | jq .

show_tree "Project structure after update"
pause

# === INSTALL - Major Version ===
log "=== INSTALL - Major Version ==="
log "Reinstalling with major version constraint..."
rm -rf .cursor .amazonq arm.json arm-lock.json

run_arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
run_arm install ai-rules/amazonq-rules@1 --include "rules/amazonq/*.md"
run_arm install ai-rules/cursor-rules@1 --include "rules/cursor/*.mdc"
run_arm install ai-rules/copilot-rules@1 --include "rules/copilot/*.instructions.md"

success "Major version installation complete"

log "arm.json with major version constraints:"
cat arm.json | jq .

log "arm-lock.json with resolved versions:"
cat arm-lock.json | jq .

show_tree "Project structure with major version constraints"
pause

# === INSTALL - Minor Version ===
log "=== INSTALL - Minor Version ==="
log "Reinstalling with minor version constraint..."
rm -rf .cursor .amazonq arm.json arm-lock.json

run_arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
run_arm install ai-rules/amazonq-rules@1.0 --include "rules/amazonq/*.md"
run_arm install ai-rules/cursor-rules@1.0 --include "rules/cursor/*.mdc"
run_arm install ai-rules/copilot-rules@1.0 --include "rules/copilot/*.instructions.md"

success "Minor version installation complete"

log "arm.json with minor version constraints:"
cat arm.json | jq .

log "arm-lock.json with resolved versions:"
cat arm-lock.json | jq .

show_tree "Project structure with minor version constraints"
pause

# === INSTALL - Branch ===
log "=== INSTALL - Branch ==="
log "Reinstalling from main branch..."
rm -rf .cursor .amazonq arm.json arm-lock.json

run_arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git
run_arm install ai-rules/amazonq-rules@main --include "rules/amazonq/*.md"
run_arm install ai-rules/cursor-rules@main --include "rules/cursor/*.mdc"
run_arm install ai-rules/copilot-rules@main --include "rules/copilot/*.instructions.md"

success "Branch installation complete"

log "arm.json with branch constraints:"
cat arm.json | jq .

log "arm-lock.json with commit hashes:"
cat arm-lock.json | jq .

show_tree "Project structure with branch tracking"
pause

# === INSTALL FROM MANIFEST ===
log "=== INSTALL FROM MANIFEST ==="
log "Removing installed files and reinstalling from manifest..."
rm -rf .cursor .amazonq .github

run_arm install

success "Install from manifest complete"

show_tree "Project structure after manifest install"
pause

# === SUMMARY ===
log "=== WORKFLOW COMPLETE ==="
success "Sample workflow completed successfully!"
echo ""
echo "This workflow demonstrated:"
echo "• Registry and sink configuration"
echo "• Installing with different version constraints"
echo "• Listing and getting info about rulesets"
echo "• Uninstalling rulesets"
echo "• Checking for outdated rulesets"
echo "• Updating rulesets"
echo "• Installing from manifest and lockfile"
echo "• File structure management"
echo ""
echo "Check the generated files:"
echo "• .armrc.json - Configuration"
echo "• arm.json - Manifest"
echo "• arm-lock.json - Lockfile"
echo "• .cursor/rules/ - Cursor rules (hierarchical)"
echo "• .amazonq/rules/ - Amazon Q rules (hierarchical)"
echo "• .github/instructions/ - GitHub Copilot instructions (flat layout)"
