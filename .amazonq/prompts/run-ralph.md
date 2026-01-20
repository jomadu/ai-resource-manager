---
description: Run Ralph autonomous agent for a specific feature
---

# Run Ralph for Feature

Run Ralph autonomous development agent for a specific feature.

## Usage

1. Tell me which feature from `specs/` you want to work on
2. I'll copy the PRD to project root and start Ralph

## What Ralph Does

Ralph is an autonomous coding agent that:
- Reads the PRD (Product Requirements Document) from `prd.json`
- Picks highest priority incomplete user story
- Implements it
- Runs quality checks (typecheck, lint, test)
- Commits changes if checks pass
- Updates PRD and progress log
- Repeats until all stories complete

## Example

User: "Run Ralph for version-constraint-interface"

I will:
1. Copy `specs/version-constraint-interface/prd.json` to `./prd.json`
2. Run `./ralph-kiro.sh` (default 10 iterations)
3. Monitor progress in `progress.txt`

## Options

- Default: 10 iterations
- Custom iterations: specify number (e.g., "run Ralph with 20 iterations")
- Custom agent: specify agent name (e.g., "run Ralph with my-agent")

## Files Created

- `prd.json` - Active task list (gitignored)
- `progress.txt` - Learnings log (gitignored)
- `.last-branch` - Branch tracking (gitignored)

## Monitoring

Check status:
```bash
cat prd.json | jq '.userStories[] | {id, title, passes}'
cat progress.txt
```
