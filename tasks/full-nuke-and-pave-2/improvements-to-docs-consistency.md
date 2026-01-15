# Documentation Conceptual Inconsistencies

This document tracks conceptual inconsistencies in the ARM documentation that need to be resolved to improve user understanding and reduce confusion.

## 1. Package vs Resource Terminology Confusion

**Issue**: Inconsistent use of "packages" and "resources" throughout documentation.

**Current State**:
- README: "packages are versioned collections of AI rules or prompts"
- Commands: Mix "install package" and "install ruleset/promptset"
- Concepts: Describes rulesets/promptsets as the actual resource types

**Impact**: Users confused about whether they install "packages" or "rulesets/promptsets"

**Resolution Strategy**:
- Use "ruleset" and "promptset" consistently throughout all documentation
- Remove confusing "package" terminology from user-facing content
- Match documentation language to actual command syntax

**Priority**: High
**Status**: Complete - Fixed README.md, concepts.md, registries.md, sinks.md, git-registry.md, gitlab-registry.md, cloudsmith-registry.md

---

## 2. Registry Structure Contradictions

**Issue**: Recommended structure conflicts with actual behavior across registry types.

**Current State**:
- Registries doc: Shows ARM resource files at root level
- Git registries: Define "packages" by file patterns within repos
- Other registries: Use actual package names

**Impact**: Users don't understand how packages are defined differently per registry type

**Resolution Strategy**:
- Create registry-specific structure examples
- Clearly explain pattern-based vs name-based package definitions
- Separate Git registry behavior from others

**Files to Update**:
- docs/registries.md - Add registry-specific structure sections
- docs/git-registry.md - Clarify how Git repos become "packages" via patterns
- docs/gitlab-registry.md - Show how actual package names work
- docs/cloudsmith-registry.md - Show single-file package model

**Key Points for Next Agent**:
- Git registries: Repository contents become packages based on file patterns (e.g., all .yml files)
- GitLab/Cloudsmith: Explicit package names with versioned artifacts
- Need clear examples showing same content structured differently per registry type

**Priority**: High
**Status**: Complete - Updated registries.md with Git vs non-Git models, clarified git-registry.md repository behavior, updated gitlab-registry.md and cloudsmith-registry.md with explicit package examples

---

## 3. Version Resolution Inconsistencies

**Issue**: Different registry types handle versioning differently without clear explanation.

**Current State**:
- Git: Uses tags AND branches (non-semantic)
- GitLab/Cloudsmith: Only semantic versions
- Docs don't explain when you get commit hash vs semantic version

**Impact**: Users surprised by version format differences

**Resolution Strategy**:
- Create dedicated versioning section explaining each type
- Add decision tree for version resolution
- Clarify semantic vs Git-based versioning upfront

**Files to Update**:
- docs/concepts.md - Add "Versioning" section explaining differences
- docs/git-registry.md - Expand version resolution section with examples
- docs/gitlab-registry.md - Clarify semantic-only versioning
- docs/cloudsmith-registry.md - Clarify semantic-only versioning

**Key Points for Next Agent**:
- Git registries: Tags (semantic) > Tags (non-semantic) > Branches (commit hash)
- GitLab/Cloudsmith: Only semantic versions, no commit hashes
- User sees different version formats depending on registry type
- Need decision tree: "What version will I get?"

**Priority**: Medium
**Status**: Complete - Added Versioning section to concepts.md with decision tree, expanded git-registry.md with detailed examples and lock file formats, clarified semantic-only versioning in gitlab-registry.md and cloudsmith-registry.md. Removed non-semantic tag support from code and tests.

---

## 4. Authentication Model Confusion

**Issue**: Inconsistent messaging about when authentication is required.

**Current State**:
- Git registries: "No additional configuration needed"
- GitLab/Cloudsmith: "Require explicit token authentication"
- Both use same `.armrc` format

**Impact**: Users confused about when tokens are needed

**Resolution Strategy**:
- Create authentication decision tree
- Separate Git auth (uses system Git) from API auth (requires tokens)
- Consolidate `.armrc` documentation

**Files to Update**:
- docs/concepts.md - Add "Authentication" section with decision tree
- docs/armrc.md - Consolidate all .armrc examples and patterns
- Registry docs - Standardize auth sections with clear "when needed" messaging

**Key Points for Next Agent**:
- Git registries: Use system Git auth (SSH keys, credential helpers, CLI tools)
- GitLab/Cloudsmith: Require explicit API tokens in .armrc
- Decision tree: "Do I need a token?" based on registry type and repo visibility
- .armrc format is same, but usage differs

**Priority**: Medium
**Status**: Complete.

---

## 5. File Pattern Behavior Inconsistencies

**Issue**: Pattern matching behavior varies and isn't clearly documented.

**Current State**:
- Default includes: "all .yml and .yaml files"
- Examples show `**/*.yml` patterns
- Archive extraction mentioned but not integrated with patterns

**Impact**: Users don't understand how patterns work with different content types

**Resolution Strategy**:
- Unify pattern documentation across all registry types
- Explain archive extraction in context of pattern matching
- Provide clear examples for each scenario

**Files to Update**:
- docs/concepts.md - Add "File Patterns" section explaining glob behavior
- All registry docs - Standardize pattern examples and archive extraction explanations
- docs/commands.md - Ensure --include/--exclude examples are consistent

**Key Points for Next Agent**:
- Default behavior: Include all .yml and .yaml files
- Patterns apply AFTER archive extraction (archives extracted first, then patterns filter)
- Same pattern syntax across all registry types
- Need examples: "Install only TypeScript rules", "Exclude experimental", etc.

**Priority**: Low
**Status**: Identified

---

## 6. Sink vs Tool Terminology

**Issue**: Unclear relationship between sink configuration and AI tool integration.

**Current State**:
- Sinks configured with `--tool` parameter
- Docs refer to "AI tools" and "tool formats"
- Purpose of sinks not clearly explained

**Impact**: Users don't understand that sinks are output destinations

**Resolution Strategy**:
- Clarify that sinks are compilation targets, not integrations
- Better explain sink purpose in concepts
- Separate sink configuration from tool usage

**Files to Update**:
- docs/concepts.md - Expand sink definition with clear purpose
- docs/sinks.md - Add "What are sinks?" section at top
- README.md - Ensure sink examples clearly show "output destination" concept

**Key Points for Next Agent**:
- Sinks are OUTPUT DESTINATIONS, not integrations
- --tool parameter specifies COMPILATION FORMAT, not which tool to integrate with
- User still needs to configure their AI tool to read from sink directory
- Analogy: Sinks are like "build targets" - same source, different output formats

**Priority**: Low
**Status**: Identified

---

## Resolution Plan

### Phase 1: Core Terminology (High Priority)
1. Fix package vs resource terminology
2. Clarify registry structure differences
3. Update command documentation

### Phase 2: Technical Clarity (Medium Priority)
4. Document version resolution per registry type
5. Consolidate authentication documentation

### Phase 3: Polish (Low Priority)
6. Unify pattern documentation
7. Clarify sink purpose and terminology

### Success Criteria
- Users can clearly distinguish packages from resources
- Registry-specific behavior is well documented
- Authentication requirements are obvious
- Examples are consistent across all docs
