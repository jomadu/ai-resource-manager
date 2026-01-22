#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; exit 1; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load environment
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
fi

# Validate required variables
if [ -z "$GIT_REGISTRY_URL" ]; then
    error "GIT_REGISTRY_URL is required. Create $SCRIPT_DIR/.env with:
  GIT_REGISTRY_URL=https://github.com/user/repo
  GIT_REGISTRY_NAME=my-registry"
fi

[ -z "$GIT_REGISTRY_NAME" ] && error "GIT_REGISTRY_NAME is required"

log "=== Git Registry Sandbox Setup ==="
log "Repository URL: $GIT_REGISTRY_URL"
log "Registry name: $GIT_REGISTRY_NAME"

# Initialize sandbox
log "Initializing sandbox..."
"$SCRIPT_DIR/init-git-sandbox.sh"

# Add registry
cd "$SCRIPT_DIR/sandbox"
log "Adding registry..."
./arm add registry git --url "$GIT_REGISTRY_URL" "$GIT_REGISTRY_NAME"

success "Git registry sandbox ready!"
echo ""
echo "Registry: $GIT_REGISTRY_NAME"
echo "URL: $GIT_REGISTRY_URL"
echo ""
echo "Next steps:"
echo "  cd $SCRIPT_DIR/sandbox"
echo "  ./arm add sink --tool cursor cursor-rules .cursor/rules"
echo "  ./arm install ruleset $GIT_REGISTRY_NAME/package-name cursor-rules"
