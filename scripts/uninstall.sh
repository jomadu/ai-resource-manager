#!/bin/bash
set -e

# ARM Uninstall Script
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="arm"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

uninstall_binary() {
    local binary_path="${INSTALL_DIR}/${BINARY_NAME}"

    if [ ! -f "$binary_path" ]; then
        log_warn "ARM binary not found at ${binary_path}"
        return 0
    fi

    log_info "Removing ARM binary from ${binary_path}..."

    if [ -w "$INSTALL_DIR" ]; then
        rm -f "$binary_path"
    else
        sudo rm -f "$binary_path"
    fi

    log_info "ARM binary removed successfully"
}

main() {
    log_info "Uninstalling ARM (AI Rules Manager)..."

    uninstall_binary

    if ! command -v "$BINARY_NAME" > /dev/null 2>&1; then
        log_info "ARM has been completely uninstalled"
    else
        log_warn "ARM may still be available in PATH from another location"
    fi
}

main "$@"
