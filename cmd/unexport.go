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

func unexportUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Unexport a program from the container.

Usage:
  apx unexport <program>

Examples:
  apx unexport <program>
`)
	return nil
}

func NewUnexportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unexport",
		Short: "Export a program from the container.",
		RunE:  unexport,
	}
	cmd.SetUsageFunc(unexportUsage)
	return cmd
}

func unexport(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	}

	return core.RemoveDesktopEntry(container, args[0])
}
