# User Documentation

## Job to be Done
Provide comprehensive, accessible documentation that helps users understand, install, configure, and use ARM effectively.

## Activities
1. Explain ARM's purpose and value proposition
2. Guide users through installation
3. Document core concepts and terminology
4. Provide complete command reference
5. Explain registry types and configuration
6. Document sink configuration and compilation
7. Guide users through publishing their own resources
8. Provide migration guides for breaking changes

## Acceptance Criteria
- [x] README.md explains what ARM is, why it exists, and key features
- [x] README.md provides installation instructions (quick install, specific version, manual)
- [x] README.md shows quick start examples for common workflows
- [x] README.md links to detailed documentation
- [x] concepts.md explains core terminology (rulesets, promptsets, registries, sinks, manifest, lock file)
- [x] commands.md provides complete reference for all CLI commands with examples
- [x] registries.md explains registry types and management
- [x] git-registry.md documents Git-based registry usage
- [x] gitlab-registry.md documents GitLab Package Registry usage
- [x] cloudsmith-registry.md documents Cloudsmith registry usage
- [x] sinks.md explains sink configuration and tool-specific formats
- [x] storage.md documents cache structure and management
- [x] resource-schemas.md defines ARM resource YAML format
- [x] publishing-guide.md guides users through creating and publishing resources
- [x] migration-v2-to-v3.md provides upgrade path from v2 to v3
- [x] armrc.md documents .armrc authentication configuration

## Data Structures

### Documentation Structure
```
docs/
├── concepts.md           # Core terminology and concepts
├── commands.md           # Complete CLI reference
├── registries.md         # Registry overview
├── git-registry.md       # Git registry details
├── gitlab-registry.md    # GitLab registry details
├── cloudsmith-registry.md # Cloudsmith registry details
├── sinks.md              # Sink configuration
├── storage.md            # Cache management
├── resource-schemas.md   # YAML format
├── publishing-guide.md   # Publishing resources
├── migration-v2-to-v3.md # Upgrade guide
├── armrc.md              # Authentication
└── examples/             # Example files
```

### README.md Sections
```
1. What is ARM?
2. Why ARM?
3. Key Features
4. Installation (Quick, Specific Version, Manual, Verify)
5. Uninstall
6. Quick Start (Resource Types, Setup, Install, Compile)
7. Documentation (links to detailed docs)
8. Buy Me a Coffee
```

## Algorithm

### Documentation Workflow
1. **User discovers ARM:**
   - Reads README.md for overview
   - Understands value proposition
   - Sees key features

2. **User installs ARM:**
   - Follows installation instructions
   - Verifies installation
   - Runs `arm help`

3. **User learns concepts:**
   - Reads concepts.md for terminology
   - Understands rulesets, promptsets, registries, sinks
   - Understands manifest and lock files

4. **User configures ARM:**
   - Reads registries.md for registry types
   - Reads specific registry docs (git, gitlab, cloudsmith)
   - Adds registries with `arm add registry`
   - Reads sinks.md for sink configuration
   - Adds sinks with `arm add sink`

5. **User installs resources:**
   - Reads commands.md for install command
   - Installs rulesets/promptsets
   - Checks installed packages with `arm list`

6. **User publishes resources:**
   - Reads publishing-guide.md
   - Creates ARM resource YAML files
   - Publishes to Git repository or package registry

7. **User upgrades ARM:**
   - Reads migration guide for breaking changes
   - Follows upgrade steps
   - Migrates configuration if needed

## Edge Cases

| Condition | Expected Behavior |
|-----------|-------------------|
| User unfamiliar with package managers | README explains concepts clearly |
| User needs specific command syntax | commands.md provides complete reference |
| User needs authentication | armrc.md explains token configuration |
| User migrating from v2 | migration-v2-to-v3.md provides step-by-step guide |
| User wants to publish | publishing-guide.md walks through process |
| User confused by terminology | concepts.md defines all terms |

## Dependencies

- Markdown rendering (GitHub, documentation sites)
- Example files in docs/examples/

## Implementation Mapping

**Source files:**
- `README.md` - Main project documentation
- `docs/concepts.md` - Core concepts (5260 bytes)
- `docs/commands.md` - CLI reference (27224 bytes)
- `docs/registries.md` - Registry overview (4416 bytes)
- `docs/git-registry.md` - Git registry (4560 bytes)
- `docs/gitlab-registry.md` - GitLab registry (3286 bytes)
- `docs/cloudsmith-registry.md` - Cloudsmith registry (3356 bytes)
- `docs/sinks.md` - Sink configuration (4021 bytes)
- `docs/storage.md` - Cache management (3182 bytes)
- `docs/resource-schemas.md` - YAML format (1576 bytes)
- `docs/publishing-guide.md` - Publishing guide (11087 bytes)
- `docs/migration-v2-to-v3.md` - Migration guide (11500 bytes)
- `docs/armrc.md` - Authentication (3452 bytes)
- `docs/examples/` - Example files

**Related specs:**
- `installation-scripts.md` - Scripts referenced in README
- `build-system.md` - Build process for contributors
- `code-quality.md` - Code standards for contributors

## Examples

### Example 1: New User Journey

**Input:**
```
User discovers ARM on GitHub
```

**Expected Output:**
```
1. Reads README.md
   - Understands ARM is a dependency manager for AI resources
   - Sees key features (versioning, reproducibility, multi-tool support)
   - Finds installation command

2. Installs ARM
   - Runs: curl -fsSL https://raw.githubusercontent.com/.../install.sh | bash
   - Verifies: arm version

3. Reads Quick Start
   - Adds registry: arm add registry git --url https://github.com/... my-registry
   - Adds sink: arm add sink --tool cursor cursor-rules .cursor/rules
   - Installs ruleset: arm install ruleset my-registry/clean-code cursor-rules

4. Explores detailed docs
   - Reads concepts.md for terminology
   - Reads commands.md for full CLI reference
   - Reads registries.md for registry options
```

**Verification:**
- User successfully installs ARM
- User understands core concepts
- User can install and use resources

### Example 2: Publishing Resources

**Input:**
```
User wants to publish their own ruleset
```

**Expected Output:**
```
1. Reads publishing-guide.md
   - Understands ARM resource YAML format
   - Sees example ruleset structure
   - Learns about Git-based publishing

2. Creates ruleset
   - Creates my-ruleset.yml with ARM schema
   - Defines rules with priorities
   - Tests locally with: arm compile ruleset my-ruleset.yml --tool cursor --output .cursor/rules/

3. Publishes to Git
   - Creates GitHub repository
   - Pushes ruleset files
   - Tags with semantic version: v1.0.0

4. Shares with others
   - Others add registry: arm add registry git --url https://github.com/user/repo my-registry
   - Others install: arm install ruleset my-registry/my-ruleset cursor-rules
```

**Verification:**
- User creates valid ARM resource
- User publishes to Git repository
- Others can install and use resource

### Example 3: Migrating from v2

**Input:**
```
User has ARM v2 installed and wants to upgrade to v3
```

**Expected Output:**
```
1. Reads migration-v2-to-v3.md
   - Understands breaking changes
   - Sees "nuke and pave" recommendation
   - Follows upgrade steps

2. Backs up configuration
   - Saves .armrc if exists
   - Notes installed packages

3. Uninstalls v2
   - Runs uninstall script
   - Removes ~/.arm/ directory

4. Installs v3
   - Runs install script
   - Verifies: arm version

5. Reconfigures
   - Adds registries again
   - Adds sinks again
   - Reinstalls packages
```

**Verification:**
- User successfully upgrades to v3
- User understands breaking changes
- User reconfigures ARM for v3

## Notes

- README.md is user-facing, focuses on getting started quickly
- docs/ provides detailed reference documentation
- Examples in docs/examples/ show real-world usage
- Migration guides document breaking changes between versions
- Documentation uses consistent terminology from concepts.md
- Commands.md is comprehensive (27KB) with examples for every command
- Publishing guide is detailed (11KB) with step-by-step instructions

## Known Issues

None - all documentation complete and accurate.

## Areas for Improvement

- Add video tutorials for common workflows
- Add interactive examples (asciinema recordings)
- Add FAQ section for common questions
- Add troubleshooting guide for common errors
- Add architecture diagrams for visual learners
- Add API documentation if programmatic usage added
- Add comparison table with other tools (npm, pip, etc.)
- Add case studies from real users
- Add search functionality for documentation
- Consider hosting documentation on dedicated site (GitHub Pages, Read the Docs)
