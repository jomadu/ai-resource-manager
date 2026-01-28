# Installation Scripts

## Job to be Done
Provide automated installation and uninstallation of ARM binaries on Linux, macOS, and Windows with platform detection and PATH configuration.

## Activities
1. Detect user's operating system and architecture
2. Download appropriate ARM binary from GitHub releases
3. Extract and install binary to system PATH
4. Verify installation and provide next steps
5. Uninstall ARM binary and clean up

## Acceptance Criteria
- [x] Detect platform (Linux, macOS, Windows) and architecture (amd64, arm64)
- [x] Download latest release by default
- [x] Support installing specific version via argument
- [x] Download and extract .tar.gz archives
- [x] Install to /usr/local/bin on Linux/macOS
- [x] Install to ~/AppData/Local/Programs/arm on Windows
- [x] Use sudo only when necessary (if install dir not writable)
- [x] Verify installation by checking if binary is in PATH
- [x] Provide PATH configuration instructions if needed
- [x] Uninstall binary from system PATH
- [x] Handle errors gracefully with colored output

## Data Structures

### Platform Detection
```bash
OS=$(uname -s)
ARCH=$(uname -m)

# Mappings:
Linux   → linux
Darwin  → darwin
CYGWIN/MINGW/MSYS → windows

x86_64/amd64 → amd64
arm64/aarch64 → arm64
```

### Installation Paths
```bash
# Linux/macOS
LINUX_INSTALL_DIR=/usr/local/bin

# Windows
WINDOWS_INSTALL_DIR=~/AppData/Local/Programs/arm
```

### GitHub Release URL
```bash
https://github.com/jomadu/ai-resource-manager/releases/download/v{version}/arm-{platform}.tar.gz
```

## Algorithm

### Install Script (scripts/install.sh)

1. **Detect Platform:**
   - Run `uname -s` to get OS
   - Run `uname -m` to get architecture
   - Map to platform string (e.g., "linux-amd64", "darwin-arm64")
   - Exit with error if unsupported

2. **Determine Version:**
   - If version argument provided: use it (strip 'v' prefix)
   - Else: fetch latest version from GitHub API
     - `curl https://api.github.com/repos/jomadu/ai-resource-manager/releases/latest`
     - Parse JSON for tag_name field
     - Strip 'v' prefix

3. **Download Binary:**
   - Construct download URL
   - Create temporary directory
   - Download .tar.gz to temp dir
   - Extract archive

4. **Install Binary:**
   - **Linux/macOS:**
     - Make binary executable
     - If /usr/local/bin writable: move directly
     - Else: use sudo to move
   - **Windows:**
     - Create install directory if not exists
     - Move .exe to install directory
     - Check if directory writable, fail if not

5. **Verify Installation:**
   - Check if binary in PATH using `command -v`
   - If found: print success message
   - If not found: print PATH configuration instructions

6. **Cleanup:**
   - Remove temporary directory

### Uninstall Script (scripts/uninstall.sh)

1. **Detect Platform:**
   - Same as install script

2. **Locate Binary:**
   - **Linux/macOS:** Check /usr/local/bin/arm
   - **Windows:** Check ~/AppData/Local/Programs/arm/arm.exe

3. **Remove Binary:**
   - If writable: remove directly
   - Else: use sudo to remove

4. **Verify Removal:**
   - Check if binary still in PATH
   - Print success or error message

5. **Cleanup:**
   - Remove empty directories if applicable

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| Unsupported OS | Exit with error message |
| Unsupported architecture | Exit with error message |
| Network failure | Exit with curl error message |
| Invalid version | Exit with 404 error from GitHub |
| Install dir not writable | Prompt for sudo or fail with error |
| Binary already installed | Overwrite existing binary |
| Binary not in PATH | Provide instructions to add to PATH |
| Windows without Git Bash | May fail (requires tar command) |

## Dependencies

- curl (for downloading)
- tar (for extracting archives)
- uname (for platform detection)
- sudo (optional, for privileged installation)
- Git Bash on Windows (for tar and bash support)

## Implementation Mapping

**Source files:**
- `scripts/install.sh` - Installation script
- `scripts/uninstall.sh` - Uninstallation script

**Related specs:**
- `ci-cd-workflows.md` - Release workflow that creates binaries
- `build-system.md` - Makefile that builds binaries

## Examples

### Example 1: Install Latest Version on macOS

**Input:**
```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash
```

**Expected Output:**
```
[INFO] Installing ARM (AI Rules Manager)...
[INFO] Installing latest version: v3.0.0
[INFO] Downloading ARM v3.0.0 for darwin-arm64...
[INFO] ARM installed to /usr/local/bin/arm
[INFO] ARM is ready! Run 'arm help' to get started
```

**Verification:**
- Binary exists at /usr/local/bin/arm
- `arm version` displays v3.0.0
- `which arm` returns /usr/local/bin/arm

### Example 2: Install Specific Version on Linux

**Input:**
```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash -s v2.5.0
```

**Expected Output:**
```
[INFO] Installing ARM (AI Rules Manager)...
[INFO] Installing specific version: v2.5.0
[INFO] Downloading ARM v2.5.0 for linux-amd64...
[INFO] ARM installed to /usr/local/bin/arm
[INFO] ARM is ready! Run 'arm help' to get started
```

**Verification:**
- Binary exists at /usr/local/bin/arm
- `arm version` displays v2.5.0

### Example 3: Install on Windows (Git Bash)

**Input:**
```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/install.sh | bash
```

**Expected Output:**
```
[INFO] Installing ARM (AI Rules Manager)...
[INFO] Installing latest version: v3.0.0
[INFO] Downloading ARM v3.0.0 for windows-amd64...
[INFO] ARM installed to ~/AppData/Local/Programs/arm/arm.exe
[INFO] Next step: Add ~/AppData/Local/Programs/arm to your PATH:
[INFO]    echo 'export PATH="~/AppData/Local/Programs/arm:$PATH"' >> ~/.bashrc
[INFO]    source ~/.bashrc
```

**Verification:**
- Binary exists at ~/AppData/Local/Programs/arm/arm.exe
- After adding to PATH: `arm version` works

### Example 4: Uninstall on macOS

**Input:**
```bash
curl -fsSL https://raw.githubusercontent.com/jomadu/ai-resource-manager/main/scripts/uninstall.sh | bash
```

**Expected Output:**
```
[INFO] Uninstalling ARM...
[INFO] Removing /usr/local/bin/arm
[INFO] ARM uninstalled successfully
```

**Verification:**
- Binary removed from /usr/local/bin/arm
- `which arm` returns nothing
- `arm version` fails with "command not found"

## Notes

- Install script uses `-fsSL` flags for curl: fail silently, show errors, follow redirects, silent progress
- Scripts use colored output (GREEN for info, YELLOW for warnings, RED for errors)
- Windows installation requires Git Bash or similar environment with tar support
- Scripts handle version strings with or without 'v' prefix
- Temporary directories are always cleaned up, even on error
- Scripts use `set -e` to exit on any error

## Known Issues

- Windows installation may fail without Git Bash (requires tar command)
- PATH configuration on Windows varies by shell (bash, PowerShell, cmd)
- Uninstall script doesn't remove configuration files (~/.armrc, ~/.arm/)

## Areas for Improvement

- Add support for PowerShell installation on Windows
- Add option to install to custom directory
- Add option to install from local binary (for testing)
- Add option to verify binary signature/checksum
- Add option to uninstall configuration files
- Add progress bar for downloads
- Add retry logic for network failures
- Add support for installing from specific commit/branch (for testing)
