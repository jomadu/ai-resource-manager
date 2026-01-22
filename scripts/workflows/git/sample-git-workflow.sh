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
    
    # Storage directory
    if [ -d "$HOME/.arm/storage" ]; then
        echo -e "${YELLOW}--- Storage Directory Tree ---${NC}"
        if command -v tree &> /dev/null; then
            tree -a "$HOME/.arm/storage" || find "$HOME/.arm/storage" | sort
        else
            find "$HOME/.arm/storage" | sort
        fi
        echo ""
        
        # Storage index files
        for storage_index in $(find "$HOME/.arm/storage" -name "*index*.json" 2>/dev/null); do
            echo -e "${YELLOW}--- $storage_index ---${NC}"
            cat "$storage_index"
            echo ""
        done
    fi
    
    echo -e "${YELLOW}=== END DEBUG OUTPUT ===${NC}"
    echo ""
}

# Get script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../../.." && pwd)"

# Load environment if available
if [ -f "$SCRIPT_DIR/.env" ]; then
    source "$SCRIPT_DIR/.env"
fi

# Load configuration from .env
REPO_URL=${REPO_URL:-"https://github.com/jomadu/ai-rules-manager-sample-git-registry"}
RULESET_NAME=${RULESET_NAME:-"grug-brained-dev"}
PROMPTSET_NAME=${PROMPTSET_NAME:-"code-review-promptset"}
RULESET_INCLUDE_PATTERNS=${RULESET_INCLUDE_PATTERNS:-"rulesets/grug-brained-dev.yml"}
PROMPTSET_INCLUDE_PATTERNS=${PROMPTSET_INCLUDE_PATTERNS:-"promptsets/code-review.yml"}
SINKS=${SINKS:-"cursor,q"}

log "=== Simple Git Workflow ==="
log "Repository: $REPO_URL"
log "Ruleset: $RULESET_NAME"
log "Promptset: $PROMPTSET_NAME"
log "Ruleset Include: $RULESET_INCLUDE_PATTERNS"
log "Promptset Include: $PROMPTSET_INCLUDE_PATTERNS"
log "Sinks: $SINKS"

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-git-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"

# Configure registry and sinks
log "Configuring registry and sinks..."
run_arm add registry git --url "$REPO_URL" sample-repo
run_arm add sink --tool cursor cursor-rules .cursor/rules
run_arm add sink --tool cursor cursor-commands .cursor/commands
run_arm add sink --tool amazonq q-rules .amazonq/rules
run_arm add sink --tool amazonq q-prompts .amazonq/prompts

# Install configured ruleset
log "Installing $RULESET_NAME..."
run_arm install ruleset sample-repo/$RULESET_NAME --include "$RULESET_INCLUDE_PATTERNS" cursor-rules q-rules

# Install configured promptset
log "Installing $PROMPTSET_NAME..."
run_arm install promptset sample-repo/$PROMPTSET_NAME --include "$PROMPTSET_INCLUDE_PATTERNS" cursor-commands q-prompts

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
echo "  ./arm uninstall ruleset sample-repo/$RULESET_NAME"
echo "  ./arm uninstall promptset sample-repo/$PROMPTSET_NAME"
echo "  ./arm update                                    # Update all resources"
echo "  ./arm update sample-repo/$RULESET_NAME          # Update specific ruleset"
echo "  ./arm update sample-repo/$PROMPTSET_NAME        # Update specific promptset"
echo ""
echo "Configuration commands:"
echo "  ./arm list registry           # Show configured registries"
echo "  ./arm list sink               # Show configured sinks"
echo "  ./arm set ruleset sample-repo/$RULESET_NAME priority 200"
echo ""
echo "Example promptset commands:"
echo "  ./arm install promptset sample-repo/$PROMPTSET_NAME --include '$PROMPTSET_INCLUDE_PATTERNS' cursor-commands q-prompts"
echo "  ./arm list promptset"
echo "  ./arm uninstall promptset sample-repo/$PROMPTSET_NAME"
echo ""
