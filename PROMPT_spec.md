Study the implementation and update specs/:

1. ULTRASTUDY the code in internal/, cmd/, test/, and project files in root directory to understand architecture, features, components, functions, and edge cases. Study user documentation in ./docs as part of the implementation.
2. Update specs/README.md with current Jobs to be Done (JTBDs) and Topics of Concern
   - Assume existing README.md is inaccurate, incomplete, and outdated
   - Verify all claims by inspection and hard study of the implementation
   - Include JTBDs and topics for both code implementation and user documentation (./docs)
   - README.md may not exist yet - create it
   - Keep the existing format (JTBDs section, Topics section, Spec docs links)
   - Update if core JTBDs or topics change
   - Allow large refactorings when structure needs to change
3. Update or create spec documents in specs/ using TEMPLATE.md
   - Assume existing specs are inaccurate, incomplete, and outdated
   - Verify all claims (especially the acceptance criteria) by inspection and hard study of the implementation
   - Specs may not exist yet - create them
   - Implementation is the source of truth
   - Deep-dive verification: Trace through actual implementation code to find exact discrepancies between spec and real behavior, then fix specs directly
   - When bugs found in code: uncheck relevant acceptance criteria and add to "Known Issues" section of spec
   - When opportunities for improvement found: add to "Areas for Improvement" section of spec
   - Reorganize/refactor specs if JTBDs or topics have evolved
   - Keep specs minimal but complete and accurate
4. When you learn something new about how to contextualize with the codebase or run the application, update AGENTS.md but keep it brief and operational only
5. Commit changes with descriptive message using the "conventional commit message" format (feat:, fix:, docs:, etc.) standard.

Focus: What does the code actually do? What are the acceptance criteria from tests?
