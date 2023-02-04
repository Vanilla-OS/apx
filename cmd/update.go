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
	"github.com/vanilla-os/orchid/cmdr"
)

func NewUpdateCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("update", apx.Trans("update.long"), apx.Trans("update.short"), update)
	/*
			Example: "apx update",
			Use:     "update",
			Short:   "Update the list of available packages",
			RunE:    update,
		}
		cmd.Flags().BoolP("all", "a", false, "Apply for all containers.")
		cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	*/
	return cmd
}

func update(cmd *cobra.Command, args []string) error {
	if cmd.Flag("all").Changed {
		var flags []string
		if cmd.Flag("assume-yes").Value.String() == "true" {
			flags = append(flags, "-y")
		}

		if err := core.ApplyForAll("update", flags); err != nil {
			return err
		}

		return nil
	}

	command := append([]string{}, container.GetPkgCommand("update")...)
	command = append(command, args...)

	if cmd.Flag("assume-yes").Value.String() == "true" {
		command = append(command, "-y")
	}

	container.Run(command...)

	return nil
}
