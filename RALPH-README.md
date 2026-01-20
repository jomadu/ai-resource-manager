# Ralph with Kiro CLI

This project uses Ralph for autonomous development.

## Quick Start

```bash
# Copy a feature's PRD to project root
cp specs/version-constraint-interface/prd.json ./prd.json

# Run Ralph (RALPH.md is the agent instructions)
./ralph-kiro.sh

# With custom max iterations (default is 10)
./ralph-kiro.sh 20

# With specific agent
./ralph-kiro.sh --agent my-agent

# Combined
./ralph-kiro.sh --agent my-agent 20
```

## Files

- **ralph-kiro.sh** - Ralph runner script
- **RALPH.md** - Agent instructions (customized for this project)
- **prd.json** - Active task list (gitignored, runtime)
- **progress.txt** - Learnings log (gitignored, auto-created)
- **.last-branch** - Branch tracking (gitignored, auto-created)

## Monitoring

```bash
# Check task status
cat prd.json | jq '.userStories[] | {id, title, passes}'

# View learnings
cat progress.txt

# Git history
git log --oneline -10
```

## References

- [Geoffrey Huntley's Ralph](https://ghuntley.com/ralph/)
- [snarktank/ralph](https://github.com/snarktank/ralph)
- [Analysis](./specs/RALPH-STRUCTURE-ANALYSIS.md)
