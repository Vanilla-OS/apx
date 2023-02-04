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

func NewRemoveCommand() *cmdr.Command {
	cmd := cmdr.NewCommand("remove <packages>",
		apx.Trans("remove.long"),
		apx.Trans("remove.short"),
		remove).WithBoolFlag(
		cmdr.NewBoolFlag(
			"assume-yes",
			"y",
			apx.Trans("remove.assumeYes"),
			false,
		))
	/*
			Example: "apx remove htop",
			Use:     "remove <packages>",
			Short:   "Remove packages inside a managed container.",
			RunE:    remove,
		}
		cmd.Flags().BoolP("assume-yes", "y", false, "Proceed without manual confirmation.")
	*/
	cmd.Example = "apx remove htop"
	cmd.Args = cobra.MinimumNArgs(1)
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
		container.RemoveBinary(pkg, true)
	}

	return nil
}
