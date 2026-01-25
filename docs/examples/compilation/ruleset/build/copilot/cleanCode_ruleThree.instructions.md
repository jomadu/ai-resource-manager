---
applyTo: "**/*"
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
    id: ruleThree
    name: "Rule 3"
    description: "Rule with optional description and enforcement (must|should|may) defined"
    enforcement: may
---

# Rule 3 (MAY)

This is the body of the rule.
