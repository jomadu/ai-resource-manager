# Cloudsmith Registry

Cloudsmith registries use Cloudsmith's package repository service to store and distribute AI rules as versioned packages.

## Configuration

Add a Cloudsmith registry:

```bash
arm add registry --type cloudsmith my-cloudsmith https://app.cloudsmith.com/myorg/ai-rules
```

## Authentication

Cloudsmith registries require explicit token authentication in `.armrc`. See the [.armrc documentation](../armrc.md) for complete details.

**Quick Setup:**

1. Create `.armrc` file:
   ```ini
   [registry https://api.cloudsmith.io/myorg/ai-rules]
   token = your-cloudsmith-api-token
   ```

2. Set file permissions:
   ```bash
   chmod 600 .armrc
   ```

3. Test installation:
   ```bash
   arm install ruleset my-cloudsmith/clean-code-ruleset cursor-rules
   ```

**Note**: Unlike Git registries which use Git's built-in authentication, Cloudsmith registries require explicit token configuration because they use Cloudsmith's Package Registry API.

## Package Structure

Cloudsmith packages are typically single-file artifacts with semantic versioning:

```
Package Name: clean-code-ruleset.yml
Version: 1.0.0, 1.1.0, 2.0.0, etc.
File: clean-code-ruleset.yml (ARM resource definition)

Package Name: security-ruleset.yml
Version: 1.0.0, 1.1.0, 2.0.0, etc.
File: security-ruleset.yml (ARM resource definition)

Package Name: code-review-promptset.yml
Version: 1.0.0, 1.1.0, 2.0.0, etc.
File: code-review-promptset.yml (ARM resource definition)
```

**Key Difference**: Unlike GitLab's generic packages which can contain multiple files, Cloudsmith packages typically contain a single file per artifact. This means each ruleset or promptset is usually published as a separate package.

## Installing Rulesets

Install from latest version:
```bash
arm install ruleset my-cloudsmith/clean-code-ruleset cursor-rules
```

Install from specific version:
```bash
arm install ruleset my-cloudsmith/clean-code-ruleset@1.0.0 cursor-rules
```

Install to multiple sinks:
```bash
arm install ruleset my-cloudsmith/clean-code-ruleset cursor-rules q-rules
```

Install with custom priority:
```bash
arm install ruleset --priority 200 my-cloudsmith/clean-code-ruleset cursor-rules
```

Install multiple rulesets:
```bash
arm install ruleset my-cloudsmith/clean-code-ruleset my-cloudsmith/security-ruleset cursor-rules
```

## Installing Promptsets

Install from latest version:
```bash
arm install promptset my-cloudsmith/code-review-promptset cursor-prompts
```

Install from specific version:
```bash
arm install promptset my-cloudsmith/code-review-promptset@1.0.0 cursor-prompts
```

Install to multiple sinks:
```bash
arm install promptset my-cloudsmith/code-review-promptset cursor-prompts q-prompts
```

Install multiple promptsets:
```bash
arm install promptset my-cloudsmith/code-review-promptset my-cloudsmith/testing-promptset cursor-prompts
```

## Version Resolution

Cloudsmith registries support semantic versioning:

- **Semantic versions**: `1.0.0`, `^1.0.0`, `~1.1.0`
- **Latest**: `latest` (resolves to highest semantic version)
- Versions are sorted by semantic version in descending order

## Publishing Packages

Use Cloudsmith CLI to publish packages:

```bash
# Upload a single ARM ruleset
cloudsmith push raw myorg/ai-rules clean-code-ruleset.yml --version 1.0.0

# Upload a single ARM promptset
cloudsmith push raw myorg/ai-rules code-review-promptset.yml --version 1.0.0

# Upload with metadata
cloudsmith push raw myorg/ai-rules security-ruleset.yml \
  --version 1.1.0 \
  --summary "Security rules for AI coding assistants" \
  --description "Comprehensive security guidelines and best practices"
```

## Best Practices

1. **Single File Per Package**: Each ruleset or promptset should be published as a separate package for better version management
2. **Semantic Versioning**: Use proper semantic versioning for compatibility tracking
3. **Descriptive Names**: Use clear package names that match your resource files
4. **Raw Format**: Use Cloudsmith's "raw" package format for maximum flexibility

## Management Commands

List all registries:
```bash
arm list registry
```

Show registry information:
```bash
arm info registry my-cloudsmith
```

Update registry configuration:
```bash
arm config registry set my-cloudsmith url https://app.cloudsmith.com/myorg/new-rules
```

Remove registry:
```bash
arm remove registry my-cloudsmith
```

## Example Workflow

```bash
# Add registry
arm add registry --type cloudsmith cloudsmith https://app.cloudsmith.com/myteam/rules

# Configure authentication
echo "[registry https://api.cloudsmith.io/myteam/rules]" >> .armrc
echo "token = your-token-here" >> .armrc
chmod 600 .armrc

# Install rulesets and promptsets
arm install ruleset cloudsmith/clean-code-ruleset cloudsmith/security-ruleset cursor-rules
arm install promptset cloudsmith/code-review-promptset cursor-prompts

# Update to latest versions
arm update ruleset
arm update promptset
```
