# Orient: Bootstrap

## R1: Identify Project Type and Tech Stack

Analyze the repository to determine:
- What programming language(s) are used?
- What framework(s) or libraries are present?
- What build system is used? (package.json, Cargo.toml, go.mod, pom.xml, etc.)
- What is the project structure? (monorepo, single app, library, etc.)
- What deployment target? (web, CLI, library, service, etc.)

## R2: Determine What Constitutes "Specification" vs "Implementation"

Based on repository structure, identify:
- Where are specifications located? (specs/, docs/, README sections, inline docs)
- What format are specs? (markdown, docstrings, API docs, etc.)
- What pattern do specs follow? (JTBD-based, feature-based, API-based)
- Where is implementation located? (src/, lib/, app/, etc.)
- What file patterns constitute implementation? (*.js, *.py, *.go, *.rs, etc.)
- What should be excluded? (tests, build artifacts, vendor, node_modules)

## R3: Identify Build/Test/Run Commands Empirically

Discover operational commands by examining:
- Package manager files (package.json scripts, Makefile targets, etc.)
- CI/CD configuration (.github/workflows, .gitlab-ci.yml, etc.)
- Documentation (README, CONTRIBUTING, etc.)
- Common conventions for the tech stack

Verify commands work:
- Try running test commands
- Try running build commands
- Try running lint commands
- Document what works, what fails, what's missing

## R4: Synthesize Operational Understanding

Combine findings into coherent operational guide:
- How should agents query work tracking?
- How should agents run tests?
- How should agents build the project?
- How should agents lint code?
- What quality criteria apply to this project?
- What conventions and patterns should be followed?
- What learnings should be captured?
