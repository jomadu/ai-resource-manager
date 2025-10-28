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
error() { echo -e "${RED}✗${NC} $1"; exit 1; }

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# Load environment
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
else
    error "Environment file not found. Create $SCRIPT_DIR/.env with:
  CLOUDSMITH_URL=https://api.cloudsmith.io
  CLOUDSMITH_OWNER=your-owner-name
  CLOUDSMITH_REPOSITORY=your-repo-name
  CLOUDSMITH_TOKEN=your-api-token"
fi

# Validate required variables
[ -z "$CLOUDSMITH_TOKEN" ] && error "CLOUDSMITH_TOKEN is required"
[ -z "$CLOUDSMITH_OWNER" ] && error "CLOUDSMITH_OWNER is required"
[ -z "$CLOUDSMITH_REPOSITORY" ] && error "CLOUDSMITH_REPOSITORY is required"

CLOUDSMITH_URL=${CLOUDSMITH_URL:-"https://api.cloudsmith.io"}

log "=== Cloudsmith Sandbox Setup ==="

# Build ARM binary
log "Building ARM binary..."
cd "$PROJECT_ROOT"
make build

# Setup sandbox
log "Setting up Cloudsmith sandbox..."
rm -rf "$SCRIPT_DIR/sandbox"
mkdir -p "$SCRIPT_DIR/sandbox"
cp ./bin/arm "$SCRIPT_DIR/sandbox/"
cd "$SCRIPT_DIR/sandbox"

# Create .armrc with authentication
cat > .armrc << EOF
[registry ${CLOUDSMITH_URL}/${CLOUDSMITH_OWNER}/${CLOUDSMITH_REPOSITORY}]
token = ${CLOUDSMITH_TOKEN}
EOF

# Configure registry and sinks
log "Configuring Cloudsmith registry..."
./arm add registry cloudsmith --url "$CLOUDSMITH_URL" --owner "$CLOUDSMITH_OWNER" --repo "$CLOUDSMITH_REPOSITORY" cloudsmith-registry
./arm add sink --type cursor cursor-rules .cursor/rules
./arm add sink --type amazonq q-rules .amazonq/rules

success "Cloudsmith sandbox ready!"
echo ""
echo "Configuration:"
echo "  Cloudsmith URL: $CLOUDSMITH_URL"
echo "  Owner: $CLOUDSMITH_OWNER"
echo "  Repository: $CLOUDSMITH_REPOSITORY"
echo ""
echo "Try these commands:"
echo "  ./arm install ruleset cloudsmith-registry/ai-rules cursor-rules q-rules"
echo "  ./arm list"
