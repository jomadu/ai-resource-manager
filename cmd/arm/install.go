package main

import (
	"context"
	"fmt"

	"github.com/jomadu/ai-rules-manager/internal/arm"
	"github.com/spf13/cobra"
)

func newInstallCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install resources",
		Long:  "Install all configured resources from manifest, or use subcommands for specific resource types.",
		RunE:  runInstallAll,
	}

	// Add subcommands
	cmd.AddCommand(newInstallRulesetCmd())
	cmd.AddCommand(newInstallPromptsetCmd())

	return cmd
}

func newInstallRulesetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ruleset <registry/ruleset[@version]> <sink...>",
		Short: "Install rulesets",
		Long:  "Install rulesets from a registry to specified sinks.",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runInstallRuleset,
	}

	cmd.Flags().StringSlice("include", nil, "Include patterns")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns")
	cmd.Flags().Int("priority", 100, "Ruleset installation priority (1-1000+)")

	return cmd
}

func newInstallPromptsetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "promptset <registry/promptset[@version]> <sink...>",
		Short: "Install promptsets",
		Long:  "Install promptsets from a registry to specified sinks.",
		Args:  cobra.MinimumNArgs(2),
		RunE:  runInstallPromptset,
	}

	cmd.Flags().StringSlice("include", nil, "Include patterns")
	cmd.Flags().StringSlice("exclude", nil, "Exclude patterns")

	return cmd
}

func runInstallAll(cmd *cobra.Command, args []string) error {
	ctx := context.Background()
	return armService.InstallAll(ctx)
}

func runInstallRuleset(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Parse ruleset argument
	ruleset, err := ParsePackageArg(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse ruleset: %w", err)
	}

	// Get sinks from remaining arguments
	sinks := args[1:]

	// Get flags
	include, _ := cmd.Flags().GetStringSlice("include")
	exclude, _ := cmd.Flags().GetStringSlice("exclude")
	priority, _ := cmd.Flags().GetInt("priority")

	// Validate priority
	if priority < 1 {
		return fmt.Errorf("priority must be a positive integer (got %d)", priority)
	}

	// Install ruleset
	req := &arm.InstallRulesetRequest{
		Registry: ruleset.Registry,
		Ruleset:  ruleset.Name,
		Version:  ruleset.Version,
		Priority: priority,
		Include:  include,
		Exclude:  exclude,
		Sinks:    sinks,
	}

	return armService.InstallRuleset(ctx, req)
}

func runInstallPromptset(cmd *cobra.Command, args []string) error {
	// Parse promptset argument
	_, err := ParsePackageArg(args[0])
	if err != nil {
		return fmt.Errorf("failed to parse promptset: %w", err)
	}

	// Get sinks from remaining arguments
	_ = args[1:]

	// Get flags
	_, _ = cmd.Flags().GetStringSlice("include")
	_, _ = cmd.Flags().GetStringSlice("exclude")

	// TODO: Implement promptset installation when service interface is updated
	return fmt.Errorf("promptset installation not yet implemented - service interface needs to be updated first")
}
