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

func showUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Show details about a package.

Usage:
  apx show <package> [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx show htop
`)
	return nil
}

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show",
		Short: "Show details about a package",
		RunE:  show,
	}
	cmd.SetUsageFunc(showUsage)
	return cmd
}

func show(cmd *cobra.Command, args []string) error {

	command := append([]string{}, core.GetPkgCommand(container, "show")...)
	command = append(command, args...)

	core.RunContainer(container, command...)

	return nil
}
