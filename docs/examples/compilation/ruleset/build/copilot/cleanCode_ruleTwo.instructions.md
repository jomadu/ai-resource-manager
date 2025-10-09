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
    id: ruleTwo
    name: "Rule 2"
    description: "Rule with optional description and enforcement (must|should|may) defined"
    enforcement: should
---

# Rule 2 (SHOULD)

This is the body of the rule.
