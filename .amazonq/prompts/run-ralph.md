---
description: Run Ralph autonomous agent for a specific feature
---

# Run Ralph for Feature

Run Ralph autonomous development agent for a specific feature.

See `RALPH-README.md` for complete Ralph documentation.

## Usage

1. Tell me which feature from `specs/` you want to work on
2. I'll copy the PRD (and progress.txt if it exists) to project root and start Ralph
3. After Ralph completes, I'll copy `prd.json` and `progress.txt` back to the feature directory

## Example

User: "Run Ralph for version-constraint-interface"

I will:
1. Copy `specs/version-constraint-interface/prd.json` to `./prd.json`
2. Copy `specs/version-constraint-interface/progress.txt` to `./progress.txt` (if it exists)
3. Run `./ralph-kiro.sh` (default 10 iterations)
4. Monitor progress in `progress.txt`
5. Copy `prd.json` and `progress.txt` back to `specs/version-constraint-interface/`
