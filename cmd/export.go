package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it work inside a container
	with support to installing packages from multiple sources without altering the root filesystem.
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
  apx export <program>
  apx export --aur <program>
  apx export --dnf <program>

Examples:
  apx export htop
  apx export --aur htop
  apx export --dnf firefox
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
