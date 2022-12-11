package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func exportUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Export/Recreate a program's desktop entry from a managed container.

Usage:
  apx export <program> [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx export htop
`)
	return nil
}

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export/Recreate a program's desktop entry from a managed container",
		RunE:  export,
	}
	cmd.SetUsageFunc(exportUsage)
	return cmd
}

func export(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	return core.ExportDesktopEntry(container, args[0])
}
