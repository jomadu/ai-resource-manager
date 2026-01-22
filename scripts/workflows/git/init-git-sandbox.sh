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

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

log "=== Git Sandbox Setup ==="

# Build ARM binary
log "Building ARM binary..."
cd "$PROJECT_ROOT"
make build

# Setup sandbox
log "Setting up Git sandbox..."
rm -rf "$SCRIPT_DIR/sandbox"
mkdir -p "$SCRIPT_DIR/sandbox"
cp ./bin/arm "$SCRIPT_DIR/sandbox/"

success "Git sandbox ready at: $SCRIPT_DIR/sandbox"
echo ""
echo "Next steps:"
echo "  cd $SCRIPT_DIR/sandbox"
echo "  ./arm help"
echo ""
echo "Try the new resource manager commands:"
echo "  ./arm install ruleset REGISTRY/RULESET SINK..."
echo "  ./arm install promptset REGISTRY/PROMPTSET SINK..."
echo "  ./arm list"
