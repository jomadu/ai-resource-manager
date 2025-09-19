package ui

import (
	"encoding/json"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/manifest"
	"github.com/jomadu/ai-rules-manager/internal/version"
	"github.com/pterm/pterm"
)

// OutdatedRuleset represents a ruleset that has newer versions available.
type OutdatedRuleset struct {
	RulesetInfo *RulesetInfo `json:"rulesetInfo"`
	Wanted      string       `json:"wanted"`
	Latest      string       `json:"latest"`
}

// ManifestInfo contains information from the manifest file.
type ManifestInfo struct {
	Constraint string   `json:"constraint"`
	Priority   int      `json:"priority"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// InstallationInfo contains information about the actual installation.
type InstallationInfo struct {
	Version        string   `json:"version"`
	InstalledPaths []string `json:"installedPaths"`
}

// RulesetInfo provides detailed information about a ruleset.
type RulesetInfo struct {
	Registry     string           `json:"registry"`
	Name         string           `json:"name"`
	Manifest     ManifestInfo     `json:"manifest"`
	Installation InstallationInfo `json:"installation"`
}

// Interface defines the UI methods needed by the service
type Interface interface {
	// Progress reporting
	InstallStep(step string)
	InstallStepWithSpinner(step string) func(result string)
	InstallComplete(registry, ruleset, version string, sinks []string)
	Success(msg string)
	Error(err error)
	Warning(msg string)

	// Display operations
	ConfigList(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig)
	RulesetList(rulesets []*RulesetInfo)
	RulesetInfoGrouped(rulesets []*RulesetInfo, detailed bool)
	OutdatedTable(outdated []OutdatedRuleset, outputFormat string)
	VersionInfo(info version.VersionInfo)
}

// UI provides pterm-based user interface functionality
type UI struct {
	debug bool
}

// New creates a new UI instance
func New(debug bool) *UI {
	return &UI{debug: debug}
}

// Success displays a success message with checkmark
func (u *UI) Success(msg string) {
	pterm.Success.Println(msg)
}

// Error displays an error message
func (u *UI) Error(err error) {
	pterm.Error.Printf("Error: %v\n", err)
}

// Warning displays a warning message
func (u *UI) Warning(msg string) {
	pterm.Warning.Println(msg)
}

// Debug displays debug information if debug mode is enabled
func (u *UI) Debug(component, msg string, fields ...interface{}) {
	if u.debug {
		if len(fields) > 0 {
			pterm.Debug.Printf("[%s] %s %v\n", component, msg, fields)
		} else {
			pterm.Debug.Printf("[%s] %s\n", component, msg)
		}
	}
}

// VersionInfo displays version information
func (u *UI) VersionInfo(info version.VersionInfo) {
	pterm.Printf("arm %s\n", info.Version)
	if info.Commit != "" {
		pterm.Printf("commit: %s\n", info.Commit)
	}
	if info.Arch != "" {
		pterm.Printf("arch: %s\n", info.Arch)
	}
}

// ConfigList displays registries and sinks configuration
func (u *UI) ConfigList(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig) {
	if len(registries) > 0 {
		registryNodes := []pterm.TreeNode{}
		for name, config := range registries {
			regType, _ := config["type"].(string)
			url, _ := config["url"].(string)

			children := []pterm.TreeNode{
				{Text: fmt.Sprintf("type: %s", regType)},
				{Text: fmt.Sprintf("url: %s", url)},
			}

			switch regType {
			case "git":
				if branches, ok := config["branches"].([]interface{}); ok {
					branchStrs := make([]string, len(branches))
					for i, b := range branches {
						branchStrs[i] = fmt.Sprintf("%v", b)
					}
					children = append(children, pterm.TreeNode{
						Text: fmt.Sprintf("branches: %v", branchStrs),
					})
				}
			case "gitlab":
				if projectID, ok := config["project_id"].(string); ok && projectID != "" {
					children = append(children, pterm.TreeNode{
						Text: fmt.Sprintf("project-id: %s", projectID),
					})
				}
				if groupID, ok := config["group_id"].(string); ok && groupID != "" {
					children = append(children, pterm.TreeNode{
						Text: fmt.Sprintf("group-id: %s", groupID),
					})
				}
				if apiVersion, ok := config["api_version"].(string); ok && apiVersion != "" {
					children = append(children, pterm.TreeNode{
						Text: fmt.Sprintf("api-version: %s", apiVersion),
					})
				}
			}

			registryNodes = append(registryNodes, pterm.TreeNode{
				Text:     name,
				Children: children,
			})
		}

		tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
			Text:     "Registries",
			Children: registryNodes,
		})
		_ = tree.Render()
		pterm.Println()
	} else {
		pterm.Info.Println("Registries: (none configured)")
	}

	if len(sinks) > 0 {
		sinkNodes := []pterm.TreeNode{}
		for name, sink := range sinks {
			layout := sink.GetLayout()
			if layout == "" {
				layout = "hierarchical"
			}

			children := []pterm.TreeNode{
				{Text: fmt.Sprintf("directory: %s", sink.Directory)},
				{Text: fmt.Sprintf("layout: %s", layout)},
				{Text: fmt.Sprintf("target: %s", string(sink.CompileTarget))},
			}

			sinkNodes = append(sinkNodes, pterm.TreeNode{
				Text:     name,
				Children: children,
			})
		}

		tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
			Text:     "Sinks",
			Children: sinkNodes,
		})
		_ = tree.Render()
	} else {
		pterm.Info.Println("Sinks: (none configured)")
	}
}

// InstallStep displays a progress step
func (u *UI) InstallStep(step string) {
	pterm.Info.Printf("%s ✓\n", step)
}

// InstallStepWithSpinner starts a spinner and returns a function to complete it
func (u *UI) InstallStepWithSpinner(step string) func(result string) {
	spinner, _ := pterm.DefaultSpinner.Start(step)
	return func(result string) {
		spinner.Info(result)
	}
}

// InstallComplete displays final installation summary
func (u *UI) InstallComplete(registry, ruleset, version string, sinks []string) {
	if len(sinks) == 1 {
		pterm.Success.Printf("Installed %s/%s@%s (installed to %s sink)\n", registry, ruleset, version, sinks[0])
	} else {
		pterm.Success.Printf("Installed %s/%s@%s (installed to %d sinks)\n", registry, ruleset, version, len(sinks))
	}
}

// RulesetList displays installed rulesets
func (u *UI) RulesetList(rulesets []*RulesetInfo) {
	if len(rulesets) == 0 {
		pterm.Info.Println("No rulesets installed")
		return
	}

	for _, ruleset := range rulesets {
		pterm.Printf("%s/%s@%s - sinks: %v, priority: %d\n",
			ruleset.Registry, ruleset.Name,
			ruleset.Installation.Version,
			ruleset.Manifest.Sinks,
			ruleset.Manifest.Priority)
	}
}

// RulesetInfo displays detailed ruleset information
func (u *UI) RulesetInfo(info *RulesetInfo, detailed bool) {
	rulesetName := fmt.Sprintf("%s@%s", info.Name, info.Installation.Version)
	if info.Manifest.Constraint != "" {
		rulesetName += fmt.Sprintf(" (%s)", info.Manifest.Constraint)
	}

	children := []pterm.TreeNode{
		{Text: fmt.Sprintf("include: %v", info.Manifest.Include)},
		{Text: fmt.Sprintf("sinks: %v", info.Manifest.Sinks)},
		{Text: fmt.Sprintf("priority: %d", info.Manifest.Priority)},
		{Text: fmt.Sprintf("files: %d installed", len(info.Installation.InstalledPaths))},
	}

	if len(info.Manifest.Exclude) > 0 {
		children = append(children, pterm.TreeNode{
			Text: fmt.Sprintf("exclude: %v", info.Manifest.Exclude),
		})
	}

	if detailed && len(info.Installation.InstalledPaths) > 0 {
		pathNodes := []pterm.TreeNode{}
		for _, path := range info.Installation.InstalledPaths {
			pathNodes = append(pathNodes, pterm.TreeNode{Text: path})
		}
		children = append(children, pterm.TreeNode{
			Text:     "installed paths:",
			Children: pathNodes,
		})
	}

	rulesetNode := pterm.TreeNode{
		Text:     rulesetName,
		Children: children,
	}

	tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
		Text:     info.Registry,
		Children: []pterm.TreeNode{rulesetNode},
	})
	_ = tree.Render()
}

// RulesetInfoGrouped displays multiple rulesets grouped by registry
func (u *UI) RulesetInfoGrouped(rulesets []*RulesetInfo, detailed bool) {
	if len(rulesets) == 0 {
		pterm.Info.Println("No rulesets installed")
		return
	}

	// Group by registry
	registryGroups := make(map[string][]*RulesetInfo)
	for _, ruleset := range rulesets {
		registryGroups[ruleset.Registry] = append(registryGroups[ruleset.Registry], ruleset)
	}

	// Display each registry group
	for registry, groupRulesets := range registryGroups {
		rulesetNodes := []pterm.TreeNode{}

		for _, ruleset := range groupRulesets {
			rulesetName := fmt.Sprintf("%s@%s", ruleset.Name, ruleset.Installation.Version)
			if ruleset.Manifest.Constraint != "" {
				rulesetName += fmt.Sprintf(" (%s)", ruleset.Manifest.Constraint)
			}

			children := []pterm.TreeNode{
				{Text: fmt.Sprintf("include: %v", ruleset.Manifest.Include)},
				{Text: fmt.Sprintf("sinks: %v", ruleset.Manifest.Sinks)},
				{Text: fmt.Sprintf("priority: %d", ruleset.Manifest.Priority)},
				{Text: fmt.Sprintf("files: %d installed", len(ruleset.Installation.InstalledPaths))},
			}

			if len(ruleset.Manifest.Exclude) > 0 {
				children = append(children, pterm.TreeNode{
					Text: fmt.Sprintf("exclude: %v", ruleset.Manifest.Exclude),
				})
			}

			if detailed && len(ruleset.Installation.InstalledPaths) > 0 {
				pathNodes := []pterm.TreeNode{}
				for _, path := range ruleset.Installation.InstalledPaths {
					pathNodes = append(pathNodes, pterm.TreeNode{Text: path})
				}
				children = append(children, pterm.TreeNode{
					Text:     "installed paths:",
					Children: pathNodes,
				})
			}

			rulesetNodes = append(rulesetNodes, pterm.TreeNode{
				Text:     rulesetName,
				Children: children,
			})
		}

		tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
			Text:     registry,
			Children: rulesetNodes,
		})
		_ = tree.Render()
		pterm.Println()
	}
}

// OutdatedTable displays outdated rulesets in specified format
func (u *UI) OutdatedTable(outdated []OutdatedRuleset, outputFormat string) {
	if len(outdated) == 0 {
		pterm.Success.Println("All rulesets are up to date!")
		return
	}

	switch outputFormat {
	case "json":
		jsonData, err := json.Marshal(outdated)
		if err != nil {
			pterm.Error.Printf("Failed to marshal JSON: %v\n", err)
			return
		}
		fmt.Println(string(jsonData))
	default: // table format
		tableData := [][]string{
			{"Registry", "Ruleset", "Constraint", "Current", "Wanted", "Latest"},
		}

		for _, r := range outdated {
			tableData = append(tableData, []string{
				r.RulesetInfo.Registry,
				r.RulesetInfo.Name,
				r.RulesetInfo.Manifest.Constraint,
				r.RulesetInfo.Installation.Version,
				r.Wanted,
				r.Latest,
			})
		}

		_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	}
}
