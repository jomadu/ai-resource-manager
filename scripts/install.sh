#!/bin/bash
set -e

# ARM Installation Script
REPO="jomadu/ai-resource-manager"
BINARY_NAME="arm"
LINUX_INSTALL_DIR=/usr/local/bin
WINDOWS_INSTALL_DIR=~/AppData/Local/Programs/arm

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

get_latest_version() {
    curl -sL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' | sed 's/^v//'
}

install_binary_windows() {
    local platform="$1" version="$2"
    local binary_name="${BINARY_NAME}-${platform}.exe"
    local download_url="https://github.com/${REPO}/releases/download/v${version}/${binary_name}.tar.gz"
    local final_binary_name="${BINARY_NAME}.exe"

    log_info "Downloading ARM v${version} for ${platform}..."

    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    # Download and extract
    curl -sL "$download_url" -o "${binary_name}.tar.gz"
    
    # Extract using tar (works on Git Bash/MSYS on Windows)
    if ! tar -xzf "${binary_name}.tar.gz"; then
        log_error "Extraction failed."
        rm -rf "$temp_dir"
        exit 1
    fi

    if [ ! -d "$WINDOWS_INSTALL_DIR" ]; then
        mkdir -p "$WINDOWS_INSTALL_DIR"
    fi

    if [ -w "$WINDOWS_INSTALL_DIR" ]; then
        mv "${binary_name}" "${WINDOWS_INSTALL_DIR}/${final_binary_name}"
    else
        log_error "Failed to install - install directory is not writable. Please run this script with appropriate permissions or install manually."
        rm -rf "$temp_dir"
        exit 1
    fi

    rm -rf "$temp_dir"
    log_info "ARM installed to ${WINDOWS_INSTALL_DIR}/${final_binary_name}"    
}

install_binary() {
    local platform="$1" version="$2"
    local binary_name="${BINARY_NAME}-${platform}"
    local download_url="https://github.com/${REPO}/releases/download/v${version}/${binary_name}.tar.gz"

    log_info "Downloading ARM v${version} for ${platform}..."

    local temp_dir=$(mktemp -d)
    cd "$temp_dir"

    curl -sL "$download_url" | tar -xz
    chmod +x "${binary_name}"

    if [ -w "$LINUX_INSTALL_DIR" ]; then
        mv "${binary_name}" "${LINUX_INSTALL_DIR}/${BINARY_NAME}"
    else
        sudo mv "${binary_name}" "${LINUX_INSTALL_DIR}/${BINARY_NAME}"
    fi

    rm -rf "$temp_dir"
    log_info "ARM installed to ${LINUX_INSTALL_DIR}/${BINARY_NAME}"    
}

main() {
    local requested_version="$1"

    log_info "Installing ARM (AI Rules Manager)..."

    local platform=$(detect_platform)
    local version

    if [ -n "$requested_version" ]; then
        version="${requested_version#v}"  # Remove 'v' prefix if present
        log_info "Installing specific version: v${version}"
    else
        version=$(get_latest_version)
        log_info "Installing latest version: v${version}"
    fi

    # Install based on OS
    if [[ "$platform" == windows* ]]; then
        install_binary_windows "$platform" "$version"
        local check_name="${BINARY_NAME}.exe"
        
        if command -v "$check_name" > /dev/null 2>&1 || command -v "$BINARY_NAME" > /dev/null 2>&1; then
            log_info "ARM is ready! Run '${BINARY_NAME} help' to get started"
        else
            if [ $(basename $SHELL) == "bash" ]; then
                log_info "Next step: Add ${WINDOWS_INSTALL_DIR} to your PATH:"
                log_info "   echo 'export PATH=\"${WINDOWS_INSTALL_DIR}:\$PATH\"' >> ~/.bashrc"
                log_info "   source ~/.bashrc"
                log_info ""
            else
                log_warn "ARM may not be in your PATH. Add it with:"
                log_warn "  export PATH=\"${WINDOWS_INSTALL_DIR}:\$PATH\""
                log_warn "Or add it permanently to your shell profile."
            fi

        fi
    else
        install_binary "$platform" "$version"
        
        if command -v "$BINARY_NAME" > /dev/null 2>&1; then
            log_info "ARM is ready! Run '${BINARY_NAME} help' to get started"
        else
            log_warn "ARM may not be in your PATH. Add ${LINUX_INSTALL_DIR} to your PATH, or run ARM directly from ${LINUX_INSTALL_DIR}/${BINARY_NAME}"
        fi
    fi
}

main "$@"
