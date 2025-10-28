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
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }
warn() { echo -e "${YELLOW}⚠${NC} $1"; }

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
    if [ "$INTERACTIVE" = true ]; then
        echo ""
        read -p "Press Enter to continue..."
        echo ""
    fi
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log "=== Comprehensive Git Workflow ==="

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-git-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

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
run_arm add registry git --url https://github.com/jomadu/ai-rules-manager-sample-git-registry ai-rules

log "Showing configuration..."
run_arm list registry
run_arm list sink
pause

# === SINK SETUP ===
log "=== SINK SETUP ==="

log "Setting up cursor sink (hierarchical)..."
run_arm add sink --type cursor cursor-rules .cursor/rules

log "Setting up Amazon Q sink (hierarchical)..."
run_arm add sink --type amazonq q-rules .amazonq/rules

log "Setting up copilot sink (flat)..."
run_arm add sink --type copilot copilot-rules .github/copilot

log "Setting up cursor prompts sink..."
run_arm add sink --type cursor cursor-commands .cursor/commands

log "Setting up Amazon Q prompts sink..."
run_arm add sink --type amazonq q-prompts .amazonq/prompts

log "Showing configuration..."
run_arm list registry
run_arm list sink
pause

# === INSTALL RULESETS ===
log "=== INSTALL RULESETS ==="

log "Installing cursor rules to cursor sink..."
run_arm install ruleset ai-rules/cursor-rules --include "rules/cursor/*.mdc" cursor-rules

log "Installing Amazon Q rules to q sink..."
run_arm install ruleset ai-rules/amazonq-rules --include "rules/amazonq/*.md" q-rules

log "Installing copilot rules to copilot sink..."
run_arm install ruleset ai-rules/copilot-rules --include "rules/copilot/*.instructions.md" copilot-rules

log "Installing grug-brained-dev ruleset to all sinks..."
run_arm install ruleset --priority 150 ai-rules/grug-brained-dev --include "rulesets/grug-brained-dev.yml" cursor-rules q-rules copilot-rules

log "Installing code-review promptset to both prompt sinks..."
run_arm install promptset ai-rules/code-review --include "promptsets/code-review.yml" cursor-commands q-prompts

log "Installing testing promptset to both prompt sinks..."
run_arm install promptset ai-rules/testing --include "promptsets/testing.yml" cursor-commands q-prompts

show_tree "Project structure after installs"
pause

# === LIST AND INFO ===
log "=== LIST AND INFO ==="

log "Running arm list (all resources)..."
run_arm list
pause

log "Running arm list ruleset..."
run_arm list ruleset
pause

log "Running arm list promptset..."
run_arm list promptset
pause

log "Running arm info (all resources)..."
run_arm info
pause

log "Running arm info on cursor ruleset..."
run_arm info ruleset ai-rules/cursor-rules
pause

log "Running arm info on code-review promptset..."
run_arm info promptset ai-rules/code-review
pause

log "Running arm info on testing promptset..."
run_arm info promptset ai-rules/testing
pause

# === UNINSTALL ALL ===
log "=== UNINSTALL ALL ==="

log "Uninstalling all rulesets..."
run_arm uninstall ruleset ai-rules/cursor-rules
run_arm uninstall ruleset ai-rules/amazonq-rules
run_arm uninstall ruleset ai-rules/copilot-rules
run_arm uninstall ruleset ai-rules/grug-brained-dev

log "Uninstalling all promptsets..."
run_arm uninstall promptset ai-rules/code-review
run_arm uninstall promptset ai-rules/testing

log "Showing empty list..."
run_arm list
pause

# === INSTALL FROM MAIN BRANCH ===
log "=== INSTALL FROM MAIN BRANCH ==="

log "Installing cursor ruleset from main branch..."
run_arm install ruleset ai-rules/cursor-rules@main --include "rules/cursor/*.mdc" cursor-rules

log "Showing info for main branch install..."
run_arm info ruleset ai-rules/cursor-rules
pause

# === OUTDATED CHECK ===
log "=== OUTDATED CHECK ==="

log "Checking for outdated rulesets..."
run_arm outdated
pause

# === RULESET CONFIG UPDATES ===
log "=== RULESET CONFIG UPDATES ==="

log "Changing cursor-rules priority to 200..."
run_arm set ruleset ai-rules/cursor-rules priority 200

log "Showing updated priority..."
run_arm info ruleset ai-rules/cursor-rules
pause

log "Changing cursor-rules version constraint to 1.0..."
run_arm set ruleset ai-rules/cursor-rules version 1.0

log "Showing updated version constraint..."
run_arm info ruleset ai-rules/cursor-rules
pause

log "Adding q sink to cursor-rules..."
run_arm set ruleset ai-rules/cursor-rules sinks cursor-rules,q-rules

log "Showing updated sinks..."
run_arm info ruleset ai-rules/cursor-rules
pause

# === VERSION CONSTRAINT DEMOS ===
log "=== VERSION CONSTRAINT DEMOS ==="

log "Installing cursor ruleset with major version 1 (should resolve to 1.1.0)..."
run_arm install ruleset ai-rules/cursor-rules@1 --include "rules/cursor/*.mdc" cursor-rules

log "Showing info (should show 1.1.0)..."
run_arm info ruleset ai-rules/cursor-rules
pause

log "Installing cursor ruleset with minor version 1.0 (should resolve to 1.0.1)..."
run_arm install ruleset ai-rules/cursor-rules@1.0 --include "rules/cursor/*.mdc" cursor-rules

log "Showing info (should show 1.0.1)..."
run_arm info ruleset ai-rules/cursor-rules
pause

log "Installing cursor ruleset with patch version 1.0.0 (should resolve to 1.0.0)..."
run_arm install ruleset ai-rules/cursor-rules@1.0.0 --include "rules/cursor/*.mdc" cursor-rules

log "Showing info (should show 1.0.0)..."
run_arm info ruleset ai-rules/cursor-rules
pause

# === SINK REMOVAL PROTECTION ===
log "=== SINK REMOVAL PROTECTION ==="

log "Attempting to remove cursor sink (should fail because ruleset is installed)..."
if run_arm remove sink cursor-rules 2>&1; then
    error "Sink removal should have failed!"
else
    success "Sink removal correctly blocked due to active ruleset"
fi
pause

# === CLEANUP ===
log "=== CLEANUP ==="

log "Removing cursor ruleset..."
run_arm uninstall ruleset ai-rules/cursor-rules

log "Now removing cursor sink (should succeed)..."
run_arm remove sink cursor-rules

success "Cleanup complete"
pause

# === SUMMARY ===
log "=== WORKFLOW COMPLETE ==="
success "Comprehensive Git workflow completed successfully!"
echo ""
echo "This workflow demonstrated:"
echo "• Sandbox setup and binary building"
echo "• Basic help and version commands"
echo "• Registry configuration"
echo "• Sink configuration (hierarchical and flat layouts)"
echo "• Installing rulesets and promptsets to specific sinks"
echo "• Installing promptsets to multiple sinks (cursor and Amazon Q)"
echo "• Listing and getting info about resources (unified and resource-specific)"
echo "• Uninstalling rulesets and promptsets"
echo "• Installing from branches"
echo "• Checking for outdated resources"
echo "• Resource configuration updates (priority, version, sinks)"
echo "• Version constraint resolution (major, minor, patch)"
echo "• Sink removal protection"
echo "• Clean teardown"
echo "• New resource manager command structure with ruleset/promptset subcommands"
