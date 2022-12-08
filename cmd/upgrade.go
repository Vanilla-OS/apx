package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func upgrasdeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Upgrade the system by installing/upgrading available packages.

Usage:
  apx upgrade
  apx --aur upgrade
  apx --dnf upgrade
`)
	return nil
}

func NewUpgradeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade the system by installing/upgrading available packages.",
		RunE:  upgrade,
	}
	cmd.SetUsageFunc(upgrasdeUsage)
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	return cmd
}

func upgrade(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "upgrade")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	core.RunContainer(container, command...)

	return nil
}
