# Cloudsmith Registry

Cloudsmith registries use Cloudsmith's package repository service to store and distribute packages as versioned packages.

## Authentication

Cloudsmith registries require explicit token authentication in `.armrc`. See the [.armrc documentation](armrc.md) for complete details.

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

**Key Concept**: Cloudsmith packages have explicit names defined in the registry. You must use the exact package name when installing. Packages are typically single-file artifacts.

```
Package: clean-code-ruleset.yml (exact name in registry)
Contents: clean-code-ruleset.yml (ARM resource definition)

Package: security-ruleset.yml (exact name in registry)
Contents: security-ruleset.yml (ARM resource definition)
```

**Install examples**:
```bash
# Must use exact package names from registry
arm install ruleset cloudsmith-registry/clean-code-ruleset.yml cursor-rules
arm install ruleset cloudsmith-registry/security-ruleset.yml q-rules

# Use --include/--exclude to filter files (default: *.yml, *.yaml)
arm install ruleset --include "**/*.yml" cloudsmith-registry/clean-code-ruleset.yml cursor-rules
```

**Key Difference**: Unlike GitLab's generic packages which can contain multiple files, Cloudsmith packages typically contain a single file per artifact.

## Version Resolution

Cloudsmith uses semantic versioning only.

```bash
# Install specific version
arm install ruleset cloudsmith-registry/clean-code-ruleset.yml@1.0.0 cursor-rules

# Install with version constraints
arm install ruleset cloudsmith-registry/clean-code-ruleset.yml@1 cursor-rules    # >= 1.0.0, < 2.0.0
arm install ruleset cloudsmith-registry/clean-code-ruleset.yml@1.1 cursor-rules  # >= 1.1.0, < 1.2.0

# Install latest version
arm install ruleset cloudsmith-registry/clean-code-ruleset.yml cursor-rules
```

**Not supported**: Branch references (`@main`), commit hashes (`@a1b2c3d`), or non-semantic tags (`@stable`).

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