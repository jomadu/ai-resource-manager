Study the implementation and update specs/:

1. Study code in internal/, cmd/, test/, docs/, and project files in root directory to understand architecture, features, components, functions, and edge cases. ULTRASTUDY.
2. Update specs/README.md with current Jobs to be Done (JTBDs) and Topics of Concern
   - README.md may not exist yet - create it
   - Existing README.md may be outdated - fix it
   - Keep the existing format (JTBDs section, Topics section, Spec docs links)
   - Update if core JTBDs or topics change
   - Allow large refactorings when structure needs to change
3. Update or create spec documents in specs/ using TEMPLATE.md
   - Specs may not exist yet - create them
   - Existing specs may be outdated - fix them
   - Implementation is the source of truth
   - Deep-dive algorithm verification: For each spec with algorithms, trace through actual implementation code to find exact discrepancies between spec and real behavior
   - Create DISCREPANCIES.md with findings (file:line references, what's wrong, what's correct)
   - Incorporate findings from DISCREPANCIES.md into specs
   - When bugs found in code: uncheck acceptance criteria + add brief note
   - Delete DISCREPANCIES.md after incorporating fixes
   - Reorganize/refactor specs if JTBDs or topics have evolved
4. Keep specs minimal but complete
5. When you learn something new about how to contextualize with the codebase or run the application, update AGENTS.md but keep it brief and operational only
6. Commit changes with descriptive message using the "conventional commit message" format (feat:, fix:, docs:, etc.) standard.

Focus: What does the code actually do? What are the acceptance criteria from tests?
