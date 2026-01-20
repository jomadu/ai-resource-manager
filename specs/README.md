# Specifications

Feature specifications and design documents.

## Structure

Each feature has its own directory:

```
specs/
└── [feature-name]/
    ├── spec.md          # Design document
    └── prd.json         # Ralph task list
```

## Using Ralph

To implement a feature with Ralph:

```bash
# Copy PRD to project root
cp specs/[feature-name]/prd.json ./prd.json

# Run Ralph (RALPH.md is the agent instructions)
./ralph-kiro.sh
```

## Existing Features

- [version-constraint-interface](./version-constraint-interface/) - Version and Constraint type refactoring
