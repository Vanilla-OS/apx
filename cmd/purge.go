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
	"github.com/vanilla-os/orchid/cmdr"
)

func NewPurgeCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("purge",
		apx.Trans("purge.long"),
		apx.Trans("purge.short"),
		purge)
	/*
			Example: "apx purge htop",
			Use:     "purge <packages>",
			Short:   "Purge packages inside a managed container",
			RunE:    purge,
		}
	*/
	cmd.Example = "apx purge htop"
	cmd.Args = cobra.MinimumNArgs(1)
	return cmd
}

func purge(cmd *cobra.Command, args []string) error {

	command := append([]string{}, container.GetPkgCommand("purge")...)
	command = append(command, args...)

	container.Run(command...)

	return nil
}
