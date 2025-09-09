package main

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
)

func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <ruleset>",
		Short: "Uninstall a ruleset",
		Args:  cobra.ExactArgs(1),
		RunE:  runUninstall,
	}
}

func runUninstall(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	ruleset, err := ParseRulesetArg(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse ruleset: %w", err)
	}

	err = armService.UninstallRuleset(ctx, ruleset.Registry, ruleset.Name)
	if err != nil {
		return fmt.Errorf("failed to uninstall %s/%s: %w", ruleset.Registry, ruleset.Name, err)
	}

	return nil
}
