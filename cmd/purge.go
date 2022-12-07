package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around apt to make it work inside a container
	with support to installing packages from multiple sources without altering the root filesystem.
*/

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vanilla-os/apx/core"
)

func purgeUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Purge packages inside a managed container.

Usage:
  apx purge <packages>
  apx --aur purge <packages>
  apx --dnf purge <packages>

Examples:
  apx purge htop
  apx --aur purge htop
  apx --dnf purge htop
`)
	return nil
}

func NewPurgeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "purge",
		Short: "Purge packages inside a managed container",
		RunE:  purge,
	}
	cmd.SetUsageFunc(purgeUsage)
	return cmd
}

func purge(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "purge")...)
	command = append(command, args...)

	core.RunContainer(container, command...)

	return nil
}
