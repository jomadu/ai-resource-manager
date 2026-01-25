# Resource Schemas

ARM uses Kubernetes-style resource definitions for packaging AI rules and prompts.

## Ruleset Schema

```yaml
apiVersion: v1                    # Required
kind: Ruleset                     # Required
metadata:                         # Required
  id: "cleanCode"                 # Required
  name: "Clean Code"              # Optional
  description: "Make clean code"  # Optional
spec:                             # Required
  rules:                          # Required
    ruleOne:                      # Required (rule key)
      name: "Rule 1"              # Optional
      description: "Rule description"  # Optional
      priority: 100               # Optional (default: 100)
      enforcement: must           # Optional (must|should|may)
      scope:                      # Optional
        - files: ["**/*.py"]      # Optional
      body: |                     # Required
        This is the body of the rule.
```

## Promptset Schema

```yaml
apiVersion: v1                    # Required
kind: Promptset                   # Required
metadata:                         # Required
  id: "code-review"               # Required
  name: "Code Review"             # Optional
  description: "Review the code!" # Optional
spec:                             # Required
  prompts:                        # Required
    promptOne:                    # Required (prompt key)
      name: "Prompt One"          # Optional
      description: "Prompt description"  # Optional
      body: |                     # Required
        This is how you review the code.
```