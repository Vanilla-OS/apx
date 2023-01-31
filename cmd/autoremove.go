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
	"github.com/vanilla-os/apx/core"
)

func NewAutoRemoveCommand() *cobra.Command {
	cmd := &cobra.Command{
		Example: "apx autoremove",
		Use:     "autoremove",
		Short:   "Remove all unused packages automatically",
		RunE:    autoRemove,
	}

	cmd.Flags().BoolP("all", "a", false, "Apply for all containers.")
	return cmd
}

func autoRemove(cmd *cobra.Command, args []string) error {
	if cmd.Flag("all").Changed {
		if err := core.ApplyForAll("autoremove", []string{}); err != nil {
			return err
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("autoremove")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
