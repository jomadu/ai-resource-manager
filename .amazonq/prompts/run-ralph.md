---
description: Run Ralph autonomous agent for a specific feature
---

# Run Ralph for Feature

Run Ralph autonomous development agent for a specific feature.

See `RALPH-README.md` for complete Ralph documentation.

## Usage

1. Tell me which feature from `specs/` you want to work on
2. I'll use `cp` to copy the PRD to project root and start Ralph

## Example

User: "Run Ralph for version-constraint-interface"

I will:
1. Copy `specs/version-constraint-interface/prd.json` to `./prd.json`
2. Run `./ralph-kiro.sh` (default 10 iterations)
3. Monitor progress in `progress.txt`
