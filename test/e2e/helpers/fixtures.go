package helpers

// Minimal test resources for E2E tests

const MinimalRuleset = `apiVersion: v1
kind: Ruleset
metadata:
  id: "testRuleset"
  name: "Test Ruleset"
  description: "A minimal test ruleset"
spec:
  rules:
    ruleOne:
      body: "This is rule one."
    ruleTwo:
      priority: 150
      body: "This is rule two with priority."
`

const SecurityRuleset = `apiVersion: v1
kind: Ruleset
metadata:
  id: "securityRuleset"
  name: "Security Ruleset"
  description: "Security best practices"
spec:
  rules:
    securityRule:
      priority: 200
      body: "Always validate user input."
`

const MinimalPromptset = `apiVersion: v1
kind: Promptset
metadata:
  id: "testPromptset"
  name: "Test Promptset"
  description: "A minimal test promptset"
spec:
  prompts:
    promptOne:
      body: "This is prompt one."
    promptTwo:
      body: "This is prompt two."
`

const CodeReviewPromptset = `apiVersion: v1
kind: Promptset
metadata:
  id: "codeReviewPromptset"
  name: "Code Review Promptset"
  description: "Code review prompts"
spec:
  prompts:
    reviewCode:
      body: "Review this code for best practices."
    suggestImprovements:
      body: "Suggest improvements for this code."
`
