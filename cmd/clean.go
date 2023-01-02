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

func NewCleanCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx clean",
		Use:     "clean",
		Short:   "Clean the apx package manager cache",
		RunE:    clean,
	}
	return cmd
}

func clean(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("clean")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
