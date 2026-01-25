# Specification Coverage Analysis

## Purpose
Map every feature, command, and concept to ensure complete specification coverage for the documentation restructuring.

## Methodology
1. Extract all commands from `specs/commands.md`
2. Extract all concepts from existing specs
3. Map to implementation in `cmd/` and `internal/`
4. Identify gaps requiring builder specs

## Command Coverage

### Core Commands (specs/commands.md)
| Command | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| `arm version` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm help` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm list` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm info` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |

### Registry Management (specs/commands.md)
| Command | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| `arm add registry git` | ✅ commands.md | ✅ registry-management.md | ✅ cmd/arm/main.go → service.go |
| `arm add registry gitlab` | ✅ commands.md | ✅ registry-management.md | ✅ cmd/arm/main.go → service.go |
| `arm add registry cloudsmith` | ✅ commands.md | ✅ registry-management.md | ✅ cmd/arm/main.go → service.go |
| `arm remove registry` | ✅ commands.md | ✅ registry-management.md | ✅ cmd/arm/main.go → service.go |
| `arm set registry` | ✅ commands.md | ✅ registry-management.md | ✅ cmd/arm/main.go → service.go |
| `arm list registry` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm info registry` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |

### Sink Management (specs/commands.md)
| Command | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| `arm add sink` | ✅ commands.md | ✅ sink-compilation.md | ✅ cmd/arm/main.go → service.go |
| `arm remove sink` | ✅ commands.md | ✅ sink-compilation.md | ✅ cmd/arm/main.go → service.go |
| `arm set sink` | ✅ commands.md | ✅ sink-compilation.md | ✅ cmd/arm/main.go → service.go |
| `arm list sink` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm info sink` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |

### Dependency Management (specs/commands.md)
| Command | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| `arm install` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm install ruleset` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm install promptset` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm uninstall` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm update` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm upgrade` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm list dependency` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm info dependency` | ✅ commands.md | ❌ (trivial) | ✅ cmd/arm/main.go |
| `arm outdated` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm set ruleset` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |
| `arm set promptset` | ✅ commands.md | ✅ package-installation.md | ✅ cmd/arm/main.go → service.go |

### Utilities (specs/commands.md)
| Command | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| `arm clean cache` | ✅ commands.md | ✅ cache-management.md | ✅ cmd/arm/main.go → service.go |
| `arm clean sinks` | ✅ commands.md | ✅ sink-compilation.md | ✅ cmd/arm/main.go → service.go |
| `arm compile` | ✅ commands.md | ✅ sink-compilation.md | ✅ cmd/arm/main.go → service.go |

## Concept Coverage

### Core Concepts (specs/concepts.md)
| Concept | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| Registries | ✅ concepts.md, registries.md | ✅ registry-management.md | ✅ internal/arm/registry/ |
| Packages | ✅ concepts.md | ✅ package-installation.md | ✅ internal/arm/service/ |
| Sinks | ✅ concepts.md, sinks.md | ✅ sink-compilation.md | ✅ internal/arm/sink/ |
| Rulesets | ✅ concepts.md, resource-schemas.md | ✅ sink-compilation.md | ✅ internal/arm/resource/ |
| Promptsets | ✅ concepts.md, resource-schemas.md | ✅ sink-compilation.md | ✅ internal/arm/resource/ |
| File Patterns | ✅ concepts.md | ✅ pattern-filtering.md | ✅ internal/arm/registry/ |
| Versioning | ✅ concepts.md | ✅ version-resolution.md | ✅ internal/arm/core/ |

### Registry Types (specs/registries.md, specs/git-registry.md, specs/gitlab-registry.md, specs/cloudsmith-registry.md)
| Registry Type | User Doc | Builder Spec Needed | Implementation |
|---------------|----------|---------------------|----------------|
| Git Registry | ✅ git-registry.md | ✅ registry-management.md | ✅ internal/arm/registry/git.go |
| GitLab Registry | ✅ gitlab-registry.md | ✅ registry-management.md | ✅ internal/arm/registry/gitlab.go |
| Cloudsmith Registry | ✅ cloudsmith-registry.md | ✅ registry-management.md | ✅ internal/arm/registry/cloudsmith.go |
| Archive Support | ✅ registries.md | ✅ pattern-filtering.md | ✅ internal/arm/core/archive.go |

### Version Resolution (specs/concepts.md, specs/git-registry.md)
| Feature | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| Semantic Versioning | ✅ concepts.md | ✅ version-resolution.md | ✅ internal/arm/core/version.go |
| Version Constraints | ✅ concepts.md | ✅ version-resolution.md | ✅ internal/arm/core/constraint.go |
| Git Tags | ✅ git-registry.md | ✅ version-resolution.md | ✅ internal/arm/registry/git.go |
| Git Branches | ✅ git-registry.md | ✅ version-resolution.md | ✅ internal/arm/registry/git.go |
| Version Priority | ✅ git-registry.md | ✅ version-resolution.md | ✅ internal/arm/core/helpers.go |

### Sink Compilation (specs/sinks.md)
| Feature | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| Hierarchical Layout | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/sink/manager.go |
| Flat Layout | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/sink/manager.go |
| Cursor Compiler | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/compiler/cursor.go |
| AmazonQ Compiler | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/compiler/amazonq.go |
| Copilot Compiler | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/compiler/copilot.go |
| Markdown Compiler | ✅ sinks.md | ✅ sink-compilation.md | ✅ internal/arm/compiler/markdown.go |
| Priority Resolution | ✅ sinks.md | ✅ priority-resolution.md | ✅ internal/arm/sink/manager.go |
| Index Generation | ✅ sinks.md | ✅ priority-resolution.md | ✅ internal/arm/sink/manager.go |

### Storage & Caching (specs/storage.md)
| Feature | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| Storage Structure | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/ |
| Registry Cache | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/registry.go |
| Package Cache | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/package.go |
| Git Repo Cache | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/repo.go |
| Cache Keys | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/storage.go |
| Metadata | ✅ storage.md | ✅ cache-management.md | ✅ internal/arm/storage/package.go |

### Authentication (specs/armrc.md)
| Feature | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| .armrc Format | ✅ armrc.md | ✅ authentication.md | ✅ internal/arm/config/manager.go |
| Token Resolution | ✅ armrc.md | ✅ authentication.md | ✅ internal/arm/config/manager.go |
| GitLab Auth | ✅ armrc.md | ✅ authentication.md | ✅ internal/arm/registry/gitlab.go |
| Cloudsmith Auth | ✅ armrc.md | ✅ authentication.md | ✅ internal/arm/registry/cloudsmith.go |
| Environment Variables | ✅ armrc.md | ✅ authentication.md | ✅ internal/arm/config/manager.go |

### File Management (specs/concepts.md)
| File | User Doc | Builder Spec Needed | Implementation |
|------|----------|---------------------|----------------|
| arm.json | ✅ concepts.md | ✅ package-installation.md | ✅ internal/arm/manifest/ |
| arm-lock.json | ✅ concepts.md | ✅ package-installation.md | ✅ internal/arm/packagelockfile/ |
| arm-index.json | ✅ concepts.md | ✅ sink-compilation.md | ✅ internal/arm/sink/manager.go |
| arm_index.* | ✅ concepts.md | ✅ priority-resolution.md | ✅ internal/arm/sink/manager.go |

### Pattern Filtering (specs/concepts.md)
| Feature | User Doc | Builder Spec Needed | Implementation |
|---------|----------|---------------------|----------------|
| Glob Patterns | ✅ concepts.md | ✅ pattern-filtering.md | ✅ internal/arm/registry/ |
| Include Patterns | ✅ concepts.md | ✅ pattern-filtering.md | ✅ internal/arm/registry/ |
| Exclude Patterns | ✅ concepts.md | ✅ pattern-filtering.md | ✅ internal/arm/registry/ |
| Archive Extraction | ✅ concepts.md | ✅ pattern-filtering.md | ✅ internal/arm/core/archive.go |

## Builder Spec Requirements

### Required Builder Specs (8 total)

1. **specs/version-resolution.md**
   - **JTBD**: Resolve package versions from registries
   - **Covers**: Semver parsing, constraint matching, tag/branch priority, version comparison
   - **Maps to**: internal/arm/core/version.go, constraint.go, helpers.go
   - **User docs**: concepts.md, git-registry.md

2. **specs/package-installation.md**
   - **JTBD**: Install, update, upgrade, uninstall packages
   - **Covers**: Install workflow, reinstall behavior, lock file updates, manifest updates
   - **Maps to**: internal/arm/service/service.go, manifest/, packagelockfile/
   - **User docs**: commands.md, concepts.md

3. **specs/registry-management.md**
   - **JTBD**: Configure and manage registries
   - **Covers**: Registry types, configuration storage, authentication integration, key generation
   - **Maps to**: internal/arm/registry/, internal/arm/manifest/
   - **User docs**: registries.md, git-registry.md, gitlab-registry.md, cloudsmith-registry.md

4. **specs/sink-compilation.md**
   - **JTBD**: Compile resources to tool-specific formats
   - **Covers**: Compilation algorithms, layout modes, filename generation, tool formats
   - **Maps to**: internal/arm/compiler/, internal/arm/sink/manager.go
   - **User docs**: sinks.md, commands.md

5. **specs/priority-resolution.md**
   - **JTBD**: Resolve conflicts between overlapping rules
   - **Covers**: Priority merging, conflict resolution, index generation, metadata embedding
   - **Maps to**: internal/arm/sink/manager.go, internal/arm/compiler/generators.go
   - **User docs**: sinks.md, concepts.md

6. **specs/cache-management.md**
   - **JTBD**: Cache packages locally to avoid redundant downloads
   - **Covers**: Storage structure, cache keys, metadata schemas, cleanup strategies
   - **Maps to**: internal/arm/storage/
   - **User docs**: storage.md, commands.md

7. **specs/pattern-filtering.md**
   - **JTBD**: Filter package files using glob patterns
   - **Covers**: Glob matching, include/exclude logic, archive extraction, path sanitization
   - **Maps to**: internal/arm/registry/, internal/arm/core/archive.go
   - **User docs**: concepts.md, registries.md

8. **specs/authentication.md**
   - **JTBD**: Authenticate with registries requiring tokens
   - **Covers**: .armrc parsing, token resolution, environment variables, security
   - **Maps to**: internal/arm/config/manager.go, internal/arm/registry/gitlab.go, cloudsmith.go
   - **User docs**: armrc.md

### Spec Template

**specs/TEMPLATE.md** - Standard structure for all builder specs

## Coverage Summary

### Complete Coverage
- ✅ All 37 commands have user documentation
- ✅ All commands are implemented
- ✅ All core concepts documented
- ✅ All features implemented

### Gaps Requiring Builder Specs
- ❌ 8 builder specs needed to document implementation details
- ❌ No algorithmic specifications for builders
- ❌ No acceptance criteria for testing
- ❌ No edge case documentation for builders

## Verification Strategy

### Phase 1: User Doc Migration
- Move all user docs to `docs/`
- Verify no broken links
- Verify all commands still documented

### Phase 2: Builder Spec Creation
- Create TEMPLATE.md first
- Create 8 builder specs following template
- Ensure each spec maps to implementation
- Ensure each spec has testable acceptance criteria

### Phase 3: Cross-Reference Validation
- Verify every command links to a spec
- Verify every concept links to a spec
- Verify every implementation links to a spec
- Create traceability matrix

## Traceability Matrix

### Command → Spec → Implementation
```
arm install ruleset
  ├─ User Doc: docs/commands.md
  ├─ Builder Spec: specs/package-installation.md
  └─ Implementation: cmd/arm/main.go → service/service.go → sink/manager.go
```

### Concept → Spec → Implementation
```
Version Resolution
  ├─ User Doc: docs/concepts.md, docs/git-registry.md
  ├─ Builder Spec: specs/version-resolution.md
  └─ Implementation: internal/arm/core/version.go, constraint.go
```

### Feature → Spec → Implementation
```
Priority Resolution
  ├─ User Doc: docs/sinks.md
  ├─ Builder Spec: specs/priority-resolution.md
  └─ Implementation: internal/arm/sink/manager.go, compiler/generators.go
```

## Success Criteria

- [ ] Every command has user documentation in `docs/`
- [ ] Every complex feature has builder spec in `specs/`
- [ ] Every builder spec maps to implementation
- [ ] Every builder spec has acceptance criteria
- [ ] Traceability matrix complete
- [ ] No orphaned documentation
- [ ] No undocumented features
