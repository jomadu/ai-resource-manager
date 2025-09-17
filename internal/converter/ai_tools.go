package converter

import (
	"bufio"
	"fmt"
	"regexp"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/urf"
	"gopkg.in/yaml.v3"
)

// AI Tool Format Converters
// These convert FROM existing AI tool formats TO URF

// CursorMDCConverter converts Cursor .mdc files to URF
type CursorMDCConverter struct {
	frontmatterPattern *regexp.Regexp
	headerPattern      *regexp.Regexp
	bulletPattern      *regexp.Regexp
	nextRuleID         int
}

// NewCursorMDCConverter creates a new Cursor .mdc converter
func NewCursorMDCConverter() Converter {
	return &CursorMDCConverter{
		frontmatterPattern: regexp.MustCompile(`^---\s*$`),
		headerPattern:      regexp.MustCompile(`^(#{1,6})\s+(.+)$`),
		bulletPattern:      regexp.MustCompile(`^[-*+]\s+(.+)$`),
		nextRuleID:         1,
	}
}

// Convert converts Cursor .mdc content to URF
func (c *CursorMDCConverter) Convert(content string, config ConversionConfig) (*URFRuleset, error) {
	// Initialize URF ruleset
	ruleset := &urf.Ruleset{
		Version: "1.0",
		Metadata: urf.Metadata{
			ID:          config.RulesetID,
			Name:        config.RulesetName,
			Version:     config.Version,
			Description: fmt.Sprintf("Rules converted from Cursor .mdc file: %s", config.InputFile),
		},
		Rules: make([]urf.Rule, 0),
	}

	// Parse frontmatter and content
	frontmatter, body, err := c.parseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Extract global metadata from frontmatter
	globalDescription := c.extractFromFrontmatter(frontmatter, "description")
	globalGlobs := c.extractArrayFromFrontmatter(frontmatter, "globs")
	alwaysApply := c.extractBoolFromFrontmatter(frontmatter, "alwaysApply")

	// Update ruleset description if available
	if globalDescription != "" {
		ruleset.Metadata.Description = globalDescription
	}

	// Parse body content for rules
	rules, err := c.parseBodyForRules(body, globalGlobs, alwaysApply)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rules: %w", err)
	}

	ruleset.Rules = rules
	return &URFRuleset{Ruleset: ruleset}, nil
}

// AmazonQMDConverter converts Amazon Q .md files to URF
type AmazonQMDConverter struct {
	frontmatterPattern *regexp.Regexp
	headerPattern      *regexp.Regexp
	bulletPattern      *regexp.Regexp
	nextRuleID         int
}

// NewAmazonQMDConverter creates a new Amazon Q .md converter
func NewAmazonQMDConverter() Converter {
	return &AmazonQMDConverter{
		frontmatterPattern: regexp.MustCompile(`^---\s*$`),
		headerPattern:      regexp.MustCompile(`^(#{1,6})\s+(.+)$`),
		bulletPattern:      regexp.MustCompile(`^[-*+]\s+(.+)$`),
		nextRuleID:         1,
	}
}

// Convert converts Amazon Q .md content to URF
func (c *AmazonQMDConverter) Convert(content string, config ConversionConfig) (*URFRuleset, error) {
	// Initialize URF ruleset
	ruleset := &urf.Ruleset{
		Version: "1.0",
		Metadata: urf.Metadata{
			ID:          config.RulesetID,
			Name:        config.RulesetName,
			Version:     config.Version,
			Description: fmt.Sprintf("Rules converted from Amazon Q .md file: %s", config.InputFile),
		},
		Rules: make([]urf.Rule, 0),
	}

	// Parse frontmatter and content (similar to Cursor but different structure)
	frontmatter, body, err := c.parseFrontmatter(content)
	if err != nil {
		return nil, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Extract global metadata
	globalDescription := c.extractFromFrontmatter(frontmatter, "description")
	if globalDescription != "" {
		ruleset.Metadata.Description = globalDescription
	}

	// Parse body content for rules
	rules, err := c.parseBodyForRules(body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse rules: %w", err)
	}

	ruleset.Rules = rules
	return &URFRuleset{Ruleset: ruleset}, nil
}

// CopilotInstructionsConverter converts GitHub Copilot instructions to URF
type CopilotInstructionsConverter struct {
	titlePattern       *regexp.Regexp
	descriptionPattern *regexp.Regexp
	pathPattern        *regexp.Regexp
	nextRuleID         int
}

// NewCopilotInstructionsConverter creates a new Copilot instructions converter
func NewCopilotInstructionsConverter() Converter {
	return &CopilotInstructionsConverter{
		titlePattern:       regexp.MustCompile(`^\*\*Title:\*\*\s*(.+)$`),
		descriptionPattern: regexp.MustCompile(`^\*\*Description:\*\*\s*(.+)$`),
		pathPattern:        regexp.MustCompile(`^\*\*Path patterns?\*\*:\s*(.+)$`),
		nextRuleID:         1,
	}
}

// Convert converts Copilot instructions content to URF
func (c *CopilotInstructionsConverter) Convert(content string, config ConversionConfig) (*URFRuleset, error) {
	// Initialize URF ruleset
	ruleset := &urf.Ruleset{
		Version: "1.0",
		Metadata: urf.Metadata{
			ID:          config.RulesetID,
			Name:        config.RulesetName,
			Version:     config.Version,
			Description: fmt.Sprintf("Rules converted from GitHub Copilot instructions: %s", config.InputFile),
		},
		Rules: make([]urf.Rule, 0),
	}

	// Parse Copilot instruction blocks
	rules := c.parseInstructionBlocks(content)

	ruleset.Rules = rules
	return &URFRuleset{Ruleset: ruleset}, nil
}

// Helper methods for CursorMDCConverter
func (c *CursorMDCConverter) parseFrontmatter(content string) (frontmatter map[string]interface{}, body string, err error) {
	lines := strings.Split(content, "\n")

	// Check if content starts with frontmatter
	if len(lines) == 0 || !c.frontmatterPattern.MatchString(lines[0]) {
		// No frontmatter, return empty map and full content
		return make(map[string]interface{}), content, nil
	}

	// Find end of frontmatter
	var frontmatterEnd int
	for i := 1; i < len(lines); i++ {
		if c.frontmatterPattern.MatchString(lines[i]) {
			frontmatterEnd = i
			break
		}
	}

	if frontmatterEnd == 0 {
		return nil, "", fmt.Errorf("unclosed frontmatter")
	}

	// Parse YAML frontmatter
	frontmatterContent := strings.Join(lines[1:frontmatterEnd], "\n")
	var frontmatterData map[string]interface{}
	if err := yaml.Unmarshal([]byte(frontmatterContent), &frontmatterData); err != nil {
		return nil, "", fmt.Errorf("invalid YAML frontmatter: %w", err)
	}

	// Return body content after frontmatter
	bodyContent := strings.Join(lines[frontmatterEnd+1:], "\n")
	return frontmatterData, bodyContent, nil
}

func (c *CursorMDCConverter) extractFromFrontmatter(frontmatter map[string]interface{}, key string) string {
	if val, ok := frontmatter[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func (c *CursorMDCConverter) extractArrayFromFrontmatter(frontmatter map[string]interface{}, key string) []string {
	if val, ok := frontmatter[key]; ok {
		if arr, ok := val.([]interface{}); ok {
			result := make([]string, 0, len(arr))
			for _, item := range arr {
				if str, ok := item.(string); ok {
					result = append(result, str)
				}
			}
			return result
		}
	}
	return nil
}

func (c *CursorMDCConverter) extractBoolFromFrontmatter(frontmatter map[string]interface{}, key string) bool {
	if val, ok := frontmatter[key]; ok {
		if b, ok := val.(bool); ok {
			return b
		}
	}
	return false
}

func (c *CursorMDCConverter) parseBodyForRules(body string, globalGlobs []string, alwaysApply bool) ([]urf.Rule, error) {
	var rules []urf.Rule
	var currentRule *urf.Rule
	var currentSection string

	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// Handle headers (potential rule titles or sections)
		if matches := c.headerPattern.FindStringSubmatch(line); matches != nil {
			// Finalize current rule
			if currentRule != nil {
				rules = append(rules, *currentRule)
			}

			// Start new rule from header
			level := len(matches[1])
			title := matches[2]

			if level <= 2 {
				// H1/H2 are sections
				currentSection = strings.ToLower(title)
				currentRule = nil
			} else {
				// H3+ are rules
				currentRule = &urf.Rule{
					ID:          c.generateRuleID(title),
					Name:        title,
					Description: fmt.Sprintf("Rule from %s section", currentSection),
					Priority:    80,
					Enforcement: c.determineEnforcement(alwaysApply),
					Scope:       c.createScope(globalGlobs),
					Body:        fmt.Sprintf("## %s\n", title),
				}
			}
			continue
		}

		// Handle bullet points
		if matches := c.bulletPattern.FindStringSubmatch(line); matches != nil {
			text := matches[1]

			// If no current rule, create one from bullet point
			if currentRule == nil {
				currentRule = &urf.Rule{
					ID:          c.generateRuleID(text),
					Name:        text,
					Description: fmt.Sprintf("Rule from %s section", currentSection),
					Priority:    80,
					Enforcement: c.determineEnforcement(alwaysApply),
					Scope:       c.createScope(globalGlobs),
					Body:        fmt.Sprintf("## %s\n", text),
				}
			} else {
				// Add to current rule body
				currentRule.Body += fmt.Sprintf("- %s\n", text)
			}
			continue
		}

		// Add regular text to current rule body
		if currentRule != nil {
			currentRule.Body += line + "\n"
		}
	}

	// Finalize last rule
	if currentRule != nil {
		rules = append(rules, *currentRule)
	}

	return rules, scanner.Err()
}

func (c *CursorMDCConverter) generateRuleID(text string) string {
	id := strings.ToLower(text)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "'", "")
	id = strings.ReplaceAll(id, "\"", "")

	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	finalID := fmt.Sprintf("%s-%d", result.String(), c.nextRuleID)
	c.nextRuleID++
	return finalID
}

func (c *CursorMDCConverter) determineEnforcement(alwaysApply bool) string {
	if alwaysApply {
		return "must"
	}
	return "should"
}

func (c *CursorMDCConverter) createScope(globs []string) []urf.Scope {
	if len(globs) == 0 {
		return []urf.Scope{{Files: []string{"**/*"}}}
	}
	return []urf.Scope{{Files: globs}}
}

// Helper methods for AmazonQMDConverter
func (c *AmazonQMDConverter) parseFrontmatter(content string) (frontmatter map[string]interface{}, body string, err error) {
	// Amazon Q files typically don't have frontmatter, just return empty map and full content
	return make(map[string]interface{}), content, nil
}

func (c *AmazonQMDConverter) extractFromFrontmatter(frontmatter map[string]interface{}, key string) string {
	// Amazon Q doesn't use frontmatter, always return empty
	return ""
}

func (c *AmazonQMDConverter) parseBodyForRules(body string) ([]urf.Rule, error) {
	// Amazon Q format is simpler - just markdown headers and content
	var rules []urf.Rule
	var currentRule *urf.Rule

	scanner := bufio.NewScanner(strings.NewReader(body))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		// Handle headers as rules
		if matches := c.headerPattern.FindStringSubmatch(line); matches != nil {
			// Finalize current rule
			if currentRule != nil {
				rules = append(rules, *currentRule)
			}

			// Start new rule
			title := matches[2]
			currentRule = &urf.Rule{
				ID:          c.generateRuleID(title),
				Name:        title,
				Description: "Rule converted from Amazon Q format",
				Priority:    80,
				Enforcement: "should",
				Scope:       []urf.Scope{{Files: []string{"**/*"}}},
				Body:        fmt.Sprintf("## %s\n", title),
			}
			continue
		}

		// Add content to current rule
		if currentRule != nil {
			currentRule.Body += line + "\n"
		}
	}

	// Finalize last rule
	if currentRule != nil {
		rules = append(rules, *currentRule)
	}

	return rules, scanner.Err()
}

func (c *AmazonQMDConverter) generateRuleID(text string) string {
	return (&CursorMDCConverter{nextRuleID: c.nextRuleID}).generateRuleID(text)
}

// Helper methods for CopilotInstructionsConverter
func (c *CopilotInstructionsConverter) parseInstructionBlocks(content string) []urf.Rule {
	var rules []urf.Rule

	blocks := c.splitIntoBlocks(content)

	for _, block := range blocks {
		rule, err := c.parseInstructionBlock(block)
		if err != nil {
			continue // Skip invalid blocks
		}
		rules = append(rules, rule)
	}

	return rules
}

func (c *CopilotInstructionsConverter) splitIntoBlocks(content string) []string {
	// Split by empty lines to get instruction blocks
	blocks := strings.Split(content, "\n\n")
	var result []string

	for _, block := range blocks {
		block = strings.TrimSpace(block)
		if block != "" {
			result = append(result, block)
		}
	}

	return result
}

func (c *CopilotInstructionsConverter) parseInstructionBlock(block string) (urf.Rule, error) {
	var title, description, pathPatterns string

	lines := strings.Split(block, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		if matches := c.titlePattern.FindStringSubmatch(line); matches != nil {
			title = matches[1]
		} else if matches := c.descriptionPattern.FindStringSubmatch(line); matches != nil {
			description = matches[1]
		} else if matches := c.pathPattern.FindStringSubmatch(line); matches != nil {
			pathPatterns = matches[1]
		}
	}

	if title == "" {
		return urf.Rule{}, fmt.Errorf("missing title in instruction block")
	}

	// Create scope from path patterns
	var scope []urf.Scope
	if pathPatterns != "" {
		patterns := strings.Split(pathPatterns, ",")
		for i, pattern := range patterns {
			patterns[i] = strings.TrimSpace(pattern)
		}
		scope = []urf.Scope{{Files: patterns}}
	} else {
		scope = []urf.Scope{{Files: []string{"**/*"}}}
	}

	rule := urf.Rule{
		ID:          c.generateRuleID(title),
		Name:        title,
		Description: description,
		Priority:    80,
		Enforcement: "should",
		Scope:       scope,
		Body:        fmt.Sprintf("## %s\n\n%s", title, description),
	}

	return rule, nil
}

func (c *CopilotInstructionsConverter) generateRuleID(text string) string {
	id := strings.ToLower(text)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "'", "")
	id = strings.ReplaceAll(id, "\"", "")

	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	finalID := fmt.Sprintf("%s-%d", result.String(), c.nextRuleID)
	c.nextRuleID++
	return finalID
}
