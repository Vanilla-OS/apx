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
)

func updateUsage(*cobra.Command) error {
	fmt.Print(`Description: 
	Update the list of available packages.

Usage:
  apx update [options]

Options:
  -h, --help            Show this help message and exit

Examples:
  apx update
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

	command := append([]string{}, container.GetPkgCommand("update")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	container.Run(command...)

	return nil
}
