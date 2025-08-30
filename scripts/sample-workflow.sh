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

show_tree "Initial sandbox structure"

# === SETUP ===
log "=== SETUP PHASE ==="

log "Adding registry configuration..."
./arm config add registry ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git

log "Adding sink configurations..."
./arm config add sink q --directories .amazonq/rules --include ai-rules/amazonq-* --exclude ai-rules/cursor-*
./arm config add sink cursor --directories .cursor/rules --include ai-rules/cursor-* --exclude ai-rules/amazonq-*

success "Configuration complete"

log "Showing .armrc.json configuration:"
cat .armrc.json | jq .
pause

# === VERSION ===
log "=== VERSION COMMAND ==="
./arm version
pause

# === HELP ===
log "=== HELP COMMAND ==="
./arm help
pause

# === INSTALL - Latest Version ===
log "=== INSTALL - Latest Version ==="
log "Installing latest version (should resolve to ^2.1.0)..."
./arm install ai-rules/amazonq-rules --include rules/amazonq/*.md
./arm install ai-rules/cursor-rules --include rules/cursor/*.mdc

success "Installation complete"

log "Generated arm.json:"
cat arm.json | jq .

log "Generated arm.lock:"
cat arm.lock | jq .

show_tree "Project structure after latest install"
pause

# === LIST ===
log "=== LIST COMMAND ==="
./arm list
pause

# === INFO - Single Ruleset ===
log "=== INFO - Single Ruleset ==="
./arm info ai-rules/amazonq-rules
pause

# === INFO - All Rulesets ===
log "=== INFO - All Rulesets ==="
./arm info
pause

# === UNINSTALL ===
log "=== UNINSTALL ==="
log "Uninstalling cursor rules..."
./arm uninstall ai-rules/cursor-rules

success "Uninstall complete"

log "Updated arm.json:"
cat arm.json | jq .

log "Updated arm.lock:"
cat arm.lock | jq .

show_tree "Project structure after uninstall"
pause

# === INSTALL - Specific Version ===
log "=== INSTALL - Specific Version ==="
log "Installing specific version 1.0.0..."
./arm install ai-rules/cursor-rules@1.0.0 --include rules/cursor/*.mdc

success "Version-specific installation complete"

log "Updated arm.json:"
cat arm.json | jq .

log "Updated arm.lock:"
cat arm.lock | jq .

show_tree "Project structure after version-specific install"
pause

# === OUTDATED ===
log "=== OUTDATED COMMAND ==="
./arm outdated
pause

# === UPDATE ===
log "=== UPDATE COMMAND ==="
log "Updating cursor rules to latest compatible version..."
./arm update ai-rules/cursor-rules

success "Update complete"

log "Updated arm.lock after update:"
cat arm.lock | jq .

show_tree "Project structure after update"
pause

# === INSTALL - Major Version ===
log "=== INSTALL - Major Version ==="
log "Reinstalling with major version constraint..."
rm -rf .cursor .amazonq arm.json arm.lock

./arm install ai-rules/amazonq-rules@1 --include rules/amazonq/*.md
./arm install ai-rules/cursor-rules@1 --include rules/cursor/*.mdc

success "Major version installation complete"

log "arm.json with major version constraints:"
cat arm.json | jq .

log "arm.lock with resolved versions:"
cat arm.lock | jq .

show_tree "Project structure with major version constraints"
pause

# === INSTALL - Minor Version ===
log "=== INSTALL - Minor Version ==="
log "Reinstalling with minor version constraint..."
rm -rf .cursor .amazonq arm.json arm.lock

./arm install ai-rules/amazonq-rules@1.0 --include rules/amazonq/*.md
./arm install ai-rules/cursor-rules@1.0 --include rules/cursor/*.mdc

success "Minor version installation complete"

log "arm.json with minor version constraints:"
cat arm.json | jq .

log "arm.lock with resolved versions:"
cat arm.lock | jq .

show_tree "Project structure with minor version constraints"
pause

# === INSTALL - Branch ===
log "=== INSTALL - Branch ==="
log "Reinstalling from main branch..."
rm -rf .cursor .amazonq arm.json arm.lock

./arm install ai-rules/amazonq-rules@main --include rules/amazonq/*.md
./arm install ai-rules/cursor-rules@main --include rules/cursor/*.mdc

success "Branch installation complete"

log "arm.json with branch constraints:"
cat arm.json | jq .

log "arm.lock with commit hashes:"
cat arm.lock | jq .

show_tree "Project structure with branch tracking"
pause

# === INSTALL FROM MANIFEST ===
log "=== INSTALL FROM MANIFEST ==="
log "Removing installed files and reinstalling from manifest..."
rm -rf .cursor .amazonq

./arm install

success "Install from manifest complete"

show_tree "Project structure after manifest install"
pause

# === INSTALL FROM LOCKFILE ONLY ===
log "=== INSTALL FROM LOCKFILE ONLY ==="
log "Removing manifest and installing from lockfile only..."
rm -rf .cursor .amazonq arm.json

./arm install

success "Install from lockfile complete"

log "Regenerated arm.json from lockfile:"
cat arm.json | jq .

show_tree "Project structure after lockfile install"
pause

# === UPDATE ALL ===
log "=== UPDATE ALL ==="
log "Updating all rulesets..."
./arm update

success "Update all complete"

log "Final arm.lock:"
cat arm.lock | jq .

show_tree "Final project structure"

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
echo "• arm.lock - Lockfile"
echo "• .cursor/rules/ - Cursor rules"
echo "• .amazonq/rules/ - Amazon Q rules"
