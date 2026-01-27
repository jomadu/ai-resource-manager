Study the implementation and update specs/:

1. Use git history to understand what's been implemented since last spec update
2. Study code in internal/, cmd/, test/, docs/, and project files in root directory to understand architecture, features, components, functions, and edge cases. ULTRASTUDY.
3. Update specs/README.md with current Jobs to be Done (JTBDs) and Topics of Concern
   - README.md may not exist yet - create it
   - Existing README.md may be outdated - fix it
   - Keep the existing format (JTBDs section, Topics section, Spec docs links)
   - Update if core JTBDs or topics change
   - Allow large refactorings when structure needs to change
4. Update or create spec documents in specs/ using TEMPLATE.md
   - Specs may not exist yet - create them
   - Existing specs may be outdated - fix them
   - Implementation is the source of truth
   - Reorganize/refactor specs if JTBDs or topics have evolved
5. Keep specs minimal but complete
6. When you learn something new about how to contextualize with the codebase or run the application, update AGENTS.md but keep it brief and operational only
7. Commit changes with descriptive message using the "conventional commit message" format (feat:, fix:, docs:, etc.) standard.

Focus: What does the code actually do? What are the acceptance criteria from tests?
