package cmd

/*	License: GPLv3
	Authors:
		Mirko Brombin <send@mirko.pm>
		Pietro di Caprio <pietro@fabricators.ltd>
	Copyright: 2023
	Description: Apx is a wrapper around multiple package managers to install packages and run commands inside a managed container.
*/

import (
	"github.com/spf13/cobra"
)

func NewRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx remove htop",
		Use:     "remove <packages>",
		Short:   "Remove packages inside a managed container.",
		RunE:    remove,
	}
	cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")

	return cmd
}

func remove(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("remove")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	err := container.Run(command...)
	if err != nil {
		return err
	}

	for _, pkg := range args {
		container.RemoveDesktopEntry(pkg)
	}

	return nil
}
