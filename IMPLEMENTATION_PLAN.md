# Implementation Plan

## Goal
Separate user documentation from builder specifications. Align specs with Ralph methodology (JTBD → activities → acceptance criteria → tasks).

See [SPECIFICATION_PHILOSOPHY.md](./SPECIFICATION_PHILOSOPHY.md) for detailed guidance on writing effective builder-oriented specifications.

## Tasks

- [ ] Move current `specs/*.md` to `docs/` (except `e2e-testing.md`)
  - [ ] Create `docs/` directory
  - [ ] Move `concepts.md` → `docs/concepts.md`
  - [ ] Move `commands.md` → `docs/commands.md`
  - [ ] Move `registries.md` → `docs/registries.md`
  - [ ] Move `git-registry.md` → `docs/git-registry.md`
  - [ ] Move `gitlab-registry.md` → `docs/gitlab-registry.md`
  - [ ] Move `cloudsmith-registry.md` → `docs/cloudsmith-registry.md`
  - [ ] Move `sinks.md` → `docs/sinks.md`
  - [ ] Move `storage.md` → `docs/storage.md`
  - [ ] Move `resource-schemas.md` → `docs/resource-schemas.md`
  - [ ] Move `armrc.md` → `docs/armrc.md`
  - [ ] Move `migration-v2-to-v3.md` → `docs/migration-v2-to-v3.md`
  - [ ] Move `specs/examples/` → `docs/examples/`
  - [ ] Keep `specs/e2e-testing.md` (already builder-oriented)

- [ ] Update README.md references
  - [ ] Update links from `specs/` to `docs/` for user documentation
  - [ ] Keep `specs/e2e-testing.md` reference as-is

- [ ] Create builder specs in `specs/`
  - [ ] Create `specs/version-resolution.md` (version resolution algorithm, edge cases)
  - [ ] Create `specs/package-installation.md` (install workflow, state transitions, reinstall behavior)
  - [ ] Create `specs/registry-management.md` (registry types, configuration, authentication)
  - [ ] Create `specs/sink-compilation.md` (compilation algorithms per tool, layout modes)
  - [ ] Create `specs/priority-resolution.md` (priority merging, conflict resolution, index generation)
  - [ ] Create `specs/cache-management.md` (storage structure, cache keys, cleanup strategies)
  - [ ] Create `specs/pattern-filtering.md` (glob matching, include/exclude logic, archive extraction)
  - [ ] Create `specs/authentication.md` (.armrc parsing, token resolution, environment variables)

- [ ] Create spec template
  - [ ] Create `specs/TEMPLATE.md` with standard structure (JTBD, activities, acceptance criteria, data structures, algorithm, edge cases, dependencies, examples)

## Spec Structure

Each builder spec should contain:
- **Job to be Done** - User outcome this enables
- **Activities** - Discrete steps/operations
- **Acceptance Criteria** - Observable, verifiable outcomes
- **Data Structures** - Schemas, types, state machines, invariants
- **Algorithm** - Step-by-step logic
- **Edge Cases** - Boundary conditions, error handling, failure modes
- **Dependencies** - Prerequisites
- **Examples** - Concrete scenarios

## Why This Aligns with Ralph

- Specs define WHAT to build (outcomes, acceptance criteria)
- Ralph decides HOW to implement (technical approach)
- Clear success signals (tests derived from acceptance criteria)
- Specs are disposable (can regenerate if wrong)
- Acceptance criteria enable test-driven development
