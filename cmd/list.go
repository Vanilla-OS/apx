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

func listUsage(*cobra.Command) error {
	fmt.Print(`Description: 
List installed packages.

Usage:
  apx list

Examples:
  apx list
`)
	return nil
}

func NewListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List installed packages.",
		RunE:  list,
	}
	cmd.SetUsageFunc(listUsage)
	cmd.Flags().SetInterspersed(false)
	cmd.Flags().BoolP("upgradable", "u", false, "List only upgradable packages.")
	return cmd
}

func list(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "list")...)

	if cmd.Flag("upgradable").Value.String() == "true" {
		command = append(command, "--upgradable")
	}

	command = append(command, args...)

	core.RunContainer(container, command...)

	return nil
}
