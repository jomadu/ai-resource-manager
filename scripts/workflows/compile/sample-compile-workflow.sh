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
pause() { 
    if [ "$INTERACTIVE" = true ]; then
        echo ""
        read -p "Press Enter to continue..."
        echo ""
    fi
}

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

log "=== Compile Workflow ==="

# Setup sandbox
log "Setting up sandbox..."
"$SCRIPT_DIR/init-compile-sandbox.sh"
cd "$SCRIPT_DIR/sandbox"
pause

# Basic compile examples
log "Compiling rulesets to Cursor format..."
run_arm compile --tool cursor example-rulesets/*.yml ./cursor-output
show_debug
pause

log "Compiling rulesets to multiple targets..."
run_arm compile --tool cursor example-rulesets/clean-code.yml ./multi-output/cursor
run_arm compile --tool amazonq example-rulesets/clean-code.yml ./multi-output/amazonq
run_arm compile --tool copilot example-rulesets/clean-code.yml ./multi-output/copilot
show_debug
pause

log "Compiling with validation only..."
run_arm compile --validate-only example-rulesets/*.yml
pause

log "Compiling promptsets..."
run_arm compile --tool cursor example-promptsets/*.yml ./promptset-output
show_debug
pause

log "Demonstrating resource-specific compilation..."
run_arm compile --tool cursor example-rulesets/clean-code.yml ./ruleset-specific-output
run_arm compile --tool cursor example-promptsets/code-review.yml ./promptset-specific-output
show_debug
pause

success "Compile workflow complete! Check outputs:"
echo ""
echo "Generated files:"
echo "  cursor-output/           - Cursor format (.mdc)"
echo "  multi-output/            - Multiple formats"
echo "  promptset-output/        - Promptset compilation"
echo "  ruleset-specific-output/ - Ruleset-specific compilation"
echo "  promptset-specific-output/ - Promptset-specific compilation"
echo ""
echo "Try more commands:"
echo "  ./arm compile --help"
echo "  ./arm compile --tool amazonq example-rulesets/*.yml ./amazonq-output"
echo "  ./arm compile --tool cursor example-promptsets/*.yml ./promptset-output"
echo ""
echo "Resource-specific compilation:"
echo "  ./arm compile --tool cursor example-rulesets/*.yml ./output"
echo "  ./arm compile --tool cursor example-promptsets/*.yml ./output"
