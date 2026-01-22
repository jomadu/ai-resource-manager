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
  GITLAB_URL=https://gitlab.com
  GITLAB_PROJECT_ID=your-project-id
  GITLAB_TOKEN=your-access-token
  GITLAB_API_VERSION=v4"
fi

# Validate required variables
[ -z "$GITLAB_TOKEN" ] && error "GITLAB_TOKEN is required"
[ -z "$GITLAB_URL" ] && error "GITLAB_URL is required"
[ -z "$GITLAB_PROJECT_ID" ] && error "GITLAB_PROJECT_ID is required"

GITLAB_API_VERSION=${GITLAB_API_VERSION:-v4}

log "=== GitLab Sandbox Setup ==="

# Build ARM binary
log "Building ARM binary..."
cd "$PROJECT_ROOT"
make build

# Setup sandbox
log "Setting up GitLab sandbox..."
rm -rf "$SCRIPT_DIR/sandbox"
mkdir -p "$SCRIPT_DIR/sandbox"
cp ./bin/arm "$SCRIPT_DIR/sandbox/"
cd "$SCRIPT_DIR/sandbox"

# Create .armrc with authentication
cat > .armrc << EOF
[registry ${GITLAB_URL}/project/${GITLAB_PROJECT_ID}]
token = ${GITLAB_TOKEN}
EOF

# Configure registry and sinks
log "Configuring GitLab registry..."
./arm add registry gitlab --url "$GITLAB_URL" --project-id "$GITLAB_PROJECT_ID" --api-version "$GITLAB_API_VERSION" gitlab-registry
./arm add sink --tool cursor cursor-rules .cursor/rules
./arm add sink --tool amazonq q-rules .amazonq/rules

success "GitLab sandbox ready!"
echo ""
echo "Configuration:"
echo "  GitLab URL: $GITLAB_URL"
echo "  Project ID: $GITLAB_PROJECT_ID"
echo ""
echo "Try these commands:"
echo "  ./arm install ruleset gitlab-registry/ai-rules cursor-rules q-rules"
echo "  ./arm list"
