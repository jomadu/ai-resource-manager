#!/bin/bash
# Usage: ./reverse-loop.sh [max_iterations]
# Generates/updates specs from implementation

MAX_ITERATIONS=${1:-0}
ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Mode:   reverse (implementation → specs)"
echo "Branch: $CURRENT_BRANCH"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:    $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

while true; do
    if [ $MAX_ITERATIONS -gt 0 ] && [ $ITERATION -ge $MAX_ITERATIONS ]; then
        echo "Reached max iterations: $MAX_ITERATIONS"
        break
    fi

    cat PROMPT_spec.md | kiro-cli chat --no-interactive --trust-all-tools

    git push origin "$CURRENT_BRANCH" || git push -u origin "$CURRENT_BRANCH"

    ITERATION=$((ITERATION + 1))
    echo -e "\n\n======================== REVERSE LOOP $ITERATION ========================\n"
done
