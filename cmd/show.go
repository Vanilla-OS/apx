package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2022
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"github.com/spf13/cobra"
)

func NewShowCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx show htop",
		Use:     "show <package>",
		Short:   "Show details about a package",
		RunE:    show,
	}
	return cmd
}

func show(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("show")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
