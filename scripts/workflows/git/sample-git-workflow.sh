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

run_arm() {
    echo -e "${BLUE}$ ./arm $*${NC}"
    ./arm "$@"
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

log "=== Sample Git Workflow ==="

# Setup sandbox with local registry
log "Setting up sandbox with local git registry..."
"$SCRIPT_DIR/init-sample-local-git-registry-sandbox.sh"

cd "$SCRIPT_DIR/sandbox"
REGISTRY_PATH="$SCRIPT_DIR/sandbox/local-registry"
REPO_URL="file://$REGISTRY_PATH"

# Configure sinks
log "Configuring sinks..."
run_arm add sink --tool cursor cursor-rules .cursor/rules
run_arm add sink --tool cursor cursor-commands .cursor/commands
run_arm add sink --tool amazonq q-rules .amazonq/rules
run_arm add sink --tool amazonq q-prompts .amazonq/prompts

# Install ruleset
log "Installing grug-brained-dev ruleset..."
run_arm install ruleset local-test/grug-rules --include "rulesets/grug-brained-dev.yml" cursor-rules q-rules

# Install promptset
log "Installing code-review promptset..."
run_arm install promptset local-test/code-review --include "promptsets/code-review.yml" cursor-commands q-prompts

# Show results
log "Listing installed resources..."
run_arm list

success "Workflow complete! Try these commands:"
echo ""
echo "Basic commands:"
echo "  cd $SCRIPT_DIR/sandbox"
echo "  ./arm list                    # Show all installed resources"
echo "  ./arm info                    # Show detailed info"
echo "  ./arm outdated                # Check for updates"
echo ""
echo "Version testing:"
echo "  ./arm install ruleset local-test/grug-rules@1.0.0 cursor-rules"
echo "  ./arm install ruleset local-test/grug-rules@1 cursor-rules"
echo "  ./arm install ruleset local-test/grug-rules@2.0.0 cursor-rules"
echo ""
echo "Branch testing:"
echo "  ./arm add registry git --url $REPO_URL --branches main,develop local-test-branches"
echo "  ./arm install ruleset local-test-branches/grug-rules@develop cursor-rules"
echo ""
echo "Pattern testing:"
echo "  ./arm install ruleset local-test/legacy --include 'rules/**/*.mdc' cursor-rules"
echo "  ./arm install ruleset local-test/yaml-only --include '**/*.yml' cursor-rules"
echo ""
