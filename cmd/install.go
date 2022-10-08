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

func installUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Install packages.

Usage:
  apx install <packages>

Examples:
  apx install htop git
`)
	return nil
}

func NewInstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install",
		Short: "Install packages",
		RunE:  install,
	}
	cmd.SetUsageFunc(installUsage)
	return cmd
}

func install(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	command := append([]string{}, core.GetPkgManager(sys)...)
	command = append(command, "install")
	command = append(command, args...)

	if sys {
		core.AlmostRun(command...)
		return nil
	}

	core.RunContainer(command...)
	return nil
}
