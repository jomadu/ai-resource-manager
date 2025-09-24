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
GITLAB_URL=${GITLAB_URL:-""}
GITLAB_PROJECT_ID=${GITLAB_PROJECT_ID:-""}
GITLAB_TOKEN=${GITLAB_TOKEN:-""}
GITLAB_API_VERSION=${GITLAB_API_VERSION:-"v4"}
RULESET_NAME=${RULESET_NAME:-"ai-rules"}
INCLUDE_PATTERNS=${INCLUDE_PATTERNS:-""}
SINKS=${SINKS:-"cursor,q"}

log "=== Simple GitLab Workflow ==="
log "GitLab URL: $GITLAB_URL"
log "Project ID: $GITLAB_PROJECT_ID"
log "Ruleset: $RULESET_NAME"
log "Includes: $INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-gitlab-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Install configured ruleset
log "Installing $RULESET_NAME..."
if [ -n "$INCLUDE_PATTERNS" ]; then
    ./arm install gitlab-registry/$RULESET_NAME --include "$INCLUDE_PATTERNS" --sinks $SINKS
else
    ./arm install gitlab-registry/$RULESET_NAME --sinks $SINKS
fi

success "Setup complete! Try these commands:"
echo ""
echo "Basic commands:"
echo "  ./arm list                    # Show installed rulesets"
echo "  ./arm info                    # Show detailed info"
echo "  ./arm outdated                # Check for updates"
echo ""
echo "Management commands:"
echo "  ./arm uninstall gitlab-registry/$RULESET_NAME"
echo "  ./arm update                  # Update all rulesets"
echo ""
echo "Configuration commands:"
echo "  ./arm config list             # Show current config"
echo "  ./arm config ruleset update gitlab-registry/$RULESET_NAME priority 200"
echo ""
