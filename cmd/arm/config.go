package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configListCmd)
}

var configAddCmd = &cobra.Command{
	Use:   "add <type> <name> [url]",
	Short: "Add registry or sink configuration",
	Long:  "Add a registry or sink to the configuration. Type must be 'registry' or 'sink'.",
	Args:  cobra.MinimumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configType := args[0]
		name := args[1]

		switch configType {
		case "registry":
			if len(args) < 3 {
				return fmt.Errorf("registry requires a URL argument")
			}
			url := args[2]
			registryType, _ := cmd.Flags().GetString("type")
			if registryType == "" {
				registryType = "git" // default
			}

			configManager := config.NewFileManager()
			return configManager.AddRegistry(context.Background(), name, url, registryType)

		case "sink":
			directories, _ := cmd.Flags().GetStringSlice("directories")
			include, _ := cmd.Flags().GetStringSlice("include")
			exclude, _ := cmd.Flags().GetStringSlice("exclude")

			configManager := config.NewFileManager()
			return configManager.AddSink(context.Background(), name, directories, include, exclude)

		default:
			return fmt.Errorf("invalid type '%s'. Must be 'registry' or 'sink'", configType)
		}
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove <type> <name>",
	Short: "Remove registry or sink configuration",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		configType := args[0]
		name := args[1]

		configManager := config.NewFileManager()
		switch configType {
		case "registry":
			return configManager.RemoveRegistry(context.Background(), name)
		case "sink":
			return configManager.RemoveSink(context.Background(), name)
		default:
			return fmt.Errorf("invalid type '%s'. Must be 'registry' or 'sink'", configType)
		}
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		configManager := config.NewFileManager()
		registries, err := configManager.GetRegistries(context.Background())
		if err != nil {
			return err
		}
		sinks, err := configManager.GetSinks(context.Background())
		if err != nil {
			return err
		}

		fmt.Println("Registries:")
		for name, reg := range registries {
			fmt.Printf("  %s: %s (%s)\n", name, reg.URL, reg.Type)
		}

		fmt.Println("Sinks:")
		for name, sink := range sinks {
			fmt.Printf("  %s:\n", name)
			fmt.Printf("    directories: %v\n", sink.Directories)
			fmt.Printf("    include: %v\n", sink.Include)
			fmt.Printf("    exclude: %v\n", sink.Exclude)
		}
		return nil
	},
}

func init() {
	configAddCmd.Flags().String("type", "git", "Registry type (git, http)")
	configAddCmd.Flags().StringSlice("directories", nil, "Sink directories")
	configAddCmd.Flags().StringSlice("include", nil, "Sink include patterns")
	configAddCmd.Flags().StringSlice("exclude", nil, "Sink exclude patterns")
}
