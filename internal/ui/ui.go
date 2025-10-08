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

// OutdatedPromptset represents a promptset that has newer versions available.
type OutdatedPromptset struct {
	PromptsetInfo *PromptsetInfo `json:"promptsetInfo"`
	Wanted        string         `json:"wanted"`
	Latest        string         `json:"latest"`
}

// OutdatedPackage represents either a ruleset or promptset that has newer versions available.
type OutdatedPackage struct {
	Package    string `json:"package"`
	Type       string `json:"type"` // "ruleset" or "promptset"
	Constraint string `json:"constraint"`
	Current    string `json:"current"`
	Wanted     string `json:"wanted"`
	Latest     string `json:"latest"`
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

// PromptsetInfo provides detailed information about a promptset.
type PromptsetInfo struct {
	Registry     string                `json:"registry"`
	Name         string                `json:"name"`
	Manifest     PromptsetManifestInfo `json:"manifest"`
	Installation InstallationInfo      `json:"installation"`
}

// PromptsetManifestInfo contains information from the manifest file for promptsets.
type PromptsetManifestInfo struct {
	Constraint string   `json:"constraint"`
	Include    []string `json:"include"`
	Exclude    []string `json:"exclude"`
	Sinks      []string `json:"sinks"`
}

// CompileStats tracks compilation statistics
type CompileStats struct {
	FilesProcessed int `json:"filesProcessed"`
	FilesCompiled  int `json:"filesCompiled"`
	FilesSkipped   int `json:"filesSkipped"`
	RulesGenerated int `json:"rulesGenerated"`
	Errors         int `json:"errors"`
}

// Interface defines the UI methods needed by the service
type Interface interface {
	// Progress reporting
	InstallStep(step string)
	InstallStepWithSpinner(step string) func(result string)
	InstallComplete(registry, resource, version, resourceType string, sinks []string)
	Success(msg string)
	Error(err error)
	Warning(msg string)

	// Display operations
	ConfigList(registries map[string]map[string]interface{}, sinks map[string]manifest.SinkConfig)
	RulesetList(rulesets []*RulesetInfo)
	PromptsetList(promptsets []*PromptsetInfo)
	PackageList(rulesets []*RulesetInfo, promptsets []*PromptsetInfo)
	RulesetInfoGrouped(rulesets []*RulesetInfo, detailed bool)
	PromptsetInfoGrouped(promptsets []*PromptsetInfo, detailed bool)
	PackageInfoGrouped(rulesets []*RulesetInfo, promptsets []*PromptsetInfo, detailed bool)
	OutdatedTable(outdated []OutdatedPackage, outputFormat string)
	VersionInfo(info version.VersionInfo)

	// Compile operations
	CompileStep(step string)
	CompileComplete(stats CompileStats, validateOnly bool)
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
func (u *UI) InstallComplete(registry, resource, version, resourceType string, sinks []string) {
	if len(sinks) == 1 {
		pterm.Success.Printf("Installed %s/%s@%s (%s installed to %s sink)\n", registry, resource, version, resourceType, sinks[0])
	} else {
		pterm.Success.Printf("Installed %s/%s@%s (%s installed to %d sinks)\n", registry, resource, version, resourceType, len(sinks))
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

// PromptsetList displays installed promptsets
func (u *UI) PromptsetList(promptsets []*PromptsetInfo) {
	if len(promptsets) == 0 {
		pterm.Info.Println("No promptsets installed")
		return
	}

	for _, promptset := range promptsets {
		pterm.Printf("%s/%s@%s - sinks: %v\n",
			promptset.Registry, promptset.Name,
			promptset.Installation.Version,
			promptset.Manifest.Sinks)
	}
}

// PackageList displays both rulesets and promptsets in a unified format
func (u *UI) PackageList(rulesets []*RulesetInfo, promptsets []*PromptsetInfo) {
	if len(rulesets) == 0 && len(promptsets) == 0 {
		pterm.Info.Println("No packages installed")
		return
	}

	// Display rulesets
	for _, ruleset := range rulesets {
		pterm.Printf("%s/%s@%s (ruleset) - sinks: %v, priority: %d\n",
			ruleset.Registry, ruleset.Name,
			ruleset.Installation.Version,
			ruleset.Manifest.Sinks,
			ruleset.Manifest.Priority)
	}

	// Display promptsets
	for _, promptset := range promptsets {
		pterm.Printf("%s/%s@%s (promptset) - sinks: %v\n",
			promptset.Registry, promptset.Name,
			promptset.Installation.Version,
			promptset.Manifest.Sinks)
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

// PromptsetInfo displays detailed promptset information
func (u *UI) PromptsetInfo(info *PromptsetInfo, detailed bool) {
	promptsetName := fmt.Sprintf("%s@%s", info.Name, info.Installation.Version)
	if info.Manifest.Constraint != "" {
		promptsetName += fmt.Sprintf(" (%s)", info.Manifest.Constraint)
	}

	children := []pterm.TreeNode{
		{Text: fmt.Sprintf("include: %v", info.Manifest.Include)},
		{Text: fmt.Sprintf("sinks: %v", info.Manifest.Sinks)},
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

	promptsetNode := pterm.TreeNode{
		Text:     promptsetName,
		Children: children,
	}

	tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
		Text:     info.Registry,
		Children: []pterm.TreeNode{promptsetNode},
	})
	_ = tree.Render()
}

// PromptsetInfoGrouped displays multiple promptsets grouped by registry
func (u *UI) PromptsetInfoGrouped(promptsets []*PromptsetInfo, detailed bool) {
	if len(promptsets) == 0 {
		pterm.Info.Println("No promptsets installed")
		return
	}

	// Group by registry
	registryGroups := make(map[string][]*PromptsetInfo)
	for _, promptset := range promptsets {
		registryGroups[promptset.Registry] = append(registryGroups[promptset.Registry], promptset)
	}

	// Display each registry group
	for registry, groupPromptsets := range registryGroups {
		promptsetNodes := []pterm.TreeNode{}

		for _, promptset := range groupPromptsets {
			promptsetName := fmt.Sprintf("%s@%s", promptset.Name, promptset.Installation.Version)
			if promptset.Manifest.Constraint != "" {
				promptsetName += fmt.Sprintf(" (%s)", promptset.Manifest.Constraint)
			}

			children := []pterm.TreeNode{
				{Text: fmt.Sprintf("include: %v", promptset.Manifest.Include)},
				{Text: fmt.Sprintf("sinks: %v", promptset.Manifest.Sinks)},
				{Text: fmt.Sprintf("files: %d installed", len(promptset.Installation.InstalledPaths))},
			}

			if len(promptset.Manifest.Exclude) > 0 {
				children = append(children, pterm.TreeNode{
					Text: fmt.Sprintf("exclude: %v", promptset.Manifest.Exclude),
				})
			}

			if detailed && len(promptset.Installation.InstalledPaths) > 0 {
				pathNodes := []pterm.TreeNode{}
				for _, path := range promptset.Installation.InstalledPaths {
					pathNodes = append(pathNodes, pterm.TreeNode{Text: path})
				}
				children = append(children, pterm.TreeNode{
					Text:     "installed paths:",
					Children: pathNodes,
				})
			}

			promptsetNodes = append(promptsetNodes, pterm.TreeNode{
				Text:     promptsetName,
				Children: children,
			})
		}

		tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
			Text:     registry,
			Children: promptsetNodes,
		})
		_ = tree.Render()
		pterm.Println()
	}
}

// PackageInfoGrouped displays both rulesets and promptsets grouped by registry
func (u *UI) PackageInfoGrouped(rulesets []*RulesetInfo, promptsets []*PromptsetInfo, detailed bool) {
	if len(rulesets) == 0 && len(promptsets) == 0 {
		pterm.Info.Println("No packages installed")
		return
	}

	// Group by registry
	registryGroups := make(map[string]struct {
		rulesets   []*RulesetInfo
		promptsets []*PromptsetInfo
	})

	for _, ruleset := range rulesets {
		group := registryGroups[ruleset.Registry]
		group.rulesets = append(group.rulesets, ruleset)
		registryGroups[ruleset.Registry] = group
	}

	for _, promptset := range promptsets {
		group := registryGroups[promptset.Registry]
		group.promptsets = append(group.promptsets, promptset)
		registryGroups[promptset.Registry] = group
	}

	// Display each registry group
	for registry, group := range registryGroups {
		packageNodes := []pterm.TreeNode{}

		// Add rulesets
		for _, ruleset := range group.rulesets {
			rulesetName := fmt.Sprintf("%s@%s (ruleset)", ruleset.Name, ruleset.Installation.Version)
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

			packageNodes = append(packageNodes, pterm.TreeNode{
				Text:     rulesetName,
				Children: children,
			})
		}

		// Add promptsets
		for _, promptset := range group.promptsets {
			promptsetName := fmt.Sprintf("%s@%s (promptset)", promptset.Name, promptset.Installation.Version)
			if promptset.Manifest.Constraint != "" {
				promptsetName += fmt.Sprintf(" (%s)", promptset.Manifest.Constraint)
			}

			children := []pterm.TreeNode{
				{Text: fmt.Sprintf("include: %v", promptset.Manifest.Include)},
				{Text: fmt.Sprintf("sinks: %v", promptset.Manifest.Sinks)},
				{Text: fmt.Sprintf("files: %d installed", len(promptset.Installation.InstalledPaths))},
			}

			if len(promptset.Manifest.Exclude) > 0 {
				children = append(children, pterm.TreeNode{
					Text: fmt.Sprintf("exclude: %v", promptset.Manifest.Exclude),
				})
			}

			if detailed && len(promptset.Installation.InstalledPaths) > 0 {
				pathNodes := []pterm.TreeNode{}
				for _, path := range promptset.Installation.InstalledPaths {
					pathNodes = append(pathNodes, pterm.TreeNode{Text: path})
				}
				children = append(children, pterm.TreeNode{
					Text:     "installed paths:",
					Children: pathNodes,
				})
			}

			packageNodes = append(packageNodes, pterm.TreeNode{
				Text:     promptsetName,
				Children: children,
			})
		}

		tree := pterm.DefaultTree.WithRoot(pterm.TreeNode{
			Text:     registry,
			Children: packageNodes,
		})
		_ = tree.Render()
		pterm.Println()
	}
}

// OutdatedTable displays outdated packages in specified format
func (u *UI) OutdatedTable(outdated []OutdatedPackage, outputFormat string) {
	if len(outdated) == 0 {
		pterm.Success.Println("All packages are up to date!")
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
	case "list":
		for _, pkg := range outdated {
			fmt.Println(pkg.Package)
		}
	default: // table format
		tableData := [][]string{
			{"Package", "Type", "Constraint", "Current", "Wanted", "Latest"},
		}

		for _, pkg := range outdated {
			tableData = append(tableData, []string{
				pkg.Package,
				pkg.Type,
				pkg.Constraint,
				pkg.Current,
				pkg.Wanted,
				pkg.Latest,
			})
		}

		_ = pterm.DefaultTable.WithHasHeader().WithData(tableData).Render()
	}
}

// CompileStep displays a compilation step
func (u *UI) CompileStep(step string) {
	pterm.Info.Printf("%s ✓\n", step)
}

// CompileComplete displays compilation results
func (u *UI) CompileComplete(stats CompileStats, validateOnly bool) {
	if validateOnly {
		pterm.Success.Printf("Validated %d files\n", stats.FilesProcessed)
	} else {
		pterm.Success.Printf("Compiled %d files, generated %d rules\n",
			stats.FilesCompiled, stats.RulesGenerated)
	}
}
