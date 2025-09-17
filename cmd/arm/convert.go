package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/converter"
	"github.com/jomadu/ai-rules-manager/internal/urf"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newConvertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "convert [input-file]",
		Short: "Convert external rule formats to Universal Rule Format (URF)",
		Long: `Convert external rule documents from various AI tools and formats to Universal Rule Format (URF).

URF files can then be installed using 'arm install' and will be automatically compiled
to the appropriate format for each configured AI tool.

Supported input formats:
- Generic Markdown (.md, .markdown) - Extracts rules from headers and bullet points
- Cursor Rules (.mdc) - Converts Cursor-specific markdown with frontmatter
- Amazon Q Rules (.md in .amazonq/rules/) - Converts Amazon Q rule files
- GitHub Copilot Instructions (.md with Title/Description format)
- Plain text (.txt) - Extracts rules from structured text (coming soon)

Examples:
  arm convert guidelines.md                      # Convert generic markdown
  arm convert .cursor/rules/style.mdc           # Convert Cursor rules
  arm convert .amazonq/rules/security.md        # Convert Amazon Q rules
  arm convert .github/copilot-instructions.md   # Convert Copilot instructions
  arm convert team-rules.md --output=rules.yml  # Specify output file
  arm convert docs/style.md --dry-run            # Preview conversion

  # Merge into existing URF files:
  arm convert new-rules.md --merge --output=existing.yml     # Merge into existing URF
  arm convert style.mdc --merge --output=team-rules.yml     # Add Cursor rules to URF
  arm convert security.md --merge --on-conflict=rename      # Handle ID conflicts by renaming`,
		Args: cobra.ExactArgs(1),
		RunE: runConvert,
	}

	cmd.Flags().String("output", "", "Output URF file (default: input filename with .yml extension)")
	cmd.Flags().Bool("dry-run", false, "Show conversion preview without creating files")
	cmd.Flags().String("ruleset-id", "", "Custom ruleset ID (default: generated from filename)")
	cmd.Flags().String("ruleset-name", "", "Custom ruleset name (default: generated from filename)")
	cmd.Flags().String("version", "1.0.0", "Ruleset version")
	cmd.Flags().String("append-to", "", "Existing URF file to append rules to (creates new ruleset in same file)")
	cmd.Flags().Bool("merge", false, "Merge rules into existing URF file specified by --output")
	cmd.Flags().String("on-conflict", "rename", "How to handle rule ID conflicts: 'rename', 'skip', 'overwrite'")

	return cmd
}

func runConvert(cmd *cobra.Command, args []string) error {
	inputFile := args[0]

	// Validate input file exists
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("input file does not exist: %s", inputFile)
	}

	// Read input file
	content, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Determine input format from extension
	inputFormat := detectInputFormat(inputFile)
	if inputFormat == "" {
		return fmt.Errorf("unsupported input format for file: %s", inputFile)
	}

	// Get flags
	outputFile, _ := cmd.Flags().GetString("output")
	dryRun, _ := cmd.Flags().GetBool("dry-run")
	rulesetID, _ := cmd.Flags().GetString("ruleset-id")
	rulesetName, _ := cmd.Flags().GetString("ruleset-name")
	version, _ := cmd.Flags().GetString("version")
	appendTo, _ := cmd.Flags().GetString("append-to")
	merge, _ := cmd.Flags().GetBool("merge")
	onConflict, _ := cmd.Flags().GetString("on-conflict")

	// Auto-generate output filename if not specified
	if outputFile == "" {
		base := strings.TrimSuffix(inputFile, filepath.Ext(inputFile))
		outputFile = base + ".yml"
	}

	// Auto-generate ruleset ID and name if not specified
	if rulesetID == "" {
		rulesetID = generateRulesetID(inputFile)
	}
	if rulesetName == "" {
		rulesetName = generateRulesetName(inputFile)
	}

	fmt.Printf("Converting %s (%s format) to URF...\n", inputFile, inputFormat)

	// Parse input to URF
	var conv converter.Converter
	switch inputFormat {
	case "markdown":
		conv = converter.NewMarkdownConverter()
	case "cursor-mdc":
		conv = converter.NewCursorMDCConverter()
	case "amazonq-md":
		conv = converter.NewAmazonQMDConverter()
	case "copilot-instructions":
		conv = converter.NewCopilotInstructionsConverter()
	case "plaintext":
		return fmt.Errorf("plain text conversion not yet implemented")
	default:
		return fmt.Errorf("unsupported input format: %s", inputFormat)
	}

	config := converter.ConversionConfig{
		RulesetID:   rulesetID,
		RulesetName: rulesetName,
		Version:     version,
		InputFile:   inputFile,
	}

	ruleset, err := conv.Convert(string(content), config)
	if err != nil {
		return fmt.Errorf("failed to convert %s: %w", inputFile, err)
	}

	// Show conversion results
	fmt.Printf("✅ Successfully parsed %d rules\n", len(ruleset.Rules))

	if dryRun {
		fmt.Println("\n🔍 DRY RUN - Conversion preview:")
		return printConversionPreview(ruleset, outputFile)
	}

	// Handle merge/append scenarios
	var finalRuleset *converter.URFRuleset
	var finalOutputFile string

	switch {
	case appendTo != "":
		// Append to existing file as new ruleset
		finalRuleset, finalOutputFile, err = appendToURFFile(ruleset, appendTo, onConflict)
		if err != nil {
			return fmt.Errorf("failed to append to URF file: %w", err)
		}
	case merge:
		// Merge rules into existing file
		finalRuleset, finalOutputFile, err = mergeWithURFFile(ruleset, outputFile, onConflict)
		if err != nil {
			return fmt.Errorf("failed to merge with URF file: %w", err)
		}
	default:
		// Create new file
		finalRuleset = ruleset
		finalOutputFile = outputFile
	}

	// Save URF file
	if err := saveURFFile(finalRuleset, finalOutputFile); err != nil {
		return fmt.Errorf("failed to save URF file: %w", err)
	}

	fmt.Printf("✅ URF file saved to: %s\n", finalOutputFile)
	if appendTo != "" || merge {
		fmt.Printf("📝 Rules merged/appended successfully\n")
	}
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  arm install %s --sinks <sink-names>\n", finalOutputFile)

	return nil
}

// detectInputFormat determines the input format from file extension and content
func detectInputFormat(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".md", ".markdown":
		// Could be generic markdown, Amazon Q, or Copilot instructions
		return detectMarkdownVariant(filename)
	case ".mdc":
		return "cursor-mdc"
	case ".txt":
		return "plaintext"
	default:
		return ""
	}
}

// detectMarkdownVariant determines the specific markdown variant
func detectMarkdownVariant(filename string) string {
	// Check file path patterns
	if strings.Contains(filename, ".amazonq") || strings.Contains(filename, "amazonq") {
		return "amazonq-md"
	}
	if strings.Contains(filename, ".github") && strings.Contains(filename, "copilot") {
		return "copilot-instructions"
	}
	if strings.Contains(filename, "instructions") || strings.Contains(filename, "copilot") {
		return "copilot-instructions"
	}

	// Default to generic markdown
	return "markdown"
}

// generateRulesetID creates a ruleset ID from filename
func generateRulesetID(filename string) string {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Convert to kebab-case
	id := strings.ToLower(name)
	id = strings.ReplaceAll(id, " ", "-")
	id = strings.ReplaceAll(id, "_", "-")

	// Remove special characters except hyphens
	var result strings.Builder
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			result.WriteRune(r)
		}
	}

	return result.String()
}

// generateRulesetName creates a human-readable name from filename
func generateRulesetName(filename string) string {
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))

	// Convert to title case
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ToUpper(name[:1]) + name[1:]

	if name == "" {
		name = "Converted Rules"
	}

	return name
}

// saveURFFile saves the URF ruleset to a YAML file
func saveURFFile(ruleset *converter.URFRuleset, outputFile string) error {
	return ruleset.SaveToFile(outputFile)
}

// appendToURFFile appends rules to an existing URF file as a new ruleset
func appendToURFFile(newRuleset *converter.URFRuleset, targetFile, onConflict string) (*converter.URFRuleset, string, error) {
	// For now, URF files contain single rulesets, so this would mean creating a new file
	// In the future, we might support multi-ruleset URF files
	return nil, "", fmt.Errorf("append-to functionality not yet implemented - URF files currently support single rulesets")
}

// mergeWithURFFile merges rules into an existing URF file
func mergeWithURFFile(newRuleset *converter.URFRuleset, outputFile, onConflict string) (*converter.URFRuleset, string, error) {
	// Check if output file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// File doesn't exist, return new ruleset
		return newRuleset, outputFile, nil
	}

	// Load existing URF file
	existingRuleset, err := loadExistingURFFile(outputFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to load existing URF file: %w", err)
	}

	// Merge rules with conflict resolution
	mergedRules, err := mergeRulesWithConflictResolution(existingRuleset.Rules, newRuleset.Rules, onConflict)
	if err != nil {
		return nil, "", fmt.Errorf("failed to merge rules: %w", err)
	}

	// Update existing ruleset with merged rules
	existingRuleset.Rules = mergedRules

	// Update metadata
	existingRuleset.Metadata.Description = fmt.Sprintf("%s (updated with additional rules)", existingRuleset.Metadata.Description)

	return existingRuleset, outputFile, nil
}

// loadExistingURFFile loads an existing URF file
func loadExistingURFFile(filename string) (*converter.URFRuleset, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var ruleset urf.Ruleset
	if err := yaml.Unmarshal(data, &ruleset); err != nil {
		return nil, fmt.Errorf("invalid URF file: %w", err)
	}

	return &converter.URFRuleset{Ruleset: &ruleset}, nil
}

// mergeRulesWithConflictResolution merges rule slices with conflict handling
func mergeRulesWithConflictResolution(existing, newRules []urf.Rule, onConflict string) ([]urf.Rule, error) {
	// Create map of existing rule IDs for quick lookup
	existingIDs := make(map[string]int)
	for i, rule := range existing {
		existingIDs[rule.ID] = i
	}

	result := make([]urf.Rule, len(existing))
	copy(result, existing)

	var conflicts []string

	for _, newRule := range newRules {
		if existingIndex, exists := existingIDs[newRule.ID]; exists {
			// Handle conflict
			switch onConflict {
			case "rename":
				// Rename the new rule ID
				newRule.ID = generateUniqueRuleID(newRule.ID, existingIDs)
				result = append(result, newRule)
				conflicts = append(conflicts, fmt.Sprintf("renamed %s", newRule.ID))
			case "skip":
				// Skip the conflicting rule
				conflicts = append(conflicts, fmt.Sprintf("skipped %s", newRule.ID))
				continue
			case "overwrite":
				// Overwrite the existing rule
				result[existingIndex] = newRule
				conflicts = append(conflicts, fmt.Sprintf("overwritten %s", newRule.ID))
			default:
				return nil, fmt.Errorf("invalid conflict resolution strategy: %s", onConflict)
			}
		} else {
			// No conflict, add the rule
			result = append(result, newRule)
		}
	}

	if len(conflicts) > 0 {
		fmt.Printf("⚠️  Resolved %d rule ID conflicts: %s\n", len(conflicts), strings.Join(conflicts, ", "))
	}

	return result, nil
}

// generateUniqueRuleID generates a unique rule ID by appending a number
func generateUniqueRuleID(baseID string, existingIDs map[string]int) string {
	counter := 1
	for {
		newID := fmt.Sprintf("%s-%d", baseID, counter)
		if _, exists := existingIDs[newID]; !exists {
			existingIDs[newID] = -1 // Mark as used
			return newID
		}
		counter++
	}
}

// printConversionPreview shows a preview of the conversion results
func printConversionPreview(ruleset *converter.URFRuleset, outputFile string) error {
	fmt.Printf("\nOutput file: %s\n", outputFile)
	fmt.Printf("Ruleset ID: %s\n", ruleset.Metadata.ID)
	fmt.Printf("Ruleset Name: %s\n", ruleset.Metadata.Name)
	fmt.Printf("Version: %s\n", ruleset.Metadata.Version)

	fmt.Printf("\nRules (%d total):\n", len(ruleset.Rules))
	for i, rule := range ruleset.Rules {
		fmt.Printf("  %d. %s (ID: %s)\n", i+1, rule.Name, rule.ID)
		fmt.Printf("     Priority: %d, Enforcement: %s\n", rule.Priority, rule.Enforcement)
		if len(rule.Scope) > 0 && len(rule.Scope[0].Files) > 0 {
			fmt.Printf("     Scope: %s\n", strings.Join(rule.Scope[0].Files, ", "))
		}
	}

	return nil
}
