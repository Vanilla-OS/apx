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

func cleanUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Clean the package manager cache.

Usage:
  apx clean

Examples:
  apx clean
`)
	return nil
}

func NewCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Clean the package manager cache",
		RunE:  clean,
	}
	cmd.SetUsageFunc(cleanUsage)
	return cmd
}

func clean(cmd *cobra.Command, args []string) error {
	command := append([]string{"sudo", "apt", "clean"}, args...)
	if cmd.Flag("sys").Value.String() == "true" {
		core.AlmostRun(command...)
	} else {
		core.RunContainer(command...)
	}

	return nil
}
