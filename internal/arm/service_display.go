package arm

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jomadu/ai-rules-manager/internal/ui"
)

func (a *ArmService) ShowVersion() error {
	a.ui.VersionInfo(a.Version())
	return nil
}

func (a *ArmService) ShowRulesetInfo(ctx context.Context, rulesets []string) error {
	if len(rulesets) == 0 {
		infos, err := a.listInstalledRulesets(ctx)
		if err != nil {
			return err
		}
		a.ui.RulesetInfoGrouped(infos, false)
		return nil
	}

	var infos []*RulesetInfo
	for _, rulesetArg := range rulesets {
		parts := strings.Split(rulesetArg, "/")
		if len(parts) != 2 {
			return fmt.Errorf("invalid ruleset format: %s (expected registry/ruleset)", rulesetArg)
		}
		info, err := a.getRulesetInfo(ctx, parts[0], parts[1])
		if err != nil {
			return err
		}
		infos = append(infos, info)
	}

	detailed := len(rulesets) == 1
	a.ui.RulesetInfoGrouped(infos, detailed)
	return nil
}

func (a *ArmService) ShowOutdated(ctx context.Context, outputFormat string, noSpinner bool) error {
	if noSpinner || outputFormat == "json" {
		outdated, err := a.getOutdatedRulesets(ctx)
		if err != nil {
			return err
		}
		// Convert to unified format
		packages := a.convertOutdatedRulesetsToPackages(outdated)
		a.ui.OutdatedTable(packages, outputFormat)
	} else {
		finishChecking := a.ui.InstallStepWithSpinner("Checking for updates...")
		outdated, err := a.getOutdatedRulesets(ctx)
		if err != nil {
			return err
		}
		finishChecking(fmt.Sprintf("Found %d outdated rulesets", len(outdated)))
		// Convert to unified format
		packages := a.convertOutdatedRulesetsToPackages(outdated)
		a.ui.OutdatedTable(packages, outputFormat)
	}
	return nil
}

func (a *ArmService) ShowList(ctx context.Context, sortByPriority bool) error {
	rulesets, err := a.listInstalledRulesets(ctx)
	if err != nil {
		return err
	}

	if sortByPriority {
		sort.Slice(rulesets, func(i, j int) bool {
			return rulesets[i].Manifest.Priority > rulesets[j].Manifest.Priority
		})
	}

	a.ui.RulesetList(rulesets)
	return nil
}

// convertOutdatedRulesetsToPackages converts OutdatedRuleset to OutdatedPackage format
func (a *ArmService) convertOutdatedRulesetsToPackages(outdated []OutdatedRuleset) []ui.OutdatedPackage {
	packages := make([]ui.OutdatedPackage, len(outdated))
	for i, ruleset := range outdated {
		packages[i] = ui.OutdatedPackage{
			Package:    fmt.Sprintf("%s/%s", ruleset.RulesetInfo.Registry, ruleset.RulesetInfo.Name),
			Type:       "ruleset",
			Constraint: ruleset.RulesetInfo.Manifest.Constraint,
			Current:    ruleset.RulesetInfo.Installation.Version,
			Wanted:     ruleset.Wanted,
			Latest:     ruleset.Latest,
		}
	}
	return packages
}
