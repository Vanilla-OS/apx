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

func unexportUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Unexport/Remove a program's desktop entry from a managed container.

Usage:
  apx unexport <program>
  apx unexport --aur <program>
  apx unexport --dnf <program>

Examples:
  apx unexport htop
  apx unexport --aur htop
  apx unexport --dnf htop
`)
	return nil
}

func NewUnexportCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "unexport",
		Short: "Unexport/Remove a program's desktop entry from a managed container",
		RunE:  unexport,
	}
	cmd.SetUsageFunc(unexportUsage)
	return cmd
}

func unexport(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	return core.RemoveDesktopEntry(container, args[0])
}
