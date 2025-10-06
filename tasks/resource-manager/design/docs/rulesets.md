# Rulesets

Rulesets are collections of AI rules packaged as versioned units, identified by names like `ai-rules/clean-code-ruleset` where `ai-rules` is the registry and `clean-code-ruleset` is the ruleset name.

## Commands

- `arm install ruleset <registry>/<ruleset>[@version] <sink>...` - Install ruleset to specific sinks
- `arm update ruleset [ruleset...]` - Update to latest compatible versions
- `arm uninstall ruleset <registry>/<ruleset>` - Remove ruleset
- `arm list ruleset` - Show installed rulesets
- `arm info ruleset [ruleset...]` - Show detailed information
- `arm config ruleset set <registry>/<ruleset> <key> <value>` - Update ruleset configuration

## ARM Resource Format

ARM uses Kubernetes-style resource definitions for rulesets:

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

## Examples

Install from latest version:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules
```

Install from specific version:
```bash
arm install ruleset ai-rules/clean-code-ruleset@v1.0.0 cursor-rules
```

Install to multiple sinks:
```bash
arm install ruleset ai-rules/clean-code-ruleset cursor-rules q-rules
```

Install with custom priority:
```bash
arm install ruleset --priority 200 ai-rules/clean-code-ruleset cursor-rules
```

Install specific files:
```bash
arm install ruleset --include "security.yml" ai-rules/clean-code-ruleset cursor-rules
```

## Priority and Conflict Resolution

ARM resource definitions compile to tool-specific formats with embedded metadata for priority resolution and conflict management. The `arm_index.*` file helps AI tools resolve conflicts between rulesets based on priority values.

See [compilation examples](examples/compilation/ruleset/) for detailed before/after examples showing how ruleset YAML files are transformed into platform-specific formats with embedded metadata.
