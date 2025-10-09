# Promptsets

Promptsets are collections of AI prompts packaged as versioned units, identified by names like `ai-rules/code-review-promptset` where `ai-rules` is the registry and `code-review-promptset` is the promptset name.

## Commands

For detailed command usage and examples, see [Promptset Management](commands.md#promptset-management) in the commands reference.

## ARM Promptset Format

ARM uses Kubernetes-style resource definitions for promptsets:

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

## Examples

Install from latest version:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands
```

Install from specific version:
```bash
arm install promptset ai-rules/code-review-promptset@v1.0.0 cursor-commands
```

Install to multiple sinks:
```bash
arm install promptset ai-rules/code-review-promptset cursor-commands q-prompts
```

Install specific files:
```bash
arm install promptset --include "review.yml" ai-rules/code-review-promptset cursor-commands
```

## Compilation

Promptsets compile to simple, content-only files across all targets:
- No metadata or frontmatter
- Pure content focus
- Universal compatibility with AI tools
- Currently identical output across all platforms

See [compilation examples](examples/compilation/promptset/) for detailed before/after examples showing how promptset YAML files are transformed into platform-specific formats.
