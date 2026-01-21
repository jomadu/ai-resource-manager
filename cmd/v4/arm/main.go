package main

import (
	"fmt"
	"os"

	"github.com/jomadu/ai-resource-manager/internal/v4/core"
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
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}
}

