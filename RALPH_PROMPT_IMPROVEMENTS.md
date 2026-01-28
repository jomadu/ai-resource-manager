# Ralph Prompt Architecture Improvements

## Core Problem

Current Ralph prompts are monolithic instruction blocks. They work but lack:
- **Composability** - Can't mix/match concerns
- **Interpretability** - Hard for agents to parse intent vs mechanics
- **Flexibility** - Each use case needs a new prompt file
- **Debuggability** - When Ralph fails, unclear which instruction caused it

## Proposed Architecture: Structured Prompt Composition

### 1. Separate Concerns into Layers

```
PROMPT.md = PHASE + TASK + GUARDRAILS + CONTEXT
```

**PHASE** - What stage of work (orient, plan, build, spec, test, deploy)
**TASK** - What to accomplish this iteration
**GUARDRAILS** - Invariants that must hold
**CONTEXT** - Project-specific knowledge

### 2. Use YAML Frontmatter for Machine-Readable Config

```yaml
---
mode: build
phase: implement
task_source: IMPLEMENTATION_PLAN.md
task_selection: most_important
subagent_budget:
  search: 500
  build: 1
  reasoning: opus
validation:
  - run_tests
  - update_plan
  - commit
context_files:
  - AGENTS.md
  - specs/*
output_format: conventional_commit
---
```

Agent parses this to understand:
- What it's supposed to do
- Resource constraints
- Success criteria
- Where to look for context

### 3. Instruction Blocks with Semantic IDs

Instead of:
```
0a. Study specs/* with up to 500 parallel Sonnet subagents...
0b. Study @IMPLEMENTATION_PLAN.md...
```

Use:
```markdown
## [ORIENT:specs]
Study `specs/*` to learn application specifications.
- Use: up to {{subagent_budget.search}} parallel subagents
- Focus: behavioral requirements, acceptance criteria

## [ORIENT:plan]
Study `{{task_source}}` to understand current work state.
- Treat as: potentially incorrect, verify claims
- Extract: next most important task

## [TASK:implement]
Implement functionality per specifications.
- Before changes: search codebase (don't assume not implemented)
- Selection: {{task_selection}} from {{task_source}}
- Subagents: {{subagent_budget.search}} for search, {{subagent_budget.build}} for build/test
```

Benefits:
- Agent can reference specific blocks: "Following [TASK:implement]..."
- Humans can debug: "Ralph failed at [ORIENT:plan]"
- Composable: Mix/match blocks for different modes
- Templatable: `{{variables}}` filled from frontmatter

### 4. Guardrails as Declarative Rules

Instead of numbered 9s:
```
99999. Important: When authoring documentation...
999999. Important: Single sources of truth...
```

Use structured rules:
```yaml
guardrails:
  documentation:
    rule: "Capture the why — tests and implementation importance"
    applies_to: [docs, comments, specs]
    severity: critical
  
  single_source_of_truth:
    rule: "No migrations/adapters. If unrelated tests fail, resolve them."
    applies_to: [implementation]
    severity: critical
  
  no_placeholders:
    rule: "Implement completely. Placeholders waste effort."
    applies_to: [implementation]
    severity: critical
    
  plan_hygiene:
    rule: "Keep {{task_source}} current. Future work depends on this."
    applies_to: [plan_updates]
    severity: critical
    when: after_task_completion
```

Agent can:
- Check which rules apply to current phase
- Prioritize by severity
- Understand when rules trigger

### 5. Mode-Specific Prompt Composition

**Base Template** (`PROMPT_base.md`):
```yaml
---
# Filled by loop.sh based on mode argument
mode: "{{MODE}}"
---

## [ORIENT:context]
{{include "blocks/orient_context.md"}}

## [PHASE:{{mode}}]
{{include "blocks/phase_{{mode}}.md"}}

## [VALIDATE]
{{include "blocks/validate_{{mode}}.md"}}

## [GUARDRAILS]
{{include "guardrails.yaml" | render}}
```

**Block Files**:
```
prompts/
├── PROMPT_base.md          # Composition template
├── config/
│   ├── build.yaml          # Build mode config
│   ├── plan.yaml           # Plan mode config
│   ├── spec.yaml           # Spec mode config
│   └── test.yaml           # Test mode config
├── blocks/
│   ├── orient_context.md   # Standard orientation
│   ├── phase_build.md      # Build-specific task
│   ├── phase_plan.md       # Planning-specific task
│   ├── phase_spec.md       # Spec-specific task
│   ├── validate_build.md   # Build validation
│   └── validate_plan.md    # Plan validation
└── guardrails.yaml         # Shared invariants
```

**loop.sh generates final prompt**:
```bash
# Render prompt from template + mode config
render_prompt() {
    local mode=$1
    local config="prompts/config/${mode}.yaml"
    local template="prompts/PROMPT_base.md"
    
    # Simple template engine (or use envsubst, mustache, etc.)
    MODE=$mode envsubst < "$template" > "PROMPT_${mode}.md"
}

render_prompt "$MODE"
cat "PROMPT_${MODE}.md" | kiro-cli chat --no-interactive --trust-all-tools
```

### 6. Self-Documenting Execution

Agent outputs structured logs:
```
[ORIENT:specs] Studying 12 spec files with 250 subagents...
[ORIENT:plan] Loaded IMPLEMENTATION_PLAN.md (47 tasks remaining)
[TASK:implement] Selected: "Add --exclude flag to install command"
[TASK:implement] Searching codebase for existing implementation...
[TASK:implement] Not found. Proceeding with implementation.
[VALIDATE:tests] Running: go test ./...
[VALIDATE:tests] ✓ All tests passed
[VALIDATE:plan] Updating IMPLEMENTATION_PLAN.md via subagent
[VALIDATE:commit] git commit -m "feat: add --exclude flag to install command"
[GUARDRAIL:plan_hygiene] ✓ Plan updated after task completion
```

Benefits:
- Trace execution to specific prompt blocks
- Identify where Ralph gets stuck
- Validate guardrails are being followed

### 7. Adaptive Subagent Budgets

Instead of hardcoded numbers:
```yaml
subagent_budget:
  search:
    default: 500
    adaptive: true
    scale_by: task_complexity
  build:
    default: 1
    reason: "Backpressure control - serial validation"
  reasoning:
    model: opus
    when: [debugging, architecture, complex_analysis]
    fallback: sonnet
```

Agent can:
- Scale resources based on task size
- Understand why limits exist
- Make intelligent model choices

### 8. Task Selection Strategies

Make selection logic explicit:
```yaml
task_selection:
  strategy: priority_first
  filters:
    - unblocked
    - has_acceptance_criteria
  sort_by:
    - priority: desc
    - dependencies: asc
  fallback: ask_user
```

Alternative strategies:
- `most_important` - Current Ralph default
- `dependency_order` - Unblock other tasks first
- `quick_wins` - Low-effort, high-value
- `risk_reduction` - Tackle unknowns early
- `user_specified` - Human picks from plan

### 9. Validation Pipelines

Instead of prose instructions:
```yaml
validation:
  steps:
    - name: run_tests
      command: "{{test_command}}"
      required: true
      on_fail: fix_and_retry
      
    - name: update_plan
      action: subagent_update
      target: "{{task_source}}"
      required: true
      
    - name: update_agents
      action: subagent_update
      target: AGENTS.md
      required: false
      when: operational_learning
      
    - name: commit
      command: "git add -A && git commit -m '{{commit_message}}'"
      required: true
      format: conventional_commit
      
    - name: tag
      command: "git tag {{next_version}}"
      required: false
      when: no_errors_and_no_tags_or_increment
```

Agent executes as pipeline, knows what's required vs optional.

### 10. Context Budget Management

Make context allocation explicit:
```yaml
context_budget:
  total_tokens: 176000
  allocation:
    prompt_structure: 5000
    specs: 15000
    plan: 5000
    agents: 2000
    code_context: 149000
  
  optimization:
    - prefer_summaries_over_full_files
    - chunk_large_specs
    - prioritize_relevant_code_only
```

Agent understands:
- Why it should be concise
- What to prioritize loading
- When to use summaries vs full content

## Implementation Strategy

### Phase 1: Backward Compatible Enhancement
Add YAML frontmatter to existing prompts without breaking current loop.sh:
```markdown
---
mode: build
# ... config ...
---

<!-- Existing prose instructions below -->
0a. Study specs/* with up to 500 parallel Sonnet subagents...
```

Agents that understand YAML get benefits, others ignore it.

### Phase 2: Hybrid Approach
Keep prose but add semantic IDs:
```markdown
## [ORIENT:specs]
Study `specs/*` with up to 500 parallel Sonnet subagents to learn the application specifications.
```

Enables better logging and debugging.

### Phase 3: Full Composition
Migrate to block-based architecture with template rendering.

## Use Case Examples

### Use Case: Test-Driven Development Mode

```yaml
---
mode: tdd
phase: red_green_refactor
task_source: IMPLEMENTATION_PLAN.md
validation:
  - write_failing_test
  - implement_minimal
  - verify_test_passes
  - refactor
  - commit
---
```

### Use Case: Documentation-First Mode

```yaml
---
mode: docs_first
phase: document_then_implement
task_source: DOCUMENTATION_PLAN.md
validation:
  - write_user_docs
  - derive_acceptance_criteria
  - implement_to_spec
  - verify_docs_accurate
  - commit
---
```

### Use Case: Refactoring Mode

```yaml
---
mode: refactor
phase: improve_without_behavior_change
task_source: REFACTORING_PLAN.md
validation:
  - snapshot_tests
  - refactor_code
  - verify_tests_unchanged
  - check_performance
  - commit
guardrails:
  behavior_preservation:
    rule: "All existing tests must pass unchanged"
    severity: critical
---
```

### Use Case: Security Audit Mode

```yaml
---
mode: security_audit
phase: identify_vulnerabilities
task_source: SECURITY_CHECKLIST.md
subagent_budget:
  analysis: 1000
  reasoning: opus
validation:
  - scan_for_vulnerabilities
  - document_findings
  - suggest_fixes
  - update_checklist
output_format: security_report
---
```

### Use Case: Performance Optimization Mode

```yaml
---
mode: performance
phase: profile_and_optimize
task_source: PERFORMANCE_PLAN.md
validation:
  - benchmark_baseline
  - implement_optimization
  - benchmark_after
  - verify_improvement
  - document_results
  - commit
guardrails:
  performance_regression:
    rule: "No optimization that degrades other metrics >5%"
    severity: critical
---
```

## Benefits Summary

**For Agents:**
- Clear structure to parse and understand
- Explicit resource constraints
- Unambiguous success criteria
- Self-documenting execution

**For Humans:**
- Debuggable (trace failures to specific blocks)
- Composable (mix/match for new use cases)
- Maintainable (change one block, affects all modes)
- Readable (YAML + semantic IDs > numbered prose)

**For the Loop:**
- Same outer loop mechanism
- Fresh context per iteration
- Deterministic setup
- Eventual consistency

## Migration Path

1. **Week 1**: Add YAML frontmatter to existing prompts (backward compatible)
2. **Week 2**: Add semantic IDs to instruction blocks
3. **Week 3**: Extract common blocks to separate files
4. **Week 4**: Implement template rendering in loop.sh
5. **Week 5**: Create 3-5 new mode configs to validate flexibility
6. **Week 6**: Update AGENTS.md with new prompt architecture

## Open Questions

1. **Template engine choice**: envsubst, mustache, jinja2, or custom?
2. **Guardrail enforcement**: Should loop.sh validate guardrails or trust agent?
3. **Config validation**: JSON Schema for YAML configs?
4. **Subagent communication**: Should subagents see the full prompt or just their block?
5. **Failure recovery**: How should agent signal "I'm stuck, need human"?

## Recommendation

Start with **Phase 1 (YAML frontmatter)** in your current PROMPT_spec.md. Test if Kiro/Claude can parse and use it. If successful, proceed to semantic IDs and composition.

The goal: **Make Ralph prompts as composable and maintainable as the code Ralph writes.**
