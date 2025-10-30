#!/bin/bash

set -e

# Parse command line arguments
INTERACTIVE=true
SHOW_DEBUG=false
while [[ $# -gt 0 ]]; do
    case $1 in
        --non-interactive|-n)
            INTERACTIVE=false
            shift
            ;;
        --show-debug|-d)
            SHOW_DEBUG=true
            shift
            ;;
        *)
            echo "Unknown option: $1"
            echo "Usage: $0 [--non-interactive|-n] [--show-debug|-d]"
            exit 1
            ;;
    esac
done

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

log() { echo -e "${BLUE}[INFO]${NC} $1"; }
success() { echo -e "${GREEN}✓${NC} $1"; }
error() { echo -e "${RED}✗${NC} $1"; }

run_arm() {
    echo -e "${BLUE}$ ./arm $*${NC}"
    ./arm "$@"
}

show_debug() {
    if [ "$SHOW_DEBUG" = false ]; then
        return
    fi
    
    echo ""
    echo -e "${YELLOW}=== DEBUG OUTPUT ===${NC}"
    
    # Directory tree
    echo -e "${YELLOW}--- Directory Tree ---${NC}"
    if command -v tree &> /dev/null; then
        tree -a -I '.git' . || find . -not -path './.git/*' | sort
    else
        find . -not -path './.git/*' | sort
    fi
    echo ""
    
    # Manifest file
    if [ -f "arm.json" ]; then
        echo -e "${YELLOW}--- arm.json (Manifest) ---${NC}"
        cat arm.json
        echo ""
    fi
    
    # Lock file
    if [ -f "arm-lock.json" ]; then
        echo -e "${YELLOW}--- arm-lock.json (Lock File) ---${NC}"
        cat arm-lock.json
        echo ""
    fi
    
    # Sink index files
    for index_file in $(find . -name "arm-index.json" -o -name "arm_index.*" 2>/dev/null); do
        echo -e "${YELLOW}--- $index_file ---${NC}"
        cat "$index_file"
        echo ""
    done
    
    # Cache directory
    if [ -d "$HOME/.arm/cache" ]; then
        echo -e "${YELLOW}--- Cache Directory Tree ---${NC}"
        if command -v tree &> /dev/null; then
            tree -a "$HOME/.arm/cache" || find "$HOME/.arm/cache" | sort
        else
            find "$HOME/.arm/cache" | sort
        fi
        echo ""
        
        # Cache index files
        for cache_index in $(find "$HOME/.arm/cache" -name "*index*.json" 2>/dev/null); do
            echo -e "${YELLOW}--- $cache_index ---${NC}"
            cat "$cache_index"
            echo ""
        done
    fi
    
    echo -e "${YELLOW}=== END DEBUG OUTPUT ===${NC}"
    echo ""
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Load environment if available
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
fi

# Try to use CLOUDSMITH_API_KEY from environment if CLOUDSMITH_TOKEN not set
if [ -z "$CLOUDSMITH_TOKEN" ] && [ -n "$CLOUDSMITH_API_KEY" ]; then
    CLOUDSMITH_TOKEN="$CLOUDSMITH_API_KEY"
fi

# Load configuration from .env
CLOUDSMITH_URL=${CLOUDSMITH_URL:-"https://api.cloudsmith.io"}
CLOUDSMITH_OWNER=${CLOUDSMITH_OWNER:-""}
CLOUDSMITH_REPOSITORY=${CLOUDSMITH_REPOSITORY:-""}
CLOUDSMITH_TOKEN=${CLOUDSMITH_TOKEN:-""}
RULESET_NAME=${RULESET_NAME:-"ai-rules"}
INCLUDE_PATTERNS=${INCLUDE_PATTERNS:-""}
SINKS=${SINKS:-"cursor,q"}

log "=== Simple Cloudsmith Workflow ==="
log "Cloudsmith URL: $CLOUDSMITH_URL"
log "Owner: $CLOUDSMITH_OWNER"
log "Repository: $CLOUDSMITH_REPOSITORY"
log "Ruleset: $RULESET_NAME"
log "Include: $INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-cloudsmith-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Install configured ruleset
log "Installing $RULESET_NAME..."
if [ -n "$INCLUDE_PATTERNS" ]; then
    run_arm install ruleset cloudsmith-registry/$RULESET_NAME --include "$INCLUDE_PATTERNS" cursor-rules q-rules
else
    run_arm install ruleset cloudsmith-registry/$RULESET_NAME cursor-rules q-rules
fi

show_debug

success "Setup complete! Try these commands:"
echo ""
echo "Basic commands:"
echo "  ./arm list                    # Show all installed resources"
echo "  ./arm list ruleset            # Show installed rulesets only"
echo "  ./arm list promptset          # Show installed promptsets only"
echo "  ./arm info                    # Show detailed info for all resources"
echo "  ./arm info ruleset            # Show detailed info for rulesets"
echo "  ./arm outdated                # Check for updates"
echo ""
echo "Management commands:"
echo "  ./arm uninstall ruleset cloudsmith-registry/$RULESET_NAME"
echo "  ./arm update                  # Update all resources"
echo "  ./arm update ruleset          # Update rulesets only"
echo ""
echo "Configuration commands:"
echo "  ./arm list registry           # Show configured registries"
echo "  ./arm list sink               # Show configured sinks"
echo "  ./arm set ruleset cloudsmith-registry/$RULESET_NAME priority 200"
echo ""
echo "Example promptset commands:"
echo "  ./arm install promptset cloudsmith-registry/code-review-promptset cursor-rules"
echo "  ./arm list promptset"
echo "  ./arm uninstall promptset cloudsmith-registry/code-review-promptset"
echo ""
