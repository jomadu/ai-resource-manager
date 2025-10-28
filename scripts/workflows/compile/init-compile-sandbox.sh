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

log "=== Compile Sandbox Setup ==="

# Build ARM binary
log "Building ARM binary..."
cd "$PROJECT_ROOT"
make build

# Setup sandbox
log "Setting up compile sandbox..."
rm -rf "$SCRIPT_DIR/sandbox"
mkdir -p "$SCRIPT_DIR/sandbox"
cp ./bin/arm "$SCRIPT_DIR/sandbox/"

# Copy example rulesets and promptsets
log "Copying example rulesets..."
cp -r "$SCRIPT_DIR/example-rulesets" "$SCRIPT_DIR/sandbox/" 2>/dev/null || true

log "Copying example promptsets..."
cp -r "$SCRIPT_DIR/example-promptsets" "$SCRIPT_DIR/sandbox/" 2>/dev/null || true

cd "$SCRIPT_DIR/sandbox"

success "Compile sandbox ready!"
echo ""
echo "Try these commands:"
echo "  ./arm compile --target cursor example-rulesets/*.yml ./output"
echo "  ./arm compile --target cursor example-promptsets/*.yml ./output"
echo "  ./arm compile --target cursor example-rulesets/clean-code.yml ./output"
echo "  ./arm compile --target cursor example-promptsets/code-review.yml ./output"
echo "  ./arm compile --help"
