package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jomadu/ai-resource-manager/internal/v4/compiler"
	"github.com/jomadu/ai-resource-manager/internal/v4/core"
	"github.com/jomadu/ai-resource-manager/internal/v4/manifest"
	"github.com/jomadu/ai-resource-manager/internal/v4/packagelockfile"
	"github.com/jomadu/ai-resource-manager/internal/v4/registry"
	"github.com/jomadu/ai-resource-manager/internal/v4/service"
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
		fmt.Fprintf(os.Stderr, "Usage: arm set <registry|sink> ...\n")
		os.Exit(1)
	}

	switch os.Args[2] {
	case "registry":
		handleSetRegistry()
	case "sink":
		handleSetSink()
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

func handleList() {
	if len(os.Args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: arm list <registry|sink>\n")
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, "Usage: arm info <registry|sink> [NAME...]\n")
		os.Exit(1)
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
