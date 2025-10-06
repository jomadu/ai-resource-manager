package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newUninstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "Uninstall resources",
		Long:  "Uninstall all configured resources, or use subcommands for specific resource types.",
		RunE:  runUninstallAll,
	}

	// Add subcommands
	cmd.AddCommand(newUninstallRulesetCmd())
	cmd.AddCommand(newUninstallPromptsetCmd())

	return cmd
}

func newUninstallRulesetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ruleset <registry/ruleset>",
		Short: "Uninstall a ruleset",
		Args:  cobra.ExactArgs(1),
		RunE:  runUninstallRuleset,
	}
}

func newUninstallPromptsetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "promptset <registry/promptset>",
		Short: "Uninstall a promptset",
		Args:  cobra.ExactArgs(1),
		RunE:  runUninstallPromptset,
	}
}

func runUninstallAll(cmd *cobra.Command, args []string) error {
	// TODO: Implement unified uninstall when service interface is updated
	return fmt.Errorf("unified uninstall not yet implemented - service interface needs to be updated first")
}

func runUninstallRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	ruleset, err := ParsePackageArg(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse ruleset: %w", err)
	}

	err = armService.UninstallRuleset(ctx, ruleset.Registry, ruleset.Name)
	if err != nil {
		return fmt.Errorf("failed to uninstall %s/%s: %w", ruleset.Registry, ruleset.Name, err)
	}

	return nil
}

func runUninstallPromptset(cmd *cobra.Command, args []string) error {
	_, err := ParsePackageArg(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse promptset: %w", err)
	}

	// TODO: Implement promptset uninstall when service interface is updated
	return fmt.Errorf("promptset uninstall not yet implemented - service interface needs to be updated first")
}
