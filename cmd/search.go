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

func searchUsage(*cobra.Command) error {
	fmt.Print(`Description: 
Search for packages in a managed container.

Usage:
  apx search <packages>
  apx --aur search <packages>
  apx --dnf search <packages>

Examples:
  apx search htop
  apx --aur search htop
  apx --dnf search dnf
`)
	return nil
}

func NewSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search for packages in a managed container.",
		RunE:  search,
	}
	cmd.SetUsageFunc(searchUsage)
	return cmd
}

func search(cmd *cobra.Command, args []string) error {
	aur := cmd.Flag("aur").Value.String() == "true"
	dnf := cmd.Flag("dnf").Value.String() == "true"

	container := "default"
	if aur {
		container = "aur"
	} else if dnf {
		container = "dnf"
	}

	command := append([]string{}, core.GetPkgCommand(container, "search")...)
	command = append(command, args...)

	core.RunContainer(container, command...)

	return nil
}
