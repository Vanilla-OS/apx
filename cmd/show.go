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

func showUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Show package details.

Usage:
  apx show <package>

Examples:
  apx show htop
`)
	return nil
}

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show package details.",
		RunE:  show,
	}
	cmd.SetUsageFunc(showUsage)
	return cmd
}

func show(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	command := append([]string{}, core.GetPkgManager(sys)...)
	command = append(command, "show")
	command = append(command, args...)

	if sys {
		core.AlmostRun(command...)
		return nil
	}

	core.RunContainer(command...)
	return nil
}
