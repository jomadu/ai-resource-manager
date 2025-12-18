package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove registries and sinks",
	Long:  "Remove registries and sinks from the ARM configuration",
}

var removeRegistryCmd = &cobra.Command{
	Use:   "registry NAME",
	Short: "Remove a registry",
	Long:  "Remove a registry from the ARM configuration by name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeRegistry(args[0])
	},
}

var removeSinkCmd = &cobra.Command{
	Use:   "sink NAME",
	Short: "Remove a sink",
	Long:  "Remove a sink from the ARM configuration by name.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		removeSink(args[0])
	},
}

func init() {
	// Add subcommands
	removeCmd.AddCommand(removeRegistryCmd)
	removeCmd.AddCommand(removeSinkCmd)
}

func removeRegistry(name string) {
	if err := armService.RemoveRegistry(ctx, name); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func removeSink(name string) {
	if err := armService.RemoveSink(ctx, name); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
