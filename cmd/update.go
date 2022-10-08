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

func updateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Update list of available packages.

Usage:
  apx update

Examples:
  apx update
`)
	return nil
}

func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update list of available packages.",
		RunE:  update,
	}
	cmd.SetUsageFunc(updateUsage)
	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	command := append([]string{"sudo", "apt", "update"}, args...)
	if cmd.Flag("sys").Value.String() == "true" {
		core.AlmostRun(command...)
	} else {
		core.RunContainer(command...)
	}

	return nil
}
