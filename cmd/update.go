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

func updateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Update the list of available packages.

Usage:
  apx update
  apx --aur update
  apx --dnf update
`)
	return nil
}

func NewUpdateCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update the list of available packages",
		RunE:  update,
	}
	cmd.SetUsageFunc(updateUsage)
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "update")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	core.RunContainer(container, command...)

	return nil
}
