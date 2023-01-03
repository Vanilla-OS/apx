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

func NewPurgeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx purge htop",
		Use:     "purge <packages>",
		Short:   "Purge packages inside a managed container",
		RunE:    purge,
	}
	return cmd
}

func purge(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("purge")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
