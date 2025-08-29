package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configListCmd)
}

var configAddCmd = &cobra.Command{
	Use:   "add <type> <name> <url>",
	Short: "Add registry or sink configuration",
	Long:  "Add a registry or sink to the configuration. Type must be 'registry' or 'sink'.",
	Args:  cobra.MinimumNArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		configType := args[0]
		name := args[1]

		switch configType {
		case "registry":
			url := args[2]
			registryType, _ := cmd.Flags().GetString("type")
			if registryType == "" {
				registryType = "git" // default
			}

			fmt.Printf("Adding registry '%s' with URL '%s' and type '%s'\n", name, url, registryType)
			// TODO: Implement registry addition via config manager
			return nil

		case "sink":
			directories, _ := cmd.Flags().GetStringSlice("directories")
			include, _ := cmd.Flags().GetStringSlice("include")
			exclude, _ := cmd.Flags().GetStringSlice("exclude")

			fmt.Printf("Adding sink '%s'\n", name)
			fmt.Printf("  Directories: %v\n", directories)
			fmt.Printf("  Include: %v\n", include)
			fmt.Printf("  Exclude: %v\n", exclude)
			// TODO: Implement sink addition via config manager
			return nil

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

		fmt.Printf("Removing %s '%s'\n", configType, name)
		// TODO: Implement removal via config manager
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Configuration:")
		// TODO: Implement listing via config manager
		return nil
	},
}

func init() {
	configAddCmd.Flags().String("type", "git", "Registry type (git, http)")
	configAddCmd.Flags().StringSlice("directories", nil, "Sink directories")
	configAddCmd.Flags().StringSlice("include", nil, "Sink include patterns")
	configAddCmd.Flags().StringSlice("exclude", nil, "Sink exclude patterns")
}
