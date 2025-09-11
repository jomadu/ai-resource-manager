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
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1"
}

usage() {
    echo "Usage: $0"
    echo ""
    echo "Sets up a clean sandbox environment for ARM testing."
    echo ""
    echo "This script:"
    echo "  1. Removes existing sandbox directory"
    echo "  2. Clears ARM cache directory"
    echo "  3. Builds ARM binary"
    echo "  4. Creates new sandbox with ARM binary"
    echo ""
    echo "Run from the project root directory."
}

main() {
    # Check for help
    if [[ "$1" == "--help" || "$1" == "-h" ]]; then
        usage
        exit 0
    fi

    # Change to project root directory
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
    cd "$PROJECT_ROOT"

    # Check if we're in the right directory
    if [[ ! -f "Makefile" ]]; then
        error "Makefile not found in project root directory."
        exit 1
    fi

    log "=== SETUP SANDBOX ==="

    log "Nuking sandbox directory..."
    rm -rf sandbox/

    log "Nuking ~/.arm/cache directory..."
    rm -rf ~/.arm/cache

    log "Building ARM binary from project root..."
    make build
    success "ARM built successfully"

    log "Creating sandbox and copying binary..."
    mkdir sandbox
    cp ./bin/arm ./sandbox

    success "Sandbox setup complete!"
    echo ""
    echo "Next steps:"
    echo "  cd sandbox"
    echo "  ./arm help"
    echo ""
    echo "Or run a workflow script from the project root:"
    echo "  ./scripts/new-workflow.sh"
}

main "$@"
