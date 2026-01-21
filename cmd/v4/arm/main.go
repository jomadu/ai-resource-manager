package main

import (
	"context"
	"fmt"
	"os"
	"strings"

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
		fmt.Println()
		fmt.Println("Flags:")
		fmt.Println("  --url          Git/GitLab/Cloudsmith repository URL (required)")
		fmt.Println("  --branches     Branches to track (git only, optional, comma-separated)")
		fmt.Println("  --project-id   GitLab project ID (gitlab only, optional)")
		fmt.Println("  --group-id     GitLab group ID (gitlab only, optional)")
		fmt.Println("  --api-version  GitLab API version (gitlab only, optional)")
		fmt.Println("  --owner        Cloudsmith owner (cloudsmith only, required)")
		fmt.Println("  --repo         Cloudsmith repository (cloudsmith only, required)")
		fmt.Println("  --force        Overwrite existing registry")
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

