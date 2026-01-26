# Publishing Rulesets to Git Repositories

This guide shows you how to create and publish a ruleset (like grug brained development) to a Git repository for use with ARM.

## Quick Start

1. Create a Git repository
2. Add ruleset YAML files
3. Tag with semantic versions
4. Users install from your repository

## Repository Structure

```
grug-brained-dev/
├── grug-brained-ruleset.yml    # Your ruleset definition
└── README.md                   # Documentation
```

That's it. ARM will find and use any `.yml` or `.yaml` files in your repository.

## Creating a Ruleset

Create a file named `grug-brained-ruleset.yml`:

```yaml
apiVersion: v1
kind: Ruleset
metadata:
  id: "grugBrained"
  name: "Grug Brained Development"
  description: "Simple, practical development principles for grug brain developer"
spec:
  rules:
    simpleCode:
      name: "Keep Code Simple"
      description: "Complexity very, very bad. Simple good."
      priority: 100
      enforcement: must
      scope:
        - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.go"]
      body: |
        grug brain developer not smart enough for complex code. 
        grug keep code simple so grug can understand later.
        
        - No clever tricks
        - No fancy abstractions unless absolutely needed
        - If grug not understand in 5 minutes, too complex
        
    avoidAbstraction:
      name: "Fear Complexity Spirit"
      description: "Abstraction is complexity spirit. Complexity spirit bad."
      priority: 95
      enforcement: must
      scope:
        - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.go"]
      body: |
        Complexity spirit sneak into code through abstraction.
        Abstraction make code hard to understand.
        
        grug say: only abstract when pain of duplication greater than 
        pain of wrong abstraction. Usually take 3 times before abstract.
        
    testingGood:
      name: "Testing Make Grug Confident"
      description: "Tests catch bugs before they catch grug"
      priority: 90
      enforcement: should
      scope:
        - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.go"]
      body: |
        grug like tests. Tests make grug feel safe to change code.
        grug not need 100% coverage. grug need tests for important parts.
        
        - Test the scary parts
        - Test the parts that break often
        - Test the parts that make money
        
    smallChanges:
      name: "Small Changes Good"
      description: "Big changes scary. Small changes safe."
      priority: 85
      enforcement: should
      scope:
        - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.go"]
      body: |
        grug make small changes. Small changes easy to understand.
        Small changes easy to review. Small changes easy to revert.
        
        Big rewrite very scary. Usually fail. grug learn this many times.
        
    sayNoToMicroservices:
      name: "Microservices Probably Not Needed"
      description: "Monolith good for most grug projects"
      priority: 80
      enforcement: may
      scope:
        - files: ["**/*.py", "**/*.js", "**/*.ts", "**/*.go"]
      body: |
        Many grug hear microservices solve all problems.
        This not true. Microservices create many new problems.
        
        Monolith good until it not good. Usually take longer than grug think.
        Start with monolith. Split later if really needed.
```

## Schema Reference

### Required Fields

```yaml
apiVersion: v1              # Always "v1"
kind: Ruleset               # Always "Ruleset"
metadata:
  id: "uniqueId"            # Unique identifier (camelCase)
spec:
  rules:
    ruleKey:                # Unique key for this rule
      body: |               # The actual rule content
        Rule text here
```

### Optional Fields

```yaml
metadata:
  name: "Display Name"                    # Human-readable name
  description: "What this ruleset does"   # Brief description

spec:
  rules:
    ruleKey:
      name: "Rule Name"                   # Human-readable rule name
      description: "What this rule does"  # Brief description
      priority: 100                       # Default: 100 (higher = more important)
      enforcement: must                   # must|should|may (default: should)
      scope:
        - files: ["**/*.py", "**/*.js"]   # File patterns this rule applies to
```

### Priority System

- Higher priority rules override lower priority rules when conflicts occur
- Default priority: 100
- Recommended ranges:
  - 200+: Team/organization standards (highest priority)
  - 100-199: General best practices
  - 1-99: Suggestions and preferences

### Enforcement Levels

- `must`: Critical rules that should always be followed
- `should`: Important rules with occasional exceptions
- `may`: Suggestions and preferences

## Versioning Your Ruleset

### Using Semantic Version Tags

Tag your releases with semantic versions:

```bash
# Initial release
git tag v1.0.0
git push origin v1.0.0

# Bug fix (typo in rule text)
git tag v1.0.1
git push origin v1.0.1

# New rule added (backwards compatible)
git tag v1.1.0
git push origin v1.1.0

# Breaking change (rule removed or significantly changed)
git tag v2.0.0
git push origin v2.0.0
```

### Using Branches

For development and testing:

```bash
# Development branch
git checkout -b develop
git push origin develop

# Feature branch
git checkout -b feature/new-rules
git push origin feature/new-rules
```

## How Users Install Your Ruleset

### Add Your Repository as a Registry

```bash
arm add registry git --url https://github.com/yourusername/grug-brained-dev grug-rules
```

### Install Specific Version

```bash
# Install latest semantic version
arm install ruleset grug-rules/grug-brained cursor-rules

# Install specific version
arm install ruleset grug-rules/grug-brained@1.0.0 cursor-rules

# Install version range
arm install ruleset grug-rules/grug-brained@1 cursor-rules      # >= 1.0.0, < 2.0.0
arm install ruleset grug-rules/grug-brained@1.1 cursor-rules    # >= 1.1.0, < 1.2.0

# Install from branch
arm install ruleset grug-rules/grug-brained@develop cursor-rules
```

### Install with Priority

```bash
# Install with high priority (overrides other rules)
arm install ruleset --priority 200 grug-rules/grug-brained cursor-rules

# Install with default priority (100)
arm install ruleset grug-rules/grug-brained cursor-rules
```

## Advanced: Multiple Rulesets in One Repository

You can publish multiple rulesets in a single repository:

```
grug-brained-dev/
├── grug-brained-ruleset.yml
├── grug-testing-ruleset.yml
├── grug-api-design-ruleset.yml
└── README.md
```

Users can install specific files:

```bash
# Install only testing rules
arm install ruleset --include "grug-testing-*.yml" grug-rules/testing cursor-rules

# Install everything except API rules
arm install ruleset --exclude "grug-api-*.yml" grug-rules/all-but-api cursor-rules
```

## Advanced: Pre-compiled Tool Files

You can include pre-compiled files for specific tools alongside ARM resources:

```
grug-brained-dev/
├── grug-brained-ruleset.yml    # ARM resource (cross-platform)
└── build/                      # Pre-compiled (optional)
    ├── cursor/
    │   └── grug-brained.mdc
    └── amazonq/
        └── grug-brained.md
```

ARM will use pre-compiled files if they exist, otherwise compile from YAML.

## Advanced: Archive Support

You can distribute rulesets as archives:

```
grug-brained-dev/
├── grug-brained-ruleset.yml
└── rules.tar.gz                # Contains additional YAML files
```

ARM automatically extracts `.zip` and `.tar.gz` files during installation.

## Version Resolution Priority

When users install without specifying a version:

1. **Semantic version tags** - Always prioritized (highest version wins)
2. **Branches** - Only used if no semver tags exist (first branch in registry config wins)

Example:
```bash
# Your repo has: v1.0.0, v1.1.0, v2.0.0, main branch
# User runs: arm install ruleset grug-rules/grug-brained cursor-rules
# Result: Installs v2.0.0 (highest semver tag)
```

## Authentication

ARM uses Git's built-in authentication:

- **Public repositories**: No authentication needed
- **Private repositories**: 
  - SSH: Use `ssh-agent` or `~/.ssh/config`
  - HTTPS: Use Git credential helpers
  - GitHub: Use `gh auth login`
  - GitLab: Use `glab auth login`

## Best Practices

1. **Start simple**: One ruleset file is enough
2. **Use semantic versioning**: Tag releases with v1.0.0, v1.1.0, etc.
3. **Write clear rule bodies**: Users will see this in their AI tools
4. **Set appropriate priorities**: Higher for critical rules, lower for suggestions
5. **Document your ruleset**: Add a README explaining the philosophy
6. **Test before tagging**: Install from your branch first
7. **Keep rules focused**: One clear principle per rule
8. **Use enforcement levels**: `must` for critical, `should` for important, `may` for suggestions

## Testing Your Ruleset

Before publishing, test locally:

```bash
# Add your local repository
arm add registry git --url /path/to/your/grug-brained-dev local-grug

# Install from your branch
arm install ruleset local-grug/grug-brained@main cursor-rules

# Check the compiled output
cat .cursor/rules/grug-brained.mdc
```

## Example: Complete Workflow

```bash
# 1. Create repository
mkdir grug-brained-dev
cd grug-brained-dev
git init

# 2. Create ruleset
cat > grug-brained-ruleset.yml << 'EOF'
apiVersion: v1
kind: Ruleset
metadata:
  id: "grugBrained"
  name: "Grug Brained Development"
  description: "Simple, practical development principles"
spec:
  rules:
    simpleCode:
      name: "Keep Code Simple"
      priority: 100
      enforcement: must
      body: |
        grug brain developer not smart enough for complex code.
        grug keep code simple so grug can understand later.
EOF

# 3. Commit and tag
git add .
git commit -m "Initial release"
git tag v1.0.0

# 4. Push to GitHub
git remote add origin https://github.com/yourusername/grug-brained-dev
git push origin main
git push origin v1.0.0

# 5. Users can now install
arm add registry git --url https://github.com/yourusername/grug-brained-dev grug-rules
arm install ruleset grug-rules/grug-brained cursor-rules
```

## Troubleshooting

**Users can't find my ruleset**
- Ensure your YAML files have `.yml` or `.yaml` extension
- Check that `apiVersion: v1` and `kind: Ruleset` are present
- Verify the file is committed and pushed

**Version not found**
- Ensure tags are pushed: `git push origin v1.0.0`
- Check tag format: Use `v1.0.0` not `1.0.0`
- Verify tag exists: `git tag -l`

**Authentication errors**
- For private repos, ensure Git authentication is configured
- Test with: `git clone <your-repo-url>`
- ARM uses the same authentication as Git

## Resources

- [ARM Resource Schemas](docs/resource-schemas.md) - Detailed YAML format
- [Git Registry Documentation](docs/git-registry.md) - Complete Git registry reference
- [Example Rulesets](docs/examples/demo/registry/) - Sample implementations
