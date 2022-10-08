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

func autoRemoveUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Remove automatically all unused packages.

Usage:
  apx autoremove

Examples:
  apx autoremove
`)
	return nil
}

func NewAutoRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "autoremove",
		Short: "Remove automatically all unused packages",
		RunE:  autoRemove,
	}
	cmd.SetUsageFunc(autoRemoveUsage)
	return cmd
}

func autoRemove(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	command := append([]string{}, core.GetPkgManager(sys)...)
	command = append(command, "autoremove")
	command = append(command, args...)

	if sys {
		core.AlmostRun(command...)
		return nil
	}

	core.RunContainer(command...)
	return nil
}
