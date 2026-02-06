#!/bin/bash
VERSION="0.1.0"
# Usage: ./rooda.sh <procedure> [-c <file>] [-m N]
#    OR: ./rooda.sh -o <file> -r <file> -d <file> -a <file> [-m N]
# Examples:
#   ./rooda.sh build
#   ./rooda.sh build -m 5
#   ./rooda.sh build -c custom-config.yml
#   ./rooda.sh -o prompts/observe_specs.md \
#             -r prompts/orient_gap.md \
#             -d prompts/decide_gap_plan.md \
#             -a prompts/act_plan.md \
#             -m 1

show_help() {
    cat <<EOF
Usage: ./rooda.sh <procedure> [--config <file>] [--max-iterations N]
   OR: ./rooda.sh --observe <file> --orient <file> --decide <file> --act <file> [--max-iterations N]

Options:
  <procedure>              Named procedure from config (bootstrap, build, etc.)
  -c, --config <file>      Path to config file (default: rooda-config.yml)
  --version                Show version number
  --list-procedures        List all available procedures from config
  -o, --observe <file>     Path to observe phase prompt
  -r, --orient <file>      Path to orient phase prompt
  -d, --decide <file>      Path to decide phase prompt
  -a, --act <file>         Path to act phase prompt
  -m, --max-iterations N   Maximum iterations (default: see below)
  --ai-cli <command>       AI CLI command to use (overrides all other settings)
  --ai-tool <preset>       AI tool preset (kiro-cli, claude, aider, or custom from config)
  --verbose                Show detailed execution including full prompt
  --quiet                  Suppress non-error output
  --help, -h               Show this help message

AI CLI Precedence (highest to lowest):
  1. --ai-cli flag (direct command override)
  2. --ai-tool preset (resolves to command via config or hardcoded)
  3. \$ROODA_AI_CLI environment variable
  4. Default: kiro-cli chat --no-interactive --trust-all-tools

Hardcoded AI Tool Presets:
  - kiro-cli: kiro-cli chat --no-interactive --trust-all-tools
  - claude: claude-cli --no-interactive
  - aider: aider --yes

Custom presets can be defined in rooda-config.yml under ai_tools section.

Max Iterations Default Behavior (three-tier system):
  1. Command-line --max-iterations takes precedence
  2. Config default_iterations used if CLI not specified
  3. Defaults to 0 (unlimited) if neither specified

Examples:
  ./rooda.sh bootstrap
  ./rooda.sh build -m 5
  ./rooda.sh build --verbose
  ./rooda.sh build --quiet
  ./rooda.sh build --ai-tool claude
  ./rooda.sh build --ai-cli "custom-cli --flags"
  ./rooda.sh --list-procedures
  ./rooda.sh -o prompts/observe_specs.md \\
            -r prompts/orient_gap.md \\
            -d prompts/decide_gap_plan.md \\
            -a prompts/act_plan.md \\
            -m 1
EOF
}

list_procedures() {
    local config_file="$1"
    
    if [ ! -f "$config_file" ]; then
        echo "Error: Configuration file not found: $config_file"
        exit 1
    fi
    
    echo "Available procedures:"
    echo ""
    
    local procedures
    procedures=$(yq eval '.procedures | keys | .[]' "$config_file")
    
    while IFS= read -r proc; do
        local display summary
        display=$(yq eval ".procedures.$proc.display // \"\"" "$config_file")
        summary=$(yq eval ".procedures.$proc.summary // \"\"" "$config_file")
        
        if [ -n "$display" ]; then
            echo "  $proc - $display"
        else
            echo "  $proc"
        fi
        
        if [ -n "$summary" ]; then
            echo "    $summary"
        fi
        echo ""
    done <<< "$procedures"
}

resolve_ai_tool_preset() {
    local preset="$1"
    local config_file="$2"
    
    # Hardcoded presets
    case "$preset" in
        kiro-cli)
            echo "kiro-cli chat --no-interactive --trust-all-tools"
            return 0
            ;;
        claude)
            echo "claude-cli --no-interactive"
            return 0
            ;;
        aider)
            echo "aider --yes"
            return 0
            ;;
    esac
    
    # Query custom preset from config
    local custom_command
    custom_command=$(yq eval ".ai_tools.$preset" "$config_file" 2>&1)
    
    if [ "$custom_command" != "null" ] && [ -n "$custom_command" ]; then
        echo "$custom_command"
        return 0
    fi
    
    # Unknown preset - return error with helpful message
    echo "Error: Unknown AI tool preset: $preset" >&2
    echo "" >&2
    echo "Available hardcoded presets:" >&2
    echo "  - kiro-cli" >&2
    echo "  - claude" >&2
    echo "  - aider" >&2
    echo "" >&2
    echo "To define custom presets, add to $config_file:" >&2
    echo "  ai_tools:" >&2
    echo "    $preset: \"your-command-here\"" >&2
    return 1
}

# Detect OS for platform-specific instructions
OS="$(uname -s)"
case "$OS" in
    Darwin*) PLATFORM="macos" ;;
    Linux*)  PLATFORM="linux" ;;
    *)       PLATFORM="unknown" ;;
esac

# Check for required dependencies
if ! command -v yq &> /dev/null; then
    echo "Error: yq is required for YAML parsing"
    if [ "$PLATFORM" = "macos" ]; then
        echo "Install with: brew install yq"
    elif [ "$PLATFORM" = "linux" ]; then
        echo "Install with: wget https://github.com/mikefarah/yq/releases/latest/download/yq_linux_amd64 -O /usr/local/bin/yq && chmod +x /usr/local/bin/yq"
    else
        echo "See: https://github.com/mikefarah/yq#install"
    fi
    exit 1
fi

# Check versions
YQ_VERSION=$(yq --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
YQ_MAJOR=$(echo "$YQ_VERSION" | cut -d. -f1)
if [ "$YQ_MAJOR" -lt 4 ]; then
    echo "Error: yq version 4.0.0 or higher required (found $YQ_VERSION)"
    echo "Upgrade with: brew upgrade yq"
    exit 1
fi

# Check for optional dependencies
if ! command -v shellcheck &> /dev/null; then
    echo "Warning: shellcheck not installed (optional but recommended for linting)"
    if [ "$PLATFORM" = "macos" ]; then
        echo "  Install with: brew install shellcheck"
    elif [ "$PLATFORM" = "linux" ]; then
        echo "  Install with: apt-get install shellcheck (or yum install shellcheck)"
    fi
fi

if ! command -v git &> /dev/null; then
    echo "Warning: git not installed (optional but recommended for version control)"
    if [ "$PLATFORM" = "macos" ]; then
        echo "  Install with: brew install git"
    elif [ "$PLATFORM" = "linux" ]; then
        echo "  Install with: apt-get install git (or yum install git)"
    fi
fi

# Validate config structure
validate_config() {
    local config_file="$1"
    local procedure="$2"
    
    # Validate YAML is parseable
    if ! yq eval '.' "$config_file" &> /dev/null; then
        echo "Error: Invalid YAML in configuration file"
        echo "  Path: $config_file"
        echo "  Run: yq eval '.' $config_file"
        echo "  To see parse errors"
        exit 1
    fi
    
    # Validate procedures key exists
    local procedures_key
    procedures_key=$(yq eval '.procedures' "$config_file")
    if [ "$procedures_key" = "null" ]; then
        echo "Error: Configuration file missing 'procedures' key"
        echo "  Path: $config_file"
        echo "  Expected structure: procedures: { procedure-name: { ... } }"
        exit 1
    fi
    
    # If procedure specified, validate it exists and has required fields
    if [ -n "$procedure" ]; then
        local proc_exists
        proc_exists=$(yq eval ".procedures.$procedure" "$config_file")
        if [ "$proc_exists" = "null" ]; then
            echo "Error: Procedure '$procedure' not found in $config_file"
            echo ""
            
            # Find closest match using simple character overlap
            local available_procs
            available_procs=$(yq eval '.procedures | keys | .[]' "$config_file")
            local best_match=""
            local best_score=0
            
            while IFS= read -r proc; do
                # Count matching characters (simple fuzzy matching)
                local score=0
                local proc_lower
                local input_lower
                proc_lower=$(echo "$proc" | tr '[:upper:]' '[:lower:]')
                input_lower=$(echo "$procedure" | tr '[:upper:]' '[:lower:]')
                
                # Substring match gets high score
                if [[ "$proc_lower" == *"$input_lower"* ]] || [[ "$input_lower" == *"$proc_lower"* ]]; then
                    score=100
                # Count common characters
                else
                    for (( i=0; i<${#input_lower}; i++ )); do
                        local char="${input_lower:$i:1}"
                        if [[ "$proc_lower" == *"$char"* ]]; then
                            ((score++))
                        fi
                    done
                fi
                
                if [ $score -gt $best_score ]; then
                    best_score=$score
                    best_match="$proc"
                fi
            done <<< "$available_procs"
            
            # Only suggest if score is reasonable (at least 3 matching chars or substring match)
            if [ -n "$best_match" ] && [ $best_score -ge 3 ]; then
                echo "Did you mean: $best_match"
                echo ""
            fi
            
            echo "Available procedures:"
            while IFS= read -r proc; do
                echo "  - $proc"
            done <<< "$available_procs"
            exit 1
        fi
        
        # Validate required OODA fields
        local observe orient decide act
        observe=$(yq eval ".procedures.$procedure.observe" "$config_file")
        orient=$(yq eval ".procedures.$procedure.orient" "$config_file")
        decide=$(yq eval ".procedures.$procedure.decide" "$config_file")
        act=$(yq eval ".procedures.$procedure.act" "$config_file")
        
        local missing_fields=()
        [ "$observe" = "null" ] && missing_fields+=("observe")
        [ "$orient" = "null" ] && missing_fields+=("orient")
        [ "$decide" = "null" ] && missing_fields+=("decide")
        [ "$act" = "null" ] && missing_fields+=("act")
        
        # Validate fields are non-empty strings
        [ -n "$observe" ] && [ "$observe" != "null" ] && [ -z "${observe// }" ] && missing_fields+=("observe (empty)")
        [ -n "$orient" ] && [ "$orient" != "null" ] && [ -z "${orient// }" ] && missing_fields+=("orient (empty)")
        [ -n "$decide" ] && [ "$decide" != "null" ] && [ -z "${decide// }" ] && missing_fields+=("decide (empty)")
        [ -n "$act" ] && [ "$act" != "null" ] && [ -z "${act// }" ] && missing_fields+=("act (empty)")
        
        if [ ${#missing_fields[@]} -gt 0 ]; then
            echo "Error: Procedure '$procedure' missing required fields"
            echo "  Missing: ${missing_fields[*]}"
            echo "  Required: observe, orient, decide, act (non-empty strings)"
            exit 1
        fi
    fi
}

# Parse arguments
OBSERVE=""
ORIENT=""
DECIDE=""
ACT=""
MAX_ITERATIONS=""  # Empty means not set, will use config default or 0
PROCEDURE=""
VERBOSE=0  # 0=default, 1=verbose, -1=quiet
AI_CLI_COMMAND="kiro-cli chat --no-interactive --trust-all-tools"  # Default AI CLI, configurable via --ai-cli or config
AI_TOOL_PRESET=""  # Preset name for --ai-tool flag
AI_CLI_FLAG=""  # Set by --ai-cli flag (highest precedence)
# Override with environment variable if set (precedence: --ai-cli flag > --ai-tool preset > $ROODA_AI_CLI > default)
[ -n "$ROODA_AI_CLI" ] && AI_CLI_COMMAND="$ROODA_AI_CLI"
# Resolve config file relative to script location
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
CONFIG_FILE="${SCRIPT_DIR}/rooda-config.yml"

# Check for help flag first
if [[ "$1" == "--help" ]] || [[ "$1" == "-h" ]]; then
    show_help
    exit 0
fi

# Check for version flag
if [[ "$1" == "--version" ]]; then
    echo "rooda.sh version $VERSION"
    exit 0
fi

# Check for list-procedures flag
if [[ "$1" == "--list-procedures" ]]; then
    list_procedures "$CONFIG_FILE"
    exit 0
fi

# First positional argument is procedure name (optional)
if [[ $# -gt 0 ]] && [[ ! "$1" =~ ^- ]]; then
    PROCEDURE="$1"
    shift
fi

while [[ $# -gt 0 ]]; do
    case $1 in
        --version)
            echo "rooda.sh version $VERSION"
            exit 0
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        --list-procedures)
            list_procedures "$CONFIG_FILE"
            exit 0
            ;;
        --config|-c)
            CONFIG_FILE="$2"
            shift 2
            ;;
        --observe|-o)
            OBSERVE="$2"
            shift 2
            ;;
        --orient|-r)
            ORIENT="$2"
            shift 2
            ;;
        --decide|-d)
            DECIDE="$2"
            shift 2
            ;;
        --act|-a)
            ACT="$2"
            shift 2
            ;;
        --max-iterations|-m)
            MAX_ITERATIONS="$2"
            shift 2
            ;;
        --ai-cli)
            AI_CLI_FLAG="$2"
            shift 2
            ;;
        --ai-tool)
            AI_TOOL_PRESET="$2"
            shift 2
            ;;
        --verbose)
            VERBOSE=1
            shift
            ;;
        --quiet)
            VERBOSE=-1
            shift
            ;;
        *)
            echo "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
done

# Resolve AI CLI command with correct precedence: --ai-cli > --ai-tool > $ROODA_AI_CLI > default
if [ -n "$AI_CLI_FLAG" ]; then
    # --ai-cli flag has highest priority
    AI_CLI_COMMAND="$AI_CLI_FLAG"
elif [ -n "$AI_TOOL_PRESET" ]; then
    # --ai-tool preset resolution
    AI_CLI_COMMAND=$(resolve_ai_tool_preset "$AI_TOOL_PRESET" "$CONFIG_FILE") || exit 1
fi
# Otherwise use $ROODA_AI_CLI (already set at line 288) or default

# If procedure specified, load from config (explicit flags override config)
if [ -n "$PROCEDURE" ]; then
    if [ ! -f "$CONFIG_FILE" ]; then
        echo "Error: Configuration file not found"
        echo "  Path: $CONFIG_FILE"
        echo "  Specify a different config with --config <file>"
        exit 1
    fi
    
    # Validate config structure and procedure
    validate_config "$CONFIG_FILE" "$PROCEDURE"
    
    # Only load from config if not already set via explicit flags
    if [ -z "$OBSERVE" ]; then
        OBSERVE=$(yq eval ".procedures.$PROCEDURE.observe" "$CONFIG_FILE" 2>&1) || {
            echo "Error: Failed to query 'observe' field from config"
            echo "  Procedure: $PROCEDURE"
            echo "  Config: $CONFIG_FILE"
            echo "  Check that procedure has valid 'observe' field"
            exit 1
        }
    fi
    if [ -z "$ORIENT" ]; then
        ORIENT=$(yq eval ".procedures.$PROCEDURE.orient" "$CONFIG_FILE" 2>&1) || {
            echo "Error: Failed to query 'orient' field from config"
            echo "  Procedure: $PROCEDURE"
            echo "  Config: $CONFIG_FILE"
            echo "  Check that procedure has valid 'orient' field"
            exit 1
        }
    fi
    if [ -z "$DECIDE" ]; then
        DECIDE=$(yq eval ".procedures.$PROCEDURE.decide" "$CONFIG_FILE" 2>&1) || {
            echo "Error: Failed to query 'decide' field from config"
            echo "  Procedure: $PROCEDURE"
            echo "  Config: $CONFIG_FILE"
            echo "  Check that procedure has valid 'decide' field"
            exit 1
        }
    fi
    if [ -z "$ACT" ]; then
        ACT=$(yq eval ".procedures.$PROCEDURE.act" "$CONFIG_FILE" 2>&1) || {
            echo "Error: Failed to query 'act' field from config"
            echo "  Procedure: $PROCEDURE"
            echo "  Config: $CONFIG_FILE"
            echo "  Check that procedure has valid 'act' field"
            exit 1
        }
    fi
    
    # Use default iterations if not specified
    if [ -z "$MAX_ITERATIONS" ]; then
        DEFAULT_ITER=$(yq eval ".procedures.$PROCEDURE.default_iterations" "$CONFIG_FILE" 2>&1) || {
            echo "Error: Failed to query 'default_iterations' field from config"
            echo "  Procedure: $PROCEDURE"
            echo "  Config: $CONFIG_FILE"
            echo "  This field is optional - check config structure"
            exit 1
        }
        if [ "$DEFAULT_ITER" != "null" ]; then
            MAX_ITERATIONS=$DEFAULT_ITER
        else
            MAX_ITERATIONS=0
        fi
    fi
fi

# Check AI CLI availability and version (only if using kiro-cli)
if [[ "$AI_CLI_COMMAND" == kiro-cli* ]]; then
    if ! command -v kiro-cli &> /dev/null; then
        echo "Error: kiro-cli is required for AI CLI integration"
        echo "Install from: https://docs.aws.amazon.com/kiro/"
        exit 1
    fi
    
    KIRO_VERSION=$(kiro-cli --version 2>&1 | grep -oE '[0-9]+\.[0-9]+\.[0-9]+' | head -1)
    KIRO_MAJOR=$(echo "$KIRO_VERSION" | cut -d. -f1)
    if [ "$KIRO_MAJOR" -lt 1 ]; then
        echo "Error: kiro-cli version 1.0.0 or higher required (found $KIRO_VERSION)"
        echo "Upgrade from: https://docs.aws.amazon.com/kiro/"
        exit 1
    fi
fi

# Validate required arguments
if [ -z "$OBSERVE" ] || [ -z "$ORIENT" ] || [ -z "$DECIDE" ] || [ -z "$ACT" ]; then
    echo "Error: All four OODA phases required"
    show_help
    exit 1
fi

# Validate files exist
for phase in OBSERVE ORIENT DECIDE ACT; do
    file="${!phase}"
    if [ ! -f "$file" ]; then
        phase_lower=$(echo "$phase" | tr '[:upper:]' '[:lower:]')
        echo "Error: ${phase} phase file not found"
        echo "  Path: $file"
        echo "  Check that the ${phase_lower} phase file exists and path is correct"
        exit 1
    fi
done

ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

if [ "$VERBOSE" -ge 0 ]; then
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    [ -n "$PROCEDURE" ] && echo "Procedure: $PROCEDURE"
    echo "Observe:   $OBSERVE"
    echo "Orient:    $ORIENT"
    echo "Decide:    $DECIDE"
    echo "Act:       $ACT"
    echo "Branch:    $CURRENT_BRANCH"
    [ -n "$MAX_ITERATIONS" ] && [ "$MAX_ITERATIONS" -gt 0 ] && echo "Max:       $MAX_ITERATIONS iterations"
    echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
fi

# Create prompt template
create_prompt() {
    # Assemble four OODA phase prompt files into single executable prompt
    # Uses heredoc (<<EOF) to create template with embedded command substitution
    # Each $(cat "$VAR") is evaluated when heredoc executes, inserting file contents
    cat <<EOF
# OODA Loop Iteration

## OBSERVE
$(cat "$OBSERVE")

## ORIENT
$(cat "$ORIENT")

## DECIDE
$(cat "$DECIDE")

## ACT
$(cat "$ACT")
EOF
}

while true; do
    if [ -n "$MAX_ITERATIONS" ] && [ "$MAX_ITERATIONS" -gt 0 ] && [ "$ITERATION" -ge "$MAX_ITERATIONS" ]; then
        [ "$VERBOSE" -ge 0 ] && echo "Reached max iterations: $MAX_ITERATIONS"
        break
    fi

    # Show full prompt in verbose mode
    if [ "$VERBOSE" -eq 1 ]; then
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        echo "Full prompt being sent to AI CLI:"
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
        create_prompt
        echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
    fi

    # Execute AI CLI - exit status intentionally ignored per ai-cli-integration.md
    # Design: Script continues to git push regardless of AI CLI success/failure
    # Rationale: Allows loop to self-correct through empirical feedback in subsequent iterations
    create_prompt | $AI_CLI_COMMAND

    PUSH_OUTPUT=$(git push origin "$CURRENT_BRANCH" 2>&1)
    PUSH_EXIT=$?
    if [ $PUSH_EXIT -ne 0 ]; then
        PUSH_OUTPUT_FALLBACK=$(git push -u origin "$CURRENT_BRANCH" 2>&1)
        PUSH_FALLBACK_EXIT=$?
        if [ $PUSH_FALLBACK_EXIT -eq 0 ]; then
            echo "Created remote branch and pushed successfully"
        else
            echo "Error: Failed to push to remote"
            
            # Parse error output to provide specific guidance
            COMBINED_OUTPUT="$PUSH_OUTPUT $PUSH_OUTPUT_FALLBACK"
            if echo "$COMBINED_OUTPUT" | grep -qi "permission denied\|authentication failed\|403"; then
                echo "Cause: Authentication failure"
                echo "Fix: Check SSH keys (ssh -T git@github.com) or credentials"
            elif echo "$COMBINED_OUTPUT" | grep -qi "could not resolve host\|connection timed out\|network is unreachable"; then
                echo "Cause: Network issue"
                echo "Fix: Check internet connectivity and try again"
            elif echo "$COMBINED_OUTPUT" | grep -qi "rejected\|non-fast-forward\|would be overwritten"; then
                echo "Cause: Merge conflict or out-of-date branch"
                echo "Fix: Run 'git pull --rebase origin $CURRENT_BRANCH' and resolve conflicts"
            else
                echo "Cause: Unknown (see error above)"
                echo "Fix: Review git output and resolve the issue"
            fi
            
            echo "The next iteration will attempt to push again"
            echo "Press Ctrl+C to stop, or Enter to continue..."
            read -r
        fi
    fi

    ITERATION=$((ITERATION + 1))
    [ "$VERBOSE" -ge 0 ] && echo -e "\n\n======================== Starting iteration $ITERATION ========================\n"
done
