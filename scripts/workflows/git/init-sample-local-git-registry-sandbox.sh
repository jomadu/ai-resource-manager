#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log "=== Sample Local Git Registry Sandbox Setup ==="

# Initialize sandbox
log "Initializing sandbox..."
"$SCRIPT_DIR/init-git-sandbox.sh"

# Create local git registry in sandbox
REGISTRY_PATH="$SCRIPT_DIR/sandbox/local-registry"
log "Creating local git registry at: $REGISTRY_PATH"
"$SCRIPT_DIR/init-local-git-registry.sh" "$REGISTRY_PATH"

# Add registry
cd "$SCRIPT_DIR/sandbox"
log "Adding registry..."
./arm add registry git --url "file://$REGISTRY_PATH" local-test

success "Sample local git registry sandbox ready!"
echo ""
echo "Registry: local-test"
echo "Registry path: $REGISTRY_PATH"
echo "Registry URL: file://$REGISTRY_PATH"
echo ""
echo "Next steps:"
echo "  cd $SCRIPT_DIR/sandbox"
echo "  ./arm add sink --tool cursor cursor-rules .cursor/rules"
echo "  ./arm install ruleset local-test/package-name cursor-rules"
