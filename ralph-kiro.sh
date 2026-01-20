#!/bin/bash
# Ralph Wiggum Loop for Kiro CLI
# Usage: ./ralph-kiro.sh [--agent AGENT_NAME] [max_iterations]

set -e

# Parse arguments
AGENT=""
MAX_ITERATIONS=10

while [[ $# -gt 0 ]]; do
  case $1 in
    --agent)
      AGENT="$2"
      shift 2
      ;;
    --agent=*)
      AGENT="${1#*=}"
      shift
      ;;
    *)
      if [[ "$1" =~ ^[0-9]+$ ]]; then
        MAX_ITERATIONS="$1"
      fi
      shift
      ;;
  esac
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PRD_FILE="prd.json"
PROGRESS_FILE="progress.txt"
PROMPT_FILE="RALPH.md"
ARCHIVE_DIR="archive"
LAST_BRANCH_FILE=".last-branch"
SPECS_DIR="specs"

# Save progress back to feature directory
save_progress() {
  if [ -n "$FEATURE_DIR" ] && [ -d "$FEATURE_DIR" ]; then
    echo "Saving progress to: $FEATURE_DIR/"
    cp "$PRD_FILE" "$FEATURE_DIR/"
    cp "$PROGRESS_FILE" "$FEATURE_DIR/"
  fi
}

# Archive previous run if branch changed
if [ -f "$PRD_FILE" ] && [ -f "$LAST_BRANCH_FILE" ]; then
  CURRENT_BRANCH=$(jq -r '.branchName // empty' "$PRD_FILE" 2>/dev/null || echo "")
  LAST_BRANCH=$(cat "$LAST_BRANCH_FILE" 2>/dev/null || echo "")
  
  if [ -n "$CURRENT_BRANCH" ] && [ -n "$LAST_BRANCH" ] && [ "$CURRENT_BRANCH" != "$LAST_BRANCH" ]; then
    DATE=$(date +%Y-%m-%d)
    FOLDER_NAME=$(echo "$LAST_BRANCH" | sed 's|^ralph/||')
    ARCHIVE_FOLDER="$ARCHIVE_DIR/$DATE-$FOLDER_NAME"
    
    echo "Archiving previous run: $LAST_BRANCH"
    mkdir -p "$ARCHIVE_FOLDER"
    [ -f "$PRD_FILE" ] && cp "$PRD_FILE" "$ARCHIVE_FOLDER/"
    [ -f "$PROGRESS_FILE" ] && cp "$PROGRESS_FILE" "$ARCHIVE_FOLDER/"
    echo "   Archived to: $ARCHIVE_FOLDER"
    
    echo "# Ralph Progress Log" > "$PROGRESS_FILE"
    echo "Started: $(date)" >> "$PROGRESS_FILE"
    echo "---" >> "$PROGRESS_FILE"
  fi
fi

# Track current branch and determine feature directory
FEATURE_DIR=""
if [ -f "$PRD_FILE" ]; then
  CURRENT_BRANCH=$(jq -r '.branchName // empty' "$PRD_FILE" 2>/dev/null || echo "")
  if [ -n "$CURRENT_BRANCH" ]; then
    echo "$CURRENT_BRANCH" > "$LAST_BRANCH_FILE"
    
    # Extract feature name from branch (remove ralph/ prefix if present)
    FEATURE_NAME=$(echo "$CURRENT_BRANCH" | sed 's|^ralph/||')
    FEATURE_DIR="$SPECS_DIR/$FEATURE_NAME"
    
    # Copy existing progress.txt from feature directory if it exists
    if [ -d "$FEATURE_DIR" ] && [ -f "$FEATURE_DIR/$PROGRESS_FILE" ]; then
      echo "Loading existing progress from: $FEATURE_DIR/$PROGRESS_FILE"
      cp "$FEATURE_DIR/$PROGRESS_FILE" "$PROGRESS_FILE"
    fi
  fi
fi

# Initialize progress file if it doesn't exist
if [ ! -f "$PROGRESS_FILE" ]; then
  echo "# Ralph Progress Log" > "$PROGRESS_FILE"
  echo "Started: $(date)" >> "$PROGRESS_FILE"
  echo "---" >> "$PROGRESS_FILE"
fi

echo "Starting Ralph with Kiro CLI - Max iterations: $MAX_ITERATIONS"
[ -n "$AGENT" ] && echo "Agent: $AGENT"

for i in $(seq 1 $MAX_ITERATIONS); do
  echo ""
  echo "==============================================================="
  echo "  Ralph Iteration $i of $MAX_ITERATIONS (Kiro CLI)"
  echo "==============================================================="

  ARGS=(chat --no-interactive --trust-all-tools)
  [ -n "$AGENT" ] && ARGS+=(--agent "$AGENT")
  
  # Capture output silently, check for promises
  OUTPUT=$(cat "$PROMPT_FILE" | kiro-cli "${ARGS[@]}" 2>&1) || true
  
  # Check for completion
  if echo "$OUTPUT" | grep -q "<promise>COMPLETE</promise>"; then
    echo ""
    echo "✓ Ralph completed all tasks!"
    echo "Completed at iteration $i of $MAX_ITERATIONS"
    save_progress
    exit 0
  fi
  
  # Check for blocked
  if echo "$OUTPUT" | grep -q "<promise>BLOCKED:"; then
    echo ""
    echo "✗ Ralph is blocked:"
    echo "$OUTPUT" | grep -o "<promise>BLOCKED:.*</promise>" | sed 's/<[^>]*>//g'
    save_progress
    exit 1
  fi
  
  # Check if no stories left (alternative completion signal)
  if echo "$OUTPUT" | grep -q "All stories complete\|No stories remaining\|All user stories.*passes: true"; then
    echo ""
    echo "✓ Ralph completed all stories!"
    echo "Completed at iteration $i of $MAX_ITERATIONS"
    save_progress
    exit 0
  fi
  
  echo ""
  echo "→ Iteration $i complete. Continuing..."
  sleep 2
done

echo ""
echo "Ralph reached max iterations ($MAX_ITERATIONS) without completing all tasks."
echo "Check $PROGRESS_FILE for status."
save_progress
exit 1
