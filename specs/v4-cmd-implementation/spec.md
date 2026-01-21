# PRD: V4 CMD Implementation

## 1. Introduction

Complete rewrite of the ARM command-line interface for v4. This implementation will create a fresh, simplified command structure based on the commands.md specification, with improved user experience and cleaner output formatting.

**Problem Statement:** The current cmd implementation needs a complete rewrite to establish a solid foundation for v4, with simplified architecture and better user feedback.

## 2. Goals

1. Implement all commands specified in commands.md with clean, simple code
2. Simplify command structure and improve output formatting
3. Create testable command handlers with unit tests where reasonable
4. Establish clear patterns for future command additions
5. Optimize for readability and maintainability over clever abstractions

## 3. User Stories

### US-001: Core Commands
**Description:** As a user, I want basic version and help commands so that I can verify installation and discover available commands.

**Acceptance Criteria:**
- [ ] `arm version` displays version, build-id, build-timestamp, build-platform
- [ ] `arm help` displays main help with all commands
- [ ] `arm help <command>` displays command-specific help
- [ ] Typecheck passes
- [ ] Tests pass

### US-002: Registry Management - Add Commands
**Description:** As a user, I want to add different registry types so that I can connect to various package sources.

**Acceptance Criteria:**
- [ ] `arm add registry git` accepts --url, --branches, --force, NAME
- [ ] `arm add registry gitlab` accepts --url, --group-id, --project-id, --api-version, --force, NAME
- [ ] `arm add registry cloudsmith` accepts --url, --owner, --repo, --force, NAME
- [ ] Commands validate required parameters
- [ ] Commands update config file correctly
- [ ] Typecheck passes
- [ ] Tests pass

### US-003: Registry Management - Other Commands
**Description:** As a user, I want to manage, view, and configure registries so that I can maintain my registry connections.

**Acceptance Criteria:**
- [ ] `arm remove registry NAME` removes registry from config
- [ ] `arm set registry NAME KEY VALUE` updates registry config
- [ ] `arm list registry` shows simple list of registry names
- [ ] `arm info registry [NAME]...` shows detailed registry info
- [ ] Typecheck passes
- [ ] Tests pass

### US-004: Sink Management - Add and Remove
**Description:** As a user, I want to add and remove sinks so that I can define output destinations for my resources.

**Acceptance Criteria:**
- [ ] `arm add sink` accepts --tool, --force, NAME, PATH
- [ ] Tool flag accepts: cursor, copilot, amazonq, markdown
- [ ] `arm remove sink NAME` removes sink from config
- [ ] Commands update config file correctly
- [ ] Typecheck passes
- [ ] Tests pass

### US-005: Sink Management - Configure and View
**Description:** As a user, I want to view and configure sinks so that I can manage my output destinations.

**Acceptance Criteria:**
- [ ] `arm set sink NAME KEY VALUE` updates sink config (tool, directory)
- [ ] `arm list sink` shows simple list of sink names
- [ ] `arm info sink [NAME]...` shows detailed sink info
- [ ] Typecheck passes
- [ ] Tests pass

### US-006: Install Ruleset
**Description:** As a user, I want to install rulesets to sinks so that I can use AI rules in my projects.

**Acceptance Criteria:**
- [ ] `arm install ruleset` accepts REGISTRY/RULESET[@VERSION] SINK...
- [ ] Accepts --priority flag (default: 100)
- [ ] Accepts --include and --exclude glob patterns
- [ ] Supports version constraints (@1, @1.1, @1.0.0)
- [ ] Uninstalls from previous sinks not in new sink list
- [ ] Typecheck passes
- [ ] Tests pass

### US-007: Install Promptset
**Description:** As a user, I want to install promptsets to sinks so that I can use AI prompts in my projects.

**Acceptance Criteria:**
- [ ] `arm install promptset` accepts REGISTRY/PROMPTSET[@VERSION] SINK...
- [ ] Accepts --include and --exclude glob patterns
- [ ] Supports version constraints (@1, @1.1, @1.0.0)
- [ ] Uninstalls from previous sinks not in new sink list
- [ ] Typecheck passes
- [ ] Tests pass

### US-008: Install All and Uninstall
**Description:** As a user, I want to install all configured dependencies and uninstall packages so that I can manage my entire dependency tree.

**Acceptance Criteria:**
- [ ] `arm install` installs all configured dependencies
- [ ] `arm uninstall` removes all packages from sinks
- [ ] Commands preserve ARM configuration
- [ ] Typecheck passes
- [ ] Tests pass

### US-009: Update and Upgrade
**Description:** As a user, I want to update packages within constraints and upgrade to latest versions so that I can keep dependencies current.

**Acceptance Criteria:**
- [ ] `arm update` updates packages within version constraints
- [ ] `arm upgrade` upgrades to latest versions ignoring constraints
- [ ] `arm upgrade` updates version constraint to major constraint (^X.0.0)
- [ ] Typecheck passes
- [ ] Tests pass

### US-010: List and Info
**Description:** As a user, I want to view all configured entities so that I can understand my ARM environment.

**Acceptance Criteria:**
- [ ] `arm list` shows registries, sinks, and dependencies grouped
- [ ] `arm info` shows detailed info for all entities
- [ ] Output is clean and readable
- [ ] Typecheck passes
- [ ] Tests pass

### US-011: Outdated Check
**Description:** As a user, I want to check for outdated dependencies so that I can identify available updates.

**Acceptance Criteria:**
- [ ] `arm outdated` shows packages with newer versions
- [ ] Displays constraint, current, wanted, and latest versions
- [ ] Accepts --output flag (table, json, list)
- [ ] Default output is table format
- [ ] Typecheck passes
- [ ] Tests pass

### US-012: Set Ruleset and Promptset
**Description:** As a user, I want to configure installed packages so that I can adjust their settings.

**Acceptance Criteria:**
- [ ] `arm set ruleset` accepts REGISTRY/RULESET KEY VALUE
- [ ] Supports keys: version, priority, sinks, include, exclude
- [ ] `arm set promptset` accepts REGISTRY/PROMPTSET KEY VALUE
- [ ] Supports keys: version, sinks, include, exclude
- [ ] Typecheck passes
- [ ] Tests pass

### US-013: Clean Cache
**Description:** As a user, I want to clean cached data so that I can free up space and remove stale data.

**Acceptance Criteria:**
- [ ] `arm clean cache` removes data older than 7 days (default)
- [ ] Accepts --max-age flag with duration (30m, 2h, 7d, 1h30m)
- [ ] Accepts --nuke flag for complete cleanup
- [ ] --nuke and --max-age are mutually exclusive
- [ ] Typecheck passes
- [ ] Tests pass

### US-014: Clean Sinks
**Description:** As a user, I want to clean sink directories so that I can remove orphaned files.

**Acceptance Criteria:**
- [ ] `arm clean sinks` removes files not in arm-index.json
- [ ] Accepts --nuke flag to remove entire ARM directory
- [ ] Typecheck passes
- [ ] Tests pass

### US-015: Compile Command
**Description:** As a user, I want to compile rulesets and promptsets so that I can test them before publishing.

**Acceptance Criteria:**
- [ ] `arm compile` accepts INPUT_PATH(s) and optional OUTPUT_PATH
- [ ] Accepts --tool flag (markdown, cursor, amazonq, copilot)
- [ ] Accepts --namespace, --force, --recursive, --validate-only flags
- [ ] Accepts --include, --exclude glob patterns
- [ ] Accepts --fail-fast flag
- [ ] Supports files, directories, and mixed inputs
- [ ] --validate-only makes OUTPUT_PATH optional
- [ ] Typecheck passes
- [ ] Tests pass

## 4. Functional Requirements

**FR-1:** All commands must follow the structure defined in commands.md exactly.

**FR-2:** Command output must be clean, simple, and easy to parse.

**FR-3:** Error messages must be clear and actionable.

**FR-4:** Commands must validate inputs before execution.

**FR-5:** Commands must update config files atomically to prevent corruption.

**FR-6:** Version constraints must support: @1.0.0 (exact), @1.1 (major.minor), @1 (major).

**FR-7:** Glob patterns must support standard glob syntax (**, *, ?).

**FR-8:** Duration parsing must support: m (minutes), h (hours), d (days), combined (1h30m).

**FR-9:** Commands must exit with appropriate status codes (0 = success, non-zero = error).

**FR-10:** Help text must be generated from command definitions, not hardcoded.

## 5. Non-Goals

- Backward compatibility with v3 CLI
- Migration tools from v3 to v4
- Interactive prompts or wizards
- Colored output or fancy formatting
- Progress bars or spinners
- Auto-completion scripts
- Shell integration

## 6. Technical Considerations

**Architecture:**
- One command handler per command
- Shared utilities for common operations (config, validation, output)
- Clear separation between CLI parsing and business logic
- Minimal abstractions - prefer explicit code

**Dependencies:**
- Use standard library where possible
- Minimal external dependencies for CLI parsing
- Reuse existing core packages (config, registry, sink, etc.)

**Testing:**
- Unit tests for command handlers where logic is non-trivial
- Focus on validation logic and error handling
- Mock external dependencies (filesystem, network)

**Code Style:**
- Follow grug-brained-dev principles
- Readable over clever
- Simple over complex
- Explicit over implicit

## 7. Success Metrics

1. All 15 user stories completed with passing tests
2. All commands from commands.md implemented
3. Zero hardcoded help text (generated from definitions)
4. Command execution time < 100ms for local operations
5. Code coverage > 70% for command handlers
6. Zero panics in production code paths
