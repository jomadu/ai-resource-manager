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

# Load environment if available
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
fi

# Load configuration from .env
CLOUDSMITH_URL=${CLOUDSMITH_URL:-"https://api.cloudsmith.io"}
CLOUDSMITH_OWNER=${CLOUDSMITH_OWNER:-""}
CLOUDSMITH_REPOSITORY=${CLOUDSMITH_REPOSITORY:-""}
CLOUDSMITH_TOKEN=${CLOUDSMITH_TOKEN:-""}
RULESET_NAME=${RULESET_NAME:-"ai-rules"}
INCLUDE_PATTERNS=${INCLUDE_PATTERNS:-""}
SINKS=${SINKS:-"cursor,q"}

log "=== Simple Cloudsmith Workflow ==="
log "Cloudsmith URL: $CLOUDSMITH_URL"
log "Owner: $CLOUDSMITH_OWNER"
log "Repository: $CLOUDSMITH_REPOSITORY"
log "Ruleset: $RULESET_NAME"
log "Include: $INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-cloudsmith-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Install configured ruleset
log "Installing $RULESET_NAME..."
if [ -n "$INCLUDE_PATTERNS" ]; then
    run_arm install ruleset cloudsmith-registry/$RULESET_NAME --include "$INCLUDE_PATTERNS" cursor-rules q-rules
else
    run_arm install ruleset cloudsmith-registry/$RULESET_NAME cursor-rules q-rules
fi

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
echo "  ./arm uninstall ruleset cloudsmith-registry/$RULESET_NAME"
echo "  ./arm update                  # Update all resources"
echo "  ./arm update ruleset          # Update rulesets only"
echo ""
echo "Configuration commands:"
echo "  ./arm list registry           # Show configured registries"
echo "  ./arm list sink               # Show configured sinks"
echo "  ./arm set ruleset cloudsmith-registry/$RULESET_NAME priority 200"
echo ""
echo "Example promptset commands:"
echo "  ./arm install promptset cloudsmith-registry/code-review-promptset cursor-rules"
echo "  ./arm list promptset"
echo "  ./arm uninstall promptset cloudsmith-registry/code-review-promptset"
echo ""
