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
log "Includes: $INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-cloudsmith-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Install configured ruleset
log "Installing $RULESET_NAME..."
if [ -n "$INCLUDE_PATTERNS" ]; then
    ./arm install cloudsmith-registry/$RULESET_NAME --include "$INCLUDE_PATTERNS" --sinks $SINKS
else
    ./arm install cloudsmith-registry/$RULESET_NAME --sinks $SINKS
fi

success "Setup complete! Try these commands:"
echo ""
echo "Basic commands:"
echo "  ./arm list                    # Show installed rulesets"
echo "  ./arm info                    # Show detailed info"
echo "  ./arm outdated                # Check for updates"
echo ""
echo "Management commands:"
echo "  ./arm uninstall cloudsmith-registry/$RULESET_NAME"
echo "  ./arm update                  # Update all rulesets"
echo ""
echo "Configuration commands:"
echo "  ./arm config list             # Show current config"
echo "  ./arm config ruleset update cloudsmith-registry/$RULESET_NAME priority 200"
echo ""
