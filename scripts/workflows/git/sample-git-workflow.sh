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
error() { echo -e "${RED}✗${NC} $1"; }

run_arm() {
    echo -e "${BLUE}$ ./arm $*${NC}"
    ./arm "$@"
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# Load environment if available
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
fi

# Load configuration from .env
REPO_URL=${REPO_URL:-"https://github.com/jomadu/ai-rules-manager-sample-git-registry"}
RULESET_NAME=${RULESET_NAME:-"grug-brained-dev"}
INCLUDE_PATTERNS=${INCLUDE_PATTERNS:-"rulesets/grug-brained-dev.yml"}
SINKS=${SINKS:-"cursor,q"}

log "=== Simple Git Workflow ==="
log "Repository: $REPO_URL"
log "Ruleset: $RULESET_NAME"
log "Includes: $INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-git-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Configure registry and sinks
log "Configuring registry and sinks..."
run_arm config registry add sample-repo "$REPO_URL" --type git
run_arm config sink add cursor .cursor/rules --type cursor
run_arm config sink add q .amazonq/rules --type amazonq

# Install configured ruleset
log "Installing $RULESET_NAME..."
run_arm install ruleset sample-repo/$RULESET_NAME --include "$INCLUDE_PATTERNS" --sinks $SINKS

success "Setup complete! Try these commands:"
echo ""
echo "Basic commands:"
echo "  ./arm list                    # Show all installed resources"
echo "  ./arm list ruleset            # Show installed rulesets only"
echo "  ./arm list promptset          # Show installed promptsets only"
echo "  ./arm info                    # Show detailed info for all resources"
echo "  ./arm info ruleset            # Show detailed info for rulesets"
echo "  ./arm outdated                # Check for updates"
echo ""
echo "Management commands:"
echo "  ./arm uninstall ruleset sample-repo/$RULESET_NAME"
echo "  ./arm update                  # Update all resources"
echo "  ./arm update ruleset          # Update rulesets only"
echo ""
echo "Configuration commands:"
echo "  ./arm config list             # Show current config"
echo "  ./arm config ruleset set sample-repo/$RULESET_NAME priority 200"
echo ""
echo "Example promptset commands:"
echo "  ./arm install promptset sample-repo/code-review-promptset --sinks cursor"
echo "  ./arm list promptset"
echo "  ./arm uninstall promptset sample-repo/code-review-promptset"
echo ""
