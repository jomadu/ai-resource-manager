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

# === SETUP SANDBOX ===
log "=== SETUP SANDBOX ==="

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

log "Running setup-sandbox script..."
"$SCRIPT_DIR/setup-sandbox.sh"

log "Entering sandbox..."
cd "$PROJECT_ROOT/sandbox"

show_tree "Initial sandbox structure"
pause

# === BASIC COMMANDS ===
log "=== BASIC COMMANDS ==="

log "Running arm help..."
run_arm help
pause

log "Running arm version..."
run_arm version
pause

# === REGISTRY SETUP ===
log "=== REGISTRY SETUP ==="

log "Setting up git registry..."
run_arm config registry add ai-rules https://github.com/jomadu/ai-rules-manager-sample-git-registry --type git

log "Showing configuration..."
run_arm config list
pause

# === SINK SETUP ===
log "=== SINK SETUP ==="

log "Setting up cursor sink (hierarchical)..."
run_arm config sink add cursor .cursor/rules

log "Setting up Amazon Q sink (hierarchical)..."
run_arm config sink add q .amazonq/rules

log "Setting up copilot sink (flat)..."
run_arm config sink add copilot .github/copilot --layout flat

log "Showing configuration..."
run_arm config list
pause

# === INSTALL RULESETS ===
log "=== INSTALL RULESETS ==="

log "Installing cursor rules to cursor sink..."
run_arm install ai-rules/cursor-rules --include "rules/cursor/*.mdc" --sinks cursor

log "Installing Amazon Q rules to q sink..."
run_arm install ai-rules/amazonq-rules --include "rules/amazonq/*.md" --sinks q

log "Installing copilot rules to copilot sink..."
run_arm install ai-rules/copilot-rules --include "rules/copilot/*.instructions.md" --sinks copilot

show_tree "Project structure after installs"
pause

# === LIST AND INFO ===
log "=== LIST AND INFO ==="

log "Running arm list..."
run_arm list
pause

log "Running arm info (all rulesets)..."
run_arm info
pause

log "Running arm info on cursor ruleset..."
run_arm info ai-rules/cursor-rules
pause

# === UNINSTALL ALL ===
log "=== UNINSTALL ALL ==="

log "Uninstalling all rulesets..."
run_arm uninstall ai-rules/cursor-rules
run_arm uninstall ai-rules/amazonq-rules
run_arm uninstall ai-rules/copilot-rules

log "Showing empty list..."
run_arm list
pause

# === INSTALL FROM MAIN BRANCH ===
log "=== INSTALL FROM MAIN BRANCH ==="

log "Installing cursor ruleset from main branch..."
run_arm install ai-rules/cursor-rules@main --include "rules/cursor/*.mdc" --sinks cursor

log "Showing info for main branch install..."
run_arm info ai-rules/cursor-rules
pause

# === OUTDATED CHECK ===
log "=== OUTDATED CHECK ==="

log "Checking for outdated rulesets..."
run_arm outdated
pause

# === VERSION CONSTRAINT DEMOS ===
log "=== VERSION CONSTRAINT DEMOS ==="

log "Installing cursor ruleset with major version 1 (should resolve to 1.1.0)..."
run_arm install ai-rules/cursor-rules@1 --include "rules/cursor/*.mdc" --sinks cursor

log "Showing info (should show 1.1.0)..."
run_arm info ai-rules/cursor-rules
pause

log "Installing cursor ruleset with minor version 1.0 (should resolve to 1.0.1)..."
run_arm install ai-rules/cursor-rules@1.0 --include "rules/cursor/*.mdc" --sinks cursor

log "Showing info (should show 1.0.1)..."
run_arm info ai-rules/cursor-rules
pause

log "Installing cursor ruleset with patch version 1.0.0 (should resolve to 1.0.0)..."
run_arm install ai-rules/cursor-rules@1.0.0 --include "rules/cursor/*.mdc" --sinks cursor

log "Showing info (should show 1.0.0)..."
run_arm info ai-rules/cursor-rules
pause

# === SINK REMOVAL PROTECTION ===
log "=== SINK REMOVAL PROTECTION ==="

log "Attempting to remove cursor sink (should fail because ruleset is installed)..."
if run_arm config sink remove cursor 2>&1; then
    error "Sink removal should have failed!"
else
    success "Sink removal correctly blocked due to active ruleset"
fi
pause

# === CLEANUP ===
log "=== CLEANUP ==="

log "Removing cursor ruleset..."
run_arm uninstall ai-rules/cursor-rules

log "Now removing cursor sink (should succeed)..."
run_arm config sink remove cursor

success "Cleanup complete"
pause

# === SUMMARY ===
log "=== WORKFLOW COMPLETE ==="
success "New workflow completed successfully!"
echo ""
echo "This workflow demonstrated:"
echo "• Sandbox setup and binary building"
echo "• Basic help and version commands"
echo "• Registry configuration"
echo "• Sink configuration (hierarchical and flat layouts)"
echo "• Installing rulesets to specific sinks"
echo "• Listing and getting info about rulesets"
echo "• Uninstalling rulesets"
echo "• Installing from branches"
echo "• Checking for outdated rulesets"
echo "• Version constraint resolution (major, minor, patch)"
echo "• Sink removal protection"
echo "• Clean teardown"
