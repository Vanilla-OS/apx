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
	"log"

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
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	return cmd
}

func upgrade(cmd *cobra.Command, args []string) error {
	sys := cmd.Flag("sys").Value.String() == "true"
	aur := cmd.Flag("aur").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	}

	command := append([]string{}, core.GetPkgCommand(sys, container, "upgrade")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	if sys {
		log.Default().Println("Performing operations on the host system.")
		core.AlmostRun(false, command...)
		return nil
	}

	core.RunContainer(container, command...)

	return nil
}
