# Cloudsmith Registry

Cloudsmith registries use Cloudsmith's package repository service to store and distribute AI rules as versioned packages.

## Configuration

Add a Cloudsmith registry:

```bash
arm config registry add my-cloudsmith https://app.cloudsmith.com/myorg/ai-rules --type cloudsmith
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
   arm install my-cloudsmith/ai-rules --sinks cursor
   ```

**Note**: Unlike Git registries which use Git's built-in authentication, Cloudsmith registries require explicit token configuration because they use Cloudsmith's Package Registry API.

## Package Structure

Cloudsmith packages are typically single-file artifacts with semantic versioning:

```
Package Name: clean-code.yml
Version: 1.0.0, 1.1.0, 2.0.0, etc.
File: clean-code.yml (URF ruleset)

Package Name: security.yml
Version: 1.0.0, 1.1.0, 2.0.0, etc.
File: security.yml (URF ruleset)
```

**Key Difference**: Unlike GitLab's generic packages which can contain multiple files, Cloudsmith packages typically contain a single file per artifact. This means each ruleset is usually published as a separate package.

## Installing Rules

Install URF rulesets:
```bash
arm install my-cloudsmith/clean-code --sinks cursor,amazonq
```

Install specific version:
```bash
arm install my-cloudsmith/clean-code@1.0.0 --sinks cursor
```

Install multiple rulesets:
```bash
arm install my-cloudsmith/clean-code my-cloudsmith/security --sinks cursor
```

## Publishing Packages

Use Cloudsmith CLI to publish packages:

```bash
# Upload a single URF ruleset
cloudsmith push raw myorg/ai-rules clean-code.yml --version 1.0.0

# Upload with metadata
cloudsmith push raw myorg/ai-rules security.yml \
  --version 1.1.0 \
  --summary "Security rules for AI coding assistants" \
  --description "Comprehensive security guidelines and best practices"
```

## Version Resolution

- **Semantic versions**: `1.0.0`, `^1.0.0`, `~1.1.0`
- **Latest**: `latest` (resolves to highest semantic version)
- Versions are sorted by semantic version in descending order

## Best Practices

1. **Single File Per Package**: Each ruleset should be published as a separate package for better version management
2. **Semantic Versioning**: Use proper semantic versioning for compatibility tracking
3. **Descriptive Names**: Use clear package names that match your ruleset files
4. **Raw Format**: Use Cloudsmith's "raw" package format for maximum flexibility

## Example Workflow

```bash
# Add registry
arm config registry add cloudsmith https://app.cloudsmith.com/myteam/rules --type cloudsmith

# Configure authentication
echo "[registry https://api.cloudsmith.io/myteam/rules]" >> .armrc
echo "token = your-token-here" >> .armrc
chmod 600 .armrc

# Install rulesets
arm install cloudsmith/clean-code cloudsmith/security --sinks cursor,amazonq

# Update to latest versions
arm update
```
