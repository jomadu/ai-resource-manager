#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}✓${NC} $1"
}

error() {
    echo -e "${RED}✗${NC} $1"
    exit 1
}

# Validate required environment variables
log "Validating environment variables..."

if [ -z "$GITLAB_TOKEN" ]; then
    error "GITLAB_TOKEN environment variable is required"
fi

if [ -z "$GITLAB_URL" ]; then
    error "GITLAB_URL environment variable is required (e.g., https://gitlab.example.com)"
fi

if [ -z "$GITLAB_PROJECT_ID" ]; then
    error "GITLAB_PROJECT_ID environment variable is required"
fi

# Set default API version if not provided
GITLAB_API_VERSION=${GITLAB_API_VERSION:-v4}

success "Environment variables validated"

# === SETUP SANDBOX ===
log "=== SETUP SANDBOX ==="

# Get script directory and project root
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

log "Running setup-sandbox script..."
"$SCRIPT_DIR/setup-sandbox.sh"

log "Entering sandbox..."
cd "$PROJECT_ROOT/sandbox"

# === CREATE .ARMRC ===
log "=== CREATING .ARMRC ==="

# Extract hostname from GITLAB_URL for .armrc
GITLAB_HOST=$(echo "$GITLAB_URL" | sed -e 's|^https\?://||' -e 's|/.*||')

log "Creating .armrc with GitLab registry authentication..."
cat > .armrc << EOF
[registry ${GITLAB_HOST}/project/${GITLAB_PROJECT_ID}]
token = ${GITLAB_TOKEN}
EOF

success ".armrc created with GitLab authentication"

# === REGISTRY SETUP ===
log "=== REGISTRY SETUP ==="

log "Adding GitLab registry..."
./arm config registry add gitlab-registry "$GITLAB_URL" --type gitlab --project-id "$GITLAB_PROJECT_ID" --api-version "$GITLAB_API_VERSION"

log "Showing configuration..."
./arm config list

success "GitLab registry setup complete!"

# === SINK SETUP ===
log "=== SINK SETUP ==="

log "Setting up cursor sink..."
./arm config sink add cursor .cursor/rules --type cursor

log "Setting up Amazon Q sink..."
./arm config sink add q .amazonq/rules --type amazonq

log "Showing updated configuration..."
./arm config list

success "Sink setup complete!"

# === SUMMARY ===
log "=== SETUP COMPLETE ==="
success "GitLab registry workflow setup completed successfully!"
echo ""
echo "Configuration summary:"
echo "• GitLab URL: $GITLAB_URL"
echo "• Project ID: $GITLAB_PROJECT_ID"
echo "• API Version: $GITLAB_API_VERSION"
echo "• Registry name: gitlab-registry"
echo ""
echo "Next steps:"
echo "• Install rulesets: ./arm install gitlab-registry/ai-rules --sinks cursor,q"
echo ""
echo "Example commands:"
echo "• ./arm install gitlab-registry/ai-rules --sinks cursor"
echo "• ./arm install gitlab-registry/ai-rules --sinks q"
echo "• ./arm install gitlab-registry/ai-rules --sinks cursor,q"
