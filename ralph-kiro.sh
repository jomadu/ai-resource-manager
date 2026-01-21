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

PROMPT_FILE="RALPH.md"

echo "Starting Ralph with Kiro CLI - Max iterations: $MAX_ITERATIONS"
[ -n "$AGENT" ] && echo "Agent: $AGENT"

for i in $(seq 1 $MAX_ITERATIONS); do
  echo ""
  echo "==============================================================="
  echo "  Ralph Iteration $i of $MAX_ITERATIONS (Kiro CLI)"
  echo "==============================================================="

  ARGS=(chat --no-interactive --trust-all-tools)
  [ -n "$AGENT" ] && ARGS+=(--agent "$AGENT")
  
  OUTPUT=$(cat "$PROMPT_FILE" | kiro-cli "${ARGS[@]}" 2>&1 | tee /dev/stderr) || true
  
  if echo "$OUTPUT" | grep -q "<promise>COMPLETE</promise>"; then
    echo ""
    echo "✓ Ralph completed all tasks!"
    echo "Completed at iteration $i of $MAX_ITERATIONS"
    exit 0
  fi
  
  echo "Iteration $i complete. Continuing..."
  sleep 2
done

echo ""
echo "Ralph reached max iterations ($MAX_ITERATIONS) without completing."
exit 1
