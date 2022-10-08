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

func upgrasdeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Upgrade the system by installing/upgrading packages.

Usage:
  apx upgrade

Examples:
  apx upgrade
`)
	return nil
}

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the system by installing/upgrading packages.",
		RunE:  upgrade,
	}
	cmd.SetUsageFunc(upgrasdeUsage)
	return cmd
}

func upgrade(cmd *cobra.Command, args []string) error {
	command := append([]string{"sudo", "apt", "upgrade"}, args...)
	if cmd.Flag("sys").Value.String() == "true" {
		core.AlmostRun(command...)
	} else {
		core.RunContainer(command...)
	}

	return nil
}
