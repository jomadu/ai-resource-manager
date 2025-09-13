package urf

import (
	"fmt"
	"strings"
)

// MetadataGenerator interface for generating metadata blocks
type MetadataGenerator interface {
	GenerateMetadata(urf *URFFile, rule *Rule, namespace string) string
}

// DefaultMetadataGenerator implements shared metadata generation
type DefaultMetadataGenerator struct{}

// NewMetadataGenerator creates a new metadata generator
func NewMetadataGenerator() MetadataGenerator {
	return &DefaultMetadataGenerator{}
}

// GenerateMetadata generates the shared metadata block
func (g *DefaultMetadataGenerator) GenerateMetadata(urf *URFFile, rule *Rule, namespace string) string {
	var content strings.Builder

	content.WriteString("---\n")
	content.WriteString(fmt.Sprintf("namespace: %s\n", namespace))
	content.WriteString("ruleset:\n")
	content.WriteString(fmt.Sprintf("  id: %s\n", urf.Metadata.ID))
	content.WriteString(fmt.Sprintf("  name: %s\n", urf.Metadata.Name))
	content.WriteString(fmt.Sprintf("  version: %s\n", urf.Metadata.Version))
	content.WriteString("  rules:\n")
	for _, r := range urf.Rules {
		content.WriteString(fmt.Sprintf("    - %s\n", r.ID))
	}
	content.WriteString("rule:\n")
	content.WriteString(fmt.Sprintf("  id: %s\n", rule.ID))
	content.WriteString(fmt.Sprintf("  name: %s\n", rule.Name))
	content.WriteString(fmt.Sprintf("  enforcement: %s\n", strings.ToUpper(rule.Enforcement)))
	content.WriteString(fmt.Sprintf("  priority: %d\n", rule.Priority))
	if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
		content.WriteString("  scope:\n")
		for _, file := range rule.Scope[0].Files {
			content.WriteString(fmt.Sprintf("    - files: %q\n", file))
		}
	}
	content.WriteString("---\n\n")

	return content.String()
}
