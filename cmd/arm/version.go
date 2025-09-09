package main

import (
	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show version",
		RunE:  runVersion,
	}
}

func runVersion(cmd *cobra.Command, args []string) error {
	version := armService.Version()
	FormatVersionInfo(version)
	return nil
}
