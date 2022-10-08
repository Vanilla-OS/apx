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
	return cmd
}

func list(cmd *cobra.Command, args []string) error {
	command := append([]string{"sudo", "apt", "list"}, args...)
	if cmd.Flag("sys").Value.String() == "true" {
		core.AlmostRun(command...)
	} else {
		core.RunContainer(command...)
	}

	return nil
}
