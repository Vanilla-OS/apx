package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it works inside a container
	from outside, directly on the host.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func exportUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Export a program from the container.

Usage:
  apx export <program>

Examples:
  apx export <program>
`)
	return nil
}

func NewExportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export a program from the container.",
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
