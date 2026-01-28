# Ralph Prompt Improvements (Grug Brain Edition)

## Grug Brain Philosophy

> grug brain developer not so smart, but grug brain developer program many long year and learn some things although mostly still confused
> 
> complexity very, very bad. say again: complexity very, very bad.

Current Ralph prompts work. Don't break what works. Make small improvements only.

## Core Problem (Simple Version)

Ralph prompts are hard to:
- **Debug** - When Ralph fails, which instruction broke?
- **Reuse** - Want same instructions across modes
- **Understand** - What is Ralph supposed to do?

## Grug Solution: Do Less, Not More

## Grug Solution: Do Less, Not More

### Improvement 1: Add Headers (No New Files)

Current:
```
0a. Study specs/* with up to 500 parallel Sonnet subagents...
0b. Study @IMPLEMENTATION_PLAN.md...
1. Your task is to implement...
```

Better:
```markdown
# ORIENT

Study specs/* with up to 500 parallel Sonnet subagents...
Study @IMPLEMENTATION_PLAN.md...

# TASK

Your task is to implement...

# VALIDATE

After implementing, run tests...

# COMMIT

When tests pass, update plan and commit...

# GUARDRAILS

- Important: Capture the why
- Important: No placeholders
- Important: Keep plan current
```

Why grug like:
- Same file, just add `#` headers
- Agent can say "doing ORIENT" or "stuck at VALIDATE"
- Human can debug: "Ralph always fails at VALIDATE"
- No templates, no variables, no sed, no complexity

### Improvement 2: Put Priority in Guardrails

Current guardrails confusing:
```
99999. Important: When authoring documentation...
999999. Important: Single sources of truth...
9999999. As soon as there are no build errors...
```

Grug not know what number mean. More 9s = more important? Why?

Better:
```markdown
# GUARDRAILS

Priority 1 (Must do every time):
- Keep IMPLEMENTATION_PLAN.md current after each task
- Commit only when tests pass
- No placeholders or stubs

Priority 2 (Important but not blocking):
- Update AGENTS.md when learning operational things
- Clean completed items from plan when large
- Add logging if needed for debugging

Priority 3 (Nice to have):
- Create git tag when no errors and ready
```

Why grug like:
- Clear priority
- Agent know what can skip if stuck
- Human can say "ignore Priority 3" to speed up

### Improvement 3: Agent Logs What It's Doing

Add to prompt:
```markdown
# LOGGING

At start of each section, output:
[ORIENT] Starting orientation...
[TASK] Starting task...
[VALIDATE] Running validation...
[COMMIT] Committing changes...

When done with section:
[ORIENT] ✓ Complete
[TASK] ✓ Complete
[VALIDATE] ✓ Complete
```

Why grug like:
- See where Ralph is
- See where Ralph gets stuck
- No fancy logging framework
- Just ask agent to print things

## What Grug NOT Do

❌ **Variables and sed** - Just write the number in prompt
❌ **Include files** - One prompt file, that's it
❌ **YAML frontmatter** - More syntax to learn, can break
❌ **Template engine** - New dependency, new complexity
❌ **Config files** - More files to manage

## Grug's Three Rules

1. **One prompt file** - No includes, no compilation, no variables
2. **If want different number, edit the file** - Simple and clear
3. **If agent confused, make words simpler** - Not more structure

## Implementation (Grug Way)

### Step 1: Add headers to existing prompts (5 minutes)

```bash
# Edit PROMPT_build.md
# Add: # ORIENT, # TASK, # VALIDATE, # COMMIT, # GUARDRAILS
# Done
```

### Step 2: Fix guardrail numbers (5 minutes)

```bash
# Replace 99999. with "Priority 1:"
# Replace 999999. with "Priority 2:"
# Done
```

### Step 3: Ask agent to log (2 minutes)

```bash
# Add "# LOGGING" section to prompt
# Ask agent to print [SECTION] messages
# Done
```

**Total time: 12 minutes**

No variables. No includes. No sed. Just edit the prompt file.

## Example: Improved PROMPT_build.md (Grug Version)

```markdown
# ORIENT

Study specs/* with up to 500 parallel Sonnet subagents to learn specifications.
Study @IMPLEMENTATION_PLAN.md to understand current work.
Study AGENTS.md to learn how to build and test.

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Follow @IMPLEMENTATION_PLAN.md and choose the most important item.

Before making changes:
- Search codebase first (don't assume not implemented)
- Use up to 500 parallel Sonnet subagents for searches
- Use only 1 subagent for build/tests
- Use Opus subagents when complex reasoning needed

Log: [TASK] Starting task: <task name>
Log when done: [TASK] ✓ Complete

# VALIDATE

Run tests for the code you changed: `go test ./...`

If tests fail:
- Fix the issues
- Run tests again
- Repeat until tests pass

If tests pass:
- Update @IMPLEMENTATION_PLAN.md with findings (use subagent)
- Mark task as complete

Log: [VALIDATE] Running tests...
Log when done: [VALIDATE] ✓ Tests passed

# COMMIT

When tests pass:
- `git add -A`
- `git commit -m "conventional commit message"`
- `git push`

Log: [COMMIT] Committing changes...
Log when done: [COMMIT] ✓ Complete

# GUARDRAILS

Priority 1 (Must do):
- Keep @IMPLEMENTATION_PLAN.md current - future work depends on this
- Commit only when tests pass
- No placeholders or stubs - implement completely
- If unrelated tests fail, fix them too

Priority 2 (Important):
- Update @AGENTS.md when learning operational things (keep brief)
- Capture the why in documentation
- Clean completed items from plan when it gets large

Priority 3 (Nice to have):
- Create git tag when no errors (start at 0.0.0, increment patch)
- Add logging if needed for debugging

# LOGGING

At start of each section above, output the log message shown.
When done with section, output the completion message.
This helps humans see where you are and where you get stuck.
```

## Example: loop.sh (Grug Version)

```bash
#!/bin/bash
# Usage: ./loop.sh [plan|spec] [max_iterations]

# Parse arguments (same as before)
if [ "$1" = "plan" ]; then
    MODE="plan"
    PROMPT_FILE="PROMPT_plan.md"
    MAX_ITERATIONS=${2:-0}
elif [ "$1" = "spec" ]; then
    MODE="spec"
    PROMPT_FILE="PROMPT_spec.md"
    MAX_ITERATIONS=${2:-0}
elif [[ "$1" =~ ^[0-9]+$ ]]; then
    MODE="build"
    PROMPT_FILE="PROMPT_build.md"
    MAX_ITERATIONS=$1
else
    MODE="build"
    PROMPT_FILE="PROMPT_build.md"
    MAX_ITERATIONS=0
fi

ITERATION=0
CURRENT_BRANCH=$(git branch --show-current)

echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "Mode:   $MODE"
echo "Prompt: $PROMPT_FILE"
echo "Branch: $CURRENT_BRANCH"
[ $MAX_ITERATIONS -gt 0 ] && echo "Max:    $MAX_ITERATIONS iterations"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

if [ ! -f "$PROMPT_FILE" ]; then
    echo "Error: $PROMPT_FILE not found"
    exit 1
fi

while true; do
    if [ $MAX_ITERATIONS -gt 0 ] && [ $ITERATION -ge $MAX_ITERATIONS ]; then
        echo "Reached max iterations: $MAX_ITERATIONS"
        break
    fi

    # Just cat the file, no sed, no variables, no complexity
    cat "$PROMPT_FILE" | kiro-cli chat --no-interactive --trust-all-tools

    git push origin "$CURRENT_BRANCH" || {
        echo "Failed to push. Creating remote branch..."
        git push -u origin "$CURRENT_BRANCH"
    }

    ITERATION=$((ITERATION + 1))
    echo -e "\n\n======================== LOOP $ITERATION ========================\n"
done
```

## When to Add Complexity (Grug Future)

Only add complexity when:
1. **Pain is real** - Not theoretical, actually hurting
2. **Simple solution tried** - Headers and priorities not enough
3. **Benefit is clear** - Know exactly what problem it solves

Examples of real pain:
- "Ralph always fails at same step" → Need better logging (add it)
- "Want to reuse instructions across modes" → Copy/paste is fine for 3 files

Examples of NOT real pain:
- "Changing subagent count in 3 files is annoying" → Just change 3 files, takes 30 seconds
- "Could make variables" → Variables add complexity, just edit the number
- "Could include common parts" → Copy/paste is simpler than includes

## Summary

**Original plan**: Variables, sed, includes, YAML, templates
**Grug plan**: Headers, priority numbers, logging

**Original time**: 27 minutes
**Grug time**: 12 minutes

**Original files**: Multiple (PROMPT_common.md, etc.)
**Grug files**: One prompt file per mode

**Original loop.sh**: sed substitution
**Grug loop.sh**: Just cat the file

Grug brain developer say: **One file. No magic. Just words.**

Ralph already work. Add headers. Add priorities. Add logging. Done.

If want different number, edit the file. If want different words, edit the file. File is source of truth.

Complexity is the enemy. Simple is the friend.

Grug done now. Go write code.
