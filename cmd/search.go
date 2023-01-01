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

func searchUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Search for packages in a managed container.

Usage:
  apx search <packages> [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx search htop
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

	command := append([]string{}, core.GetPkgCommand(container, "search")...)
	command = append(command, args...)

	core.RunContainer(container, command...)

	return nil
}
