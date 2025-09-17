package converter

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/urf"
)

// MarkdownConverter converts markdown text to URF format
type MarkdownConverter struct {
	// Patterns for detecting different markdown elements
	headerPattern    *regexp.Regexp
	bulletPattern    *regexp.Regexp
	codeBlockPattern *regexp.Regexp

	// Current parsing state
	currentSection string
	currentRule    *urf.Rule
	inCodeBlock    bool
	codeContent    strings.Builder
	nextRuleID     int
}

// NewMarkdownConverter creates a new markdown converter
func NewMarkdownConverter() Converter {
	return &MarkdownConverter{
		headerPattern:    regexp.MustCompile(`^(#{1,6})\s+(.+)$`),
		bulletPattern:    regexp.MustCompile(`^[-*+]\s+(.+)$`),
		codeBlockPattern: regexp.MustCompile(`^` + "```" + `(\w*)?$`),
		nextRuleID:       1,
	}
}

// Convert converts markdown content to URF format
func (c *MarkdownConverter) Convert(content string, config ConversionConfig) (*URFRuleset, error) {
	// Initialize URF ruleset
	ruleset := &urf.Ruleset{
		Version: "1.0",
		Metadata: urf.Metadata{
			ID:          config.RulesetID,
			Name:        config.RulesetName,
			Version:     config.Version,
			Description: fmt.Sprintf("Rules converted from %s", config.InputFile),
		},
		Rules: make([]urf.Rule, 0),
	}

	// Reset parser state
	c.currentSection = "general"
	c.currentRule = nil
	c.inCodeBlock = false
	c.nextRuleID = 1

	// Parse content line by line
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		if err := c.parseLine(line, ruleset); err != nil {
			return nil, fmt.Errorf("error parsing line %d: %w", lineNum, err)
		}
	}

	// Finalize any pending rule
	c.finalizeCurrentRule(ruleset)

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading content: %w", err)
	}

	return &URFRuleset{Ruleset: ruleset}, nil
}

// parseLine processes a single line of markdown
func (c *MarkdownConverter) parseLine(line string, ruleset *urf.Ruleset) error {
	// Skip empty lines
	if line == "" {
		return nil
	}

	// Handle code blocks
	if c.codeBlockPattern.MatchString(line) {
		return c.handleCodeBlock(line)
	}

	// If we're inside a code block, collect content
	if c.inCodeBlock {
		c.codeContent.WriteString(line + "\n")
		return nil
	}

	// Handle headers (section markers)
	if matches := c.headerPattern.FindStringSubmatch(line); matches != nil {
		return c.handleHeader(matches[1], matches[2], ruleset)
	}

	// Handle bullet points (rules)
	if matches := c.bulletPattern.FindStringSubmatch(line); matches != nil {
		return c.handleBulletPoint(matches[1], ruleset)
	}

	// Handle regular text (could be description or additional content)
	if c.currentRule != nil {
		// Append to current rule body if we have one
		if c.currentRule.Body != "" {
			c.currentRule.Body += "\n" + line
		}
	}

	return nil
}

// handleHeader processes markdown headers as section markers
func (c *MarkdownConverter) handleHeader(_, text string, ruleset *urf.Ruleset) error {
	// Finalize previous rule if any
	c.finalizeCurrentRule(ruleset)

	// Extract section name from header text
	section := c.extractSectionName(text)
	c.currentSection = section

	return nil
}

// handleBulletPoint processes markdown bullet points as rules
func (c *MarkdownConverter) handleBulletPoint(text string, ruleset *urf.Ruleset) error {
	// Finalize previous rule if any
	c.finalizeCurrentRule(ruleset)

	// Create new rule from bullet point
	rule := urf.Rule{
		ID:          c.generateRuleID(text),
		Name:        text,
		Description: fmt.Sprintf("Rule from %s section", c.currentSection),
		Priority:    c.inferPriority(text),
		Enforcement: c.inferEnforcement(text),
		Scope:       c.inferScope(text),
		Body:        c.generateRuleBody(text),
	}

	c.currentRule = &rule
	return nil
}

// handleCodeBlock processes markdown code blocks
func (c *MarkdownConverter) handleCodeBlock(_ string) error {
	if !c.inCodeBlock {
		// Starting a code block
		c.inCodeBlock = true
		c.codeContent.Reset()
	} else {
		// Ending a code block
		c.inCodeBlock = false

		// Attach code to current rule if available
		if c.currentRule != nil {
			codeContent := strings.TrimSpace(c.codeContent.String())
			if codeContent != "" {
				// Add code block to rule body
				c.currentRule.Body += "\n\n```\n" + codeContent + "\n```"
			}
		}
	}

	return nil
}

// finalizeCurrentRule adds the current rule to the ruleset
func (c *MarkdownConverter) finalizeCurrentRule(ruleset *urf.Ruleset) {
	if c.currentRule != nil {
		ruleset.Rules = append(ruleset.Rules, *c.currentRule)
		c.currentRule = nil
	}
}

// extractSectionName extracts a clean section name from header text
func (c *MarkdownConverter) extractSectionName(text string) string {
	section := strings.ToLower(text)
	section = strings.TrimSpace(section)

	// Remove common words
	section = strings.ReplaceAll(section, "rules", "")
	section = strings.ReplaceAll(section, "guidelines", "")
	section = strings.ReplaceAll(section, "practices", "")
	section = strings.ReplaceAll(section, "best", "")
	section = strings.TrimSpace(section)

	// Default section if empty
	if section == "" {
		section = "general"
	}

	return section
}

// generateRuleID creates a unique rule ID
func (c *MarkdownConverter) generateRuleID(text string) string {
	// Convert text to kebab-case ID
	id := strings.ToLower(text)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "'", "")
	id = strings.ReplaceAll(id, "\"", "")

	// Remove special characters except hyphens
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	cleanID := result.String()

	// Ensure uniqueness by adding a number
	finalID := fmt.Sprintf("%s-%d", cleanID, c.nextRuleID)
	c.nextRuleID++

	return finalID
}

// inferPriority attempts to infer rule priority from text
func (c *MarkdownConverter) inferPriority(text string) int {
	text = strings.ToLower(text)

	// Critical priority keywords (90-100)
	if containsAny(text, []string{"must", "never", "always", "critical", "security", "required"}) {
		return 95
	}

	// High priority keywords (80-89)
	if containsAny(text, []string{"should", "prefer", "recommend", "important"}) {
		return 85
	}

	// Medium priority keywords (70-79)
	if containsAny(text, []string{"consider", "may", "optional", "style"}) {
		return 75
	}

	// Default priority
	return 80
}

// inferEnforcement attempts to infer enforcement level from text
func (c *MarkdownConverter) inferEnforcement(text string) string {
	text = strings.ToLower(text)

	// Must enforcement
	if containsAny(text, []string{"must", "never", "always", "critical", "required"}) {
		return "must"
	}

	// Should enforcement
	if containsAny(text, []string{"should", "prefer", "recommend", "important"}) {
		return "should"
	}

	// May enforcement
	if containsAny(text, []string{"consider", "may", "optional"}) {
		return "may"
	}

	// Default to should
	return "should"
}

// inferScope attempts to infer file scope from text
func (c *MarkdownConverter) inferScope(text string) []urf.Scope {
	text = strings.ToLower(text)
	var patterns []string

	// Language-specific patterns
	if containsAny(text, []string{"typescript", "ts"}) {
		patterns = append(patterns, "**/*.ts", "**/*.tsx")
	}
	if containsAny(text, []string{"javascript", "js"}) {
		patterns = append(patterns, "**/*.js", "**/*.jsx")
	}
	if containsAny(text, []string{"python", "py"}) {
		patterns = append(patterns, "**/*.py")
	}
	if containsAny(text, []string{"go", "golang"}) {
		patterns = append(patterns, "**/*.go")
	}
	if containsAny(text, []string{"java"}) {
		patterns = append(patterns, "**/*.java")
	}

	// Default scope if no specific patterns found
	if len(patterns) == 0 {
		patterns = []string{"**/*"}
	}

	return []urf.Scope{{Files: patterns}}
}

// generateRuleBody creates the rule body content
func (c *MarkdownConverter) generateRuleBody(text string) string {
	return fmt.Sprintf("## %s\n\n%s", text, "Follow this guideline to maintain code quality and consistency.")
}

// containsAny checks if text contains any of the given substrings
func containsAny(text string, substrings []string) bool {
	for _, substring := range substrings {
		if strings.Contains(text, substring) {
			return true
		}
	}
	return false
}
