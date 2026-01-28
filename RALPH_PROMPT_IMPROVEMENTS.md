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
- No templates, no YAML, no complexity

### Improvement 2: Put Numbers in Guardrails (Not Everywhere)

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

### Improvement 3: One File for Common Things

Make `PROMPT_common.md`:
```markdown
# ORIENT

Study specs/* with up to 500 parallel Sonnet subagents to learn specifications.
Study @IMPLEMENTATION_PLAN.md to understand current work.
Study AGENTS.md to learn how to build/test.

# GUARDRAILS

Priority 1:
- Keep IMPLEMENTATION_PLAN.md current
- Commit only when tests pass
- No placeholders

Priority 2:
- Update AGENTS.md for operational learnings
- Clean completed items when plan gets large
```

Then mode-specific prompts just include it:
```markdown
# MODE: BUILD

{{include PROMPT_common.md}}

# TASK

Follow IMPLEMENTATION_PLAN.md and choose most important item.
Search codebase first (don't assume not implemented).
Implement using parallel subagents.

# VALIDATE

Run tests for the code you changed.
If tests pass, update plan and commit.
```

Why grug like:
- Don't repeat same instructions in every file
- Change common thing once, affects all modes
- Still just markdown files, no fancy template engine
- Can use simple `cat` or `sed` to include

### Improvement 4: Make Variables Obvious

Current:
```
Study specs/* with up to 500 parallel Sonnet subagents...
```

What if want different number? Edit every prompt? Grug not like.

Better in loop.sh:
```bash
# At top of loop.sh
SUBAGENT_SEARCH=500
SUBAGENT_BUILD=1
TEST_COMMAND="go test ./..."

# Replace in prompt before sending
cat "$PROMPT_FILE" | \
    sed "s/{{SUBAGENT_SEARCH}}/$SUBAGENT_SEARCH/g" | \
    sed "s/{{SUBAGENT_BUILD}}/$SUBAGENT_BUILD/g" | \
    sed "s/{{TEST_COMMAND}}/$TEST_COMMAND/g" | \
    kiro-cli chat --no-interactive --trust-all-tools
```

Prompt uses simple placeholders:
```markdown
Study specs/* with up to {{SUBAGENT_SEARCH}} parallel Sonnet subagents...
Use only {{SUBAGENT_BUILD}} subagent for build/tests.
Run: {{TEST_COMMAND}}
```

Why grug like:
- One place to change numbers
- No YAML, no config files
- Just bash variables and sed
- Works today, no new tools

### Improvement 5: Agent Logs What It's Doing

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

❌ **YAML frontmatter** - More syntax to learn, can break
❌ **Template engine** - New dependency, new complexity
❌ **Config files** - More files to manage
❌ **Validation pipelines** - Fancy but not needed yet
❌ **Adaptive budgets** - Sounds smart but grug not understand when needed
❌ **Multiple strategies** - One strategy work fine, why add more?

## Grug's Three Rules

1. **If current thing work, change small** - Don't rewrite everything
2. **If can do in bash, do in bash** - No new tools
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

### Step 3: Extract common parts (10 minutes)

```bash
# Create PROMPT_common.md with ORIENT and GUARDRAILS
# Update loop.sh to cat PROMPT_common.md and PROMPT_build.md together
# Done
```

### Step 4: Add variables to loop.sh (5 minutes)

```bash
# Add SUBAGENT_SEARCH=500 at top
# Add sed to replace {{SUBAGENT_SEARCH}} before piping
# Done
```

### Step 5: Ask agent to log (2 minutes)

```bash
# Add "# LOGGING" section to prompt
# Ask agent to print [SECTION] messages
# Done
```

**Total time: 27 minutes**

Compare to original plan: weeks of work, new tools, YAML, templates, configs.

Grug way: half hour, no new dependencies, works today.

## Example: Improved PROMPT_build.md (Grug Version)

```markdown
# ORIENT

Study specs/* with up to {{SUBAGENT_SEARCH}} parallel Sonnet subagents to learn specifications.
Study @IMPLEMENTATION_PLAN.md to understand current work.
Study AGENTS.md to learn how to build and test.

Log: [ORIENT] Starting orientation...
Log when done: [ORIENT] ✓ Complete

# TASK

Follow @IMPLEMENTATION_PLAN.md and choose the most important item.

Before making changes:
- Search codebase first (don't assume not implemented)
- Use up to {{SUBAGENT_SEARCH}} parallel Sonnet subagents for searches
- Use only {{SUBAGENT_BUILD}} subagent for build/tests
- Use Opus subagents when complex reasoning needed

Log: [TASK] Starting task: <task name>
Log when done: [TASK] ✓ Complete

# VALIDATE

Run tests for the code you changed: {{TEST_COMMAND}}

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
```

## Example: loop.sh (Grug Version)

```bash
#!/bin/bash

# Configuration (one place to change things)
SUBAGENT_SEARCH=500
SUBAGENT_BUILD=1
TEST_COMMAND="go test ./..."

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

    # Simple variable substitution (grug way)
    cat "$PROMPT_FILE" | \
        sed "s/{{SUBAGENT_SEARCH}}/$SUBAGENT_SEARCH/g" | \
        sed "s/{{SUBAGENT_BUILD}}/$SUBAGENT_BUILD/g" | \
        sed "s/{{TEST_COMMAND}}/$TEST_COMMAND/g" | \
        kiro-cli chat --no-interactive --trust-all-tools

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
2. **Simple solution tried** - Headers and variables not enough
3. **Benefit is clear** - Know exactly what problem it solves

Examples of real pain:
- "Ralph always fails at same step" → Need better logging (add it)
- "Changing subagent count in 5 files is annoying" → Need variables (add it)
- "Want to reuse instructions across modes" → Need includes (add it)

Examples of theoretical pain:
- "Might want 10 different modes someday" → Don't have 10 modes yet, don't build for it
- "Could make adaptive budgets" → Current fixed budgets working fine
- "YAML would be more structured" → Markdown working fine

## Summary

**Original plan**: YAML, templates, configs, validation pipelines, adaptive budgets, multiple strategies
**Grug plan**: Headers, priority numbers, simple variables, logging

**Original time**: Weeks
**Grug time**: 27 minutes

**Original complexity**: High
**Grug complexity**: Low

**Original risk**: Break everything
**Grug risk**: Change small, test each step

Grug brain developer say: **Start simple. Add complexity only when pain is real.**

Ralph already work. Make small improvements. Ship it. See what breaks. Fix that. Repeat.

Complexity is the enemy. Simple is the friend.

Grug done now. Go write code.
