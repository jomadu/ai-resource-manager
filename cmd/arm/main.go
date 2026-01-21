package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jomadu/ai-resource-manager/internal/arm/compiler"
	"github.com/jomadu/ai-resource-manager/internal/arm/core"
	"github.com/jomadu/ai-resource-manager/internal/arm/manifest"
	"github.com/jomadu/ai-resource-manager/internal/arm/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/arm/registry"
	"github.com/jomadu/ai-resource-manager/internal/arm/service"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		return
	}

	switch os.Args[1] {
	case "version":
		printVersion()
	case "help":
		if len(os.Args) > 2 {
			printCommandHelp(os.Args[2])
		} else {
			printHelp()
		}
	case "add":
		handleAdd()
	case "remove":
		handleRemove()
	case "set":
		handleSet()
	case "list":
		handleList()
	case "info":
		handleInfo()
	case "install":
		handleInstall()
	case "uninstall":
		handleUninstall()
	case "update":
		handleUpdate()
	case "upgrade":
		handleUpgrade()
	case "outdated":
		handleOutdated()
	case "clean":
		handleClean()
	case "compile":
		handleCompile()
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		fmt.Fprintf(os.Stderr, "Run 'arm help' for usage.\n")
		os.Exit(1)
	}
}

func printVersion() {
	info := core.GetBuildInfo()
	fmt.Printf("arm %s\n", info.Version.Version)
	fmt.Printf("build-id: %s\n", info.Commit)
	fmt.Printf("build-timestamp: %s\n", info.BuildTime)
	fmt.Printf("build-platform: %s\n", info.Arch)
}

func printHelp() {
	fmt.Println("AI Resource Manager (ARM) - Manage AI resources for coding assistants")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  arm <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  version              Display version information")
	fmt.Println("  help [command]       Display help for a command")
	fmt.Println("  add                  Add registries or sinks")
	fmt.Println("  remove               Remove registries or sinks")
	fmt.Println("  set                  Configure registries or sinks")
	fmt.Println("  list                 List registries or sinks")
	fmt.Println("  info                 Show detailed information")
	fmt.Println("  install              Install rulesets or promptsets")
	fmt.Println("  uninstall            Uninstall packages")
	fmt.Println("  update               Update packages within version constraints")
	fmt.Println("  upgrade              Upgrade packages to latest versions")
	fmt.Println("  outdated             Check for outdated dependencies")
	fmt.Println("  clean                Clean cache or sinks")
	fmt.Println("  compile              Compile rulesets and promptsets")
	fmt.Println()
	fmt.Println("Run 'arm help <command>' for more information on a command.")
}

func printCommandHelp(command string) {
	switch command {
	case "version":
		fmt.Println("Display version information")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm version")
		fmt.Println()
		fmt.Println("Displays:")
		fmt.Println("  - Version number")
		fmt.Println("  - Build ID (commit hash)")
		fmt.Println("  - Build timestamp")
		fmt.Println("  - Build platform (OS/architecture)")
	case "add":
		fmt.Println("Add registries or sinks")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm add registry git --url URL [--branches BRANCH...] [--force] NAME")
		fmt.Println("  arm add registry gitlab --url URL [--project-id ID] [--group-id ID] [--api-version VERSION] [--force] NAME")
		fmt.Println("  arm add registry cloudsmith --url URL --owner OWNER --repo REPO [--force] NAME")
		fmt.Println("  arm add sink --tool TOOL [--force] NAME PATH")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --url          Git/GitLab/Cloudsmith repository URL (required)")
		fmt.Println("  --branches     Branches to track (git only, optional, comma-separated)")
		fmt.Println("  --project-id   GitLab project ID (gitlab only, optional)")
		fmt.Println("  --group-id     GitLab group ID (gitlab only, optional)")
		fmt.Println("  --api-version  GitLab API version (gitlab only, optional)")
		fmt.Println("  --owner        Cloudsmith owner (cloudsmith only, required)")
		fmt.Println("  --repo         Cloudsmith repository (cloudsmith only, required)")
		fmt.Println("  --tool         Sink tool: cursor, copilot, amazonq, markdown (required)")
		fmt.Println("  --force        Overwrite existing registry or sink")
	case "remove":
		fmt.Println("Remove registries or sinks")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm remove registry NAME")
		fmt.Println()
		fmt.Println("Removes the specified registry from the configuration.")
	case "set":
		fmt.Println("Configure registries or sinks")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm set registry NAME KEY VALUE")
		fmt.Println()
		fmt.Println("Supported keys:")
		fmt.Println("  name           Rename the registry")
		fmt.Println("  url            Update the registry URL")
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  arm set registry my-registry url https://github.com/new/repo")
	case "list":
		fmt.Println("List registries or sinks")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm list registry")
		fmt.Println()
		fmt.Println("Displays a simple list of configured registry names.")
	case "info":
		fmt.Println("Show detailed information")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm info registry [NAME...]")
		fmt.Println()
		fmt.Println("Displays detailed information about registries.")
		fmt.Println("If no names are provided, shows all registries.")
	case "install":
		fmt.Println("Install rulesets or promptsets")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm install                                                                    # Install all dependencies")
		fmt.Println("  arm install ruleset [--priority N] [--include PATTERN] [--exclude PATTERN] REGISTRY/RULESET[@VERSION] SINK...")
		fmt.Println("  arm install promptset [--include PATTERN] [--exclude PATTERN] REGISTRY/PROMPTSET[@VERSION] SINK...")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --priority     Priority for ruleset (default: 100)")
		fmt.Println("  --include      Include glob pattern (can be specified multiple times)")
		fmt.Println("  --exclude      Exclude glob pattern (can be specified multiple times)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  arm install")
		fmt.Println("  arm install ruleset --priority 200 my-registry/clean-code@1.0.0 cursor-rules")
		fmt.Println("  arm install promptset my-registry/code-review cursor-commands")
	case "uninstall":
		fmt.Println("Uninstall packages")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm uninstall [packages...]")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  packages       Package names in format registry/package (optional)")
		fmt.Println()
		fmt.Println("If no packages specified, uninstalls all packages.")
		fmt.Println("Removes specified packages from sinks and dependency configuration.")
	case "update":
		fmt.Println("Update packages within version constraints")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm update [packages...]")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  packages       Package names in format registry/package (optional)")
		fmt.Println()
		fmt.Println("If no packages specified, updates all packages.")
		fmt.Println("Updates packages to the latest versions that satisfy the version constraints")
		fmt.Println("specified in the manifest file. Updates the lock file with new versions.")
	case "upgrade":
		fmt.Println("Upgrade packages to latest versions")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm upgrade [registry/package ...]")
		fmt.Println()
		fmt.Println("Upgrades packages to the latest versions, ignoring version constraints.")
		fmt.Println("Updates the manifest file with new major version constraints (^X.0.0).")
		fmt.Println()
		fmt.Println("Arguments:")
		fmt.Println("  registry/package  One or more packages to upgrade (omit to upgrade all)")
		fmt.Println("Updates the lock file with new versions.")
	case "outdated":
		fmt.Println("Check for outdated dependencies")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm outdated [--output FORMAT]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --output       Output format: table (default), json, list")
		fmt.Println()
		fmt.Println("Displays packages with newer versions available, showing:")
		fmt.Println("  - Constraint: version constraint from manifest")
		fmt.Println("  - Current: currently installed version")
		fmt.Println("  - Wanted: latest version satisfying constraint")
		fmt.Println("  - Latest: latest available version")
	case "clean":
		fmt.Println("Clean cache or sinks")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm clean cache [--max-age DURATION] [--nuke]")
		fmt.Println("  arm clean sinks [--nuke]")
		fmt.Println()
		fmt.Println("Cache Flags:")
		fmt.Println("  --max-age      Remove cache older than duration (default: 7d)")
		fmt.Println("                 Examples: 30m, 2h, 7d, 1h30m")
		fmt.Println("  --nuke         Remove all cache (mutually exclusive with --max-age)")
		fmt.Println()
		fmt.Println("Sinks Flags:")
		fmt.Println("  --nuke         Remove entire ARM directory from sinks")
		fmt.Println()
		fmt.Println("Removes cached data or orphaned files from sinks.")
	case "compile":
		fmt.Println("Compile rulesets and promptsets")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  arm compile INPUT_PATH... [OUTPUT_PATH] [flags]")
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --tool         Target tool: markdown, cursor, amazonq, copilot")
		fmt.Println("  --namespace    Namespace for compiled resources")
		fmt.Println("  --force        Overwrite existing files")
		fmt.Println("  --recursive    Process directories recursively")
		fmt.Println("  --validate-only Validate without writing output (OUTPUT_PATH optional)")
		fmt.Println("  --include      Include glob patterns (can be specified multiple times)")
		fmt.Println("  --exclude      Exclude glob patterns (can be specified multiple times)")
		fmt.Println("  --fail-fast    Stop on first error")
		fmt.Println()
		fmt.Println("Compiles ARM resources to tool-specific formats.")
		fmt.Println("Supports files, directories, and mixed inputs.")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

func handleAdd() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: arm add <registry|sink> ...\n")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "registry":
		handleAddRegistry()
	case "sink":
		handleAddSink()
	default:
		fmt.Fprintf(os.Stderr, "Unknown add target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleAddRegistry() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: arm add registry <git|gitlab|cloudsmith> ...\n")
		os.Exit(1)
	}

	switch os.Args[3] {
	case "git":
		handleAddGitRegistry()
	case "gitlab":
		handleAddGitLabRegistry()
	case "cloudsmith":
		handleAddCloudsmithRegistry()
	default:
		fmt.Fprintf(os.Stderr, "Unknown registry type: %s\n", os.Args[3])
		os.Exit(1)
	}
}

func handleAddGitRegistry() {
	var url string
	var branches []string
	var force bool
	var name string

	// Parse flags and positional args
	i := 4
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--url" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--url requires a value\n")
				os.Exit(1)
			}
			url = os.Args[i+1]
			i += 2
		} else if arg == "--branches" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--branches requires a value\n")
				os.Exit(1)
			}
			branches = strings.Split(os.Args[i+1], ",")
			i += 2
		} else if arg == "--force" {
			force = true
			i++
		} else if !strings.HasPrefix(arg, "--") {
			name = arg
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if url == "" {
		fmt.Fprintf(os.Stderr, "--url is required\n")
		os.Exit(1)
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "NAME is required\n")
		os.Exit(1)
	}

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.AddGitRegistry(ctx, name, url, branches, force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added git registry '%s'\n", name)
}

func handleAddGitLabRegistry() {
	var url string
	var projectID string
	var groupID string
	var apiVersion string
	var force bool
	var name string

	// Parse flags and positional args
	i := 4
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--url" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--url requires a value\n")
				os.Exit(1)
			}
			url = os.Args[i+1]
			i += 2
		} else if arg == "--project-id" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--project-id requires a value\n")
				os.Exit(1)
			}
			projectID = os.Args[i+1]
			i += 2
		} else if arg == "--group-id" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--group-id requires a value\n")
				os.Exit(1)
			}
			groupID = os.Args[i+1]
			i += 2
		} else if arg == "--api-version" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--api-version requires a value\n")
				os.Exit(1)
			}
			apiVersion = os.Args[i+1]
			i += 2
		} else if arg == "--force" {
			force = true
			i++
		} else if !strings.HasPrefix(arg, "--") {
			name = arg
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if url == "" {
		fmt.Fprintf(os.Stderr, "--url is required\n")
		os.Exit(1)
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "NAME is required\n")
		os.Exit(1)
	}

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.AddGitLabRegistry(ctx, name, url, projectID, groupID, apiVersion, force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added gitlab registry '%s'\n", name)
}

func handleAddCloudsmithRegistry() {
	var url string
	var owner string
	var repo string
	var force bool
	var name string

	// Parse flags and positional args
	i := 4
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--url" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--url requires a value\n")
				os.Exit(1)
			}
			url = os.Args[i+1]
			i += 2
		} else if arg == "--owner" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--owner requires a value\n")
				os.Exit(1)
			}
			owner = os.Args[i+1]
			i += 2
		} else if arg == "--repo" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--repo requires a value\n")
				os.Exit(1)
			}
			repo = os.Args[i+1]
			i += 2
		} else if arg == "--force" {
			force = true
			i++
		} else if !strings.HasPrefix(arg, "--") {
			name = arg
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if url == "" {
		fmt.Fprintf(os.Stderr, "--url is required\n")
		os.Exit(1)
	}
	if owner == "" {
		fmt.Fprintf(os.Stderr, "--owner is required\n")
		os.Exit(1)
	}
	if repo == "" {
		fmt.Fprintf(os.Stderr, "--repo is required\n")
		os.Exit(1)
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "NAME is required\n")
		os.Exit(1)
	}

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.AddCloudsmithRegistry(ctx, name, url, owner, repo, force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added cloudsmith registry '%s'\n", name)
}

func handleAddSink() {
	var tool string
	var force bool
	var name string
	var path string

	// Parse flags and positional args
	i := 3
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--tool" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--tool requires a value\n")
				os.Exit(1)
			}
			tool = os.Args[i+1]
			i += 2
		} else if arg == "--force" {
			force = true
			i++
		} else if !strings.HasPrefix(arg, "--") {
			if name == "" {
				name = arg
			} else if path == "" {
				path = arg
			} else {
				fmt.Fprintf(os.Stderr, "Too many positional arguments\n")
				os.Exit(1)
			}
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if tool == "" {
		fmt.Fprintf(os.Stderr, "--tool is required\n")
		os.Exit(1)
	}
	if name == "" {
		fmt.Fprintf(os.Stderr, "NAME is required\n")
		os.Exit(1)
	}
	if path == "" {
		fmt.Fprintf(os.Stderr, "PATH is required\n")
		os.Exit(1)
	}

	// Validate tool
	var compilerTool compiler.Tool
	switch tool {
	case "cursor":
		compilerTool = compiler.Cursor
	case "copilot":
		compilerTool = compiler.Copilot
	case "amazonq":
		compilerTool = compiler.AmazonQ
	case "markdown":
		compilerTool = compiler.Markdown
	default:
		fmt.Fprintf(os.Stderr, "Invalid tool: %s (must be cursor, copilot, amazonq, or markdown)\n", tool)
		os.Exit(1)
	}

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.AddSink(ctx, name, path, compilerTool, force); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Added sink '%s'\n", name)
}



func handleRemove() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: arm remove <registry|sink> ...\n")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "registry":
		handleRemoveRegistry()
	case "sink":
		handleRemoveSink()
	default:
		fmt.Fprintf(os.Stderr, "Unknown remove target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleRemoveRegistry() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: arm remove registry NAME\n")
		os.Exit(1)
	}

	name := os.Args[3]

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.RemoveRegistry(ctx, name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed registry '%s'\n", name)
}

func handleRemoveSink() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "Usage: arm remove sink NAME\n")
		os.Exit(1)
	}

	name := os.Args[3]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	if err := svc.RemoveSink(ctx, name); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Removed sink '%s'\n", name)
}

func handleSet() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: arm set <registry|sink|ruleset|promptset> ...\n")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "registry":
		handleSetRegistry()
	case "sink":
		handleSetSink()
	case "ruleset":
		handleSetRuleset()
	case "promptset":
		handleSetPromptset()
	default:
		fmt.Fprintf(os.Stderr, "Unknown set target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleSetRegistry() {
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: arm set registry NAME KEY VALUE\n")
		os.Exit(1)
	}

	name := os.Args[3]
	key := os.Args[4]
	value := os.Args[5]

	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	var err error

	switch key {
	case "name":
		err = svc.SetRegistryName(ctx, name, value)
	case "url":
		err = svc.SetRegistryURL(ctx, name, value)
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s (valid: name, url)\n", key)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated registry '%s' %s\n", name, key)
}

func handleSetSink() {
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: arm set sink NAME KEY VALUE\n")
		os.Exit(1)
	}

	name := os.Args[3]
	key := os.Args[4]
	value := os.Args[5]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	var err error

	switch key {
	case "tool":
		var tool compiler.Tool
		switch value {
		case "cursor":
			tool = compiler.Cursor
		case "copilot":
			tool = compiler.Copilot
		case "amazonq":
			tool = compiler.AmazonQ
		case "markdown":
			tool = compiler.Markdown
		default:
			fmt.Fprintf(os.Stderr, "Invalid tool: %s (valid: cursor, copilot, amazonq, markdown)\n", value)
			os.Exit(1)
		}
		err = svc.SetSinkTool(ctx, name, tool)
	case "directory":
		err = svc.SetSinkDirectory(ctx, name, value)
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s (valid: tool, directory)\n", key)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated sink '%s' %s\n", name, key)
}

func handleSetRuleset() {
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: arm set ruleset REGISTRY/RULESET KEY VALUE\n")
		os.Exit(1)
	}

	packageSpec := os.Args[3]
	key := os.Args[4]
	value := os.Args[5]

	parts := strings.Split(packageSpec, "/")
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid package spec: %s (expected REGISTRY/RULESET)\n", packageSpec)
		os.Exit(1)
	}
	registryName := parts[0]
	ruleset := parts[1]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	var err error

	switch key {
	case "version":
		err = svc.SetRulesetVersion(ctx, registryName, ruleset, value)
	case "priority":
		var priority int
		_, scanErr := fmt.Sscanf(value, "%d", &priority)
		if scanErr != nil {
			fmt.Fprintf(os.Stderr, "Invalid priority value: %s (must be integer)\n", value)
			os.Exit(1)
		}
		err = svc.SetRulesetPriority(ctx, registryName, ruleset, priority)
	case "sinks":
		sinks := strings.Split(value, ",")
		err = svc.SetRulesetSinks(ctx, registryName, ruleset, sinks)
	case "include":
		include := strings.Split(value, ",")
		err = svc.SetRulesetInclude(ctx, registryName, ruleset, include)
	case "exclude":
		exclude := strings.Split(value, ",")
		err = svc.SetRulesetExclude(ctx, registryName, ruleset, exclude)
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s (valid: version, priority, sinks, include, exclude)\n", key)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated ruleset '%s/%s' %s\n", registryName, ruleset, key)
}

func handleSetPromptset() {
	if len(os.Args) < 6 {
		fmt.Fprintf(os.Stderr, "Usage: arm set promptset REGISTRY/PROMPTSET KEY VALUE\n")
		os.Exit(1)
	}

	packageSpec := os.Args[3]
	key := os.Args[4]
	value := os.Args[5]

	parts := strings.Split(packageSpec, "/")
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "Invalid package spec: %s (expected REGISTRY/PROMPTSET)\n", packageSpec)
		os.Exit(1)
	}
	registryName := parts[0]
	promptset := parts[1]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	var err error

	switch key {
	case "version":
		err = svc.SetPromptsetVersion(ctx, registryName, promptset, value)
	case "sinks":
		sinks := strings.Split(value, ",")
		err = svc.SetPromptsetSinks(ctx, registryName, promptset, sinks)
	case "include":
		include := strings.Split(value, ",")
		err = svc.SetPromptsetInclude(ctx, registryName, promptset, include)
	case "exclude":
		exclude := strings.Split(value, ",")
		err = svc.SetPromptsetExclude(ctx, registryName, promptset, exclude)
	default:
		fmt.Fprintf(os.Stderr, "Unknown key: %s (valid: version, sinks, include, exclude)\n", key)
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated promptset '%s/%s' %s\n", registryName, promptset, key)
}

func handleList() {
	if len(os.Args) < 3 {
		handleListAll()
		return
	}

	switch os.Args[2] {
	case "registry":
		handleListRegistry()
	case "sink":
		handleListSink()
	default:
		fmt.Fprintf(os.Stderr, "Unknown list target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleListAll() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()

	// List registries
	registries, err := svc.GetAllRegistriesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Registries:")
	if len(registries) == 0 {
		fmt.Println("  (none)")
	} else {
		for name := range registries {
			fmt.Printf("  %s\n", name)
		}
	}

	// List sinks
	sinks, err := svc.GetAllSinkConfigs(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nSinks:")
	if len(sinks) == 0 {
		fmt.Println("  (none)")
	} else {
		for name := range sinks {
			fmt.Printf("  %s\n", name)
		}
	}

	// List dependencies
	rulesets, err := manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	promptsets, err := manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDependencies:")
	if len(rulesets) == 0 && len(promptsets) == 0 {
		fmt.Println("  (none)")
	} else {
		for key := range rulesets {
			fmt.Printf("  %s (ruleset)\n", key)
		}
		for key := range promptsets {
			fmt.Printf("  %s (promptset)\n", key)
		}
	}
}

func handleListRegistry() {
	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	registries, err := svc.GetAllRegistriesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(registries) == 0 {
		fmt.Println("No registries configured")
		return
	}

	for name := range registries {
		fmt.Println(name)
	}
}

func handleInfo() {
	if len(os.Args) < 3 {
		handleInfoAll()
		return
	}

	switch os.Args[2] {
	case "registry":
		handleInfoRegistry()
	case "sink":
		handleInfoSink()
	default:
		fmt.Fprintf(os.Stderr, "Unknown info target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleInfoAll() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()

	// Show registries
	registries, err := svc.GetAllRegistriesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Registries:")
	if len(registries) == 0 {
		fmt.Println("  (none)")
	} else {
		for name, config := range registries {
			fmt.Printf("\n  %s:\n", name)
			regType, _ := config["type"].(string)
			fmt.Printf("    type: %s\n", regType)
			if url, ok := config["url"].(string); ok {
				fmt.Printf("    url: %s\n", url)
			}
			if regType == "git" {
				if branches, ok := config["branches"].([]interface{}); ok && len(branches) > 0 {
					fmt.Printf("    branches: %v\n", branches)
				}
			} else if regType == "gitlab" {
				if projectID, ok := config["projectId"].(string); ok && projectID != "" {
					fmt.Printf("    projectId: %s\n", projectID)
				}
				if groupID, ok := config["groupId"].(string); ok && groupID != "" {
					fmt.Printf("    groupId: %s\n", groupID)
				}
				if apiVersion, ok := config["apiVersion"].(string); ok && apiVersion != "" {
					fmt.Printf("    apiVersion: %s\n", apiVersion)
				}
			} else if regType == "cloudsmith" {
				if owner, ok := config["owner"].(string); ok {
					fmt.Printf("    owner: %s\n", owner)
				}
				if repo, ok := config["repository"].(string); ok {
					fmt.Printf("    repository: %s\n", repo)
				}
			}
		}
	}

	// Show sinks
	sinks, err := svc.GetAllSinkConfigs(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nSinks:")
	if len(sinks) == 0 {
		fmt.Println("  (none)")
	} else {
		for name, config := range sinks {
			fmt.Printf("\n  %s:\n", name)
			fmt.Printf("    tool: %s\n", config.Tool)
			fmt.Printf("    directory: %s\n", config.Directory)
		}
	}

	// Show dependencies
	rulesets, err := manifestMgr.GetAllRulesetDependenciesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	promptsets, err := manifestMgr.GetAllPromptsetDependenciesConfig(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nDependencies:")
	if len(rulesets) == 0 && len(promptsets) == 0 {
		fmt.Println("  (none)")
	} else {
		for key, config := range rulesets {
			fmt.Printf("\n  %s:\n", key)
			fmt.Printf("    type: ruleset\n")
			fmt.Printf("    version: %s\n", config.Version)
			fmt.Printf("    priority: %d\n", config.Priority)
			if len(config.Sinks) > 0 {
				fmt.Printf("    sinks: %v\n", config.Sinks)
			}
			if len(config.Include) > 0 {
				fmt.Printf("    include: %v\n", config.Include)
			}
			if len(config.Exclude) > 0 {
				fmt.Printf("    exclude: %v\n", config.Exclude)
			}
		}
		for key, config := range promptsets {
			fmt.Printf("\n  %s:\n", key)
			fmt.Printf("    type: promptset\n")
			fmt.Printf("    version: %s\n", config.Version)
			if len(config.Sinks) > 0 {
				fmt.Printf("    sinks: %v\n", config.Sinks)
			}
			if len(config.Include) > 0 {
				fmt.Printf("    include: %v\n", config.Include)
			}
			if len(config.Exclude) > 0 {
				fmt.Printf("    exclude: %v\n", config.Exclude)
			}
		}
	}
}

func handleInfoRegistry() {
	// Get manifest path from env or use default
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()

	// Get names from args or all registries
	var names []string
	if len(os.Args) > 3 {
		names = os.Args[3:]
	} else {
		// Get all registry names
		registries, err := svc.GetAllRegistriesConfig(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for name := range registries {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		fmt.Println("No registries configured")
		return
	}

	// Display info for each registry
	for i, name := range names {
		if i > 0 {
			fmt.Println()
		}

		config, err := svc.GetRegistryConfig(ctx, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting registry '%s': %v\n", name, err)
			continue
		}

		fmt.Printf("Registry: %s\n", name)
		
		// Display type
		if regType, ok := config["type"].(string); ok {
			fmt.Printf("  Type: %s\n", regType)
		}

		// Display URL
		if url, ok := config["url"].(string); ok {
			fmt.Printf("  URL: %s\n", url)
		}

		// Display type-specific fields
		if branches, ok := config["branches"].([]interface{}); ok && len(branches) > 0 {
			fmt.Printf("  Branches: ")
			for j, b := range branches {
				if j > 0 {
					fmt.Printf(", ")
				}
				fmt.Printf("%v", b)
			}
			fmt.Println()
		}

		if projectID, ok := config["projectId"].(string); ok && projectID != "" {
			fmt.Printf("  Project ID: %s\n", projectID)
		}

		if groupID, ok := config["groupId"].(string); ok && groupID != "" {
			fmt.Printf("  Group ID: %s\n", groupID)
		}

		if apiVersion, ok := config["apiVersion"].(string); ok && apiVersion != "" {
			fmt.Printf("  API Version: %s\n", apiVersion)
		}

		if owner, ok := config["owner"].(string); ok && owner != "" {
			fmt.Printf("  Owner: %s\n", owner)
		}

		if repo, ok := config["repository"].(string); ok && repo != "" {
			fmt.Printf("  Repository: %s\n", repo)
		}
	}
}

func handleListSink() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()
	sinks, err := svc.GetAllSinkConfigs(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(sinks) == 0 {
		fmt.Println("No sinks configured")
		return
	}

	for name := range sinks {
		fmt.Println(name)
	}
}

func handleInfoSink() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm.json"
	}

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManager()
	registryFactory := &registry.DefaultFactory{}
	svc := service.NewArmService(manifestMgr, lockfileMgr, registryFactory)

	ctx := context.Background()

	// Get names from args or all sinks
	var names []string
	if len(os.Args) > 3 {
		names = os.Args[3:]
	} else {
		// Get all sink names
		sinks, err := svc.GetAllSinkConfigs(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		for name := range sinks {
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		fmt.Println("No sinks configured")
		return
	}

	// Display info for each sink
	for i, name := range names {
		if i > 0 {
			fmt.Println()
		}

		config, err := svc.GetSinkConfig(ctx, name)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error getting sink '%s': %v\n", name, err)
			continue
		}

		fmt.Printf("Sink: %s\n", name)
		fmt.Printf("  Tool: %s\n", config.Tool)
		fmt.Printf("  Directory: %s\n", config.Directory)
	}
}

func handleInstall() {
	if len(os.Args) < 3 {
		handleInstallAll()
		return
	}

	switch os.Args[2] {
	case "ruleset":
		handleInstallRuleset()
	case "promptset":
		handleInstallPromptset()
	default:
		fmt.Fprintf(os.Stderr, "Unknown install target: %s\n", os.Args[2])
		os.Exit(1)
	}
}

func handleInstallAll() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	if err := svc.InstallAll(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All dependencies installed successfully")
}

func handleInstallRuleset() {
	var priority int = 100
	var include []string
	var exclude []string
	var packageSpec string
	var sinks []string

	// Parse flags and positional args
	i := 3
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--priority" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--priority requires a value\n")
				os.Exit(1)
			}
			fmt.Sscanf(os.Args[i+1], "%d", &priority)
			i += 2
		} else if arg == "--include" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--include requires a value\n")
				os.Exit(1)
			}
			include = append(include, os.Args[i+1])
			i += 2
		} else if arg == "--exclude" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--exclude requires a value\n")
				os.Exit(1)
			}
			exclude = append(exclude, os.Args[i+1])
			i += 2
		} else if !strings.HasPrefix(arg, "--") {
			if packageSpec == "" {
				packageSpec = arg
			} else {
				sinks = append(sinks, arg)
			}
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if packageSpec == "" {
		fmt.Fprintf(os.Stderr, "Package spec required (REGISTRY/RULESET[@VERSION])\n")
		os.Exit(1)
	}

	if len(sinks) == 0 {
		fmt.Fprintf(os.Stderr, "At least one sink required\n")
		os.Exit(1)
	}

	// Parse package spec
	registryName, err := parseRegistry(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	ruleset, err := parsePackage(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	version, err := parseVersion(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize service
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	// Call service
	if err := svc.InstallRuleset(ctx, registryName, ruleset, version, priority, include, exclude, sinks); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installed %s/%s to sinks: %s\n", registryName, ruleset, strings.Join(sinks, ", "))
}

func handleInstallPromptset() {
	var include []string
	var exclude []string
	var packageSpec string
	var sinks []string

	// Parse flags and positional args
	i := 3
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--include" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--include requires a value\n")
				os.Exit(1)
			}
			include = append(include, os.Args[i+1])
			i += 2
		} else if arg == "--exclude" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--exclude requires a value\n")
				os.Exit(1)
			}
			exclude = append(exclude, os.Args[i+1])
			i += 2
		} else if !strings.HasPrefix(arg, "--") {
			if packageSpec == "" {
				packageSpec = arg
			} else {
				sinks = append(sinks, arg)
			}
			i++
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	if packageSpec == "" {
		fmt.Fprintf(os.Stderr, "Package spec required (REGISTRY/PROMPTSET[@VERSION])\n")
		os.Exit(1)
	}

	if len(sinks) == 0 {
		fmt.Fprintf(os.Stderr, "At least one sink required\n")
		os.Exit(1)
	}

	// Parse package spec
	registryName, err := parseRegistry(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	promptset, err := parsePackage(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	version, err := parseVersion(packageSpec)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Initialize service
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	// Call service
	if err := svc.InstallPromptset(ctx, registryName, promptset, version, include, exclude, sinks); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Installed %s/%s to sinks: %s\n", registryName, promptset, strings.Join(sinks, ", "))
}

func parseRegistry(input string) (string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: %s (expected registry/package)", input)
	}
	if parts[0] == "" {
		return "", fmt.Errorf("registry name cannot be empty")
	}
	return parts[0], nil
}

func parsePackage(input string) (string, error) {
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid format: %s (expected registry/package)", input)
	}

	pkgWithVersion := parts[1]
	if pkgWithVersion == "" {
		return "", fmt.Errorf("package name cannot be empty")
	}

	if strings.Contains(pkgWithVersion, "@") {
		pkg := strings.SplitN(pkgWithVersion, "@", 2)[0]
		if pkg == "" {
			return "", fmt.Errorf("package name cannot be empty")
		}
		return pkg, nil
	}

	return pkgWithVersion, nil
}

func parseVersion(input string) (string, error) {
	if !strings.Contains(input, "@") {
		return "", nil
	}

	parts := strings.SplitN(input, "@", 2)
	if len(parts) != 2 {
		return "", nil
	}

	version := parts[1]
	if version == "" {
		return "", fmt.Errorf("version cannot be empty after @")
	}

	return version, nil
}

func handleUninstall() {
	// Parse package arguments (everything after "uninstall")
	packages := os.Args[2:]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	// If no packages specified, uninstall all
	if len(packages) == 0 {
		if err := svc.UninstallAll(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All packages uninstalled successfully")
		return
	}

	// Uninstall specific packages
	if err := svc.UninstallPackages(ctx, packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Packages uninstalled successfully")
}

func handleUpdate() {
	// Parse package arguments (everything after "update")
	packages := os.Args[2:]

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	// If no packages specified, update all
	if len(packages) == 0 {
		if err := svc.UpdateAll(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All packages updated successfully")
		return
	}

	// Update specific packages
	if err := svc.UpdatePackages(ctx, packages); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Packages updated successfully")
}

func handleUpgrade() {
	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	packages := os.Args[2:]
	if len(packages) == 0 {
		if err := svc.UpgradeAll(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("All packages upgraded successfully")
	} else {
		if err := svc.UpgradePackages(ctx, packages); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Specified packages upgraded successfully")
	}
}

func handleOutdated() {
	outputFormat := "table"

	// Parse flags
	i := 2
	for i < len(os.Args) {
		arg := os.Args[i]
		if arg == "--output" {
			if i+1 >= len(os.Args) {
				fmt.Fprintf(os.Stderr, "--output requires a value\n")
				os.Exit(1)
			}
			outputFormat = os.Args[i+1]
			if outputFormat != "table" && outputFormat != "json" && outputFormat != "list" {
				fmt.Fprintf(os.Stderr, "Invalid output format: %s (must be table, json, or list)\n", outputFormat)
				os.Exit(1)
			}
			i += 2
		} else {
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", arg)
			os.Exit(1)
		}
	}

	manifestPath := os.Getenv("ARM_MANIFEST_PATH")
	if manifestPath == "" {
		manifestPath = "arm-manifest.json"
	}

	lockfilePath := strings.TrimSuffix(manifestPath, ".json") + "-lock.json"

	manifestMgr := manifest.NewFileManagerWithPath(manifestPath)
	lockfileMgr := packagelockfile.NewFileManagerWithPath(lockfilePath)

	svc := service.NewArmService(manifestMgr, lockfileMgr, nil)
	ctx := context.Background()

	outdated, err := svc.ListOutdated(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if len(outdated) == 0 {
		fmt.Println("All packages are up to date")
		return
	}

	switch outputFormat {
	case "json":
		printOutdatedJSON(outdated)
	case "list":
		printOutdatedList(outdated)
	default:
		printOutdatedTable(outdated)
	}
}

func printOutdatedTable(outdated []*service.OutdatedDependency) {
	// Calculate column widths
	maxPackage := len("Package")
	maxConstraint := len("Constraint")
	maxCurrent := len("Current")
	maxWanted := len("Wanted")
	maxLatest := len("Latest")

	for _, dep := range outdated {
		pkgName := dep.Current.RegistryName + "/" + dep.Current.Name
		if len(pkgName) > maxPackage {
			maxPackage = len(pkgName)
		}
		if len(dep.Constraint) > maxConstraint {
			maxConstraint = len(dep.Constraint)
		}
		if len(dep.Current.Version.Version) > maxCurrent {
			maxCurrent = len(dep.Current.Version.Version)
		}
		if len(dep.Wanted.Version.Version) > maxWanted {
			maxWanted = len(dep.Wanted.Version.Version)
		}
		if len(dep.Latest.Version.Version) > maxLatest {
			maxLatest = len(dep.Latest.Version.Version)
		}
	}

	// Print header
	fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s\n",
		maxPackage, "Package",
		maxConstraint, "Constraint",
		maxCurrent, "Current",
		maxWanted, "Wanted",
		maxLatest, "Latest")

	// Print separator
	for i := 0; i < maxPackage; i++ {
		fmt.Print("-")
	}
	fmt.Print("  ")
	for i := 0; i < maxConstraint; i++ {
		fmt.Print("-")
	}
	fmt.Print("  ")
	for i := 0; i < maxCurrent; i++ {
		fmt.Print("-")
	}
	fmt.Print("  ")
	for i := 0; i < maxWanted; i++ {
		fmt.Print("-")
	}
	fmt.Print("  ")
	for i := 0; i < maxLatest; i++ {
		fmt.Print("-")
	}
	fmt.Println()

	// Print rows
	for _, dep := range outdated {
		pkgName := dep.Current.RegistryName + "/" + dep.Current.Name
		fmt.Printf("%-*s  %-*s  %-*s  %-*s  %-*s\n",
			maxPackage, pkgName,
			maxConstraint, dep.Constraint,
			maxCurrent, dep.Current.Version.Version,
			maxWanted, dep.Wanted.Version.Version,
			maxLatest, dep.Latest.Version.Version)
	}
}

func printOutdatedJSON(outdated []*service.OutdatedDependency) {
	fmt.Println("[")
	for i, dep := range outdated {
		pkgName := dep.Current.RegistryName + "/" + dep.Current.Name
		fmt.Printf("  {\n")
		fmt.Printf("    \"package\": \"%s\",\n", pkgName)
		fmt.Printf("    \"constraint\": \"%s\",\n", dep.Constraint)
		fmt.Printf("    \"current\": \"%s\",\n", dep.Current.Version.Version)
		fmt.Printf("    \"wanted\": \"%s\",\n", dep.Wanted.Version.Version)
		fmt.Printf("    \"latest\": \"%s\"\n", dep.Latest.Version.Version)
		if i < len(outdated)-1 {
			fmt.Printf("  },\n")
		} else {
			fmt.Printf("  }\n")
		}
	}
	fmt.Println("]")
}

func printOutdatedList(outdated []*service.OutdatedDependency) {
	for _, dep := range outdated {
		pkgName := dep.Current.RegistryName + "/" + dep.Current.Name
		fmt.Printf("%s: %s -> %s (latest: %s)\n",
			pkgName,
			dep.Current.Version.Version,
			dep.Wanted.Version.Version,
			dep.Latest.Version.Version)
	}
}

func handleClean() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Error: clean requires a subcommand (cache, sinks)\n")
		fmt.Fprintf(os.Stderr, "Run 'arm help clean' for usage.\n")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "cache":
		handleCleanCache()
	case "sinks":
		handleCleanSinks()
	default:
		fmt.Fprintf(os.Stderr, "Unknown clean subcommand: %s\n", os.Args[2])
		fmt.Fprintf(os.Stderr, "Run 'arm help clean' for usage.\n")
		os.Exit(1)
	}
}

func handleCleanCache() {
	var maxAge string
	var nuke bool

	// Parse flags
	args := os.Args[3:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--max-age":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --max-age requires a value\n")
				os.Exit(1)
			}
			maxAge = args[i+1]
			i++
		case "--nuke":
			nuke = true
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", args[i])
			os.Exit(1)
		}
	}

	// Validate mutual exclusivity
	if maxAge != "" && nuke {
		fmt.Fprintf(os.Stderr, "Error: --max-age and --nuke are mutually exclusive\n")
		os.Exit(1)
	}

	// Default max-age to 7d
	if maxAge == "" && !nuke {
		maxAge = "7d"
	}

	svc := service.NewArmService(nil, nil, nil)
	ctx := context.Background()

	if nuke {
		if err := svc.NukeCache(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Cache nuked successfully")
		return
	}

	// Parse duration
	duration, err := parseDuration(maxAge)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: invalid duration format: %v\n", err)
		os.Exit(1)
	}

	if err := svc.CleanCacheByAge(ctx, duration); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cache cleaned (removed items older than %s)\n", maxAge)
}

func parseDuration(s string) (time.Duration, error) {
	// Support simple formats: 30m, 2h, 7d, 1h30m
	// Convert days to hours since time.ParseDuration doesn't support days
	if strings.HasSuffix(s, "d") {
		days := strings.TrimSuffix(s, "d")
		var d int
		if _, err := fmt.Sscanf(days, "%d", &d); err != nil {
			return 0, fmt.Errorf("invalid duration: %s", s)
		}
		return time.Duration(d) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}

func handleCleanSinks() {
	var nuke bool

	// Parse flags
	args := os.Args[3:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--nuke":
			nuke = true
		default:
			fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", args[i])
			os.Exit(1)
		}
	}

	manifestMgr := manifest.NewFileManager()
	svc := service.NewArmService(manifestMgr, nil, nil)
	ctx := context.Background()

	if nuke {
		if err := svc.NukeSinks(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Sinks nuked successfully")
		return
	}

	if err := svc.CleanSinks(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Sinks cleaned successfully")
}

func handleCompile() {
	var tool string
	var namespace string
	var force bool
	var recursive bool
	var validateOnly bool
	var include []string
	var exclude []string
	var failFast bool
	var paths []string
	var outputPath string

	// Parse flags and arguments
	args := os.Args[2:]
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--tool":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --tool requires a value\n")
				os.Exit(1)
			}
			tool = args[i+1]
			i++
		case "--namespace":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --namespace requires a value\n")
				os.Exit(1)
			}
			namespace = args[i+1]
			i++
		case "--force":
			force = true
		case "--recursive":
			recursive = true
		case "--validate-only":
			validateOnly = true
		case "--include":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --include requires a value\n")
				os.Exit(1)
			}
			include = append(include, args[i+1])
			i++
		case "--exclude":
			if i+1 >= len(args) {
				fmt.Fprintf(os.Stderr, "Error: --exclude requires a value\n")
				os.Exit(1)
			}
			exclude = append(exclude, args[i+1])
			i++
		case "--fail-fast":
			failFast = true
		default:
			if strings.HasPrefix(args[i], "--") {
				fmt.Fprintf(os.Stderr, "Unknown flag: %s\n", args[i])
				os.Exit(1)
			}
			paths = append(paths, args[i])
		}
	}

	// Validate required arguments
	if len(paths) == 0 {
		fmt.Fprintf(os.Stderr, "Error: at least one INPUT_PATH is required\n")
		fmt.Fprintf(os.Stderr, "Run 'arm help compile' for usage.\n")
		os.Exit(1)
	}

	// Determine output path
	if !validateOnly {
		if len(paths) < 2 {
			fmt.Fprintf(os.Stderr, "Error: OUTPUT_PATH is required (or use --validate-only)\n")
			fmt.Fprintf(os.Stderr, "Run 'arm help compile' for usage.\n")
			os.Exit(1)
		}
		outputPath = paths[len(paths)-1]
		paths = paths[:len(paths)-1]
	}

	// Create service and compile
	svc := service.NewArmService(nil, nil, nil)
	ctx := context.Background()

	req := &service.CompileRequest{
		Paths:        paths,
		Tool:         tool,
		OutputDir:    outputPath,
		Namespace:    namespace,
		Force:        force,
		Recursive:    recursive,
		ValidateOnly: validateOnly,
		Include:      include,
		Exclude:      exclude,
		FailFast:     failFast,
	}

	if err := svc.CompileFiles(ctx, req); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if validateOnly {
		fmt.Println("Validation successful")
	} else {
		fmt.Println("Compilation successful")
	}
}
