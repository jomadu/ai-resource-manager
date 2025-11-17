#!/bin/bash
set -e

# ARM Uninstall Script
LINUX_INSTALL_DIR=/usr/local/bin
WINDOWS_INSTALL_DIR=~/AppData/Local/Programs/arm
BINARY_NAME="arm"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

log_info() { echo -e "${GREEN}[INFO]${NC} $1"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

detect_platform() {
    local os arch
    case "$(uname -s)" in
        Linux*)     os="linux" ;;
        Darwin*)    os="darwin" ;;
        CYGWIN*|MINGW*|MSYS*) os="windows" ;;
        *)          log_error "Unsupported OS: $(uname -s)"; exit 1 ;;
    esac

    case "$(uname -m)" in
        x86_64|amd64)   arch="amd64" ;;
        arm64|aarch64)  arch="arm64" ;;
        *)              log_error "Unsupported architecture: $(uname -m)"; exit 1 ;;
    esac

    echo "${os}-${arch}"
}

uninstall_binary_windows() {
    local binary_path="${WINDOWS_INSTALL_DIR}/${BINARY_NAME}.exe"

    if [ ! -f "$binary_path" ]; then
        log_warn "ARM binary not found at ${binary_path}"
        return 0
    fi

    log_info "Removing ${WINDOWS_INSTALL_DIR}..."

    if [ -w "$WINDOWS_INSTALL_DIR" ]; then
        rm -rf "$WINDOWS_INSTALL_DIR"
    else
        log_error "Failed to remove ${WINDOWS_INSTALL_DIR}. Please run this script with appropriate permissions or remove manually."
        return 1
    fi

    log_info "ARM removed successfully"
}

uninstall_binary() {
    local binary_path="${LINUX_INSTALL_DIR}/${BINARY_NAME}"

    if [ ! -f "$binary_path" ]; then
        log_warn "ARM binary not found at ${binary_path}"
        return 0
    fi

    log_info "Removing ARM binary from ${binary_path}..."

    if [ -w "$LINUX_INSTALL_DIR" ]; then
        rm -f "$binary_path"
    else
        sudo rm -f "$binary_path"
    fi

    log_info "ARM binary removed successfully"
}

main() {
    log_info "Uninstalling ARM (AI Rules Manager)..."

    local platform=$(detect_platform)

    if [[ "$platform" == windows* ]]; then
        uninstall_binary_windows
    else
        uninstall_binary
    fi

    if ! command -v "$BINARY_NAME" > /dev/null 2>&1; then
        log_info "ARM has been completely uninstalled"
    else
        log_warn "ARM may still be available in PATH from another location"
    fi
}

main "$@"
