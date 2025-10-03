---
applyTo: "**/*.py"
---

---
namespace: someNamespace
ruleset:
    id: "cleanCode"
    name: "Clean Code"
    description: "Make clean code"
    rules:
        - ruleOne
        - ruleTwo
        - ruleThree
        - ruleFour
        - ruleFive
rule:
    id: ruleOne
    name: "Rule 1"
    description: "Rule with optional description, priority, enforcement (must|should|may), scope[0].files attributes defined"
    priority: 100
    enforcement: must
    scope:
        - files: ["**/*.py"]
---

# Rule 1 (MUST)

This is the body of the rule.
