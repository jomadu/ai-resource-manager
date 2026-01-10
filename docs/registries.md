# Registries

Registries are remote sources where rulesets and promptsets are stored and versioned, similar to npm registries. ARM supports:

- **Git registries**: GitHub repositories, GitLab projects, or any Git remote
- **GitLab Package registries**: GitLab's Generic Package Registry for versioned packages
- **Cloudsmith registries**: Cloudsmith's package repository service for single-file artifacts

## Commands

For detailed command usage and examples, see [Registry Management](commands.md#registry-management) in the commands reference.

## Package Structure

**Recommended structure:**
```
clean-code-ruleset.yml      # ARM resource definitions
security-ruleset.yml
code-review-promptset.yml   # ARM resource definitions
build/                      # Pre-compiled rules (optional)
├── cursor/
│   ├── clean-code.mdc
│   └── security.mdc
└── amazonq/
    ├── clean-code.md
    └── security.md
```

This structure works for all registry types, with ARM resource files at the root level and pre-compiled rules organized under `build/` by AI tool. ARM defaults to installing resource files (`*.yml, *.yaml`) when no `--include` patterns are specified.

## Archive Support

ARM automatically extracts and processes **zip** and **tar.gz** archives during installation:

- **Supported formats**: `.zip` and `.tar.gz` files
- **Automatic extraction**: Archives are detected by extension and extracted transparently
- **Merge behavior**: Extracted files are merged with loose files, with archives taking precedence in case of path conflicts
- **Security**: Path sanitization prevents directory traversal attacks
- **Pattern filtering**: `--include` patterns are applied to the merged content after extraction

## Registry Types

### Git Registry
Uses Git repositories with tags and branches for versioning. See [Git Registry](storage/registries/git-registry.md) for details.

### GitLab Registry
Uses GitLab's Generic Package Registry for versioned packages. See [GitLab Registry](storage/registries/gitlab-registry.md) for details.

### Cloudsmith Registry
Uses Cloudsmith's package repository service for single-file artifacts. See [Cloudsmith Registry](storage/registries/cloudsmith-registry.md) for details.

## Examples

**Community registries:**
- [PatrickJS/awesome-cursorrules](https://github.com/PatrickJS/awesome-cursorrules) - Community collection of Cursor rules
- [snarktank/ai-dev-tasks](https://github.com/snarktank/ai-dev-tasks) - AI development task templates
- [steipete/agent-rules](https://github.com/steipete/agent-rules) - Agent configuration rules
