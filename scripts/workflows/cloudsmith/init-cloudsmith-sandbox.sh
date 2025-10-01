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
    error "Environment file not found. Copy .env.cloudsmith.example to .env and configure it."
fi

# Validate required variables
[ -z "$CLOUDSMITH_TOKEN" ] && error "CLOUDSMITH_TOKEN is required"
[ -z "$CLOUDSMITH_OWNER" ] && error "CLOUDSMITH_OWNER is required"
[ -z "$CLOUDSMITH_REPOSITORY" ] && error "CLOUDSMITH_REPOSITORY is required"

CLOUDSMITH_URL=${CLOUDSMITH_URL:-"https://api.cloudsmith.io"}
CLOUDSMITH_APP_URL="https://app.cloudsmith.com/${CLOUDSMITH_OWNER}/${CLOUDSMITH_REPOSITORY}"

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
[registry ${CLOUDSMITH_APP_URL}]
token = ${CLOUDSMITH_TOKEN}
EOF

# Configure registry and sinks
log "Configuring Cloudsmith registry..."
./arm config registry add cloudsmith-registry "$CLOUDSMITH_APP_URL" --type cloudsmith
./arm config sink add cursor .cursor/rules --type cursor
./arm config sink add q .amazonq/rules --type amazonq

success "Cloudsmith sandbox ready!"
echo ""
echo "Configuration:"
echo "  Cloudsmith URL: $CLOUDSMITH_URL"
echo "  Owner: $CLOUDSMITH_OWNER"
echo "  Repository: $CLOUDSMITH_REPOSITORY"
echo ""
echo "Try these commands:"
echo "  ./arm install cloudsmith-registry/ai-rules --sinks cursor,q"
echo "  ./arm list"
