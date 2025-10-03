## ARM Resource Schema (v1)

### Ruleset

```yml
apiVersion: v1
kind: Ruleset
metadata:
    id: "clean-code"
    name: "Clean Code"
    description: "Make clean code"
spec:
    rules:
        ruleOne:
            name: "Rule 1"
            description: "Rule with optional description, priority, enforcement (must|should|may), scope[0].files attributes defined"
            priority: 100
            enforcement: must
            scope:
                - files: ["**/*.py"]
            body: |
                This is the body of the rule.
        ruleTwo:
            name: "Rule 2"
            description: "Rule with optional description and enforcement (must|should|may) defined"
            enforcement: should
            body: |
                This is the body of the rule.
        ruleThree:
            name: "Rule 3"
            description: "Rule with optional description and enforcement (must|should|may) defined"
            enforcement: may
            body: |
                This is the body of the rule.
        ruleFour:
            name: "Rule 4"
            description: "Rule with optional description defined"
            body: |
                This is the body of the rule.
        ruleFive:
            name: "Rule with no description"
            body: |
                This is the body of the rule.
```

### Promptset

```yml
apiVersion: v1
kind: Promptset
metadata:
    id: "code-review"
    name: "Code Review"
    description: "Review the code!"
spec:
    prompts:
        promptOne:
            name: "Prompt One"
            description: "Prompt with optional description defined"
            body: |
                This is how you review the code.
        promptTwo:
            name: "Prompt Two"
            description: "Prompt with optional description defined"
            body: |
                This is how you review the code.
```

### Future Resources

**Workflow** - Multi-step AI processes with dependencies and state passing between prompts.

**Template** - Parameterized code scaffolding that generates multiple files from user inputs.

## ARM

### Cache

```
~/.arm/cache/
    registries/
        <registry-key>/
            index.json
            repository/ # if a git type registry
                .git/
                <repository files>
                ...
            packages/
                <package-key>/
                    <version>/
                        <package files>
                        ...
```

~/.arm/cache/registries/<registry-key>/index.json

```
{
    "registry_metadata": {
        <properties that make up the registry key>
        ...
    },
    "created_on": "2025-09-08T23:10:43.984784Z",
    "last_updated_on": "2025-09-08T23:10:43.984784Z",
    "last_accessed_on": "2025-09-08T23:10:43.984784Z",
    "packages": {
        "<package-key>": {
            "package_metadata": {
                <properties that make up the package key>
                ...
            },
            "created_on": "2025-09-08T23:10:43.984784Z",
            "last_updated_on": "2025-09-08T23:10:43.984784Z",
            "last_accessed_on": "2025-09-08T23:10:43.984784Z",
            "versions": {
                "<version>": {
                    "created_on": "2025-09-08T23:10:43.984784Z",
                    "last_updated_on": "2025-09-08T23:10:43.984784Z",
                    "last_accessed_on": "2025-09-08T23:10:43.984784Z",
                }
            }
        }

    }

}
```

### Installation to Sink

```
<sink-dir>/arm/
    arm-index.json
    arm_index.mdc # if packages are installed as rulesets
    <registry-name>
        <package-name>
            <version>
                <package-files>
                ...

<project-dir>/
    arm-lock.json
    arm.json
```

<project-dir>/arm.json

```
{
    "registries": {
        "<registry-name>": {
            "url": "<registry url>",
            "type": "<registry type (git|gitlab|cloudsmith)>
            # additional Registry config
            ...
        }
    },
    "packages": {
        "rulesets": {
            "<registry-name>": {
                "<package-name>": {
                    "version": "<version constraint>",
                    "priority": <priority 0-1000+>,
                    "include": [...],
                    "exclude": [...],
                    "sinks": [
                        "cursor",
                        "q"
                    ]
                }
            }
        },
        "promptsets": {
            "<registry-name>": {
                "<package-name>": {
                    "version": "<version constraint>",
                    "include": [...],
                    "exclude": [...],
                    "sinks": [
                        "cursor",
                        "q"
                    ]
                }
            }
        }
    },
    "sinks": {
        "cursor": {
            "directory": ".cursor/rules",
            "layout": "hierarchical",
            "compileTarget": "cursor"
        },
        "q": {
            "directory": ".amazonq/rules",
            "layout": "hierarchical",
            "compileTarget": "amazonq"
        }
    }
}
```

<project-dir>/arm-lock.json

```
{
    "rulesets": {
        "<registry-name>": {
            "<package-name>": {
                "version": "<resolved version>",
                "display": "<display name of resolved version>",
                "checksum": "sha256:......."
            }
        }
    },
    "prompsets": {
        "<registry-name>": {
            "<package-name>": {
                "version": "<resolved version>",
                "display": "<display name of resolved version>",
                "checksum": "sha256:......."
            }
        }
    }
}
```

<sink-dir>/arm/arm-index.json

```
{
    "rulesets": {
        "<registry-name>": {
            "<package-name>": {
                "version": "<resolved version>",
                "priority": <priority 1-1000+>,
                "file_paths": [
                    "arm/<registry-name>/<package-name>/<version>/<package-file>,
                    ...
                ]
            }
        }
    },
    "promptsets": {
        "<registry-name>": {
            "<package-name>": {
                "version": "<resolved version>",
                "file_paths": [
                    "arm/<registry-name>/<package-name>/<version>/<package-file>,
                    ...
                ]
            }
        }
    }
}
```

<sink-dir>/arm/arm_index.*

```
# ARM Rulesets

This file defines the installation priorities for rulesets managed by ARM.

## Priority Rules

**This index is the authoritative source of truth for ruleset priorities.** When conflicts arise between rulesets, follow this priority order:

1. **Higher priority numbers take precedence** over lower priority numbers
2. **Rules from higher priority rulesets override** conflicting rules from lower priority rulesets
3. **Always consult this index** to resolve any ambiguity about which rules to follow

## Installed Rulesets

### sample-repo/grug-brained-dev@v2.1.0
- **Priority:** 100
- **Rules:**
  - arm/sample-repo/grug-brained-dev/v2.1.0/ai-rules-sample_grug-testing.mdc
  - arm/sample-repo/grug-brained-dev/v2.1.0/ai-rules-sample_grug-documentation.mdc
  - arm/sample-repo/grug-brained-dev/v2.1.0/ai-rules-sample_grug-simplicity.mdc
```
